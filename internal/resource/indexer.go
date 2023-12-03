package resource

import (
	"context"
	"fmt"
)

// Indexer is a very simple wrapper around a Storer and a Searcher.
// It makes it harder to forget to update one of the two when adding a new resource.
type Indexer interface {
	// Index indexes the given resource.
	Index(ctx context.Context, r Resource) error

	// Deindex deindexes the given resource.
	Deindex(ctx context.Context, r Resource) error
}

type Index struct {
	Store  Storer
	Search Searcher
}

func NewIndex(store Storer, search Searcher) *Index {
	return &Index{
		Store:  store,
		Search: search,
	}
}

func (i *Index) Index(ctx context.Context, r Resource) error {
	if err := i.Store.Store(ctx, r); err != nil {
		return fmt.Errorf("indexer: failed to store resource: %w", err)
	}
	if err := i.Search.Index(r); err != nil {
		return fmt.Errorf("indexer: failed to index resource: %w", err)
	}
	return nil
}

func (i *Index) Deindex(ctx context.Context, r Resource) error {
	if err := i.Store.DeleteResource(ctx, r); err != nil {
		return fmt.Errorf("indexer: failed to delete resource: %w", err)
	}
	if err := i.Search.Deindex(r); err != nil {
		return fmt.Errorf("indexer: failed to deindex resource: %w", err)
	}
	return nil
}
