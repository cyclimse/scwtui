package bleve

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	sdk "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/stretchr/testify/require"
)

func TestSearch_Search(t *testing.T) {
	projectID := gofakeit.UUID()
	projectName := gofakeit.Name()

	search, err := NewSearch(map[string]string{
		projectID: projectName,
	})
	require.NoError(t, err)

	functionNamespaceID := gofakeit.UUID()
	functionID := gofakeit.UUID()

	resources := []resource.Resource{
		&scaleway.Project{
			ID:   projectID,
			Name: projectName,
		},
		&scaleway.FunctionNamespace{
			ID:        functionNamespaceID,
			ProjectID: projectID,
		},
		&scaleway.Function{
			Function: sdk.Function{
				ID:   functionID,
				Name: "function0",
			},
			Namespace: sdk.Namespace{
				ID:        functionNamespaceID,
				ProjectID: projectID,
			},
		},
	}

	for _, r := range resources {
		err := search.Index(r)
		require.NoError(t, err)
	}

	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"getting all resources in project via project id", projectID, []string{projectID, functionNamespaceID, functionID}},
		{"getting a function via function id", functionID, []string{functionID}},
		{"getting a function via function name", "function0", []string{functionID}},
		{"getting a function via function name (case insensitive)", "FUNCTION0", []string{functionID}},
		{"getting a function via its namespace id", functionNamespaceID, []string{functionID, functionNamespaceID}},
	}
	for _, test := range tests {
		test := test
		ctx := context.Background()
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := search.Search(ctx, test.query)
			require.NoError(t, err)
			require.ElementsMatch(t, test.want, got)
		})
	}
}
