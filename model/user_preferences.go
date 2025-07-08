package model

import "github.com/google/uuid"

type UserPreferences struct {
	IdentifiedModel
	SystemUserId uuid.UUID              `gorm:"index;column:system_user" json:"-"`
	SystemUser   SystemUser             `json:"systemUser"`
	Preferences  map[string]interface{} `gorm:"type:json"`
}
