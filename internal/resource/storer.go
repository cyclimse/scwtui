package resource

import (
	"context"
	"errors"
)

// ErrResourceNotFound is returned when a resource is not found.
var ErrResourceNotFound = errors.New("store: resource not found")

type Predicate func(r Resource) bool

type Storer interface {
	// Store stores a resource.
	Store(ctx context.Context, r Resource) error

	// Queries

	// ListAllResources returns all resources.
	ListAllResources(ctx context.Context) ([]Resource, error)

	// DeleteResource deletes a resource.
	DeleteResource(ctx context.Context, r Resource) error
}
