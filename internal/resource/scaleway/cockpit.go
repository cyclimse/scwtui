package scaleway

import (
	"context"

	"github.com/cyclimse/scaleway-dangling/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Cockpit sdk.Cockpit

func (c Cockpit) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          c.ProjectID,
		Name:        c.ProjectID,
		ProjectID:   c.ProjectID,
		Status:      statusPtr(c.Status),
		Description: nil,
		Tags:        nil,
		Type:        resource.TypeCockpit,
		Locality:    resource.Global,
	}
}

func (c Cockpit) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs: false,
	}
}

func (c Cockpit) Delete(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	_, err := api.DeactivateCockpit(&sdk.DeactivateCockpitRequest{
		ProjectID: c.ProjectID,
	})
	if err != nil {
		return err
	}

	_, err = api.WaitForCockpit(&sdk.WaitForCockpitRequest{
		ProjectID: c.ProjectID,
	})
	if err != nil {
		return err
	}

	return s.DeleteResource(ctx, c)
}
