package model

// Organization represents a base organization model.
// Users should extend this struct to provide additional details.
type Organization struct {
	IdentifiedModel
	Name string `gorm:"type:varchar(255)" json:"name"`
}
