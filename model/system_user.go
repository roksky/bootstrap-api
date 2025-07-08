package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SystemUser struct {
	UserId              uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"userId"`
	UserName            string         `gorm:"index" json:"userName"`
	DateCreated         time.Time      `json:"dateCreated"`
	DateUpdated         time.Time      `json:"dateUpdated"`
	DateDeleted         gorm.DeletedAt `gorm:"index" json:"dateDeleted"`
	PrimaryOrganization uuid.UUID      `json:"primaryOrganization"`
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
