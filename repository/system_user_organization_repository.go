package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/roksky/bootstrap-api/helper"
	"github.com/roksky/bootstrap-api/model"
	"gorm.io/gorm"
)

type SystemUserOrganizationRepository struct {
	Db *gorm.DB
}

func NewSystemUserOrganizationRepository(Db *gorm.DB) BaseRepository[model.SystemUserOrganization, uuid.UUID, SystemUserOrganizationSearch] {
	return &SystemUserOrganizationRepository{Db: Db}
}

type SystemUserOrganizationSearch struct {
	OrganizationId uuid.UUID
	SystemUser     string
	PageSize       int
	PageNumber     int
	OrderBy        string
}

func (e *SystemUserOrganizationRepository) GetDB() *gorm.DB {
	return e.Db
}

func (e *SystemUserOrganizationRepository) Save(tx *gorm.DB, filterContext *SystemUserOrganizationSearch, item *model.SystemUserOrganization) (*model.SystemUserOrganization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	result := db.Create(item)
	return item, result.Error
}

func (e *SystemUserOrganizationRepository) SaveMany(tx *gorm.DB, filterContext *SystemUserOrganizationSearch, items []*model.SystemUserOrganization) ([]*model.SystemUserOrganization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	result := db.Create(items)
	return items, result.Error
}

func (e *SystemUserOrganizationRepository) Update(tx *gorm.DB, filterContext *SystemUserOrganizationSearch, item *model.SystemUserOrganization) (*model.SystemUserOrganization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	result := db.Model(item).Updates(item)
	if result.Error == nil {
		return e.FindById(tx, filterContext, item.Id)
	} else {
		return nil, result.Error
	}
}

func (e *SystemUserOrganizationRepository) UpdateMany(tx *gorm.DB, filterContext *SystemUserOrganizationSearch, items []*model.SystemUserOrganization) ([]*model.SystemUserOrganization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	itemIds := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		result := db.Model(item).Where("id = ?", item.Id).Updates(item)
		itemIds = append(itemIds, item.Id)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return e.FindByIds(tx, filterContext, itemIds)
}

func (e *SystemUserOrganizationRepository) Delete(tx *gorm.DB, filterContext *SystemUserOrganizationSearch, itemId uuid.UUID) error {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var SystemUserOrganization model.SystemUserOrganization
	result := db.Where("id = ?", itemId).Delete(&SystemUserOrganization)
	return result.Error
}

func (e *SystemUserOrganizationRepository) DeleteByIds(tx *gorm.DB, filterContext *SystemUserOrganizationSearch, itemIds []uuid.UUID) error {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var SystemUserOrganization model.SystemUserOrganization
	result := db.Where("id IN ?", itemIds).Delete(&SystemUserOrganization)
	return result.Error
}

func (e *SystemUserOrganizationRepository) FindById(tx *gorm.DB, filterContext *SystemUserOrganizationSearch, itemId uuid.UUID) (*model.SystemUserOrganization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var SystemUserOrganization model.SystemUserOrganization
	result := db.Find(&SystemUserOrganization, itemId)
	if result.Error == nil {
		return &SystemUserOrganization, nil
	} else {
		return nil, errors.New("SystemUserOrganization is not found")
	}
}

func (e *SystemUserOrganizationRepository) FindByIds(tx *gorm.DB, filterContext *SystemUserOrganizationSearch, itemIds []uuid.UUID) ([]*model.SystemUserOrganization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var entities []model.SystemUserOrganization
	result := db.Where("id IN ?", itemIds).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return helper.ConvertSliceToReference(entities), nil
	}
}

func (e *SystemUserOrganizationRepository) FindAll(tx *gorm.DB, filterContext *SystemUserOrganizationSearch, pageSize int, page int) ([]*model.SystemUserOrganization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var entities []*model.SystemUserOrganization
	result := db.Limit(pageSize).Offset(page * pageSize).Find(&entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}

func (e *SystemUserOrganizationRepository) Search(tx *gorm.DB, searchParams *SystemUserOrganizationSearch) ([]*model.SystemUserOrganization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var entities []*model.SystemUserOrganization

	tx2 := db.Joins("SystemUser")
	if searchParams.OrganizationId != uuid.Nil {
		tx2 = tx2.Where("organization = ?", searchParams.OrganizationId)
	}
	if searchParams.SystemUser != "" {
		tx2 = tx2.Where("system_user = ?", searchParams.SystemUser)
	}
	if searchParams.OrderBy != "" {
		tx2 = tx2.Order(searchParams.OrderBy)
	}

	result := tx2.Limit(searchParams.PageSize).Offset(searchParams.PageNumber * searchParams.PageSize).Find(&entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}

func (e *SystemUserOrganizationRepository) Count(tx *gorm.DB, searchParams *SystemUserOrganizationSearch) (int64, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var count int64

	tx2 := db
	if searchParams.OrganizationId != uuid.Nil {
		tx2 = tx2.Where("organization = ?", searchParams.OrganizationId)
	}
	if searchParams.SystemUser != "" {
		tx2 = tx2.Where("system_user = ?", searchParams.SystemUser)
	}

	result := tx2.Model(&model.SystemUserOrganization{}).Count(&count)
	if result.Error != nil {
		return count, result.Error
	} else {
		return count, nil
	}
}

func (e *SystemUserOrganizationRepository) Deleted(tx *gorm.DB, searchParams *SystemUserOrganizationSearch) ([]string, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var entities []string

	tx2 := db.Unscoped().Model(&model.SystemUserOrganization{}).Where("date_deleted IS NOT NULL").Where("organization_id = ?", searchParams.OrganizationId)
	if searchParams.OrganizationId != uuid.Nil {
		tx2 = tx2.Where("organization = ?", searchParams.OrganizationId)
	}
	if searchParams.SystemUser != "" {
		tx2 = tx2.Where("system_user = ?", searchParams.SystemUser)
	}

	result := tx2.Limit(searchParams.PageSize).Pluck("id", &entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}
