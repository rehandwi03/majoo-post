package criteria

import "github.com/rehandwi03/test-case-backend-majoo/util"

type OutletCriteria struct {
	Name       string `json:"name"`
	Location   string `json:"location"`
	Pagination util.Pagination
}
