package repository

import "gorm.io/gorm"

// BaseRepository defines a generic repository interface for CRUD operations.
// T represents the type of the entity.
// K represents the type of the entity's ID.
// S represents the type of the search parameters.
type BaseRepository[T any, K comparable, S any] interface {
	// GetDB returns the gorm.DB instance.
	GetDB() *gorm.DB

	// Save saves a single entity.
	// filterContext provides additional context for the operation.
	// item is the entity to be saved.
	Save(filterContext *S, item *T) (*T, error)

	// SaveMany saves multiple entities.
	// filterContext provides additional context for the operation.
	// item is the list of entities to be saved.
	SaveMany(filterContext *S, item []*T) ([]*T, error)

	// Update updates a single entity.
	// filterContext provides additional context for the operation.
	// item is the entity to be updated.
	Update(filterContext *S, item *T) (*T, error)

	// UpdateMany updates multiple entities.
	// filterContext provides additional context for the operation.
	// item is the list of entities to be updated.
	UpdateMany(filterContext *S, item []*T) ([]*T, error)

	// Delete deletes a single entity by its ID.
	// searchParams provides additional context for the operation.
	// itemId is the ID of the entity to be deleted.
	Delete(searchParams *S, itemId K) error

	// DeleteByIds deletes multiple entities by their IDs.
	// searchParams provides additional context for the operation.
	// itemIds is the list of IDs of the entities to be deleted.
	DeleteByIds(searchParams *S, itemIds []K) error

	// FindById finds a single entity by its ID.
	// searchParams provides additional context for the operation.
	// itemId is the ID of the entity to be found.
	FindById(searchParams *S, itemId K) (*T, error)

	// FindByIds finds multiple entities by their IDs.
	// searchParams provides additional context for the operation.
	// itemIds is the list of IDs of the entities to be found.
	FindByIds(searchParams *S, itemIds []K) ([]*T, error)

	// FindAll finds all entities with pagination.
	// searchParams provides additional context for the operation.
	// pageSize is the number of entities per page.
	// page is the current page number.
	FindAll(searchParams *S, pageSize int, page int) ([]*T, error)

	// Search searches for entities based on search parameters.
	// searchParams provides the criteria for the search.
	Search(searchParams *S) ([]*T, error)

	// Count counts the number of entities based on search parameters.
	// searchParams provides the criteria for the count.
	Count(searchParams *S) (int64, error)

	// Deleted returns a list of IDs of deleted entities based on search parameters.
	// searchParams provides the criteria for the search.
	Deleted(searchParams *S) ([]string, error)
}
