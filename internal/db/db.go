package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"go-server/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(conf config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DBHost,
		conf.DBPort,
		conf.DBUser,
		conf.DBPassword,
		conf.DBDatabase)

	newLogger := logger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Disconnect(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func Migrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}
