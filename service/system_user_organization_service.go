package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/roksky/bootstrap-api/data/request"
	"github.com/roksky/bootstrap-api/data/response"
	"github.com/roksky/bootstrap-api/model"
	"github.com/roksky/bootstrap-api/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (e *SystemUserOrganizationService) Register(item *request.RegisterOrgSystemUser) error {
	err := e.Validate.Struct(item)
	if err != nil {
		return err
	}

	err = e.repository.GetDB().Transaction(func(tx *gorm.DB) error {
		// create a user
		systemUser := &model.SystemUser{
			UserId:              uuid.New(),
			FullNames:           item.FullNames,
			UserName:            item.UserName,
			Password:            item.Password,
			NeedsPasswordChange: item.NeedsPasswordChange,
			OrganizationManaged: true,
			DateCreated:         time.Now(),
			DateUpdated:         time.Now(),
			PrimaryOrganization: item.PrimaryOrganization,
		}

		// if the password is null set the password to the username
		if systemUser.Password == "" {
			systemUser.Password = systemUser.UserName
		}
		// Encrypt the password using bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(systemUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		systemUser.Password = string(hashedPassword)

		if err := tx.Create(systemUser).Error; err != nil {
			return err // Rollback the transaction
		}

		// create a user organization
		systemUserOrganization := &model.SystemUserOrganization{
			SystemUser:     *systemUser,
			OrganizationId: item.PrimaryOrganization,
			UserRole:       item.UserRole,
		}
		if err := tx.Create(systemUserOrganization).Error; err != nil {
			return err // Rollback the transaction
		}

		return nil
	})

	return err
}

func (e *SystemUserOrganizationService) Create(filterContext *repository.SystemUserOrganizationSearch, item *model.SystemUserOrganization) (*model.SystemUserOrganization, error) {
	err := e.Validate.Struct(item)
	if err != nil {
		return nil, err
	}
	return e.repository.Save(filterContext, item)
}

func (e *SystemUserOrganizationService) CreateMany(filterContext *repository.SystemUserOrganizationSearch, items []*model.SystemUserOrganization) ([]*model.SystemUserOrganization, error) {
	for _, item := range items {
		err := e.Validate.Struct(item)
		if err != nil {
			return nil, err
		}
	}

	return e.repository.SaveMany(filterContext, items)
}

func (e *SystemUserOrganizationService) Update(filterContext *repository.SystemUserOrganizationSearch, item *model.SystemUserOrganization) (*model.SystemUserOrganization, error) {
	err := e.Validate.Struct(item)
	if err != nil {
		return nil, err
	}
	if item.Id == uuid.Nil {
		return nil, errors.New("entity id is missing")
	}
	return e.repository.Update(filterContext, item)
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

	return e.repository.UpdateMany(filterContext, items)
}

func (e *SystemUserOrganizationService) Delete(filterContext *repository.SystemUserOrganizationSearch, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("entity id is missing")
	}

	return e.repository.Delete(filterContext, id)
}

func (e *SystemUserOrganizationService) DeleteMany(filterContext *repository.SystemUserOrganizationSearch, ids []uuid.UUID) error {
	for _, id := range ids {
		if id == uuid.Nil {
			return errors.New("invalid id")
		}
	}

	return e.repository.DeleteByIds(filterContext, ids)
}

func (e *SystemUserOrganizationService) FindById(filterContext *repository.SystemUserOrganizationSearch, id uuid.UUID) (*model.SystemUserOrganization, error) {
	if id == uuid.Nil {
		return nil, errors.New("entity id is missing")
	}
	return e.repository.FindById(filterContext, id)
}

func (e *SystemUserOrganizationService) FindByIds(filterContext *repository.SystemUserOrganizationSearch, ids []uuid.UUID) ([]*model.SystemUserOrganization, error) {
	for _, id := range ids {
		if id == uuid.Nil {
			return nil, errors.New("invalid id")
		}
	}

	return e.repository.FindByIds(filterContext, ids)
}

func (e *SystemUserOrganizationService) FindAll(filterContext *repository.SystemUserOrganizationSearch, pageSize int, page int) (response.PagedResult[*model.SystemUserOrganization], error) {
	var result response.PagedResult[*model.SystemUserOrganization]

	items, err := e.repository.FindAll(filterContext, pageSize, page)
	if err != nil {
		return result, err
	}

	count, err := e.repository.Count(nil)

	result.Items = items
	result.TotalItems = count
	result.PageSize = count
	result.PageNumber = 0

	return result, nil
}

func (e *SystemUserOrganizationService) Search(searchParams *repository.SystemUserOrganizationSearch) (response.PagedResult[*model.SystemUserOrganization], error) {
	var result response.PagedResult[*model.SystemUserOrganization]

	items, err := e.repository.Search(searchParams)
	if err != nil {
		return result, err
	}

	count, err := e.repository.Count(searchParams)

	result.Items = items
	result.TotalItems = count
	result.PageSize = int64(searchParams.PageSize)
	result.PageNumber = searchParams.PageNumber

	return result, nil
}

func (e *SystemUserOrganizationService) Deleted(searchParams *repository.SystemUserOrganizationSearch) ([]string, error) {
	return e.repository.Deleted(searchParams)
}
