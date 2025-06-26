package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SystemUser struct {
	UserId              uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"userId"`
	FullNames           string         `json:"fullNames"`
	UserName            string         `gorm:"index" json:"userName"`
	Password            string         `json:"password" binding:"required"`
	DateCreated         time.Time      `json:"dateCreated"`
	DateUpdated         time.Time      `json:"dateUpdated"`
	DateDeleted         gorm.DeletedAt `gorm:"index" json:"dateDeleted"`
	PrimaryOrganization uuid.UUID      `json:"primaryOrganization"`
	NeedsPasswordChange bool           `json:"needsPasswordChange"`
	OrganizationManaged bool           `json:"organizationManaged"`
}

func (t *SystemUser) TableName() string {
	return "system_users"
}

func (t *SystemUser) BeforeCreate(tx *gorm.DB) (err error) {
	t.DateCreated = time.Now()
	t.DateUpdated = time.Now()
	return
}

func (t *SystemUser) BeforeUpdate(tx *gorm.DB) (err error) {
	t.DateUpdated = time.Now()
	return
}
