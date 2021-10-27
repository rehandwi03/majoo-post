package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/rehandwi03/test-case-backend-majoo/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(ctx context.Context, user model.User) (uuid.UUID, error)
	GetByParam(ctx context.Context, params map[string]interface{}) (res model.User, err error)
	GetByParams(ctx context.Context, params map[string]interface{}) (res []model.User, err error)
	Fetch(ctx context.Context, params map[string]interface{}) (res []model.User, count int64, err error)
	Delete(ctx context.Context, data *model.User) error
}

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &userRepository{conn: conn}
}

func (u userRepository) Fetch(ctx context.Context, params map[string]interface{}) (
	res []model.User, count int64, err error,
) {
	res, err = u.GetByParams(ctx, params)
	if err != nil {
		return res, count, err
	}

	done := make(chan bool, 1)
	u.countRecords(ctx, model.User{}, done, &count, params)

	<-done

	return res, count, nil
}

func (u userRepository) Save(ctx context.Context, user model.User) (uuid.UUID, error) {
	err := u.conn.WithContext(ctx).Save(&user).Error
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (u userRepository) GetByParam(ctx context.Context, params map[string]interface{}) (res model.User, err error) {
	query := u.conn.WithContext(ctx)
	if params["where"] != nil && params["where"].(map[string]interface{})["default"] != nil {
		for field, value := range params["where"].(map[string]interface{})["default"].(map[string]interface{}) {
			query = query.Where(field, value)
		}
	}

	if params["where"] != nil && params["where"].(map[string]interface{})["or"] != nil {
		for field, value := range params["where"].(map[string]interface{})["or"].(map[string]interface{}) {
			query = query.Or(field, value)
		}
	}

	err = query.First(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u userRepository) GetByParams(ctx context.Context, params map[string]interface{}) (res []model.User, err error) {
	query := u.conn.WithContext(ctx)
	if params["where"] != nil && params["where"].(map[string]interface{})["default"] != nil {
		for field, value := range params["where"].(map[string]interface{})["default"].(map[string]interface{}) {
			query = query.Where(field, value)
		}
	}

	if params["where"] != nil && params["where"].(map[string]interface{})["or"] != nil {
		for field, value := range params["where"].(map[string]interface{})["or"].(map[string]interface{}) {
			query = query.Or(field, value)
		}
	}

	if params["where"] != nil && params["where"].(map[string]interface{})["pagination"] != nil {
		page := params["where"].(map[string]interface{})["pagination"].(map[string]interface{})["page"].(int)
		limit := params["where"].(map[string]interface{})["pagination"].(map[string]interface{})["limit"].(int)
		sort := params["where"].(map[string]interface{})["pagination"].(map[string]interface{})["sort"].(string)

		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset).Order(sort)
	}

	err = query.Find(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func (u userRepository) Delete(ctx context.Context, data *model.User) error {
	err := u.conn.WithContext(ctx).Delete(data).Error
	if err != nil {
		return err
	}

	return nil
}

func (p userRepository) countRecords(
	ctx context.Context, countDataSource interface{}, done chan bool,
	count *int64, params map[string]interface{},
) {
	query := p.conn.WithContext(ctx)
	if params["where"] != nil && params["where"].(map[string]interface{})["default"] != nil {
		for whereKey, whereValue := range params["where"].(map[string]interface{})["default"].(map[string]interface{}) {
			query = query.Where(whereKey, whereValue)
		}
	}

	query.Model(countDataSource).Count(count)
	done <- true
}
