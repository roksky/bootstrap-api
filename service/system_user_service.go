package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/roksky/bootstrap-api/model"
	"github.com/roksky/bootstrap-api/repository"
	"golang.org/x/crypto/bcrypt"
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

// RegisterUserAndOrg a user and org
func (e *SystemUserService) RegisterUserAndOrg(user *model.SystemUser, org *model.Organization) (*model.SystemUserOrganization, error) {
	organization, err := e.organizationRepo.Save(nil, org)
	if err != nil {
		return nil, err
	}

	user.PrimaryOrganization = organization.Id

	registeredUser, err := e.Register(user)
	if err != nil {
		return nil, err
	}

	systemUser := &model.SystemUserOrganization{
		SystemUserId:   registeredUser.UserId,
		OrganizationId: organization.Id,
		UserRole:       model.Owner,
	}

	return e.systemUserOrganizationRepository.Save(nil, systemUser)
}

// RegisterOrganization creates an organization
func (e *SystemUserService) RegisterOrganization(org *model.Organization) (*model.Organization, error) {
	organization, err := e.organizationRepo.Save(nil, org)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

// Register a user
func (e *SystemUserService) Register(item *model.SystemUser) (*model.SystemUser, error) {
	err := e.Validate.Struct(item)
	if err != nil {
		return nil, err
	}
	// check if id is present
	if item.UserId == uuid.Nil {
		newId, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		item.UserId = newId
	}

	// Encrypt the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	item.Password = string(hashedPassword)

	err = e.systemUserRepo.Save(item)
	return item, err
}

func (e *SystemUserService) DeleteUserById(userId uuid.UUID) error {
	return e.systemUserRepo.Delete(&model.SystemUser{
		UserId: userId,
	})
}

func (e *SystemUserService) DeleteOrgById(orgId uuid.UUID) error {
	return e.organizationRepo.Delete(nil, orgId)
}

func (e *SystemUserService) DeleteOrgSystemUserById(orgId uuid.UUID, userId uuid.UUID) error {
	search := &repository.SystemUserOrganizationSearch{
		SystemUser:     userId.String(),
		OrganizationId: orgId,
	}

	users, err := e.systemUserOrganizationRepository.Search(search)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return errors.New("user not found")
	}

	return e.systemUserOrganizationRepository.Delete(nil, users[0].Id)
}

func (e *SystemUserService) GetUserByUserName(userName string) (*model.SystemUser, error) {
	return e.systemUserRepo.FindByUserName(userName)
}

func (e *SystemUserService) GetByUserId(userId uuid.UUID) (*model.SystemUser, error) {
	return e.systemUserRepo.FindById(userId)
}

func (e *SystemUserService) VerifyPassword(username, password string) error {
	if username == "" {
		return errors.New("user not found")
	}
	if password == "" {
		return errors.New("user not found")
	}
	user, err := e.systemUserRepo.FindByUserName(username)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Compare the stored hashed password, with the password provided
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
