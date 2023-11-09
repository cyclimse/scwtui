package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	sdk "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (d *ResourceDiscover) discoverKapsuleClustersInRegion(ctx context.Context, region scw.Region) ([]resource.Resource, error) {
	api := sdk.NewAPI(d.client)

	clusters, err := api.ListClusters(&sdk.ListClustersRequest{
		Region: region,
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if handleRequestError(err) != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0, len(clusters.Clusters))

	for _, cluster := range clusters.Clusters {
		if cluster == nil {
			continue
		}

		resources = append(resources, scaleway.KapsuleCluster(*cluster))
	}

	return resources, nil
}
