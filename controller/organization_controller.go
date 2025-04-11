package controller

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/roksky/bootstrap-api/data/response"
	"github.com/roksky/bootstrap-api/helper"
	"github.com/roksky/bootstrap-api/model"
	"github.com/roksky/bootstrap-api/repository"
	"github.com/roksky/bootstrap-api/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type OrganizationController struct {
	service service.BaseService[model.Organization, uuid.UUID, repository.OrganizationSearch]
}

func NewOrganizationController(service service.BaseService[model.Organization, uuid.UUID, repository.OrganizationSearch]) *OrganizationController {
	return &OrganizationController{
		service: service,
	}
}

func (controller *OrganizationController) Create(ctx *gin.Context) {
	log.Info().Msg("create organization")

	tokenInfo, err := GetTokenInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	createItem := &model.Organization{}
	createItem.CreatedBy = tokenInfo.GetUserID()
	createItem.UpdatedBy = tokenInfo.GetUserID()
	err = ctx.ShouldBindJSON(createItem)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	search := &repository.OrganizationSearch{}

	ctx.Header("Content-Type", "application/json")

	item, err := controller.service.Create(search, createItem)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Code: "1", Message: err.Error()})
	} else {
		ctx.JSON(http.StatusCreated, item)
	}
}

func (controller *OrganizationController) CreateMany(ctx *gin.Context) {
	log.Info().Msg("create many organizations")

	tokenInfo, err := GetTokenInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	createItems, err := helper.ReadJsonAsType[*model.Organization](ctx.Request.Body)

	search := &repository.OrganizationSearch{}

	ctx.Header("Content-Type", "application/json")

	for _, item := range createItems {
		item.CreatedBy = tokenInfo.GetUserID()
		item.UpdatedBy = tokenInfo.GetUserID()
	}

	item, err := controller.service.CreateMany(search, createItems)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusCreated, item)
	}
}

func (controller *OrganizationController) Update(ctx *gin.Context) {
	log.Info().Msg("update organization")

	tokenInfo, err := GetTokenInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	id, err := helper.ParamAsUUId(ctx, "orgId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Code: "1", Message: err.Error()})
		return
	}

	createItem := &model.Organization{}
	createItem.Id = id
	createItem.UpdatedBy = tokenInfo.GetUserID()
	err = ctx.ShouldBindJSON(createItem)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	search := &repository.OrganizationSearch{}

	ctx.Header("Content-Type", "application/json")

	item, err := controller.service.Update(search, createItem)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Code: "1", Message: err.Error()})
	} else {
		ctx.JSON(http.StatusOK, item)
	}
}

func (controller *OrganizationController) UpdateMany(ctx *gin.Context) {
	log.Info().Msg("update many organizations")

	tokenInfo, err := GetTokenInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	createItems, err := helper.ReadJsonAsType[*model.Organization](ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	for _, item := range createItems {
		item.UpdatedBy = tokenInfo.GetUserID()
	}

	search := &repository.OrganizationSearch{}

	ctx.Header("Content-Type", "application/json")

	item, err := controller.service.UpdateMany(search, createItems)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusCreated, item)
	}
}

func (controller *OrganizationController) Delete(ctx *gin.Context) {
	log.Info().Msg("delete organization")

	idUuid, err := helper.ParamAsUUId(ctx, "orgId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if idUuid == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	search := &repository.OrganizationSearch{}

	ctx.Header("Content-Type", "application/json")

	err = controller.service.Delete(search, idUuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, "deleted")
	}
}

func (controller *OrganizationController) DeleteMany(ctx *gin.Context) {
	log.Info().Msg("delete organizations")

	ids, err := helper.ReadJsonAsType[uuid.UUID](ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	search := &repository.OrganizationSearch{}

	ctx.Header("Content-Type", "application/json")

	err = controller.service.DeleteMany(search, ids)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, "deleted")
	}
}

func (controller *OrganizationController) FindById(ctx *gin.Context) {
	log.Info().Msg("findbyid organizations")

	id, err := helper.ParamAsUUId(ctx, "orgId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	search := &repository.OrganizationSearch{}

	ctx.Header("Content-Type", "application/json")

	item, err := controller.service.FindById(search, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, item)
	}
}

func (controller *OrganizationController) FindByIds(ctx *gin.Context) {
	log.Info().Msg("findbyids organizations")

	ids, err := helper.ReadJsonAsType[uuid.UUID](ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	search := &repository.OrganizationSearch{}

	ctx.Header("Content-Type", "application/json")

	item, err := controller.service.FindByIds(search, ids)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, item)
	}
}

func (controller *OrganizationController) SearchAll(ctx *gin.Context) {
	log.Info().Msg("search all organizations")

	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "100"))
	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("pageNumber", "0"))
	orderBy := helper.GetSortString(ctx.DefaultQuery("orderBy", ""))

	search := &repository.OrganizationSearch{
		PageSize:         pageSize,
		PageNumber:       pageNumber,
		OrganizationType: ctx.Query("organizationType"),
		OrderBy:          orderBy,
	}

	ctx.Header("Content-Type", "application/json")

	item, err := controller.service.Search(search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, item)
	}
}

func (controller *OrganizationController) GetDeleted(ctx *gin.Context) {
	log.Info().Msg("search all deleted organizations")

	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "100"))
	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("pageNumber", "0"))

	search := &repository.OrganizationSearch{
		PageSize:         pageSize,
		PageNumber:       pageNumber,
		OrganizationType: ctx.Query("organizationType"),
	}

	ctx.Header("Content-Type", "application/json")

	item, err := controller.service.Deleted(search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, item)
	}
}

func (controller *OrganizationController) GroupName() string {
	return "/org"
}

func (controller *OrganizationController) Handlers() []*HttpFunc {
	return []*HttpFunc{
		NewHttpFunc(GET, "", controller.SearchAll),
		NewHttpFunc(GET, "/:orgId", controller.FindById),
		NewHttpFunc(GET, "s/:orgIds", controller.FindByIds),
		NewHttpFunc(POST, "", controller.Create),
		NewHttpFunc(POST, "s", controller.CreateMany),
		NewHttpFunc(PATCH, "/:orgId", controller.Update),
		NewHttpFunc(PATCH, "s", controller.UpdateMany),
		NewHttpFunc(DELETE, "/:orgId", controller.Delete),
		NewHttpFunc(DELETE, "s", controller.DeleteMany),
		NewHttpFunc(GET, "/deleted", controller.GetDeleted),
	}
}

func (controller *OrganizationController) IsAuthEnabled() bool {
	return true
}
