package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/roksky/bootstrap-api/auth"
	"github.com/roksky/bootstrap-api/constants"
	"github.com/roksky/bootstrap-api/controller"
	"github.com/roksky/bootstrap-api/model"
	"github.com/roksky/bootstrap-api/repository"
	"github.com/roksky/bootstrap-api/service"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	ginoauth2 "github.com/go-oauth2/gin-server"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteHandler interface {
	RegisterRoutes(controllers []controller.Controller)
	RegisterRoute(controller controller.Controller)
	GetEngine() *gin.Engine
	EnableAuth(clientId string, clientSecret string) error
	AllowCORS()
	EnableSentry()
}

type Router struct {
	baseUrl                       string
	engine                        *gin.Engine
	server                        *server.Server
	authEnabled                   bool
	authDatabase                  *auth.Database
	systemUserService             *service.SystemUserService
	systemUserOrganizationService service.BaseService[model.SystemUserOrganization, uuid.UUID, repository.SystemUserOrganizationSearch]
}

func NewRouteHandler(baseUrl string, db *gorm.DB, authDB *auth.Database) (RouteHandler, error) {
	return &Router{
		baseUrl:                       baseUrl,
		engine:                        gin.Default(),
		authEnabled:                   false,
		authDatabase:                  authDB,
		systemUserService:             service.NewSystemUserService(repository.NewSystemUserRepo(db), repository.NewOrganizationRepository(db), repository.NewSystemUserOrganizationRepository(db), validator.New()),
		systemUserOrganizationService: service.NewSystemUserOrganizationService(repository.NewSystemUserOrganizationRepository(db), validator.New()),
	}, nil
}

func (r *Router) RegisterRoutes(controllers []controller.Controller) {
	for _, cnt := range controllers {
		r.RegisterRoute(cnt)
	}
}

func (r *Router) RegisterRoute(cnt controller.Controller) {
	baseRouter := r.engine.Group(r.baseUrl)
	routerConfig := ginoauth2.Config{
		ErrorHandleFunc: func(ctx *gin.Context, err error) {
			errMsg := fmt.Sprint(err)
			if errMsg == "invalid access token" || errMsg == "expired access token" {
				ctx.AbortWithError(401, err)
			} else {
				ctx.AbortWithError(500, err)
			}
		},
		TokenKey: constants.TokenKey,
		Skipper: func(_ *gin.Context) bool {
			return false
		},
	}
	if r.authEnabled && cnt.IsAuthEnabled() {
		baseRouter.Use(ginoauth2.HandleTokenVerify(routerConfig))
	}

	controllerRouter := baseRouter.Group(cnt.GroupName())

	if r.authEnabled && cnt.IsAuthEnabled() {
		controllerRouter.Use(ginoauth2.HandleTokenVerify(routerConfig))
	}
	for _, route := range cnt.Handlers() {
		switch route.GetHttpMethod() {
		case controller.GET:
			controllerRouter.GET(route.GetUrlTemplate(), route.GetHandlerFunc())
		case controller.POST:
			controllerRouter.POST(route.GetUrlTemplate(), route.GetHandlerFunc())
		case controller.PUT:
			controllerRouter.PUT(route.GetUrlTemplate(), route.GetHandlerFunc())
		case controller.PATCH:
			controllerRouter.PATCH(route.GetUrlTemplate(), route.GetHandlerFunc())
		case controller.DELETE:
			controllerRouter.DELETE(route.GetUrlTemplate(), route.GetHandlerFunc())
		}
	}
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

func (r *Router) AllowCORS() {
	r.engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

func (r *Router) EnableSentry() {
	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           "https://4edd0d587e0fdee6b69867ab46c2cb78@o4507588780425216.ingest.us.sentry.io/4507588782260224",
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	// Once it's done, you can attach the handler as one of your middleware
	r.engine.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
}

func (r *Router) EnableAuth(clientId string, clientSecret string) error {
	// Initialize the OAuth2 manager with the in-memory token store
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(auth.NewDatabaseTokenStore(r.authDatabase))
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// Set up the client store
	dbStore := auth.NewDatabaseClientStore(r.authDatabase)

	err := dbStore.Set(clientId, &auth.Client{
		ID:                   clientId,
		Secret:               clientSecret,
		Domain:               "http://localhost:9119",
		Scope:                "read,write",
		AuthorizedGrantTypes: "password,refresh_token,client_credentials,authorization_code,implicit",
	})
	if err != nil {
		return err
	}

	manager.MapClientStorage(dbStore)

	manager.SetPasswordTokenCfg(&manage.Config{AccessTokenExp: time.Hour * 24, RefreshTokenExp: time.Hour * 24 * 30 * 12, IsGenerateRefresh: true})

	// Initialize and configure the OAuth2 server
	r.server = ginoauth2.InitServer(manager)
	ginoauth2.SetAllowGetAccessRequest(true)
	ginoauth2.SetClientInfoHandler(server.ClientFormHandler)
	ginoauth2.SetClientScopeHandler(r.clientScopeHandler)
	ginoauth2.SetClientAuthorizedHandler(r.clientAuthorizedHandler)
	ginoauth2.SetPasswordAuthorizationHandler(r.passwordAuthorizationHandler)
	ginoauth2.SetExtensionFieldsHandler(r.extensionFieldsHandler)

	r.engine.Any("/oauth2/token", ginoauth2.HandleTokenRequest)
	r.engine.Any("/oauth2/refresh_token", r.refreshTokenHandler(manager))

	// set all the url to check token
	r.authEnabled = true
	return nil
}

func (r *Router) refreshTokenHandler(manager *manage.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken := c.PostForm("refresh_token")
		clientID := c.PostForm("client_id")
		clientSecret := c.PostForm("client_secret")

		if refreshToken == "" || clientID == "" || clientSecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
			return
		}

		tgr := &oauth2.TokenGenerateRequest{
			Refresh:      refreshToken,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}

		token, err := manager.RefreshAccessToken(context.Background(), tgr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"access_token":  token.GetAccess(),
			"refresh_token": token.GetRefresh(),
			"expires_in":    token.GetAccessExpiresIn(),
			"token_type":    "Bearer",
		}

		extensions := r.extensionFieldsHandler(token)
		// merge ro response
		for k, v := range extensions {
			response[k] = v
		}

		c.JSON(http.StatusOK, response)
	}
}

