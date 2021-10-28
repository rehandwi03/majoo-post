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

type userHandler struct {
	userSvc service.UserService
}

func NewUserHandler(app fiber.Router, userService service.UserService) {
	handler := userHandler{userSvc: userService}

	app.Post("/users", middleware.JwtProtected(), handler.saveUser)
	app.Get("/users/:id", middleware.JwtProtected(), handler.getByID)
	app.Put("/users", middleware.JwtProtected(), handler.updateUser)
	app.Delete("/users/:id", middleware.JwtProtected(), handler.deleteByID)
	app.Get("/users", middleware.JwtProtected(), handler.fetch)
}

func (u *userHandler) fetch(c *fiber.Ctx) error {
	pagination := util.GeneratePaginationFromRequest(c)

	userCriteria := criteria.UserCriteria{
		Pagination: pagination,
	}

	userCriteria.Email = c.Query("email")
	userCriteria.PhoneNumber = c.Query("phoneNumber")

	res, err := u.userSvc.Fetch(c.Context(), userCriteria)
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

func (u *userHandler) deleteByID(c *fiber.Ctx) error {
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

	err := u.userSvc.DeleteUser(c.Context(), params)
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

func (u *userHandler) getByID(c *fiber.Ctx) error {
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

	res, err := u.userSvc.GetByParam(c.Context(), params)
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

func (u *userHandler) updateUser(c *fiber.Ctx) error {
	request := new(request2.UserUpdateRequest)

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

	res, err := u.userSvc.UpdateUser(c.Context(), request)
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
					"user_id": res,
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

func (u *userHandler) saveUser(c *fiber.Ctx) error {
	request := new(request2.UserAddRequest)

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

	res, err := u.userSvc.SaveUser(c.Context(), request)
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
					"user_id": res,
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
