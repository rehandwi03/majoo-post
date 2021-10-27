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

type UserService interface {
	SaveUser(ctx context.Context, request *request.UserAddRequest) (uuid.UUID, error)
	UpdateUser(ctx context.Context, request *request.UserUpdateRequest) (uuid.UUID, error)
	DeleteUser(ctx context.Context, params map[string]interface{}) error
	GetByParam(ctx context.Context, params map[string]interface{}) (*response.UserResponse, error)
	Fetch(ctx context.Context, userCriteria criteria.UserCriteria) (*util.PaginationResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepo: userRepository}
}

func (u *userService) SaveUser(ctx context.Context, request *request.UserAddRequest) (uuid.UUID, error) {
	param := map[string]interface{}{
		"where": map[string]interface{}{
			"or": map[string]interface{}{},
			"default": map[string]interface{}{
				"email = ?": request.Email,
			},
		},
	}

	_, err := u.userRepo.GetByParam(ctx, param)
	if err == nil {
		return uuid.Nil, &custom_error.BadRequest{Message: "email already exist"}
	}

	userModel := model.User{
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
	}

	if err := userModel.EncryptPassword(); err != nil {
		return uuid.Nil, err
	}
	res, err := u.userRepo.Save(
		ctx, userModel,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
}

func (u *userService) UpdateUser(ctx context.Context, request *request.UserUpdateRequest) (uuid.UUID, error) {
	param := map[string]interface{}{
		"where": map[string]interface{}{
			"or": map[string]interface{}{},
			"default": map[string]interface{}{
				"id = ?": request.ID,
			},
		},
	}

	emailParam := map[string]interface{}{
		"where": map[string]interface{}{
			"or": map[string]interface{}{},
			"default": map[string]interface{}{
				"id = ?": request.ID,
			},
		},
	}

	userData, err := u.userRepo.GetByParam(ctx, param)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, &custom_error.NotFoundError{Message: err.Error()}
		}

		return uuid.Nil, err
	}

	if request.Email != userData.Email {
		_, err = u.userRepo.GetByParam(ctx, emailParam)
		if err == nil {
			return uuid.Nil, &custom_error.BadRequest{Message: "email already exist"}
		}
	}

	userModel := model.User{
		ID:          userData.ID,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		Email:       request.Email,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
	}
	err = userModel.EncryptPassword()
	if err != nil {
		return uuid.Nil, err
	}

	res, err := u.userRepo.Save(
		ctx, userModel,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return res, nil
}

func (u *userService) DeleteUser(ctx context.Context, params map[string]interface{}) error {
	userData, err := u.userRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &custom_error.NotFoundError{Message: err.Error()}
		}

		return err
	}

	err = u.userRepo.Delete(ctx, &userData)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) GetByParam(ctx context.Context, params map[string]interface{}) (*response.UserResponse, error) {
	userData, err := u.userRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &custom_error.NotFoundError{Message: err.Error()}
		}

		return nil, err
	}

	response := new(response.UserResponse)
	response.ID = userData.ID
	response.FirstName = userData.FirstName
	response.LastName = userData.LastName
	response.Email = userData.Email
	response.PhoneNumber = userData.PhoneNumber
	response.CreatedAt = userData.CreatedAt.Time

	return response, nil
}

func (u *userService) Fetch(ctx context.Context, criteria criteria.UserCriteria) (*util.PaginationResponse, error) {
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

	if criteria.Email != "" {
		params["where"].(map[string]interface{})["default"] = map[string]interface{}{
			"email ILIKE ?": "%" + criteria.
				Email + "%",
		}
	}

	if criteria.PhoneNumber != "" {
		params["where"].(map[string]interface{})["default"] = map[string]interface{}{
			"phone_number LIKE" +
				" ?": "%" + criteria.
				PhoneNumber + "%",
		}
	}

	res, rowCount, err := u.userRepo.Fetch(ctx, params)
	if err != nil {
		return nil, err
	}

	var responseData []response.UserResponse
	for _, val := range res {
		var data response.UserResponse

		data.ID = val.ID
		data.FirstName = val.FirstName
		data.LastName = val.LastName
		data.Email = val.Email
		data.PhoneNumber = val.PhoneNumber
		data.CreatedAt = val.CreatedAt.Time

		responseData = append(responseData, data)
	}

	resPagination := util.BuildPagination(criteria.Pagination, responseData, rowCount)

	return &resPagination, nil
}
