package pluggy

type paginatedResponse[T any] struct {
	Page       int `json:"page"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
	Results    []T `json:"results"`
}
