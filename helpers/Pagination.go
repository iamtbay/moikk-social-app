package helpers

import (
	"math"
)

func Pagination(page, totalCount, postPerPage int) (int, int) {
	totalPage := int(math.Ceil(float64(totalCount) / float64(postPerPage)))
	//
	if page > totalPage {
		return totalPage, totalPage
	} else if page < 1 {
		return 1, totalPage
	}
	return page, totalPage
}
