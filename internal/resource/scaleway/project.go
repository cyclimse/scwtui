package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/account/v3"
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

func (p Project) Delete(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := sdk.NewProjectAPI(client)
	err := api.DeleteProject(&sdk.ProjectAPIDeleteProjectRequest{
		ProjectID: p.ID,
	})
	if err != nil {
		return err
	}

	return s.DeleteResource(ctx, p)
}
