package params

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"per_page"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}
