package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/rehandwi03/test-case-backend-majoo/criteria"
	custom_error "github.com/rehandwi03/test-case-backend-majoo/internal/error"
	"github.com/rehandwi03/test-case-backend-majoo/model"
	"github.com/rehandwi03/test-case-backend-majoo/repository"
	"github.com/rehandwi03/test-case-backend-majoo/request"
	"github.com/rehandwi03/test-case-backend-majoo/response"
	"github.com/rehandwi03/test-case-backend-majoo/util"
	"gorm.io/gorm"
)

type ProductService interface {
	SaveProduct(ctx context.Context, request *request.ProductAddRequest) (uuid.UUID, error)
	UpdateProduct(ctx context.Context, request *request.ProductUpdateRequest) (uuid.UUID, error)
	DeleteProduct(ctx context.Context, params map[string]interface{}) error
	GetByParam(ctx context.Context, params map[string]interface{}) (*response.ProductResponse, error)
	Fetch(ctx context.Context, ProductCriteria criteria.ProductCriteria) (*util.PaginationResponse, error)
	SaveProductIDImage(ctx context.Context, productId string, fileName string) error
}

type productService struct {
	productRepo repository.ProductRepository
	outletRepo  repository.OutletRepository
}

func NewProductService(
	productRepository repository.ProductRepository, outletRepository repository.OutletRepository,
) ProductService {
	return &productService{productRepo: productRepository, outletRepo: outletRepository}
}

func (p *productService) SaveProductIDImage(ctx context.Context, productId string, fileName string) error {
	params := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?": productId,
			},
		},
	}

	product, err := p.productRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &custom_error.NotFoundError{Message: "product not found"}
		}
	}

	product.Image = fileName

	_, err = p.productRepo.Save(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (p *productService) SaveProduct(ctx context.Context, request *request.ProductAddRequest) (uuid.UUID, error) {
	params := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?": request.OutletID,
			},
		},
	}
	_, err := p.outletRepo.GetByParams(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: "user not found"}
		}

		return uuid.Nil, err
	}

	res, err := p.productRepo.Save(
		ctx, model.Product{
			OutletID:    request.OutletID,
			Name:        request.Name,
			Description: request.Description,
			Stock:       request.Stock,
			Price:       request.Price,
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
}

func (p *productService) UpdateProduct(ctx context.Context, request *request.ProductUpdateRequest) (
	uuid.UUID, error,
) {
	param := map[string]interface{}{
		"where": map[string]interface{}{
			"or": map[string]interface{}{},
			"default": map[string]interface{}{
				"id = ?": request.ID,
			},
		},
	}

	outletParam := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?": request.OutletID,
			},
		},
	}

	_, err := p.outletRepo.GetByParam(ctx, outletParam)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: "user not found"}
		}

		return uuid.Nil, err
	}

	ProductData, err := p.productRepo.GetByParam(ctx, param)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: err.Error()}
		}

		return uuid.Nil, err
	}

	res, err := p.productRepo.Save(
		ctx, model.Product{
			ID:          ProductData.ID,
			OutletID:    request.OutletID,
			Name:        request.Name,
			Description: request.Description,
			Stock:       request.Stock,
			Price:       request.Price,
			Audit: model.Audit{
				CreatedAt: ProductData.CreatedAt,
			},
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
}

func (p *productService) DeleteProduct(ctx context.Context, params map[string]interface{}) error {
	ProductData, err := p.productRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &custom_error.NotFoundError{Message: err.Error()}
		}

		return err
	}

	err = p.productRepo.Delete(ctx, &ProductData)
	if err != nil {
		return err
	}

	return nil
}

func (p *productService) GetByParam(ctx context.Context, params map[string]interface{}) (
	*response.ProductResponse, error,
) {
	productData, err := p.productRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &custom_error.NotFoundError{Message: "product not found"}
		}

		return nil, err
	}

	response := new(response.ProductResponse)
	response.ID = productData.ID
	response.OutletID = productData.OutletID
	response.Name = productData.Name
	response.Description = productData.Description
	response.Stock = productData.Stock
	response.Price = productData.Price
	response.CreatedAt = productData.CreatedAt.Time

	return response, nil
}

func (p *productService) Fetch(ctx context.Context, criteria criteria.ProductCriteria) (
	*util.PaginationResponse, error,
) {
	params := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
			},
			"pagination": map[string]interface{}{
				"page":  criteria.Pagination.Page,
				"sort":  criteria.Pagination.Sort,
				"limit": criteria.Pagination.Limit,
			},
		},
	}

	if criteria.Name != "" {
		params["where"].(map[string]interface{})["default"] = map[string]interface{}{
			"name ILIKE ?": "%" + criteria.Name + "%",
		}
	}
	if criteria.Stock != "" {
		params["where"].(map[string]interface{})["default"] = map[string]interface{}{
			"stock = ?": criteria.Stock,
		}
	}
	if criteria.Price != "" {
		params["where"].(map[string]interface{})["default"] = map[string]interface{}{
			"price = ?": criteria.Price,
		}
	}

	res, rowCount, err := p.productRepo.Fetch(ctx, params)
	if err != nil {
		return nil, err
	}

	var responseData []response.ProductResponse
	for _, val := range res {
		var data response.ProductResponse

		data.ID = val.ID
		data.OutletID = val.OutletID
		data.Name = val.Name
		data.Description = val.Description
		data.Stock = val.Stock
		data.Price = val.Price
		data.CreatedAt = val.CreatedAt.Time

		responseData = append(responseData, data)
	}

	resPagination := util.BuildPagination(criteria.Pagination, responseData, rowCount)

	return &resPagination, nil
}
