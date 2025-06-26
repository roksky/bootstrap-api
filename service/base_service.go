package service

import "github.com/roksky/bootstrap-api/data/response"

// BaseService defines a generic interface for basic CRUD operations.
// T represents the type of the item.
// K represents the type of the item's identifier.
// S represents the type of the search parameters.
type BaseService[T any, K comparable, S any] interface {
	// Create adds a new item.
	// filterContext provides additional context for the operation.
	// item is the item to be created.
	// Returns the created item and an error if any.
	Create(filterContext *S, item *T) (*T, error)

	// CreateMany adds multiple new items.
	// filterContext provides additional context for the operation.
	// items is the list of items to be created.
	// Returns the list of created items and an error if any.
	CreateMany(filterContext *S, items []*T) ([]*T, error)

	// Update modifies an existing item.
	// filterContext provides additional context for the operation.
	// item is the item to be updated.
	// Returns the updated item and an error if any.
	Update(filterContext *S, item *T) (*T, error)

	// UpdateMany modifies multiple existing items.
	// filterContext provides additional context for the operation.
	// items is the list of items to be updated.
	// Returns the list of updated items and an error if any.
	UpdateMany(filterContext *S, items []*T) ([]*T, error)

	// Delete removes an item by its identifier.
	// searchParams provides the search parameters.
	// id is the identifier of the item to be deleted.
	// Returns an error if any.
	Delete(searchParams *S, id K) error

	// DeleteMany removes multiple items by their identifiers.
	// searchParams provides the search parameters.
	// ids is the list of identifiers of the items to be deleted.
	// Returns an error if any.
	DeleteMany(searchParams *S, ids []K) error

	// FindById retrieves an item by its identifier.
	// searchParams provides the search parameters.
	// id is the identifier of the item to be retrieved.
	// Returns the found item and an error if any.
	FindById(searchParams *S, id K) (*T, error)

	// FindByIds retrieves multiple items by their identifiers.
	// searchParams provides the search parameters.
	// ids is the list of identifiers of the items to be retrieved.
	// Returns the list of found items and an error if any.
	FindByIds(searchParams *S, ids []K) ([]*T, error)

	// FindAll retrieves all items with pagination.
	// searchParams provides the search parameters.
	// pageSize is the number of items per page.
	// page is the page number.
	// Returns a paged result of found items and an error if any.
	FindAll(searchParams *S, pageSize int, page int) (response.PagedResult[*T], error)

	// Search performs a search based on the search parameters.
	// searchParams provides the search parameters.
	// Returns a paged result of found items and an error if any.
	Search(searchParams *S) (response.PagedResult[*T], error)

	// Deleted retrieves the identifiers of deleted items.
	// searchParams provides the search parameters.
	// Returns the list of identifiers of deleted items and an error if any.
	Deleted(searchParams *S) ([]string, error)
}
