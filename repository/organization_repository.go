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

func (e *OrganizationRepository) Save(tx *gorm.DB, filterContext *OrganizationSearch, item *model.Organization) (*model.Organization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	result := db.Create(item)
	return item, result.Error
}

func (e *OrganizationRepository) SaveMany(tx *gorm.DB, filterContext *OrganizationSearch, item []*model.Organization) ([]*model.Organization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	result := db.Create(item)
	return item, result.Error
}

func (e *OrganizationRepository) Update(tx *gorm.DB, filterContext *OrganizationSearch, item *model.Organization) (*model.Organization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	result := db.Model(item).Updates(item)
	if result.Error == nil {
		return e.FindById(db, filterContext, item.Id)
	} else {
		return nil, result.Error
	}
}

func (e *OrganizationRepository) UpdateMany(tx *gorm.DB, filterContext *OrganizationSearch, items []*model.Organization) ([]*model.Organization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	itemIds := make([]uuid.UUID, len(items))
	for _, item := range items {
		result := db.Model(item).Where("id = ?", item.Id).Updates(item)
		itemIds = append(itemIds, item.Id)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return e.FindByIds(db, filterContext, itemIds)
}

func (e *OrganizationRepository) Delete(tx *gorm.DB, filterContext *OrganizationSearch, itemId uuid.UUID) error {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var organization model.Organization
	result := db.Where("id = ?", itemId).Delete(&organization)
	return result.Error
}

func (e *OrganizationRepository) DeleteByIds(tx *gorm.DB, filterContext *OrganizationSearch, itemIds []uuid.UUID) error {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var organization model.Organization
	result := db.Where("id IN ?", itemIds).Delete(&organization)
	return result.Error
}

func (e *OrganizationRepository) FindById(tx *gorm.DB, filterContext *OrganizationSearch, itemId uuid.UUID) (*model.Organization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var organization model.Organization
	result := db.Find(&organization, itemId)
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

func (e *OrganizationRepository) FindByIds(tx *gorm.DB, filterContext *OrganizationSearch, itemIds []uuid.UUID) ([]*model.Organization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var entities []model.Organization
	result := db.Where("id IN ?", itemIds).Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return helper.ConvertSliceToReference(entities), nil
	}
}

func (e *OrganizationRepository) FindAll(tx *gorm.DB, filterContext *OrganizationSearch, pageSize int, page int) ([]*model.Organization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var entities []*model.Organization
	result := db.Limit(pageSize).Offset(page * pageSize).Find(&entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}

func (e *OrganizationRepository) Search(tx *gorm.DB, searchParams *OrganizationSearch) ([]*model.Organization, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var entities []*model.Organization

	if searchParams.OrganizationType != "" {
		db = db.Where("organization_type = ?", searchParams.OrganizationType)
	}
	if searchParams.OrderBy != "" {
		db = db.Order(searchParams.OrderBy)
	}

	result := db.Limit(searchParams.PageSize).Offset(searchParams.PageNumber * searchParams.PageSize).Find(&entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}

func (e *OrganizationRepository) Count(tx *gorm.DB, searchParams *OrganizationSearch) (int64, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var count int64

	if searchParams.OrganizationType != "" {
		db = db.Where("organization_type = ?", searchParams.OrganizationType)
	}
	result := db.Model(&model.Organization{}).Count(&count)
	if result.Error != nil {
		return count, result.Error
	} else {
		return count, nil
	}
}

func (e *OrganizationRepository) Deleted(tx *gorm.DB, searchParams *OrganizationSearch) ([]string, error) {
	db := e.Db
	if tx != nil {
		db = tx
	}
	var entities []string

	db = db.Unscoped().Model(&model.Organization{}).Where("date_deleted IS NOT NULL")
	if searchParams.OrganizationType != "" {
		db = db.Where("organization_type = ?", searchParams.OrganizationType)
	}

	result := db.Limit(searchParams.PageSize).Pluck("id", &entities)
	helper.ErrorPanic(result.Error)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return entities, nil
	}
}
