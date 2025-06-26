package request

import (
	"github.com/roksky/bootstrap-api/model"

	"github.com/google/uuid"
)

type RegisterOrgSystemUser struct {
	FullNames           string               `json:"fullNames"`
	UserName            string               `json:"userName"`
	Password            string               `json:"password"`
	PrimaryOrganization uuid.UUID            `json:"primaryOrganization"`
	NeedsPasswordChange bool                 `json:"needsPasswordChange"`
	CreatedBy           string               `json:"createdBy"`
	UserRole            model.SystemUserRole `json:"userRole"`
}
