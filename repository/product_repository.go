package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/rehandwi03/test-case-backend-majoo/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Save(ctx context.Context, product model.Product) (uuid.UUID, error)
	GetByParam(ctx context.Context, params map[string]interface{}) (model.Product, error)
	GetByParams(ctx context.Context, params map[string]interface{}) ([]model.Product, error)
	Delete(ctx context.Context, data *model.Product) error
	Fetch(ctx context.Context, params map[string]interface{}) (res []model.Product, count int64, err error)
}

type productRepository struct {
	conn *gorm.DB
}

func NewProductRepository(conn *gorm.DB) ProductRepository {
	return &productRepository{conn: conn}
}

func (p productRepository) Fetch(ctx context.Context, params map[string]interface{}) (
	res []model.Product, count int64, err error,
) {
	res, err = p.GetByParams(ctx, params)
	if err != nil {
		return res, count, err
	}

	done := make(chan bool, 1)
	p.countRecords(ctx, model.Product{}, done, &count, params)

	<-done

	return res, count, nil
}

func (p productRepository) Save(ctx context.Context, product model.Product) (uuid.UUID, error) {
	err := p.conn.WithContext(ctx).Save(&product).Error
	if err != nil {
		return uuid.Nil, err
	}

	return product.ID, nil
}

func (p productRepository) GetByParam(ctx context.Context, params map[string]interface{}) (
	res model.Product, err error,
) {
	query := p.conn.WithContext(ctx)
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

func (p productRepository) GetByParams(ctx context.Context, params map[string]interface{}) (
	res []model.Product, err error,
) {
	query := p.conn.WithContext(ctx)
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

func (p productRepository) Delete(ctx context.Context, data *model.Product) error {
	err := p.conn.WithContext(ctx).Delete(data).Error
	if err != nil {
		return err
	}

	return nil
}

func (p productRepository) countRecords(
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
