package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/roksky/bootstrap-api/helper"
	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DatabaseConnection() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	dbConf := EnvConfigs.Database

	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Name)
	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{
		Logger: newLogger,
	})
	helper.ErrorPanic(err)

	return db
}
