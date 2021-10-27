package util

import (
	"github.com/gofiber/fiber/v2"
	"math"
	"strconv"
	"strings"
)

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

type PaginationResponse struct {
	Data   interface{} `json:"data"`
	Paging Paging      `json:"paging"`
}

type Paging struct {
	TotalRecord int    `json:"total_record"`
	TotalPage   int    `json:"total_page"`
	Page        int    `json:"page"`
	OrderBy     string `json:"order_by"`
	SortBy      string `json:"sort_by"`
	Size        int    `json:"size"`
}

func GeneratePaginationFromRequest(c *fiber.Ctx) Pagination {
	limit, _ := strconv.Atoi(c.Query("limit", "1"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	sort := c.Query("sort", "created_at asc")

	return Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

func BuildPagination(
	pagination Pagination, data interface{}, totalRow int64,
) PaginationResponse {

	splitOrderAndSort := strings.Fields(pagination.Sort)

	var response PaginationResponse

	response.Data = data
	response.Paging.TotalRecord = int(totalRow)
	response.Paging.Page = pagination.Page
	response.Paging.Size = pagination.Limit
	totalPage := int(math.Ceil(float64(totalRow)) / float64(pagination.Limit))
	if (((pagination.Limit * totalPage) - int(totalRow)) * -1) > 0 {
		totalPage++
	}
	if totalRow == 1 {
		totalPage = 1
	}

	response.Paging.TotalPage = totalPage
	response.Paging.OrderBy = splitOrderAndSort[0]
	response.Paging.SortBy = splitOrderAndSort[1]

	return response
}
