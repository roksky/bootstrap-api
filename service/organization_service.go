package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/roksky/bootstrap-api/data/response"
	"github.com/roksky/bootstrap-api/model"
	"github.com/roksky/bootstrap-api/repository"
)

type OrganizationService struct {
	repository repository.BaseRepository[model.Organization, uuid.UUID, repository.OrganizationSearch]
	Validate   *validator.Validate
}

func NewOrganizationService(repository repository.BaseRepository[model.Organization, uuid.UUID, repository.OrganizationSearch], validate *validator.Validate) BaseService[model.Organization, uuid.UUID, repository.OrganizationSearch] {
	return &OrganizationService{
		repository: repository,
		Validate:   validate,
	}
}

func (e *OrganizationService) Create(filterContext *repository.OrganizationSearch, item *model.Organization) (*model.Organization, error) {
	err := e.Validate.Struct(item)
	if err != nil {
		return nil, err
	}
	return e.repository.Save(nil, filterContext, item)
}

func (e *OrganizationService) CreateMany(filterContext *repository.OrganizationSearch, items []*model.Organization) ([]*model.Organization, error) {
	for _, item := range items {
		err := e.Validate.Struct(item)
		if err != nil {
			return nil, err
		}
	}

	return e.repository.SaveMany(nil, filterContext, items)
}

func (e *OrganizationService) Update(filterContext *repository.OrganizationSearch, item *model.Organization) (*model.Organization, error) {
	err := e.Validate.Struct(item)
	if err != nil {
		return nil, err
	}
	if item.Id == uuid.Nil {
		return nil, errors.New("entity id is missing")
	}
	return e.repository.Update(nil, filterContext, item)
}

func (e *OrganizationService) UpdateMany(filterContext *repository.OrganizationSearch, items []*model.Organization) ([]*model.Organization, error) {
	for _, item := range items {
		err := e.Validate.Struct(item)
		if err != nil {
			return nil, err
		}
		if item.Id == uuid.Nil {
			return nil, errors.New("entity id is missing")
		}
	}

	return e.repository.UpdateMany(nil, filterContext, items)
}

func (e *OrganizationService) Delete(filterContext *repository.OrganizationSearch, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("entity id is missing")
	}

	return e.repository.Delete(nil, filterContext, id)
}

func (e *OrganizationService) DeleteMany(filterContext *repository.OrganizationSearch, ids []uuid.UUID) error {
	for _, id := range ids {
		if id == uuid.Nil {
			return errors.New("invalid id")
		}
	}

	return e.repository.DeleteByIds(nil, filterContext, ids)
}

func (e *OrganizationService) FindById(filterContext *repository.OrganizationSearch, id uuid.UUID) (*model.Organization, error) {
	if id == uuid.Nil {
		return nil, errors.New("entity id is missing")
	}
	return e.repository.FindById(nil, filterContext, id)
}

func (e *OrganizationService) FindByIds(filterContext *repository.OrganizationSearch, ids []uuid.UUID) ([]*model.Organization, error) {
	for _, id := range ids {
		if id == uuid.Nil {
			return nil, errors.New("invalid id")
		}
	}

	return e.repository.FindByIds(nil, filterContext, ids)
}

func (e *OrganizationService) FindAll(filterContext *repository.OrganizationSearch, pageSize int, page int) (response.PagedResult[*model.Organization], error) {
	var result response.PagedResult[*model.Organization]

	items, err := e.repository.FindAll(nil, filterContext, pageSize, page)
	if err != nil {
		return result, err
	}

	count, err := e.repository.Count(nil, filterContext)

	result.Items = items
	result.TotalItems = count
	result.PageSize = count
	result.PageNumber = 0

	return result, nil
}

func (e *OrganizationService) Search(searchParams *repository.OrganizationSearch) (response.PagedResult[*model.Organization], error) {
	var result response.PagedResult[*model.Organization]

	items, err := e.repository.Search(nil, searchParams)
	if err != nil {
		return result, err
	}

	count, err := e.repository.Count(nil, searchParams)

	result.Items = items
	result.TotalItems = count
	result.PageSize = int64(searchParams.PageSize)
	result.PageNumber = searchParams.PageNumber

	return result, nil
}

func (e *OrganizationService) Deleted(searchParams *repository.OrganizationSearch) ([]string, error) {
	return e.repository.Deleted(nil, searchParams)
}
