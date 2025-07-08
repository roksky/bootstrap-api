package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/roksky/bootstrap-api/data/response"
	"github.com/roksky/bootstrap-api/model"
	"github.com/roksky/bootstrap-api/repository"
)

type SystemUserOrganizationService struct {
	repository repository.BaseRepository[model.SystemUserOrganization, uuid.UUID, repository.SystemUserOrganizationSearch]
	Validate   *validator.Validate
}

func NewSystemUserOrganizationService(repository repository.BaseRepository[model.SystemUserOrganization, uuid.UUID, repository.SystemUserOrganizationSearch], validate *validator.Validate) *SystemUserOrganizationService {
	return &SystemUserOrganizationService{
		repository: repository,
		Validate:   validate,
	}
}

type SystemUserOrganizationSearch struct {
	OrganizationType string
	PageSize         int
	PageNumber       int
}

func (e *SystemUserOrganizationService) Create(filterContext *repository.SystemUserOrganizationSearch, item *model.SystemUserOrganization) (*model.SystemUserOrganization, error) {
	err := e.Validate.Struct(item)
	if err != nil {
		return nil, err
	}
	return e.repository.Save(nil, filterContext, item)
}

func (e *SystemUserOrganizationService) CreateMany(filterContext *repository.SystemUserOrganizationSearch, items []*model.SystemUserOrganization) ([]*model.SystemUserOrganization, error) {
	for _, item := range items {
		err := e.Validate.Struct(item)
		if err != nil {
			return nil, err
		}
	}

	return e.repository.SaveMany(nil, filterContext, items)
}

func (e *SystemUserOrganizationService) Update(filterContext *repository.SystemUserOrganizationSearch, item *model.SystemUserOrganization) (*model.SystemUserOrganization, error) {
	err := e.Validate.Struct(item)
	if err != nil {
		return nil, err
	}
	if item.Id == uuid.Nil {
		return nil, errors.New("entity id is missing")
	}
	return e.repository.Update(nil, filterContext, item)
}

func (e *SystemUserOrganizationService) UpdateMany(filterContext *repository.SystemUserOrganizationSearch, items []*model.SystemUserOrganization) ([]*model.SystemUserOrganization, error) {
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

func (e *SystemUserOrganizationService) Delete(filterContext *repository.SystemUserOrganizationSearch, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("entity id is missing")
	}

	return e.repository.Delete(nil, filterContext, id)
}

func (e *SystemUserOrganizationService) DeleteMany(filterContext *repository.SystemUserOrganizationSearch, ids []uuid.UUID) error {
	for _, id := range ids {
		if id == uuid.Nil {
			return errors.New("invalid id")
		}
	}

	return e.repository.DeleteByIds(nil, filterContext, ids)
}

func (e *SystemUserOrganizationService) FindById(filterContext *repository.SystemUserOrganizationSearch, id uuid.UUID) (*model.SystemUserOrganization, error) {
	if id == uuid.Nil {
		return nil, errors.New("entity id is missing")
	}
	return e.repository.FindById(nil, filterContext, id)
}

func (e *SystemUserOrganizationService) FindByIds(filterContext *repository.SystemUserOrganizationSearch, ids []uuid.UUID) ([]*model.SystemUserOrganization, error) {
	for _, id := range ids {
		if id == uuid.Nil {
			return nil, errors.New("invalid id")
		}
	}

	return e.repository.FindByIds(nil, filterContext, ids)
}

func (e *SystemUserOrganizationService) FindAll(filterContext *repository.SystemUserOrganizationSearch, pageSize int, page int) (response.PagedResult[*model.SystemUserOrganization], error) {
	var result response.PagedResult[*model.SystemUserOrganization]

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

func (e *SystemUserOrganizationService) Search(searchParams *repository.SystemUserOrganizationSearch) (response.PagedResult[*model.SystemUserOrganization], error) {
	var result response.PagedResult[*model.SystemUserOrganization]

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

func (e *SystemUserOrganizationService) Deleted(searchParams *repository.SystemUserOrganizationSearch) ([]string, error) {
	return e.repository.Deleted(nil, searchParams)
}
