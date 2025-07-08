package repository

import "gorm.io/gorm"

// BaseRepository defines a generic repository interface for CRUD operations.
// T represents the type of the entity.
// K represents the type of the entity's ID.
// S represents the type of the search parameters.
// Many methods accept an optional *gorm.DB transaction (tx). If tx is nil, the default DB connection is used.
type BaseRepository[T any, K comparable, S any] interface {
	// GetDB returns the gorm.DB instance.
	GetDB() *gorm.DB

	// Save saves a single entity.
	// filterContext provides additional context for the operation.
	// item is the entity to be saved.
	// tx is an optional transaction. If nil, the default DB is used.
	Save(tx *gorm.DB, filterContext *S, item *T) (*T, error)

	// SaveMany saves multiple entities.
	// filterContext provides additional context for the operation.
	// item is the list of entities to be saved.
	// tx is an optional transaction. If nil, the default DB is used.
	SaveMany(tx *gorm.DB, filterContext *S, item []*T) ([]*T, error)

	// Update updates a single entity.
	// filterContext provides additional context for the operation.
	// item is the entity to be updated.
	// tx is an optional transaction. If nil, the default DB is used.
	Update(tx *gorm.DB, filterContext *S, item *T) (*T, error)

	// UpdateMany updates multiple entities.
	// filterContext provides additional context for the operation.
	// item is the list of entities to be updated.
	// tx is an optional transaction. If nil, the default DB is used.
	UpdateMany(tx *gorm.DB, filterContext *S, item []*T) ([]*T, error)

	// Delete deletes a single entity by its ID.
	// searchParams provides additional context for the operation.
	// itemId is the ID of the entity to be deleted.
	// tx is an optional transaction. If nil, the default DB is used.
	Delete(tx *gorm.DB, searchParams *S, itemId K) error

	// DeleteByIds deletes multiple entities by their IDs.
	// searchParams provides additional context for the operation.
	// itemIds is the list of IDs of the entities to be deleted.
	// tx is an optional transaction. If nil, the default DB is used.
	DeleteByIds(tx *gorm.DB, searchParams *S, itemIds []K) error

	// FindById finds a single entity by its ID.
	// searchParams provides additional context for the operation.
	// itemId is the ID of the entity to be found.
	// tx is an optional transaction. If nil, the default DB is used.
	FindById(tx *gorm.DB, searchParams *S, itemId K) (*T, error)

	// FindByIds finds multiple entities by their IDs.
	// searchParams provides additional context for the operation.
	// itemIds is the list of IDs of the entities to be found.
	// tx is an optional transaction. If nil, the default DB is used.
	FindByIds(tx *gorm.DB, searchParams *S, itemIds []K) ([]*T, error)

	// FindAll finds all entities with pagination.
	// searchParams provides additional context for the operation.
	// pageSize is the number of entities per page.
	// page is the current page number.
	// tx is an optional transaction. If nil, the default DB is used.
	FindAll(tx *gorm.DB, searchParams *S, pageSize int, page int) ([]*T, error)

	// Search searches for entities based on search parameters.
	// searchParams provides the criteria for the search.
	// tx is an optional transaction. If nil, the default DB is used.
	Search(tx *gorm.DB, searchParams *S) ([]*T, error)

	// Count counts the number of entities based on search parameters.
	// searchParams provides the criteria for the count.
	// tx is an optional transaction. If nil, the default DB is used.
	Count(tx *gorm.DB, searchParams *S) (int64, error)

	// Deleted returns a list of IDs of deleted entities based on search parameters.
	// searchParams provides the criteria for the search.
	// tx is an optional transaction. If nil, the default DB is used.
	Deleted(tx *gorm.DB, searchParams *S) ([]string, error)
}
