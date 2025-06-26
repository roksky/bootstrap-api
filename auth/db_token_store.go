package auth

import (
	"context"
	"github.com/go-oauth2/oauth2/v4"
)

// NewDatabaseTokenStore create a token store instance based on a database
func NewDatabaseTokenStore(authDB *Database) (oauth2.TokenStore, error) {
	return &TokenStore{db: authDB}, nil
}

// TokenStore implementation of a DB
type TokenStore struct {
	db *Database
}

// Create and store the new token information
func (ts *TokenStore) Create(_ context.Context, info oauth2.TokenInfo) error {
	dbToken := NewDBTokenFromTokenInfo(info)
	return ts.db.tokenRepo.Save(dbToken)
}

// RemoveByCode use the authorization code to delete the token information
func (ts *TokenStore) RemoveByCode(_ context.Context, code string) error {
	return ts.db.tokenRepo.DeleteByCode(code)
}

// RemoveByAccess use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(_ context.Context, access string) error {
	return ts.db.tokenRepo.DeleteByAccess(access)
}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(_ context.Context, refresh string) error {
	return ts.db.tokenRepo.DeleteByRefresh(refresh)
}

// GetByCode use the authorization code for token information data
func (ts *TokenStore) GetByCode(_ context.Context, code string) (oauth2.TokenInfo, error) {
	return ts.db.tokenRepo.FindByCode(code)
}

// GetByAccess use the access token for token information data
func (ts *TokenStore) GetByAccess(_ context.Context, access string) (oauth2.TokenInfo, error) {
	return ts.db.tokenRepo.FindByAccess(access)
}

// GetByRefresh use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(_ context.Context, refresh string) (oauth2.TokenInfo, error) {
	return ts.db.tokenRepo.FindByRefresh(refresh)
}
