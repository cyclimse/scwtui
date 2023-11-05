package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/store/sqlite/db"
)

func (s *Store) GetResource(ctx context.Context, resourceID string) (resource.Resource, error) {
	r, err := s.queries.GetResource(ctx, resourceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, resource.ErrResourceNotFound
		}
		return nil, fmt.Errorf("store: failed to get resource with id %s: %w", resourceID, err)
	}

	return s.unmarshaller.UnmarshalResource(resource.Type(r.Type), r.Data.(string))
}

// ListAllResources implements resource.Storer.
func (s *Store) ListAllResources(ctx context.Context) ([]resource.Resource, error) {
	rows, err := s.queries.ListAllResources(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: failed to list all resources: %w", err)
	}

	resources := make([]resource.Resource, 0, len(rows))

	for _, row := range rows {
		r, err := s.unmarshaller.UnmarshalResource(resource.Type(row.Type), row.Data.(string))
		if err != nil {
			return nil, err
		}

		resources = append(resources, r)
	}

	return resources, nil
}

// FindTypedByPredicateInProject finds all the resources of a given type in a project that match a predicate.
func (s *Store) FindTypedByPredicateInProject(ctx context.Context, resourceType resource.Type, projectID string, predicate resource.Predicate) ([]resource.Resource, error) {
	rows, err := s.queries.ListTypedResourcesInProject(ctx, db.ListTypedResourcesInProjectParams{
		ProjectID: projectID,
		Type:      int64(resourceType),
	})
	if err != nil {
		return nil, fmt.Errorf("store: failed to list resources in project %s: %w", projectID, err)
	}

	resources := make([]resource.Resource, 0, len(rows))

	for _, row := range rows {
		r, err := s.unmarshaller.UnmarshalResource(resource.Type(row.Type), row.Data.(string))
		if err != nil {
			return nil, err
		}

		if predicate(r) {
			resources = append(resources, r)
		}
	}

	return resources, nil
}

// ListResourcesByIDs implements resource.Storer.
func (s *Store) ListResourcesByIDs(ctx context.Context, ids []string) ([]resource.Resource, error) {
	queryIds := make([]interface{}, len(ids))
	for i, id := range ids {
		queryIds[i] = id
	}

	rows, err := s.queries.ListResourcesByIDs(ctx, queryIds)
	if err != nil {
		return nil, fmt.Errorf("store: failed to list resources by ids: %w", err)
	}

	resources := make([]resource.Resource, 0, len(rows))

	for _, row := range rows {
		r, err := s.unmarshaller.UnmarshalResource(resource.Type(row.Type), row.Data.(string))
		if err != nil {
			return nil, err
		}

		resources = append(resources, r)
	}

	return resources, nil
}
