package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	sdk "github.com/scaleway/scaleway-sdk-go/api/account/v3"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func ListProjects(ctx context.Context, client *scw.Client) ([]resource.Resource, error) {
	api := sdk.NewProjectAPI(client)

	projects, err := api.ListProjects(&sdk.ProjectAPIListProjectsRequest{}, scw.WithAllPages(), scw.WithContext(ctx))
	if handleRequestError(err) != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0, len(projects.Projects))

	for _, project := range projects.Projects {
		if project == nil {
			continue
		}

		resources = append(resources, scaleway.Project(*project))
	}

	return resources, nil
}
