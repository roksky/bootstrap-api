package model

import (
	"github.com/google/uuid"
)

type SystemUserOrganization struct {
	IdentifiedModel
	SystemUserId   uuid.UUID      `gorm:"index;column:system_user" json:"-"`
	SystemUser     SystemUser     `json:"systemUser"`
	OrganizationId uuid.UUID      `gorm:"index;column:organization" json:"-"`
	Organization   Organization   `json:"organization"`
	UserRole       SystemUserRole `gorm:"index" json:"userRole"`
}
