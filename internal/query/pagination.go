package query

type PaginationInput struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"limit"`
}

type PaginationOutput struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"limit"`
	Total uint64 `json:"total"`
}
