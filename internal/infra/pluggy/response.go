package pluggy

type paginatedResponse[T any] struct {
	Page       float64 `json:"page"`
	Total      float64 `json:"total"`
	TotalPages float64 `json:"totalPages"`
	Results    []T     `json:"results"`
}
