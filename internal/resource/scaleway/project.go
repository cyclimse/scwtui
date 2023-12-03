package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/account/v3"
	cockpit_sdk "github.com/scaleway/scaleway-sdk-go/api/cockpit/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Project sdk.Project

func (p Project) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          p.ID,
		Name:        p.Name,
		ProjectID:   p.ID,
		Status:      nil,
		Description: &p.Description,
		Tags:        nil,
		Type:        resource.TypeProject,
		Locality:    resource.Global,
	}
}

func (p Project) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs: false,
	}
}

func (p Project) Delete(ctx context.Context, index resource.Indexer, client *scw.Client) error {
	api := sdk.NewProjectAPI(client)
	err := api.DeleteProject(&sdk.ProjectAPIDeleteProjectRequest{
		ProjectID: p.ID,
	})
	if err != nil {
		return err
	}

	return index.Deindex(ctx, p)
}

func (p Project) Actions() []resource.Action {
	return []resource.Action{
		{
			Name: "Activate Cockpit",
			Do: func(ctx context.Context, index resource.Indexer, client *scw.Client) error {
				api := cockpit_sdk.NewAPI(client)
				_, err := api.ActivateCockpit(&cockpit_sdk.ActivateCockpitRequest{
					ProjectID: p.ID,
				})
				if err != nil {
					return err
				}

				cockpit, err := api.WaitForCockpit(&cockpit_sdk.WaitForCockpitRequest{
					ProjectID: p.ID,
				})
				if err != nil {
					return err
				}

				return index.Index(ctx, Cockpit(*cockpit))
			},
		},
	}
}
