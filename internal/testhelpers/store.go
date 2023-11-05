package testhelpers

import (
	"context"
	"testing"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/store/sqlite"
	"github.com/stretchr/testify/require"
)

func NewStoreFromResources(t *testing.T, resources []resource.Resource) *sqlite.Store {
	t.Helper()

	ctx := context.Background()
	store, err := sqlite.NewStore(ctx, t.TempDir())
	require.NoError(t, err)

	for _, r := range resources {
		err := store.Store(ctx, r)
		require.NoError(t, err)
	}

	return store
}
