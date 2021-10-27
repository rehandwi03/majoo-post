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

type OutletService interface {
	SaveOutlet(ctx context.Context, request *request.OutletAddRequest) (uuid.UUID, error)
	UpdateOutlet(ctx context.Context, request *request.OutletUpdateRequest) (uuid.UUID, error)
	DeleteOutlet(ctx context.Context, params map[string]interface{}) error
	GetByParam(ctx context.Context, params map[string]interface{}) (*response.OutletResponse, error)
	Fetch(ctx context.Context, OutletCriteria criteria.OutletCriteria) (*util.PaginationResponse, error)
}

type outletService struct {
	outletRepo   repository.OutletRepository
	merchantRepo repository.MerchantRepository
}

func NewOutletService(
	outletRepository repository.OutletRepository, merchantRepository repository.MerchantRepository,
) OutletService {
	return &outletService{outletRepo: outletRepository, merchantRepo: merchantRepository}
}

func (o *outletService) SaveOutlet(ctx context.Context, request *request.OutletAddRequest) (uuid.UUID, error) {
	params := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?": request.MerchantID,
			},
		},
	}
	_, err := o.merchantRepo.GetByParams(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: "merchant not found"}
		}

		return uuid.Nil, err
	}

	res, err := o.outletRepo.Save(
		ctx, model.Outlet{
			MerchantID:  request.MerchantID,
			Name:        request.Name,
			Location:    request.Location,
			PhoneNumber: request.PhoneNumber,
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
}

func (o *outletService) UpdateOutlet(ctx context.Context, request *request.OutletUpdateRequest) (
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

	merchantParam := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?": request.MerchantID,
			},
		},
	}

	_, err := o.merchantRepo.GetByParam(ctx, merchantParam)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: "merchant not found"}
		}

		return uuid.Nil, err
	}

	outletData, err := o.outletRepo.GetByParam(ctx, param)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: "outlet not found"}
		}

		return uuid.Nil, err
	}

	res, err := o.outletRepo.Save(
		ctx, model.Outlet{
			ID:          request.ID,
			MerchantID:  request.MerchantID,
			Name:        request.Name,
			Location:    request.Location,
			PhoneNumber: request.PhoneNumber,
			Audit: model.Audit{
				CreatedAt: outletData.CreatedAt,
			},
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
}

func (o *outletService) DeleteOutlet(ctx context.Context, params map[string]interface{}) error {
	OutletData, err := o.outletRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &custom_error.NotFoundError{Message: "outlet not found"}
		}

		return err
	}

	err = o.outletRepo.Delete(ctx, &OutletData)
	if err != nil {
		return err
	}

	return nil
}

func (o *outletService) GetByParam(ctx context.Context, params map[string]interface{}) (
	*response.OutletResponse, error,
) {
	outletData, err := o.outletRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &custom_error.NotFoundError{Message: err.Error()}
		}

		return nil, err
	}

	response := new(response.OutletResponse)
	response.ID = outletData.ID
	response.MerchantID = outletData.MerchantID
	response.Name = outletData.Name
	response.Location = outletData.Location
	response.PhoneNumber = outletData.PhoneNumber
	response.CreatedAt = outletData.CreatedAt.Time

	return response, nil
}

func (o *outletService) Fetch(ctx context.Context, criteria criteria.OutletCriteria) (
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
			"email ILIKE ?": "%" + criteria.
				Name + "%",
		}
	}

	if criteria.Location != "" {
		params["where"].(map[string]interface{})["default"] = map[string]interface{}{
			"location ILIKE ?": "%" + criteria.
				Location + "%",
		}
	}

	res, rowCount, err := o.outletRepo.Fetch(ctx, params)
	if err != nil {
		return nil, err
	}

	var responseData []response.OutletResponse
	for _, val := range res {
		var data response.OutletResponse

		data.ID = val.ID
		data.MerchantID = val.MerchantID
		data.Name = val.Name
		data.Location = val.Location
		data.PhoneNumber = val.PhoneNumber
		data.CreatedAt = val.CreatedAt.Time

		responseData = append(responseData, data)
	}

	resPagination := util.BuildPagination(criteria.Pagination, responseData, rowCount)

	return &resPagination, nil
}
