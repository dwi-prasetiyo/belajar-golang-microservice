package repository

import (
	"context"
	"net/http"
	"product-service/src/common/dto/request"
	"product-service/src/common/errors"
	"product-service/src/common/model"
	"product-service/src/common/util"
	"strings"

	pb "github.com/dwi-prasetiyo/protobuf/protogen/product"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product interface {
	Create(ctx context.Context, req *model.Product) error
	FindMany(ctx context.Context, req *request.FindManyProduct) ([]*model.Product, error)
	ReduceStocks(ctx context.Context, req []*pb.ProductOrder) error
	RollbackStocks(ctx context.Context, req []*pb.ProductOrder) error
}

type productImpl struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) Product {
	return &productImpl{db: db}
}

func (r *productImpl) Create(ctx context.Context, req *model.Product) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *productImpl) FindMany(ctx context.Context, req *request.FindManyProduct) ([]*model.Product, error) {
	limit, offset := util.CreateLimitAndOffset(req.Page)

	var args []any

	query := "SELECT * FROM products "
	if req.Search != "" {
		query += "WHERE to_tsvector('indonesian', description) @@ to_tsquery('indonesian', ?) "
		args = append(args, strings.Join(strings.Fields(req.Search), " & "))
	}

	query += "LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var res []*model.Product

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&res).Error
	return res, err
}

func (r *productImpl) ReduceStocks(ctx context.Context, req []*pb.ProductOrder) error {
	var productIDs []int
	existIDs := make(map[int]bool)
	
	for _, p := range req {
		if existIDs[int(p.ProductId)] {
			continue
		}

		productIDs = append(productIDs, int(p.ProductId))
		existIDs[int(p.ProductId)] = true
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var products []*model.Product
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&products, "id IN (?)", productIDs).Error
		if err != nil {
			return err
		}

		for _, p := range req {
			result := tx.Model(&model.Product{}).
				Where("id = ? AND stock >= ?", p.ProductId, p.Quantity).
				Update("stock", gorm.Expr("stock - ?", p.Quantity))

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return &errors.Response{
					Message:  "Product not found or stock is not enough",
					HttpCode: http.StatusBadRequest,
					GrpcCode: codes.FailedPrecondition,
				}
			}
		}
		return nil
	})
}

func (r *productImpl) RollbackStocks(ctx context.Context, req []*pb.ProductOrder) error {
	var productIDs []int
	existIDs := make(map[int]bool)
	
	for _, p := range req {
		if existIDs[int(p.ProductId)] {
			continue
		}

		productIDs = append(productIDs, int(p.ProductId))
		existIDs[int(p.ProductId)] = true
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var products []*model.Product
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&products, "id IN (?)", productIDs).Error
		if err != nil {
			return err
		}

		for _, p := range req {
			result := tx.Model(&model.Product{}).
				Where("id = ?", p.ProductId).
				Update("stock", gorm.Expr("stock + ?", p.Quantity))

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return &errors.Response{
					Message:  "Product not found or stock is not enough",
					HttpCode: http.StatusBadRequest,
					GrpcCode: codes.FailedPrecondition,
				}
			}
		}
		return nil
	})
}
