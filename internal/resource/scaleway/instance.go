package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Instance sdk.Server

func (i Instance) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          i.ID,
		Name:        i.Name,
		ProjectID:   i.Project,
		Status:      statusPtr(i.State),
		Description: nil,
		Tags:        i.Tags,
		Type:        resource.TypeInstance,
		Locality:    resource.Zone(i.Zone),
	}
}

func (i Instance) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs: false,
	}
}

func (i Instance) Delete(ctx context.Context, index resource.Indexer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	err := api.DeleteServer(&sdk.DeleteServerRequest{
		ServerID: i.ID,
		Zone:     i.Zone,
	})
	if err != nil {
		return err
	}

	return index.Deindex(ctx, i)
}
