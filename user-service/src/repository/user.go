package repository

import (
	"context"
	"user-service/src/common/model"

	"gorm.io/gorm"
)

type User interface {
	Find(ctx context.Context, whereClause string, args []any) (*model.User, error)
	Create(ctx context.Context, user *model.User, password string) error 
}

type userImpl struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	return &userImpl{db: db}
}

func (r *userImpl) Find(ctx context.Context, whereClause string, args []any) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where(whereClause, args...).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if user.ID == "" {
		return nil, nil
	}

	return &user, nil
}

func (r *userImpl) Create(ctx context.Context, user *model.User, password string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		err := tx.Create(&model.Credential{
			UserID:   user.ID,
			Password: password,
		}).Error

		return err
	})

	return err
}
