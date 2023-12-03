package resource

import "context"

type Searcher interface {
	// Search searches for resources. Returns a list of resource IDs that match the query.
	Search(ctx context.Context, query string) ([]string, error)

	// Index indexes a resource.
	Index(r Resource) error

	// Deindex deindexes a resource.
	Deindex(r Resource) error
}
