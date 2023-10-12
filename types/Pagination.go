package mytypes

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	TotalCount  int `json:"total_count"`
}
