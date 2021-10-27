package criteria

import "github.com/rehandwi03/test-case-backend-majoo/util"

type MerchantCriteria struct {
	Name            string `json:"name"`
	InstitutionName string `json:"institution_name"`
	Pagination      util.Pagination
}
