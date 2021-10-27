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

type MerchantService interface {
	SaveMerchant(ctx context.Context, request *request.MerchantAddRequest) (uuid.UUID, error)
	UpdateMerchant(ctx context.Context, request *request.MerchantUpdateRequest) (uuid.UUID, error)
	DeleteMerchant(ctx context.Context, params map[string]interface{}) error
	GetByParam(ctx context.Context, params map[string]interface{}) (*response.MerchantResponse, error)
	Fetch(ctx context.Context, MerchantCriteria criteria.MerchantCriteria) (*util.PaginationResponse, error)
}

type merchantService struct {
	merchantRepo repository.MerchantRepository
	userRepo     repository.UserRepository
	userID       uuid.UUID
}

func NewMerchantService(
	merchantRepository repository.MerchantRepository, userRepository repository.UserRepository,
) MerchantService {
	return &merchantService{merchantRepo: merchantRepository, userRepo: userRepository}
}

func (m *merchantService) SaveMerchant(ctx context.Context, request *request.MerchantAddRequest) (uuid.UUID, error) {
	userId, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, &custom_error.NotFoundError{Message: "user id not found"}
	}

	params := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?": userId,
			},
		},
	}
	_, err := m.userRepo.GetByParams(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: "user not found"}
		}

		return uuid.Nil, err
	}

	res, err := m.merchantRepo.Save(
		ctx, model.Merchant{
			UserID:          userId,
			Name:            request.Name,
			InstitutionName: request.InstitutionName,
			PhoneNumber:     request.PhoneNumber,
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
}

func (m *merchantService) UpdateMerchant(ctx context.Context, request *request.MerchantUpdateRequest) (
	uuid.UUID, error,
) {
	userId, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, &custom_error.NotFoundError{Message: "user id not found"}
	}
	merchantParams := map[string]interface{}{
		"where": map[string]interface{}{
			"or": map[string]interface{}{},
			"default": map[string]interface{}{
				"id = ?":      request.ID,
				"user_id = ?": userId,
			},
		},
	}

	userParams := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"id = ?": userId,
			},
		},
	}

	_, err := m.userRepo.GetByParam(ctx, userParams)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: "user not found"}
		}

		return uuid.Nil, err
	}

	merchantData, err := m.merchantRepo.GetByParam(ctx, merchantParams)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: err.Error()}
		}

		return uuid.Nil, err
	}

	res, err := m.merchantRepo.Save(
		ctx, model.Merchant{
			ID:              merchantData.ID,
			UserID:          userId,
			Name:            request.Name,
			InstitutionName: request.InstitutionName,
			PhoneNumber:     request.PhoneNumber,
			Audit: model.Audit{
				CreatedAt: merchantData.CreatedAt,
			},
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
}

func (m *merchantService) DeleteMerchant(ctx context.Context, params map[string]interface{}) error {
	MerchantData, err := m.merchantRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &custom_error.NotFoundError{Message: err.Error()}
		}

		return err
	}

	err = m.merchantRepo.Delete(ctx, &MerchantData)
	if err != nil {
		return err
	}

	return nil
}

func (m *merchantService) GetByParam(ctx context.Context, params map[string]interface{}) (
	*response.MerchantResponse, error,
) {
	merchantData, err := m.merchantRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &custom_error.NotFoundError{Message: "merchant not found"}
		}

		return nil, err
	}

	response := new(response.MerchantResponse)
	response.ID = merchantData.ID
	response.UserID = merchantData.UserID
	response.Name = merchantData.Name
	response.InstitutionName = merchantData.InstitutionName
	response.PhoneNumber = merchantData.PhoneNumber
	response.CreatedAt = merchantData.CreatedAt.Time

	return response, nil
}

func (m *merchantService) Fetch(ctx context.Context, criteria criteria.MerchantCriteria) (
	*util.PaginationResponse, error,
) {
	userId, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return nil, &custom_error.NotFoundError{Message: "user id not found"}
	}

	params := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"user_id = ?": userId,
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
	if criteria.InstitutionName != "" {
		params["where"].(map[string]interface{})["default"] = map[string]interface{}{
			"institution_name ILIKE ?": "%" + criteria.InstitutionName + "%",
		}
	}

	res, rowCount, err := m.merchantRepo.Fetch(ctx, params)
	if err != nil {
		return nil, err
	}

	var responseData []response.MerchantResponse
	for _, val := range res {
		var data response.MerchantResponse

		data.ID = val.ID
		data.UserID = val.UserID
		data.Name = val.Name
		data.InstitutionName = val.InstitutionName
		data.PhoneNumber = val.PhoneNumber
		data.CreatedAt = val.CreatedAt.Time

		responseData = append(responseData, data)
	}

	resPagination := util.BuildPagination(criteria.Pagination, responseData, rowCount)

	return &resPagination, nil
}
