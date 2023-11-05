package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	sdk "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (d *ResourceDiscover) discoverContainersInRegion(ctx context.Context, region scw.Region) ([]resource.Resource, error) {
	api := sdk.NewAPI(d.client)

	resources := make([]resource.Resource, 0)

	for _, project := range d.projects {
		nss, err := api.ListNamespaces(&sdk.ListNamespacesRequest{
			Region:    region,
			ProjectID: &[]string{project.Metadata().ID}[0],
		}, scw.WithAllPages(), scw.WithContext(ctx))
		if handleRequestError(err) != nil {
			return nil, err
		}

		for _, ns := range nss.Namespaces {
			if ns == nil {
				continue
			}

			resources = append(resources, scaleway.ContainerNamespace(*ns))

			fs, err := api.ListContainers(&sdk.ListContainersRequest{
				Region:      region,
				NamespaceID: ns.ID,
			}, scw.WithAllPages(), scw.WithContext(ctx))
			if handleRequestError(err) != nil {
				return nil, err
			}

			for _, f := range fs.Containers {
				if f == nil {
					continue
				}

				resources = append(resources, scaleway.Container{
					Container: *f,
					Namespace: *ns,
				})
			}
		}
	}

	return resources, nil
}
