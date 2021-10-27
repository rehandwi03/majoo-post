package http

import (
	"github.com/gofiber/fiber/v2"
	custom_error "github.com/rehandwi03/test-case-backend-majoo/internal/error"
	"github.com/rehandwi03/test-case-backend-majoo/internal/helper"
	request2 "github.com/rehandwi03/test-case-backend-majoo/request"
	"github.com/rehandwi03/test-case-backend-majoo/service"
	"log"
)

type authHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(app fiber.Router, authService service.AuthService) {
	handler := authHandler{authSvc: authService}
	app.Post("/login", handler.Login)
}

func (a *authHandler) Login(c *fiber.Ctx) error {
	request := new(request2.LoginRequest)

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

	res, err := a.authSvc.Login(c.Context(), request)
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
		return c.Status(fiber.StatusOK).JSON(
			helper.SuccessResponse{
				Status: "success", Message: "success login", Data: res,
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
