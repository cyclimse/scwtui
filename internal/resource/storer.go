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

	// ListResourcesByIDs returns all resources with the given IDs.
	// Used in combination of the search API to get the resources that match a query.I
	ListResourcesByIDs(ctx context.Context, ids []string) ([]Resource, error)

	// FindTypedByPredicateInProject returns all resources of a given type in a project that match a predicate.
	FindTypedByPredicateInProject(ctx context.Context, resourceType Type, projectID string, predicate Predicate) ([]Resource, error)

	// DeleteResource deletes a resource.
	DeleteResource(ctx context.Context, r Resource) error
}
