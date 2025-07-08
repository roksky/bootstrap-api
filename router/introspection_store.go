package router

import (
	"encoding/json"
	"errors"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strings"
	"time"
)

// IntrospectionResponse models what most OAuth2 servers return at /introspect
type IntrospectionResponse struct {
	Active    bool   `json:"active"`
	ClientID  string `json:"client_id"`
	Username  string `json:"username,omitempty"`
	Scope     string `json:"scope,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	TokenType string `json:"token_type,omitempty"`
}

// IntrospectionTokenStore implements oauth2.TokenStore by calling /introspect
type IntrospectionTokenStore struct {
	IntrospectURL string
	ClientID      string
	ClientSecret  string
	HTTPClient    *http.Client
}

func NewIntrospectionTokenStore(url, clientID, clientSecret string) *IntrospectionTokenStore {
	return &IntrospectionTokenStore{
		IntrospectURL: url,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		HTTPClient:    &http.Client{Timeout: 5 * time.Second},
	}
}

// GetByAccess calls the introspection endpoint and, if active, returns a TokenInfo
func (s *IntrospectionTokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	payload := strings.NewReader("token=" + access)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.IntrospectURL, payload)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.ClientID, s.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("introspection failed: " + string(body))
	}

	var introspect IntrospectionResponse
	if err := json.Unmarshal(body, &introspect); err != nil {
		return nil, err
	}
	if !introspect.Active {
		return nil, errors.New("token is not active")
	}

	// Map introspection response into a TokenInfo
	ti := &models.Token{
		ClientID:        introspect.ClientID,
		UserID:          introspect.Username,
		Scope:           introspect.Scope,
		Access:          access,
		AccessCreateAt:  time.Now(),
		AccessExpiresIn: time.Until(time.Unix(introspect.Exp, 0)),
	}
	return ti, nil
}

// We don’t need refresh tokens here.
func (s *IntrospectionTokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return nil
}
func (s *IntrospectionTokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	return nil
}
func (s *IntrospectionTokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return nil
}
func (s *IntrospectionTokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	return nil, errors.New("refresh not supported")
}

func (s *IntrospectionTokenStore) RemoveByCode(ctx context.Context, code string) error {
	return nil
}

func (s *IntrospectionTokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return s.GetByAccess(ctx, code)
}
