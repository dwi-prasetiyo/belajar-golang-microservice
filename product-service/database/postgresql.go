package database

import (
	"product-service/env"
	"product-service/src/common/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgreSQL() *gorm.DB {
	gormDB, err := gorm.Open(postgres.Open(env.Conf.Postgres.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	return gormDB
}

func ClosePostgreSQL(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
