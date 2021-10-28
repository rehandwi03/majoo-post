package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rehandwi03/test-case-backend-majoo/criteria"
	custom_error "github.com/rehandwi03/test-case-backend-majoo/internal/error"
	"github.com/rehandwi03/test-case-backend-majoo/internal/helper"
	"github.com/rehandwi03/test-case-backend-majoo/internal/middleware"
	request2 "github.com/rehandwi03/test-case-backend-majoo/request"
	"github.com/rehandwi03/test-case-backend-majoo/service"
	"github.com/rehandwi03/test-case-backend-majoo/util"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type productHandler struct {
	productSvc service.ProductService
}

func NewProductHandler(app fiber.Router, productService service.ProductService) {
	handler := productHandler{productSvc: productService}

	app.Post("/products", middleware.JwtProtected(), handler.saveProduct)
	app.Get("/products/:id", middleware.JwtProtected(), handler.getByID)
	app.Put("/products", middleware.JwtProtected(), handler.updateProduct)
	app.Delete("/products/:id", middleware.JwtProtected(), handler.deleteByID)
	app.Get("/products", middleware.JwtProtected(), handler.fetch)
	app.Post("/products/image", middleware.JwtProtected(), handler.uploadImage)

}

func (p *productHandler) uploadImage(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Message: "StatusBadRequest",
				Status:  "failed",
				Errors:  "image value is null",
			},
		)
	}

	productId := form.Value["product_id"]

	if len(productId) < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.ErrorResponse{
				Message: "StatusBadRequest",
				Status:  "failed",
				Errors:  "product id not found",
			},
		)
	}

	files := form.File["image"]

	for _, file := range files {
		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])

		header := file.Header["Content-Type"][0]

		splitHeader := strings.Split(header, "/")
		if splitHeader[0] != "image" {
			return &custom_error.BadRequest{Message: "file format isn't image"}

		}

		fileName := strconv.Itoa(rand.Int()) + file.Filename

		if err := c.SaveFile(
			file, fmt.Sprintf("./internal/file/%s", fileName),
		); err != nil {
			return err
		}

		err = p.productSvc.SaveProductIDImage(c.Context(), productId[0], fileName)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				helper.ErrorResponse{
					Message: "StatusInternalServerError",
					Status:  "failed",
					Errors:  "internal server error",
				},
			)
		}
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.SuccessResponse{
			Status:  "success",
			Message: "success upload image",
		},
	)

}

func (p *productHandler) fetch(c *fiber.Ctx) error {
	pagination := util.GeneratePaginationFromRequest(c)

	productCriteria := criteria.ProductCriteria{
		Pagination: pagination,
	}

	productCriteria.Name = c.Query("name")
	productCriteria.Stock = c.Query("stock")
	productCriteria.Price = c.Query("price")

	res, err := p.productSvc.Fetch(c.Context(), productCriteria)
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

func (p *productHandler) deleteByID(c *fiber.Ctx) error {
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

	err := p.productSvc.DeleteProduct(c.Context(), params)
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

func (p *productHandler) getByID(c *fiber.Ctx) error {
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

	res, err := p.productSvc.GetByParam(c.Context(), params)
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

func (p *productHandler) updateProduct(c *fiber.Ctx) error {
	request := new(request2.ProductUpdateRequest)

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

	res, err := p.productSvc.UpdateProduct(c.Context(), request)
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
					"product_id": res,
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

func (p *productHandler) saveProduct(c *fiber.Ctx) error {
	request := new(request2.ProductAddRequest)

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

	res, err := p.productSvc.SaveProduct(c.Context(), request)
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
					"product_id": res,
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
