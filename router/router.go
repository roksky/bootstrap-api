package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	ginoauth2 "github.com/go-oauth2/gin-server"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/roksky/bootstrap-api/constants"
	"github.com/roksky/bootstrap-api/controller"
)

type RouteHandler interface {
	RegisterRoutes(controllers []controller.Controller)
	RegisterRoute(controller controller.Controller)
	GetEngine() *gin.Engine
	EnableAuth(introspectURL string, clientId string, clientSecret string) error
	AllowCORS()
	EnableSentry(dsn string)
}

type Router struct {
	baseUrl     string
	engine      *gin.Engine
	server      *server.Server
	authEnabled bool
}

func NewRouteHandler(baseUrl string) (RouteHandler, error) {
	return &Router{
		baseUrl:     baseUrl,
		engine:      gin.Default(),
		authEnabled: false,
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

	// Middleware to enforce Bearer token validation
	authMiddleware := func(c *gin.Context) {
		ti, err := r.server.ValidationBearerToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		// you can pull claims from ti.GetUserID(), ti.GetScope(), etc.
		c.Set("TokenInfo", ti)
		c.Next()
	}
	baseRouter.Use(authMiddleware)

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

func (r *Router) EnableSentry(dsn string) {
	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           dsn,
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

func (r *Router) EnableAuth(introspectURL string, clientId string, clientSecret string) error {
	// Initialize the OAuth2 manager with the in-memory token store
	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(NewIntrospectionTokenStore(introspectURL, clientId, clientSecret))
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	manager.SetPasswordTokenCfg(&manage.Config{AccessTokenExp: time.Hour * 24, RefreshTokenExp: time.Hour * 24 * 30 * 12, IsGenerateRefresh: true})

	// Initialize and configure the OAuth2 server
	r.server = ginoauth2.InitServer(manager)
	ginoauth2.SetAllowGetAccessRequest(true)

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

		//extensions := r.extensionFieldsHandler(token)
		// merge ro response
		//for k, v := range extensions {
		//response[k] = v
		//}

		c.JSON(http.StatusOK, response)
	}
}
