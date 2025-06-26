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

func (e *SystemUserOrganizationRepository) Save(filterContext *SystemUserOrganizationSearch, item *model.SystemUserOrganization) (*model.SystemUserOrganization, error) {
	result := e.Db.Create(item)
	return item, result.Error
}

func (e *SystemUserOrganizationRepository) SaveMany(filterContext *SystemUserOrganizationSearch, item []*model.SystemUserOrganization) ([]*model.SystemUserOrganization, error) {
	result := e.Db.Create(item)
	return item, result.Error
}

func (e *SystemUserOrganizationRepository) Update(filterContext *SystemUserOrganizationSearch, item *model.SystemUserOrganization) (*model.SystemUserOrganization, error) {
	result := e.Db.Model(item).Updates(item)
	if result.Error == nil {
		return e.FindById(filterContext, item.Id)
	} else {
		return nil, result.Error
	}
}

func (e *SystemUserOrganizationRepository) UpdateMany(filterContext *SystemUserOrganizationSearch, items []*model.SystemUserOrganization) ([]*model.SystemUserOrganization, error) {
	itemIds := make([]uuid.UUID, len(items))
	for _, item := range items {
		result := e.Db.Model(item).Where("id = ?", item.Id).Updates(item)
		itemIds = append(itemIds, item.Id)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return e.FindByIds(filterContext, itemIds)
}

func (e *SystemUserOrganizationRepository) Delete(filterContext *SystemUserOrganizationSearch, itemId uuid.UUID) error {
	var SystemUserOrganization model.SystemUserOrganization
	result := e.Db.Where("id = ?", itemId).Delete(&SystemUserOrganization)
	return result.Error
}

func (e *SystemUserOrganizationRepository) DeleteByIds(filterContext *SystemUserOrganizationSearch, itemIds []uuid.UUID) error {
	var SystemUserOrganization model.SystemUserOrganization
	result := e.Db.Where("id IN ?", itemIds).Delete(&SystemUserOrganization)
	return result.Error
}

func (e *SystemUserOrganizationRepository) FindById(filterContext *SystemUserOrganizationSearch, itemId uuid.UUID) (*model.SystemUserOrganization, error) {
	var SystemUserOrganization model.SystemUserOrganization
	result := e.Db.Find(&SystemUserOrganization, itemId)
	if result.Error == nil {
		return &SystemUserOrganization, nil
	} else {
		return nil, errors.New("SystemUserOrganization is not found")
	}
}

func (e *SystemUserOrganizationRepository) FindByIds(filterContext *SystemUserOrganizationSearch, itemIds []uuid.UUID) ([]*model.SystemUserOrganization, error) {
	var entities []model.SystemUserOrganization
	result := e.Db.Where("id IN ?", itemIds).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return helper.ConvertSliceToReference(entities), nil
	}
}

func (e *SystemUserOrganizationRepository) FindAll(filterContext *SystemUserOrganizationSearch, pageSize int, page int) ([]*model.SystemUserOrganization, error) {
	var entities []*model.SystemUserOrganization
	result := e.Db.Limit(pageSize).Offset(page * pageSize).Find(&entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}

func (e *SystemUserOrganizationRepository) Search(searchParams *SystemUserOrganizationSearch) ([]*model.SystemUserOrganization, error) {
	var entities []*model.SystemUserOrganization

	tx := e.Db.Joins("SystemUser")
	if searchParams.OrganizationId != uuid.Nil {
		tx = tx.Where("organization = ?", searchParams.OrganizationId)
	}
	if searchParams.SystemUser != "" {
		tx = tx.Where("system_user = ?", searchParams.SystemUser)
	}
	if searchParams.OrderBy != "" {
		tx = tx.Order(searchParams.OrderBy)
	}

	result := tx.Limit(searchParams.PageSize).Offset(searchParams.PageNumber * searchParams.PageSize).Find(&entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}

func (e *SystemUserOrganizationRepository) Count(searchParams *SystemUserOrganizationSearch) (int64, error) {
	var count int64

	tx := e.Db
	if searchParams.OrganizationId != uuid.Nil {
		tx = tx.Where("organization = ?", searchParams.OrganizationId)
	}
	if searchParams.SystemUser != "" {
		tx = tx.Where("system_user = ?", searchParams.SystemUser)
	}

	result := tx.Model(&model.SystemUserOrganization{}).Count(&count)
	if result.Error != nil {
		return count, result.Error
	} else {
		return count, nil
	}
}

func (e *SystemUserOrganizationRepository) Deleted(searchParams *SystemUserOrganizationSearch) ([]string, error) {
	var entities []string

	tx := e.Db.Unscoped().Model(&model.SystemUserOrganization{}).Where("date_deleted IS NOT NULL").Where("organization_id = ?", searchParams.OrganizationId)
	if searchParams.OrganizationId != uuid.Nil {
		tx = tx.Where("organization = ?", searchParams.OrganizationId)
	}
	if searchParams.SystemUser != "" {
		tx = tx.Where("system_user = ?", searchParams.SystemUser)
	}

	result := tx.Limit(searchParams.PageSize).Pluck("id", &entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}
