package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/roksky/bootstrap-api/model"
	"github.com/roksky/bootstrap-api/repository"
)

type SystemUserService struct {
	systemUserRepo                   *repository.SystemUserRepo
	organizationRepo                 repository.BaseRepository[model.Organization, uuid.UUID, repository.OrganizationSearch]
	systemUserOrganizationRepository repository.BaseRepository[model.SystemUserOrganization, uuid.UUID, repository.SystemUserOrganizationSearch]
	Validate                         *validator.Validate
}

func NewSystemUserService(repository *repository.SystemUserRepo, organizationRepo repository.BaseRepository[model.Organization, uuid.UUID, repository.OrganizationSearch], systemUserOrganizationRepository repository.BaseRepository[model.SystemUserOrganization, uuid.UUID, repository.SystemUserOrganizationSearch], validate *validator.Validate) *SystemUserService {
	return &SystemUserService{
		systemUserRepo:                   repository,
		organizationRepo:                 organizationRepo,
		systemUserOrganizationRepository: systemUserOrganizationRepository,
		Validate:                         validate,
	}
}

// SearchUserByEmailOrMobile finds a user by email of mobile
func (e *SystemUserService) SearchUserByEmailOrMobile(email string, mobileNumber string) (*model.SystemUser, error) {
	if email != "" {
		return e.systemUserRepo.FindByUserName(email)
	}
	if mobileNumber != "" {
		return e.systemUserRepo.FindByUserName(email)
	}
	return nil, errors.New("email and mobile are both nil")
}

// RegisterOrganization creates an organization
func (e *SystemUserService) RegisterOrganization(org *model.Organization) (*model.Organization, error) {
	organization, err := e.organizationRepo.Save(nil, nil, org)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

func (e *SystemUserService) DeleteUserById(userId uuid.UUID) error {
	return e.systemUserRepo.Delete(&model.SystemUser{
		UserId: userId,
	})
}

func (e *SystemUserService) DeleteOrgById(orgId uuid.UUID) error {
	return e.organizationRepo.Delete(nil, nil, orgId)
}

func (e *SystemUserService) DeleteOrgSystemUserById(orgId uuid.UUID, userId uuid.UUID) error {
	search := &repository.SystemUserOrganizationSearch{
		SystemUser:     userId.String(),
		OrganizationId: orgId,
	}

	users, err := e.systemUserOrganizationRepository.Search(nil, search)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return errors.New("user not found")
	}

	return e.systemUserOrganizationRepository.Delete(nil, nil, users[0].Id)
}

func (e *SystemUserService) GetUserByUserName(userName string) (*model.SystemUser, error) {
	return e.systemUserRepo.FindByUserName(userName)
}

func (e *SystemUserService) GetByUserId(userId uuid.UUID) (*model.SystemUser, error) {
	return e.systemUserRepo.FindById(userId)
}
