package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rehandwi03/test-case-backend-majoo/criteria"
	custom_error "github.com/rehandwi03/test-case-backend-majoo/internal/error"
	"github.com/rehandwi03/test-case-backend-majoo/internal/helper"
	"github.com/rehandwi03/test-case-backend-majoo/internal/middleware"
	request2 "github.com/rehandwi03/test-case-backend-majoo/request"
	"github.com/rehandwi03/test-case-backend-majoo/service"
	"github.com/rehandwi03/test-case-backend-majoo/util"
	"log"
)

type merchantHandler struct {
	merchantSvc service.MerchantService
}

func NewMerchantHandler(app fiber.Router, merchantService service.MerchantService) {
	handler := merchantHandler{merchantSvc: merchantService}

	app.Post("/merchants", middleware.JwtProtected(), handler.saveMerchant)
	app.Get("/merchants/:id", middleware.JwtProtected(), handler.getByID)
	app.Put("/merchants", middleware.JwtProtected(), handler.updateMerchant)
	app.Delete("/merchants/:id", middleware.JwtProtected(), handler.deleteByID)
	app.Get("/merchants", middleware.JwtProtected(), handler.fetch)
}

func (m *merchantHandler) fetch(c *fiber.Ctx) error {
	pagination := util.GeneratePaginationFromRequest(c)

	merchantCriteria := criteria.MerchantCriteria{
		Pagination: pagination,
	}

	// MerchantCriteria.Email = c.Query("email")
	// MerchantCriteria.PhoneNumber = c.Query("phoneNumber")

	res, err := m.merchantSvc.Fetch(c.Context(), merchantCriteria)
	switch err.(type) {
	case nil:
		return c.Status(fiber.StatusOK).JSON(
			res,
		)
	default:
		log.Printf("error internal server error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusInternalServerError",
				Errors:  err,
			},
		)
	}
}

func (m *merchantHandler) deleteByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		log.Printf("error id is null")
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusNotFound",
				Errors:  "id param is null",
			},
		)
	}

	userId, ok := c.Context().Value("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusNotFound",
				Errors:  "user id not found",
			},
		)
	}

	params := map[string]interface{}{
		"or": map[string]interface{}{},
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?":      id,
				"user_id = ?": userId,
			},
		},
	}

	err := m.merchantSvc.DeleteMerchant(c.Context(), params)
	switch err.(type) {
	case *custom_error.NotFoundError:
		log.Printf("error not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusNotFound",
				Errors:  err,
			},
		)
	case nil:
		return c.Status(fiber.StatusOK).JSON(
			helper.SuccessResponse{
				Status:  "success",
				Message: "success delete data",
			},
		)
	default:
		log.Printf("error internal server error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusInternalServerError",
				Errors:  err,
			},
		)
	}
}

func (m *merchantHandler) getByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		log.Printf("error id is null")
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusNotFound",
				Errors:  "id param is null",
			},
		)
	}

	userId, ok := c.Context().Value("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusNotFound",
				Errors:  "user id not found",
			},
		)
	}

	params := map[string]interface{}{
		"or": map[string]interface{}{},
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?":      id,
				"user_id = ?": userId,
			},
		},
	}

	res, err := m.merchantSvc.GetByParam(c.Context(), params)
	switch err.(type) {
	case *custom_error.NotFoundError:
		log.Printf("error not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusNotFound",
				Errors:  err,
			},
		)
	case nil:
		return c.Status(fiber.StatusOK).JSON(
			helper.SuccessResponse{
				Status:  "success",
				Message: "success get data",
				Data:    res,
			},
		)
	default:
		log.Printf("error internal server error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusInternalServerError",
				Errors:  err,
			},
		)
	}
}

func (m *merchantHandler) updateMerchant(c *fiber.Ctx) error {
	request := new(request2.MerchantUpdateRequest)

	err := c.BodyParser(&request)
	if err != nil {
		log.Printf("error parsing request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusBadRequest",
				Errors:  err.Error(),
			},
		)
	}

	errors := helper.ValidateRequest(*request)
	if errors != nil {
		log.Printf("error validate request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusBadRequest",
				Errors:  errors,
			},
		)
	}

	res, err := m.merchantSvc.UpdateMerchant(c.Context(), request)
	switch err.(type) {
	case *custom_error.NotFoundError:
		log.Printf("error not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusNotFound",
				Errors:  err,
			},
		)
	case *custom_error.BadRequest:
		log.Printf("error bad request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusBadRequest",
				Errors:  err,
			},
		)
	case nil:
		return c.Status(fiber.StatusOK).JSON(
			helper.SuccessResponse{
				Status:  "success",
				Message: "success update data",
				Data: map[string]interface{}{
					"merchant_id": res,
				},
			},
		)
	default:
		log.Printf("error internal server error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusInternalServerError",
				Errors:  err,
			},
		)
	}
}

func (m *merchantHandler) saveMerchant(c *fiber.Ctx) error {
	request := new(request2.MerchantAddRequest)

	err := c.BodyParser(&request)
	if err != nil {
		log.Printf("error parsing request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusBadRequest",
				Errors:  err.Error(),
			},
		)
	}

	errors := helper.ValidateRequest(*request)
	if errors != nil {
		log.Printf("error validate request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusBadRequest",
				Errors:  errors,
			},
		)
	}

	res, err := m.merchantSvc.SaveMerchant(c.Context(), request)
	switch err.(type) {
	case *custom_error.BadRequest:
		log.Printf("error bad request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusBadRequest",
				Errors:  err,
			},
		)
	case nil:
		return c.Status(fiber.StatusCreated).JSON(
			helper.SuccessResponse{
				Status:  "success",
				Message: "success add data",
				Data: map[string]interface{}{
					"merchant_id": res,
				},
			},
		)
	default:
		log.Printf("error internal server error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			helper.ErrorResponse{
				Status:  "failed",
				Message: "StatusInternalServerError",
				Errors:  err,
			},
		)
	}
}
