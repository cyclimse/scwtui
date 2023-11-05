package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	sdk "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (d *ResourceDiscover) discoverRdbInstancesInRegion(ctx context.Context, region scw.Region) ([]resource.Resource, error) {
	api := sdk.NewAPI(d.client)

	resources := make([]resource.Resource, 0)

	for _, project := range d.projects {
		projectID := project.Metadata().ID

		instances, err := api.ListInstances(&sdk.ListInstancesRequest{
			Region:    region,
			ProjectID: &projectID,
		}, scw.WithAllPages(), scw.WithContext(ctx))
		if handleRequestError(err) != nil {
			return nil, err
		}

		for _, i := range instances.Instances {
			if i == nil {
				continue
			}

			resources = append(resources, scaleway.RdbInstance(*i))
		}
	}

	return resources, nil
}
