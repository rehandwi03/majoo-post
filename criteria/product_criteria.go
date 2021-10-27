package criteria

import "github.com/rehandwi03/test-case-backend-majoo/util"

type ProductCriteria struct {
	Name       string `json:"name"`
	Stock      string `json:"stock"`
	Price      string `json:"price"`
	Pagination util.Pagination
}
