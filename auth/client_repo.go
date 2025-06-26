package auth

import "gorm.io/gorm"

type ClientRepo struct {
	db *gorm.DB
}

func NewClientRepo(db *gorm.DB) *ClientRepo {
	return &ClientRepo{
		db: db,
	}
}

func (m *ClientRepo) Save(item *Client) error {
	result := m.db.Create(item)
	return result.Error
}

func (m *ClientRepo) FindById(id string) (*Client, error) {
	var client Client
	result := m.db.Where("ID = ?", id).Find(&client)
	if result.Error == nil {
		return &client, nil
	} else {
		return nil, nil
	}
}

func (m *ClientRepo) DeleteById(code string) error {
	var client Client
	result := m.db.Where("ID = ?", code).Delete(&client)
	return result.Error
}
