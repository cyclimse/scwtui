package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	sdk "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (d *ResourceDiscover) discoverCockpits(ctx context.Context) ([]resource.Resource, error) {
	api := sdk.NewAPI(d.client)

	cockpits := make([]resource.Resource, 0, len(d.projects))

	for _, project := range d.projects {
		cockpit, err := api.GetCockpit(&sdk.GetCockpitRequest{
			ProjectID: project.Metadata().ID,
		}, scw.WithContext(ctx))
		if handleRequestError(err) != nil {
			return nil, err
		}

		if cockpit == nil {
			continue
		}

		cockpits = append(cockpits, scaleway.Cockpit(*cockpit))
	}

	return cockpits, nil
}
