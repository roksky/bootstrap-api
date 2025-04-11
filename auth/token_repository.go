package auth

import "gorm.io/gorm"

type TokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (m *TokenRepo) Save(item *DbToken) error {
	result := m.db.Create(item)
	return result.Error
}

func (m *TokenRepo) FindByCode(code string) (*DbToken, error) {
	var token DbToken
	result := m.db.Where("Code = ?", code).Find(&token)
	if result.Error == nil {
		return &token, nil
	} else {
		return nil, nil
	}
}

func (m *TokenRepo) DeleteByCode(code string) error {
	var token DbToken
	result := m.db.Where("Code = ?", code).Delete(&token)
	return result.Error
}

func (m *TokenRepo) FindByAccess(access string) (*DbToken, error) {
	var token DbToken
	result := m.db.Where("Access = ?", access).Find(&token)
	if result.Error == nil {
		return &token, nil
	} else {
		return nil, nil
	}
}

func (m *TokenRepo) DeleteByAccess(access string) error {
	var token DbToken
	result := m.db.Where("Access = ?", access).Delete(&token)
	return result.Error
}

func (m *TokenRepo) FindByRefresh(refresh string) (*DbToken, error) {
	var token DbToken
	result := m.db.Where("Refresh = ?", refresh).Find(&token)
	if result.Error == nil {
		return &token, nil
	} else {
		return nil, nil
	}
}

func (m *TokenRepo) DeleteByRefresh(refresh string) error {
	var token DbToken
	result := m.db.Where("Refresh = ?", refresh).Delete(&token)
	return result.Error
}
