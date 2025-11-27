package repository

import (
	"context"
	"user-service/src/common/model"

	"gorm.io/gorm"
)

type Credential interface {
	Find(ctx context.Context, whereClause string, args []any) (*model.Credential, error)
}

type credentialImpl struct {
	db *gorm.DB
}

func NewCredential(db *gorm.DB) Credential {
	return &credentialImpl{db: db}
}

func (r *credentialImpl) Find(ctx context.Context, whereClause string, args []any) (*model.Credential, error) {
	var credential model.Credential

	if err := r.db.WithContext(ctx).Table("credentials AS c").Select("c.*").Joins("JOIN users AS u ON u.id = c.user_id").Where(whereClause, args...).Scan(&credential).Error; err != nil {
		return nil, err
	}

	if credential.UserID == "" {
		return nil, nil
	}

	return &credential, nil
}
