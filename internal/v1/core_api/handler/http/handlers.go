package http

import (
	"context"
	"main/config"
	"main/internal/v1/core_api/domain/ports"
	"main/pkg/logger"
)

// handlerHttp Core handlers
type HandlerHttp struct {
	ctx     context.Context
	cfg     *config.Config
	service ports.IService
	logger  logger.Logger
}

// NewHttpHandler Core Domain HTTP handlers constructor
func NewHttpHandler(ctx context.Context, cfg *config.Config, service ports.IService, logger logger.Logger) ports.IHandlers {
	return &HandlerHttp{
		ctx:     ctx,
		cfg:     cfg,
		service: service,
		logger:  logger,
	}
}

// Example Handler:
/*
// GetSellRepAcList Func
// @Description Func lists seller rep
// @Summary Func list sell rep by jwt token org_id
// @Tags Common
// @Accept json
// @Produce json
// @Param search query string false "search"
// @Success 200 {object} entities.HandlerResponse
// @Failure 400 {object} entities.HandlerResponse
// @Failure 500 {object} entities.HandlerResponse
// @Router /common/filters/sell-rep [GET]
func (h HandlerHttp) GetSellRepAcList(c *fiber.Ctx) error {
	var responser fiber.Map
	var statusCode int
	var userData ent.UserData
	var record int64
	var data interface{}
	dat := ent.SearchRequestDto{}
	userData = cm.GetAuthDataFromToken(c)

	errParser := c.QueryParser(&dat)
	if errParser != nil {
		statusCode, responser = fiber.StatusBadGateway, cm.HTTPResponser(nil, true, errParser.Error())
	}

	errValidate := validator.Validate.Struct(dat)
	if errValidate != nil {
		statusCode, responser = fiber.StatusBadRequest, cm.HTTPResponser(nil, true, errValidate.Error())
	}

	if errParser == nil && errValidate == nil {
		record, data = h.service.GetSellRepAcList(userData.CurOrgId, dat)
		if record == -1 {
			statusCode, responser = fiber.StatusInternalServerError, cm.HTTPResponser(nil, true, cm.Translate(userData.DefLang, "General.ERROR_DB"))
		} else if record == 0 {
			statusCode, responser = fiber.StatusOK, cm.HTTPResponser(nil, true, cm.Translate(userData.DefLang, "General.ERROR_NO_ROW"))
		} else {
			statusCode, responser = fiber.StatusOK, cm.HTTPResponser(data, false, cm.Translate(userData.DefLang, "General.OK"))
		}
	}

	return c.Status(statusCode).JSON(responser)
}
*/
