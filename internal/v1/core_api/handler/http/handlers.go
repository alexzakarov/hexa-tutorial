package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"main/config"
	ent "main/internal/v1/core_api/domain/entities"
	"main/internal/v1/core_api/domain/ports"
	"main/pkg/logger"
	cm "main/pkg/utils/common"
	"main/pkg/utils/validator"
)

// handlerHttp Core handlers
type HandlerHttp struct {
	ctx     context.Context
	cfg     *config.Config
	service ports.IService
	logger  logger.Logger
}

func (h HandlerHttp) CreateUser(c *fiber.Ctx) error {
	var responser fiber.Map
	var statusCode int
	var data interface{}
	dat := ent.UserReqDto{}

	errParser := c.BodyParser(&dat)
	if errParser != nil {
		statusCode, responser = fiber.StatusBadGateway, cm.HTTPResponser(nil, true, errParser.Error())
		return c.Status(statusCode).JSON(responser) // Hemen yanıtı döndür
	}

	errValidate := validator.Validate.Struct(dat)
	if errValidate != nil {
		statusCode, responser = fiber.StatusBadRequest, cm.HTTPResponser(nil, true, errValidate.Error())
		return c.Status(statusCode).JSON(responser) // Hemen yanıtı döndür
	}

	if errParser == nil && errValidate == nil {
		err := h.service.CreateUser(dat)
		if err != nil {
			statusCode, responser = fiber.StatusInternalServerError, cm.HTTPResponser(nil, true, cm.Translate("tr", "General.ERROR_DB"))
		} else {
			statusCode, responser = fiber.StatusOK, cm.HTTPResponser(data, false, cm.Translate("tr", "General.OK"))
		}
	}

	return c.Status(statusCode).JSON(responser)

}

func (h HandlerHttp) GetUserById(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (h HandlerHttp) UpdateUser(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (h HandlerHttp) DeleteUser(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
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
