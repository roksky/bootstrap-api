package model

type Organization struct {
	IdentifiedModel
	Name string `gorm:"type:varchar(255)" json:"name"`
}
