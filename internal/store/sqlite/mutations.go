package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/store/sqlite/db"
)

// Store implements resource.Storer.
func (s *Store) Store(ctx context.Context, r resource.Resource) error {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("store: failed to marshal resource: %w", err)
	}

	m := r.Metadata()
	_, err = s.queries.UpsertResource(ctx, db.UpsertResourceParams{
		ID:          m.ID,
		Name:        m.Name,
		ProjectID:   m.ProjectID,
		Description: m.Description,
		Tags:        strings.Join(m.Tags, ","),
		Type:        int64(m.Type),
		Locality:    m.Locality,
		Data:        string(jsonData),
	})
	if err != nil {
		return fmt.Errorf("store: failed to store resource: %w", err)
	}

	return nil
}

// DeleteResource implements resource.Storer.
func (s *Store) DeleteResource(ctx context.Context, r resource.Resource) error {
	_, err := s.queries.DeleteResource(ctx, r.Metadata().ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resource.ErrResourceNotFound
		}
		return fmt.Errorf("store: failed to delete resource: %w", err)
	}

	return nil
}
