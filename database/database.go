package database

import (
	"github.com/roksky/bootstrap-api/model"
	"gorm.io/gorm"
)

type Database interface {
	AutoMigrate() error
	Database() *gorm.DB
}

type MyDb struct {
	db *gorm.DB
}

func (m *MyDb) AutoMigrate() error {
	return m.db.AutoMigrate(getModels()...)
}

func (m *MyDb) Database() *gorm.DB {
	return m.db
}

func NewDatabase(db *gorm.DB) Database {
	return &MyDb{
		db: db,
	}
}

// getModels returns a list of all the models to be migrated
// Add a list of all the models to be auto migrated here
func getModels() []any {
	return []any{
		&model.Organization{},
	}
}
