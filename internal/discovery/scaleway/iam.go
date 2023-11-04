package scaleway

import (
	"context"

	"github.com/cyclimse/scaleway-dangling/internal/resource"
	"github.com/cyclimse/scaleway-dangling/internal/resource/scaleway"
	"github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (d *ResourceDiscover) discoverIAMApplications(ctx context.Context) ([]resource.Resource, error) {
	organizationID, ok := d.client.GetDefaultOrganizationID()
	if !ok {
		return nil, nil
	}

	api := iam.NewAPI(d.client)

	apps, err := api.ListApplications(&iam.ListApplicationsRequest{
		OrganizationID: &organizationID,
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0, len(apps.Applications))

	for _, app := range apps.Applications {
		if app == nil {
			continue
		}

		resources = append(resources, scaleway.IAMApplication(*app))
	}

	return resources, nil
}
