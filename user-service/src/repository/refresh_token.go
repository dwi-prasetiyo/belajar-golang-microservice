package repository

import (
	"context"
	"user-service/src/common/model"

	"gorm.io/gorm"
)

type RefreshToken interface {
	Create(ctx context.Context, data *model.RefreshToken) error
	Delete(ctx context.Context, refreshToken string) error
	Find(ctx context.Context, refreshToken string) (*model.RefreshTokenWithRole, error)
	Update(ctx context.Context, data *model.RefreshToken, oldToken string) error
}

type refreshTokenImpl struct {
	db *gorm.DB
}

func NewRefreshToken(db *gorm.DB) RefreshToken {
	return &refreshTokenImpl{db: db}
}

func (r *refreshTokenImpl) Create(ctx context.Context, data *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *refreshTokenImpl) Delete(ctx context.Context, refreshToken string) error {
	return r.db.WithContext(ctx).Delete(&model.RefreshToken{}, "token = ?", refreshToken).Error
}

func (r *refreshTokenImpl) Find(ctx context.Context, refreshToken string) (*model.RefreshTokenWithRole, error) {
	res := new(model.RefreshTokenWithRole)

	err := r.db.WithContext(ctx).Table("refresh_tokens AS rt").Select("rt.token, rt.user_id, u.role").Joins("INNER JOIN users AS u ON rt.user_id = u.id").Where("rt.token = ?", refreshToken).First(&res).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if res.Token == "" {
		return nil, nil
	}

	return res, nil
}

func (r *refreshTokenImpl) Update(ctx context.Context, data *model.RefreshToken, oldToken string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.RefreshToken{}, "token = ?", oldToken).Error; err != nil {
			return err
		}

		if err := tx.Create(data).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}