package model

import (
	"database/sql/driver"
	"errors"
)

type SystemUserRole string

const (
	Owner    SystemUserRole = "owner"
	Employee SystemUserRole = "employee"
	Admin    SystemUserRole = "admin"
)

// Implement the Scanner interface for the enum type
func (s *SystemUserRole) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}
	*s = SystemUserRole(strValue)
	return nil
}

// Implement the Valuer interface for the enum type
func (s SystemUserRole) Value() (driver.Value, error) {
	return string(s), nil
}
