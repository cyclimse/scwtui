package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	sdk "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (d *ResourceDiscover) discoverInstancesInZone(ctx context.Context, zone scw.Zone) ([]resource.Resource, error) {
	api := sdk.NewAPI(d.client)

	servers, err := api.ListServers(&sdk.ListServersRequest{
		Zone: zone,
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if handleRequestError(err) != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0, len(servers.Servers))

	for _, server := range servers.Servers {
		if server == nil {
			continue
		}

		resources = append(resources, scaleway.Instance(*server))
	}

	return resources, nil
}
