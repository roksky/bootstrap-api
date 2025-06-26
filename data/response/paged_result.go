package response

type PagedResult[T any] struct {
	TotalItems int64 `json:"totalItems"`
	PageNumber int   `json:"pageNumber"`
	PageSize   int64 `json:"pageSize"`
	Items      []T   `json:"items"`
}
