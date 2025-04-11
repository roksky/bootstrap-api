package request

import "github.com/roksky/bootstrap-api/model"

type RegisterUserAndOrg struct {
	User         *model.SystemUser
	Organization *model.Organization
}

type CreateOrganizationRequest struct {
	Name string `validate:"required,min=1,max=200" json:"name"`
}

type UpdateOrganizationRequest struct {
	Id   int    `validate:"required"`
	Name string `validate:"required,max=200,min=1" json:"name"`
}
