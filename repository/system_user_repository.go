package repository

import (
	"github.com/google/uuid"
	"github.com/roksky/bootstrap-api/model"
	"gorm.io/gorm"
)

type SystemUserRepo struct {
	db *gorm.DB
}

func NewSystemUserRepo(db *gorm.DB) *SystemUserRepo {
	return &SystemUserRepo{
		db: db,
	}
}

func (m *SystemUserRepo) Save(item *model.SystemUser) error {
	result := m.db.Create(item)
	return result.Error
}

func (m *SystemUserRepo) Delete(item *model.SystemUser) error {
	result := m.db.Delete(item)
	return result.Error
}

func (m *SystemUserRepo) FindByUserName(userName string) (*model.SystemUser, error) {
	var systemUser model.SystemUser
	result := m.db.Where("user_name = ?", userName).Find(&systemUser)
	if result.Error == nil {
		return &systemUser, nil
	} else {
		return nil, result.Error
	}
}

func (m *SystemUserRepo) FindById(userId uuid.UUID) (*model.SystemUser, error) {
	var systemUser model.SystemUser
	result := m.db.Where("user_id = ?", userId).Find(&systemUser)
	if result.Error == nil {
		return &systemUser, nil
	} else {
		return nil, result.Error
	}
}
