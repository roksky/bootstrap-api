package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IdentifiedModel struct {
	Id          uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	DateCreated time.Time      `json:"dateCreated"`
	DateUpdated time.Time      `json:"dateUpdated"`
	DateDeleted gorm.DeletedAt `gorm:"index" json:"dateDeleted"`

	CreatedBy string `gorm:"type:varchar(255)" json:"createdBy"`
	UpdatedBy string `gorm:"type:varchar(255)" json:"updatedBy"`
	DeletedBy string `gorm:"type:varchar(255)" json:"deletedBy"`
}

func (u *IdentifiedModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.DateCreated = time.Now()
	u.DateUpdated = time.Now()
	return
}

func (u *IdentifiedModel) BeforeUpdate(tx *gorm.DB) (err error) {
	u.DateUpdated = time.Now()
	return
}
