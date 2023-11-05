package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (d *ResourceDiscover) discoverRegistryNamespacesInRegion(ctx context.Context, region scw.Region) ([]resource.Resource, error) {
	api := registry.NewAPI(d.client)

	nss, err := api.ListNamespaces(&registry.ListNamespacesRequest{
		Region: region,
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if handleRequestError(err) != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0, len(nss.Namespaces))

	for _, ns := range nss.Namespaces {
		if ns == nil {
			continue
		}

		resources = append(resources, scaleway.RegistryNamespace(*ns))
	}

	return resources, nil
}
