package auth

import (
	"github.com/roksky/bootstrap-api/model"
	"github.com/roksky/bootstrap-api/repository"
	"gorm.io/gorm"
)

type Database struct {
	db         *gorm.DB
	tokenRepo  *TokenRepo
	clientRepo *ClientRepo
	userRepo   *repository.SystemUserRepo
}

func (m *Database) AutoMigrate() error {
	return m.db.AutoMigrate(getModels()...)
}

func (m *Database) Database() *gorm.DB {
	return m.db
}

func NewAuthDatabase(db *gorm.DB) *Database {
	return &Database{
		db:         db,
		tokenRepo:  NewTokenRepo(db),
		clientRepo: NewClientRepo(db),
		userRepo:   repository.NewSystemUserRepo(db),
	}
}

func (m *Database) GetTokenRepo() *TokenRepo {
	return m.tokenRepo
}

func (m *Database) GetClientRepo() *ClientRepo {
	return m.clientRepo
}

func (m *Database) GetUserRepo() *repository.SystemUserRepo {
	return m.userRepo
}

// getModels returns a list of all the models to be migrated
// Add a list of all the models to be auto migrated here
func getModels() []any {
	return []any{
		&DbToken{},
		&Client{},
		&model.SystemUser{},
	}
}
