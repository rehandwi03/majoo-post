package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rehandwi03/test-case-backend-majoo/criteria"
	custom_error "github.com/rehandwi03/test-case-backend-majoo/internal/error"
	"github.com/rehandwi03/test-case-backend-majoo/internal/helper"
	"github.com/rehandwi03/test-case-backend-majoo/internal/middleware"
	request2 "github.com/rehandwi03/test-case-backend-majoo/request"
	"github.com/rehandwi03/test-case-backend-majoo/service"
	"github.com/rehandwi03/test-case-backend-majoo/util"
	"log"
)

type outletHandler struct {
	outletSvc service.OutletService
}

func NewOutletHandler(app fiber.Router, outletService service.OutletService) {
	handler := outletHandler{outletSvc: outletService}

	app.Post("/outlets", middleware.JwtProtected(), handler.saveOutlet)
	app.Get("/outlets/:id", middleware.JwtProtected(), handler.getByID)
	app.Put("/outlets", middleware.JwtProtected(), handler.updateOutlet)
	app.Delete("/outlets/:id", middleware.JwtProtected(), handler.deleteByID)
	app.Get("/outlets", middleware.JwtProtected(), handler.fetch)
}

func (o *outletHandler) fetch(c *fiber.Ctx) error {
	pagination := util.GeneratePaginationFromRequest(c)

	OutletCriteria := criteria.OutletCriteria{
		Pagination: pagination,
	}

	OutletCriteria.Name = c.Query("name")
	OutletCriteria.Location = c.Query("location")

	res, err := o.outletSvc.Fetch(c.Context(), OutletCriteria)
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

func (o *outletHandler) deleteByID(c *fiber.Ctx) error {
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

	params := map[string]interface{}{
		"or": map[string]interface{}{},
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?": id,
			},
		},
	}

	err := o.outletSvc.DeleteOutlet(c.Context(), params)
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

func (o *outletHandler) getByID(c *fiber.Ctx) error {
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

	params := map[string]interface{}{
		"or": map[string]interface{}{},
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?": id,
			},
		},
	}

	res, err := o.outletSvc.GetByParam(c.Context(), params)
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

func (o *outletHandler) updateOutlet(c *fiber.Ctx) error {
	request := new(request2.OutletUpdateRequest)

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

	res, err := o.outletSvc.UpdateOutlet(c.Context(), request)
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
					"outlet_id": res,
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

func (o *outletHandler) saveOutlet(c *fiber.Ctx) error {
	request := new(request2.OutletAddRequest)

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

	res, err := o.outletSvc.SaveOutlet(c.Context(), request)
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
					"outlet_id": res,
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
