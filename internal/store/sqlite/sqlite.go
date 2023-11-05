package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/cyclimse/scaleway-dangling/internal/store/sqlite/db"
)

//go:embed schema.sql
var ddl string

// NewStore creates a new store.
// The store is used to save the resources in a database, to avoid querying the API every time.
// It also allows doing some analysis on the resources, like finding dangling resources.
func NewStore(ctx context.Context, shard string) (*Store, error) {
	sqlDB, err := sql.Open("sqlite3", fmt.Sprintf("file:%s:?mode=memory&cache=shared", shard))
	if err != nil {
		return nil, fmt.Errorf("store: failed to open database: %w", err)
	}

	_, err = sqlDB.ExecContext(ctx, ddl)
	if err != nil {
		return nil, fmt.Errorf("store: failed to create schema: %w", err)
	}

	queries := db.New(sqlDB)

	return &Store{
		DB:           sqlDB,
		queries:      queries,
		unmarshaller: ScalewayResourceUnmarshal{},
	}, nil
}

type Store struct {
	DB           *sql.DB
	queries      *db.Queries
	unmarshaller ResourceUnmarshaler
}

// Close closes the store.
func (s *Store) Close() error {
	return s.DB.Close()
}

// SetUnmarshaller sets the unmarshaller used to unmarshal the resources.
// In demo mode, the unmarshaller is set to a demo unmarshaller that returns a demo resource.
func (s *Store) SetUnmarshaller(u ResourceUnmarshaler) {
	s.unmarshaller = u
}
