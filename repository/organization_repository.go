package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/roksky/bootstrap-api/helper"
	"github.com/roksky/bootstrap-api/model"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	Db *gorm.DB
}

func NewOrganizationRepository(Db *gorm.DB) BaseRepository[model.Organization, uuid.UUID, OrganizationSearch] {
	return &OrganizationRepository{Db: Db}
}

type OrganizationSearch struct {
	OrganizationType string
	PageSize         int
	PageNumber       int
	OrderBy          string
}

func (e *OrganizationRepository) GetDB() *gorm.DB {
	return e.Db
}

func (e *OrganizationRepository) Save(filterContext *OrganizationSearch, item *model.Organization) (*model.Organization, error) {
	result := e.Db.Create(item)
	return item, result.Error
}

func (e *OrganizationRepository) SaveMany(filterContext *OrganizationSearch, item []*model.Organization) ([]*model.Organization, error) {
	result := e.Db.Create(item)
	return item, result.Error
}

func (e *OrganizationRepository) Update(filterContext *OrganizationSearch, item *model.Organization) (*model.Organization, error) {
	result := e.Db.Model(item).Updates(item)
	if result.Error == nil {
		return e.FindById(filterContext, item.Id)
	} else {
		return nil, result.Error
	}
}

func (e *OrganizationRepository) UpdateMany(filterContext *OrganizationSearch, items []*model.Organization) ([]*model.Organization, error) {
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

func (e *OrganizationRepository) Delete(filterContext *OrganizationSearch, itemId uuid.UUID) error {
	var organization model.Organization
	result := e.Db.Where("id = ?", itemId).Delete(&organization)
	return result.Error
}

func (e *OrganizationRepository) DeleteByIds(filterContext *OrganizationSearch, itemIds []uuid.UUID) error {
	var organization model.Organization
	result := e.Db.Where("id IN ?", itemIds).Delete(&organization)
	return result.Error
}

func (e *OrganizationRepository) FindById(filterContext *OrganizationSearch, itemId uuid.UUID) (*model.Organization, error) {
	var organization model.Organization
	result := e.Db.Find(&organization, itemId)
	if result.Error == nil {
		if organization.Id != uuid.Nil {
			return &organization, nil
		} else {
			return nil, errors.New("entity is not found")
		}
	} else {
		return nil, errors.New("organization is not found")
	}
}

func (e *OrganizationRepository) FindByIds(filterContext *OrganizationSearch, itemIds []uuid.UUID) ([]*model.Organization, error) {
	var entities []model.Organization
	result := e.Db.Where("id IN ?", itemIds).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return helper.ConvertSliceToReference(entities), nil
	}
}

func (e *OrganizationRepository) FindAll(filterContext *OrganizationSearch, pageSize int, page int) ([]*model.Organization, error) {
	var entities []*model.Organization
	result := e.Db.Limit(pageSize).Offset(page * pageSize).Find(&entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}

func (e *OrganizationRepository) Search(searchParams *OrganizationSearch) ([]*model.Organization, error) {
	var entities []*model.Organization

	tx := e.Db
	if searchParams.OrganizationType != "" {
		tx = tx.Where("organization_type = ?", searchParams.OrganizationType)
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

func (e *OrganizationRepository) Count(searchParams *OrganizationSearch) (int64, error) {
	var count int64

	tx := e.Db
	if searchParams.OrganizationType != "" {
		tx = tx.Where("organization_type = ?", searchParams.OrganizationType)
	}
	result := tx.Model(&model.Organization{}).Count(&count)
	if result.Error != nil {
		return count, result.Error
	} else {
		return count, nil
	}
}

func (e *OrganizationRepository) Deleted(searchParams *OrganizationSearch) ([]string, error) {
	var entities []string

	tx := e.Db.Unscoped().Model(&model.Organization{}).Where("date_deleted IS NOT NULL")
	if searchParams.OrganizationType != "" {
		tx = tx.Where("organization_type = ?", searchParams.OrganizationType)
	}

	result := tx.Limit(searchParams.PageSize).Pluck("id", &entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}
