package resource

import "context"

type SetOfIDs = map[string]struct{}

type Searcher interface {
	// Search searches for resources. Returns a list of resource IDs that match the query.
	Search(ctx context.Context, query string) (SetOfIDs, error)

	// Index indexes a resource.
	Index(r Resource) error

	// Deindex deindexes a resource.
	Deindex(r Resource) error
}