func (r *Router) clientAuthorizedHandler(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
	client, err := r.authDatabase.GetClientRepo().FindById(clientID)
	if err != nil {
		return false, err
	}
	if client == nil {
		return false, nil
	}

	err = client.VerifyGrantTypes(string(grant))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Router) clientScopeHandler(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
	client, err := r.authDatabase.GetClientRepo().FindById(tgr.ClientID)
	if err != nil {
		return false, err
	}
	if client == nil {
		return false, nil
	}

	err = client.VerifyScopes(tgr.Scope)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Router) passwordAuthorizationHandler(ctx context.Context, clientID, username, password string) (userID string, err error) {
	err = r.systemUserService.VerifyPassword(username, password)
	if err != nil {
		return "", err
	}

	return username, nil
}

func (r *Router) extensionFieldsHandler(ti oauth2.TokenInfo) (fieldsValue map[string]interface{}) {
	// add user  orgs to token
	fieldsValue = make(map[string]interface{})

	systemUser, err := r.systemUserService.GetUserByUserName(ti.GetUserID())
	if err != nil {
		panic(err)
	}
	if systemUser.PrimaryOrganization == uuid.Nil {
		fieldsValue["primary_org"] = nil
	} else {
		fieldsValue["primary_org"] = systemUser.PrimaryOrganization.String()
	}

	searchParams := &repository.SystemUserOrganizationSearch{
		PageSize:   100,
		PageNumber: 0,
		SystemUser: systemUser.UserId.String(),
	}
	orgs, err := r.systemUserOrganizationService.Search(searchParams)
	if err != nil {
		panic(err)
	}
	fieldsValue["organizations"] = make([]string, len(orgs.Items))
	for i, org := range orgs.Items {
		fieldsValue["organizations"].([]string)[i] = org.OrganizationId.String()
	}

	// get the token users orgs and primary org
	return fieldsValue
}
