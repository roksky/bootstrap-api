package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/roksky/bootstrap-api/constants"
)

type HttpMethod int

const (
	GET HttpMethod = iota
	POST
	PUT
	PATCH
	DELETE
)

type HttpFunc struct {
	method      HttpMethod
	httpFunc    gin.HandlerFunc
	urlTemplate string
}

func NewHttpFunc(method HttpMethod, urlTemplate string, httpFunc gin.HandlerFunc) *HttpFunc {
	return &HttpFunc{
		method:      method,
		httpFunc:    httpFunc,
		urlTemplate: urlTemplate,
	}
}

func (h *HttpFunc) GetHttpMethod() HttpMethod {
	return h.method
}

func (h *HttpFunc) GetHandlerFunc() gin.HandlerFunc {
	return h.httpFunc
}

func (h *HttpFunc) GetUrlTemplate() string {
	return h.urlTemplate
}

type Controller interface {
	GroupName() string
	Handlers() []*HttpFunc
	IsAuthEnabled() bool
}

func GetTokenInfo(ctx *gin.Context) (oauth2.TokenInfo, error) {
	ti, exists := ctx.Get(constants.TokenKey)
	if !exists {
		return nil, errors.New("no token info found")
	}

	tokenInfo := ti.(oauth2.TokenInfo)
	return tokenInfo, nil
}
