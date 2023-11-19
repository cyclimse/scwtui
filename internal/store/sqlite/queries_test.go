package sqlite_test

import (
	"context"
	"testing"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	"github.com/cyclimse/scwtui/internal/store/sqlite"
	cockpit_sdk "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"
	fnc_sdk "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stringPtr(s string) *string {
	return &s
}

// nolint:funlen // this is a test
func TestStore(t *testing.T) {
	ctx := context.Background()
	store, err := sqlite.NewStore(ctx, t.TempDir())
	require.NoError(t, err)

	testCases := []struct {
		name     string
		resource resource.Resource
	}{
		{
			name: "scaleway function",
			resource: scaleway.Function{
				Function: fnc_sdk.Function{
					ID:          "function-id",
					Name:        "function-name",
					Description: stringPtr("function-description"),
					Region:      scw.RegionFrPar,
					Status:      fnc_sdk.FunctionStatusReady,
					Runtime:     fnc_sdk.FunctionRuntimeGo121,
					Privacy:     fnc_sdk.FunctionPrivacyPublic,
					HTTPOption:  fnc_sdk.FunctionHTTPOptionEnabled,
				},
				Namespace: fnc_sdk.Namespace{
					ID:        "namespace-id",
					Name:      "namespace-name",
					ProjectID: "project-id",
					Region:    scw.RegionFrPar,
					Status:    fnc_sdk.NamespaceStatusReady,
				},
			},
		},
		{
			name: "scaleway function namespace",
			resource: scaleway.FunctionNamespace{
				ID:        "namespace-id",
				Name:      "namespace-name",
				ProjectID: "project-id",
				Region:    scw.RegionFrPar,
				Status:    fnc_sdk.NamespaceStatusReady,
			},
		},
		{
			name: "scaleway cockpit",
			resource: scaleway.Cockpit{
				ProjectID: "project-id",
				Status:    cockpit_sdk.CockpitStatusCreating,
			},
		},
		{
			name: "scaleway project",
			resource: scaleway.Project{
				ID:   "project-id",
				Name: "project-name",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := store.Store(ctx, tc.resource)
			require.NoError(t, err)

			// retrieve the stored resource
			retrieved, err := store.GetResource(ctx, tc.resource.Metadata().ID)
			require.NoError(t, err)

			assert.Equal(t, tc.resource, retrieved)
		})
	}
}

// nolint:funlen // this is a test
func TestStore_ListResourcesByIDs(t *testing.T) {
	ctx := context.Background()
	store, err := sqlite.NewStore(ctx, t.TempDir())
	require.NoError(t, err)

	resources := []resource.Resource{
		&scaleway.Function{
			Function: fnc_sdk.Function{
				ID:          "function-id",
				Name:        "function-name",
				Description: stringPtr("function-description"),
				Region:      scw.RegionFrPar,
				Status:      fnc_sdk.FunctionStatusReady,
				Runtime:     fnc_sdk.FunctionRuntimeGo121,
				Privacy:     fnc_sdk.FunctionPrivacyPublic,
				HTTPOption:  fnc_sdk.FunctionHTTPOptionEnabled,
			},
			Namespace: fnc_sdk.Namespace{
				ID:        "namespace-id",
				Name:      "namespace-name",
				ProjectID: "project-id",
				Region:    scw.RegionFrPar,
				Status:    fnc_sdk.NamespaceStatusReady,
			},
		},
		&scaleway.FunctionNamespace{
			ID:        "namespace-id",
			Name:      "namespace-name",
			ProjectID: "project-id",
			Region:    scw.RegionFrPar,
			Status:    fnc_sdk.NamespaceStatusReady,
		},
		&scaleway.ContainerNamespace{
			ID:        "container-namespace-id",
			Name:      "container-namespace-name",
			ProjectID: "project-id",
			Region:    scw.RegionFrPar,
		},
	}

	for _, r := range resources {
		err := store.Store(ctx, r)
		require.NoError(t, err)
	}

	ids := []string{
		"function-id",
		"namespace-id",
		"container-namespace-id",
	}

	// retrieve the stored resources
	retrieved, err := store.ListResourcesByIDs(ctx, ids)
	require.NoError(t, err)
	assert.Len(t, retrieved, 3)

	retrievedIDs := make([]string, len(retrieved))
	for i, r := range retrieved {
		retrievedIDs[i] = r.Metadata().ID
	}

	assert.ElementsMatch(t, ids, retrievedIDs)
}
