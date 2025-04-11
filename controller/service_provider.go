package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/roksky/bootstrap-api/job"
	"github.com/roksky/bootstrap-api/repository"
	"github.com/roksky/bootstrap-api/service"
	"gorm.io/gorm"
)

type Provider interface {
	GetControllers() []Controller
	GetJobs() []job.Job
}

func NewProvider(db *gorm.DB, validate *validator.Validate) Provider {
	return &ControllersRegistry{
		db:       db,
		validate: validate,
	}
}

type ControllersRegistry struct {
	db       *gorm.DB
	validate *validator.Validate
}

func (c *ControllersRegistry) GetControllers() []Controller {
	organizationRepository := repository.NewOrganizationRepository(c.db)

	organizationService := service.NewOrganizationService(organizationRepository, c.validate)

	return []Controller{
		NewOrganizationController(organizationService),
	}
}

func (c *ControllersRegistry) GetJobs() []job.Job {
	return []job.Job{}
}
