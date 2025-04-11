package auth

import (
	"context"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"strings"
)

type Client struct {
	ID                   string `gorm:"primaryKey"`
	Secret               string
	Domain               string
	Public               bool
	UserID               string
	Scope                string
	AuthorizedGrantTypes string
}

// GetID client id
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret client secret
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain client domain
func (c *Client) GetDomain() string {
	return c.Domain
}

// IsPublic public
func (c *Client) IsPublic() bool {
	return c.Public
}

// GetUserID user id
func (c *Client) GetUserID() string {
	return c.UserID
}

func (c *Client) VerifyScopes(scopes string) error {
	authorizedScopes := splitAndClean(c.Scope, ",")
	scopesToCheck := splitAndClean(scopes, ",")
	if !checkIfSubArrayExists(authorizedScopes, scopesToCheck) {
		return errors.ErrInvalidScope
	}
	return nil
}

func (c *Client) VerifyGrantTypes(grantTypes string) error {
	authorizedGrantTypes := splitAndClean(c.AuthorizedGrantTypes, ",")
	grantTypesToCheck := splitAndClean(grantTypes, ",")
	if !checkIfSubArrayExists(authorizedGrantTypes, grantTypesToCheck) {
		return errors.ErrUnauthorizedClient
	}
	return nil
}

func checkIfSubArrayExists(array []string, subArray []string) bool {
	for _, subArrayItem := range subArray {
		found := false
		for _, arrayItem := range array {
			if subArrayItem == arrayItem {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func splitAndClean(values string, sep string) []string {
	parts := strings.Split(values, sep)
	cleanedValues := make([]string, len(parts))
	// Iterate through the parts
	for i, part := range parts {
		cleanedValues[i] = strings.ToLower(strings.TrimSpace(part))
	}
	return cleanedValues
}

// NewDatabaseClientStore create client store
func NewDatabaseClientStore(authDB *Database) *DatabaseClientStore {
	return &DatabaseClientStore{
		db: authDB,
	}
}

// DatabaseClientStore client information store
type DatabaseClientStore struct {
	db *Database
}

// GetByID according to the ID for the client information
func (cs *DatabaseClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return cs.db.clientRepo.FindById(id)
}

// Set client information
func (cs *DatabaseClientStore) Set(id string, client *Client) (err error) {
	existing, _ := cs.db.clientRepo.FindById(id)
	if existing != nil {
		_ = cs.db.clientRepo.DeleteById(id)
	}
	client.ID = id

	return cs.db.clientRepo.Save(client)
}
