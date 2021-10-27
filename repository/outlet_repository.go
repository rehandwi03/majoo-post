package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/rehandwi03/test-case-backend-majoo/model"
	"gorm.io/gorm"
)

type OutletRepository interface {
	Save(ctx context.Context, user model.Outlet) (uuid.UUID, error)
	GetByParam(ctx context.Context, params map[string]interface{}) (model.Outlet, error)
	GetByParams(ctx context.Context, params map[string]interface{}) ([]model.Outlet, error)
	Delete(ctx context.Context, data *model.Outlet) error
	Fetch(ctx context.Context, params map[string]interface{}) (res []model.Outlet, count int64, err error)
}

type outletRepository struct {
	conn *gorm.DB
}

func NewOutletRepository(conn *gorm.DB) OutletRepository {
	return &outletRepository{conn: conn}
}

func (o outletRepository) Fetch(ctx context.Context, params map[string]interface{}) (
	res []model.Outlet, count int64, err error,
) {
	res, err = o.GetByParams(ctx, params)
	if err != nil {
		return res, count, err
	}

	done := make(chan bool, 1)
	o.countRecords(ctx, model.Outlet{}, done, &count, params)

	<-done

	return res, count, nil
}

func (o outletRepository) Save(ctx context.Context, outlet model.Outlet) (uuid.UUID, error) {
	err := o.conn.WithContext(ctx).Save(&outlet).Error
	if err != nil {
		return uuid.Nil, err
	}

	return outlet.ID, nil
}

func (o outletRepository) GetByParam(ctx context.Context, params map[string]interface{}) (
	res model.Outlet, err error,
) {
	query := o.conn.WithContext(ctx)
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

func (o outletRepository) GetByParams(ctx context.Context, params map[string]interface{}) (
	res []model.Outlet, err error,
) {
	query := o.conn.WithContext(ctx)
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

func (o outletRepository) Delete(ctx context.Context, data *model.Outlet) error {
	err := o.conn.WithContext(ctx).Delete(data).Error
	if err != nil {
		return err
	}

	return nil
}

func (o outletRepository) countRecords(
	ctx context.Context, countDataSource interface{}, done chan bool,
	count *int64, params map[string]interface{},
) {
	query := o.conn.WithContext(ctx)
	if params["where"] != nil && params["where"].(map[string]interface{})["default"] != nil {
		for whereKey, whereValue := range params["where"].(map[string]interface{})["default"].(map[string]interface{}) {
			query = query.Where(whereKey, whereValue)
		}
	}

	query.Model(countDataSource).Count(count)
	done <- true
}
