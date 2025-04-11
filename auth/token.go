package auth

import (
	"github.com/go-oauth2/oauth2/v4"
	"time"
)

// NewDBToken create to token model instance
func NewDBToken() *DbToken {
	return &DbToken{}
}

// NewDBTokenFromTokenInfo create to token model instance
func NewDBTokenFromTokenInfo(info oauth2.TokenInfo) *DbToken {
	return &DbToken{
		ClientID:            info.GetClientID(),
		UserID:              info.GetUserID(),
		RedirectURI:         info.GetRedirectURI(),
		Scope:               info.GetScope(),
		Code:                info.GetCode(),
		CodeChallenge:       info.GetCodeChallenge(),
		CodeChallengeMethod: string(info.GetCodeChallengeMethod()),
		CodeCreateAt:        info.GetCodeCreateAt(),
		CodeExpiresIn:       info.GetCodeExpiresIn(),
		Access:              info.GetAccess(),
		AccessCreateAt:      info.GetAccessCreateAt(),
		AccessExpiresIn:     info.GetAccessExpiresIn(),
		Refresh:             info.GetRefresh(),
		RefreshCreateAt:     info.GetRefreshCreateAt(),
		RefreshExpiresIn:    info.GetRefreshExpiresIn(),
	}
}

// DbToken token model
type DbToken struct {
	ClientID            string        `bson:"ClientID"`
	UserID              string        `bson:"UserID"`
	RedirectURI         string        `bson:"RedirectURI"`
	Scope               string        `bson:"Scope"`
	Code                string        `gorm:"index" bson:"Code"`
	CodeChallenge       string        `bson:"CodeChallenge"`
	CodeChallengeMethod string        `bson:"CodeChallengeMethod"`
	CodeCreateAt        time.Time     `bson:"CodeCreateAt"`
	CodeExpiresIn       time.Duration `bson:"CodeExpiresIn"`
	Access              string        `gorm:"index" bson:"Access"`
	AccessCreateAt      time.Time     `bson:"AccessCreateAt"`
	AccessExpiresIn     time.Duration `bson:"AccessExpiresIn"`
	Refresh             string        `gorm:"index" bson:"Refresh"`
	RefreshCreateAt     time.Time     `bson:"RefreshCreateAt"`
	RefreshExpiresIn    time.Duration `bson:"RefreshExpiresIn"`
}

func (t *DbToken) TableName() string {
	return "tokens"
}

func (t *DbToken) New() oauth2.TokenInfo {
	return NewDBToken()
}

// GetClientID the client id
func (t *DbToken) GetClientID() string {
	return t.ClientID
}

// SetClientID the client id
func (t *DbToken) SetClientID(clientID string) {
	t.ClientID = clientID
}

// GetUserID the user id
func (t *DbToken) GetUserID() string {
	return t.UserID
}

// SetUserID the user id
func (t *DbToken) SetUserID(userID string) {
	t.UserID = userID
}

// GetRedirectURI redirect URI
func (t *DbToken) GetRedirectURI() string {
	return t.RedirectURI
}

// SetRedirectURI redirect URI
func (t *DbToken) SetRedirectURI(redirectURI string) {
	t.RedirectURI = redirectURI
}

// GetScope get scope of authorization
func (t *DbToken) GetScope() string {
	return t.Scope
}

// SetScope get scope of authorization
func (t *DbToken) SetScope(scope string) {
	t.Scope = scope
}

// GetCode authorization code
func (t *DbToken) GetCode() string {
	return t.Code
}

// SetCode authorization code
func (t *DbToken) SetCode(code string) {
	t.Code = code
}

// GetCodeCreateAt create Time
func (t *DbToken) GetCodeCreateAt() time.Time {
	return t.CodeCreateAt
}

// SetCodeCreateAt create Time
func (t *DbToken) SetCodeCreateAt(createAt time.Time) {
	t.CodeCreateAt = createAt
}

// GetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *DbToken) GetCodeExpiresIn() time.Duration {
	return t.CodeExpiresIn
}

// SetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *DbToken) SetCodeExpiresIn(exp time.Duration) {
	t.CodeExpiresIn = exp
}

// GetCodeChallenge challenge code
func (t *DbToken) GetCodeChallenge() string {
	return t.CodeChallenge
}

// SetCodeChallenge challenge code
func (t *DbToken) SetCodeChallenge(code string) {
	t.CodeChallenge = code
}

// GetCodeChallengeMethod challenge method
func (t *DbToken) GetCodeChallengeMethod() oauth2.CodeChallengeMethod {
	return oauth2.CodeChallengeMethod(t.CodeChallengeMethod)
}

// SetCodeChallengeMethod challenge method
func (t *DbToken) SetCodeChallengeMethod(method oauth2.CodeChallengeMethod) {
	t.CodeChallengeMethod = string(method)
}

// GetAccess access DbToken
func (t *DbToken) GetAccess() string {
	return t.Access
}

// SetAccess access DbToken
func (t *DbToken) SetAccess(access string) {
	t.Access = access
}

// GetAccessCreateAt create Time
func (t *DbToken) GetAccessCreateAt() time.Time {
	return t.AccessCreateAt
}

// SetAccessCreateAt create Time
func (t *DbToken) SetAccessCreateAt(createAt time.Time) {
	t.AccessCreateAt = createAt
}

// GetAccessExpiresIn the lifetime in seconds of the access token
func (t *DbToken) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresIn
}

// SetAccessExpiresIn the lifetime in seconds of the access token
func (t *DbToken) SetAccessExpiresIn(exp time.Duration) {
	t.AccessExpiresIn = exp
}

// GetRefresh refresh DbToken
func (t *DbToken) GetRefresh() string {
	return t.Refresh
}

// SetRefresh refresh DbToken
func (t *DbToken) SetRefresh(refresh string) {
	t.Refresh = refresh
}

// GetRefreshCreateAt create Time
func (t *DbToken) GetRefreshCreateAt() time.Time {
	return t.RefreshCreateAt
}

// SetRefreshCreateAt create Time
func (t *DbToken) SetRefreshCreateAt(createAt time.Time) {
	t.RefreshCreateAt = createAt
}

// GetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *DbToken) GetRefreshExpiresIn() time.Duration {
	return t.RefreshExpiresIn
}

// SetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *DbToken) SetRefreshExpiresIn(exp time.Duration) {
	t.RefreshExpiresIn = exp
}
