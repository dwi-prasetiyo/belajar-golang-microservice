package repository

import (
	"context"
	"user-service/src/common/model"

	"gorm.io/gorm"
)

type UserBlockLog interface {
	Find(ctx context.Context, userID string) (*model.UserBlockLog, error)
}

type userBlockLogImpl struct {
	db *gorm.DB
}

func NewUserBlockLog(db *gorm.DB) UserBlockLog {
	return &userBlockLogImpl{db: db}
}

func (r *userBlockLogImpl) Find(ctx context.Context, userID string) (*model.UserBlockLog, error) {
	var userBlockLog model.UserBlockLog

	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).First(&userBlockLog).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if userBlockLog.ID == 0 {
		return nil, nil
	}

	return &userBlockLog, nil
}
