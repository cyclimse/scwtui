package scaleway_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/cyclimse/scwtui/internal/resource"
// 	"github.com/cyclimse/scwtui/internal/resource/scaleway"
// 	"github.com/cyclimse/scwtui/internal/testhelpers"
// 	sdk "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
// 	"github.com/scaleway/scaleway-sdk-go/scw"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// nolint:funlen // this is a test
// func TestFunctionNamespace_DanglingStatus(t *testing.T) {
// 	t.Parallel()
// 	ctx := context.Background()

// 	t.Run("no functions in namespace", func(t *testing.T) {
// 		t.Parallel()

// 		resources := []resource.Resource{
// 			&scaleway.FunctionNamespace{
// 				ID:        "ns1",
// 				Name:      "Namespace 1",
// 				ProjectID: "project1",
// 			},
// 		}

// 		store := testhelpers.NewStoreFromResources(t, resources)
// 		defer store.Close()

// 		ns := scaleway.FunctionNamespace{
// 			ID:        "ns1",
// 			Name:      "Namespace 1",
// 			ProjectID: "project1",
// 		}
// 		expectedStatus := resource.DanglingStatus{
// 			IsDangling:     true,
// 			DanglingReason: scaleway.DanglingReasonNoFunctionsInNamespace,
// 			SafeToDelete:   true,
// 		}
// 		status, err := ns.DanglingStatus(ctx, store)
// 		require.NoError(t, err)
// 		assert.Equal(t, expectedStatus, status)
// 	})

// 	t.Run("functions in namespace", func(t *testing.T) {
// 		t.Parallel()

// 		resources := []resource.Resource{
// 			&scaleway.FunctionNamespace{
// 				ID:        "ns1",
// 				Name:      "Namespace 1",
// 				ProjectID: "project1",
// 				Region:    scw.RegionFrPar,
// 			},
// 			&scaleway.Function{
// 				Function: sdk.Function{
// 					ID:          "f1",
// 					Name:        "Function 1",
// 					NamespaceID: "ns1",
// 					Region:      scw.RegionFrPar,
// 				},
// 				Namespace: sdk.Namespace{
// 					ID:        "ns1",
// 					Name:      "Namespace 1",
// 					ProjectID: "project1",
// 					Region:    scw.RegionFrPar,
// 				},
// 			},
// 		}

// 		store := testhelpers.NewStoreFromResources(t, resources)
// 		defer store.Close()

// 		ns := scaleway.FunctionNamespace{
// 			ID:        "ns1",
// 			Name:      "Namespace 1",
// 			ProjectID: "project1",
// 		}
// 		expectedStatus := resource.DanglingStatus{
// 			IsDangling:   false,
// 			SafeToDelete: false,
// 		}
// 		status, err := ns.DanglingStatus(ctx, store)
// 		require.NoError(t, err)
// 		assert.Equal(t, expectedStatus, status)
// 	})
// }
