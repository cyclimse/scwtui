package sqlite_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/store/sqlite"
	"github.com/cyclimse/scwtui/internal/testhelpers"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint:funlen // test multiple methods
func TestStore_Store(t *testing.T) {
	ctx := context.Background()

	store, err := sqlite.NewStore(ctx, t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, store.Close())
	})

	t.Run("store resource that does not exist", func(t *testing.T) {
		t.Parallel()

		r := &testhelpers.MockResource{
			MetadataValue: resource.Metadata{
				Name:        "name",
				ID:          "id",
				ProjectID:   "project-id",
				Description: &[]string{"description"}[0],
				Tags:        []string{"tag1", "tag2"},
				Type:        resource.TypeFunction,
				Locality:    resource.Region(scw.RegionFrPar),
			},
		}

		err = store.Store(ctx, r)
		require.NoError(t, err)

		// We will directly query the database to check if the resource was correctly stored.
		// This is not ideal but it's the only way to check if the resource was correctly stored.
		var (
			id          string
			name        string
			projectID   string
			description string
			tags        string
			typ         int
			locality    string
			data        string
		)

		err = store.DB.QueryRowContext(ctx, "SELECT * FROM resources WHERE id = ?", r.Metadata().ID).Scan(
			&id,
			&name,
			&projectID,
			&description,
			&tags,
			&typ,
			&locality,
			&data,
		)
		require.NoError(t, err)

		meta := r.Metadata()

		assert.Equal(t, meta.ID, id)
		assert.Equal(t, meta.Name, name)
		assert.Equal(t, meta.ProjectID, projectID)
		assert.Equal(t, *meta.Description, description)
		assert.Equal(t, meta.Tags, strings.Split(tags, ",")) // tags are stored as a comma separated string
		assert.Equal(t, meta.Type, resource.Type(typ))
		assert.Equal(t, meta.Locality, resource.Region(scw.RegionFrPar))
	})

	t.Run("store resource that already exists", func(t *testing.T) {
		t.Parallel()

		r := &testhelpers.MockResource{
			MetadataValue: resource.Metadata{
				ID:       "exists",
				Locality: resource.Region(scw.RegionFrPar),
			},
		}

		err = store.Store(ctx, r)
		require.NoError(t, err)

		// Update the resource
		r.MetadataValue = resource.Metadata{
			ID:       "exists",
			Locality: resource.Region(scw.RegionPlWaw),
		}

		err = store.Store(ctx, r)
		require.NoError(t, err)

		// We will directly query the database to check if the resource was correctly stored.
		// This is not ideal but it's the only way to check if the resource was correctly stored.
		var region string

		err = store.DB.QueryRowContext(ctx, "SELECT locality FROM resources WHERE id = ?", r.Metadata().ID).Scan(
			&region,
		)
		require.NoError(t, err)

		assert.Equal(t, "pl-waw", region)
	})
}

func TestStore_DeleteResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("delete resource that does not exist", func(t *testing.T) {
		t.Parallel()

		store, err := sqlite.NewStore(ctx, t.TempDir())
		require.NoError(t, err)
		defer store.Close()

		err = store.DeleteResource(ctx, &testhelpers.MockResource{
			MetadataValue: resource.Metadata{
				ID: "does-not-exist",
			},
		})
		require.Error(t, err)
		require.ErrorIs(t, err, resource.ErrResourceNotFound)
	})

	t.Run("delete resource that exists", func(t *testing.T) {
		t.Parallel()

		r := &testhelpers.MockResource{
			MetadataValue: resource.Metadata{
				ID:       "exists",
				Locality: resource.Region(scw.RegionFrPar),
			},
		}

		store := testhelpers.NewStoreFromResources(t, []resource.Resource{
			r,
		})
		defer store.Close()

		err := store.DeleteResource(ctx, r)
		require.NoError(t, err)

		_, err = store.GetResource(ctx, r.Metadata().ID)
		require.Error(t, err)
		require.ErrorIs(t, err, resource.ErrResourceNotFound)
	})
}
