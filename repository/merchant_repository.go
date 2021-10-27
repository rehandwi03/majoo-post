package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/rehandwi03/test-case-backend-majoo/model"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	Save(ctx context.Context, user model.Merchant) (uuid.UUID, error)
	GetByParam(ctx context.Context, params map[string]interface{}) (model.Merchant, error)
	GetByParams(ctx context.Context, params map[string]interface{}) ([]model.Merchant, error)
	Delete(ctx context.Context, data *model.Merchant) error
	Fetch(ctx context.Context, params map[string]interface{}) (res []model.Merchant, count int64, err error)
}

type merchantRepository struct {
	conn *gorm.DB
}

func NewMerchantRepository(conn *gorm.DB) MerchantRepository {
	return &merchantRepository{conn: conn}
}

func (m merchantRepository) Fetch(ctx context.Context, params map[string]interface{}) (
	res []model.Merchant, count int64, err error,
) {
	res, err = m.GetByParams(ctx, params)
	if err != nil {
		return res, count, err
	}

	done := make(chan bool, 1)
	m.countRecords(ctx, model.Merchant{}, done, &count, params)

	<-done

	return res, count, nil
}

func (m merchantRepository) Save(ctx context.Context, merchant model.Merchant) (uuid.UUID, error) {
	err := m.conn.WithContext(ctx).Save(&merchant).Error
	if err != nil {
		return uuid.Nil, err
	}

	return merchant.ID, nil
}

func (m merchantRepository) GetByParam(ctx context.Context, params map[string]interface{}) (
	res model.Merchant, err error,
) {
	query := m.conn.WithContext(ctx)
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

func (m merchantRepository) GetByParams(ctx context.Context, params map[string]interface{}) (
	res []model.Merchant, err error,
) {
	query := m.conn.WithContext(ctx)
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

func (m merchantRepository) Delete(ctx context.Context, data *model.Merchant) error {
	err := m.conn.WithContext(ctx).Delete(data).Error
	if err != nil {
		return err
	}

	return nil
}

func (m merchantRepository) countRecords(
	ctx context.Context, countDataSource interface{}, done chan bool,
	count *int64, params map[string]interface{},
) {
	query := m.conn.WithContext(ctx)
	if params["where"] != nil && params["where"].(map[string]interface{})["default"] != nil {
		for whereKey, whereValue := range params["where"].(map[string]interface{})["default"].(map[string]interface{}) {
			query = query.Where(whereKey, whereValue)
		}
	}

	query.Model(countDataSource).Count(count)
	done <- true
}
