package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type IAMApplication iam.Application

func (app IAMApplication) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          app.ID,
		Name:        app.Name,
		ProjectID:   "",
		Status:      nil,
		Description: &app.Description,
		Tags:        nil,
		Type:        resource.TypeIAMApplication,
		Locality:    resource.Global,
	}
}

func (app IAMApplication) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs: false,
	}
}

func (app IAMApplication) Delete(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := iam.NewAPI(client)
	err := api.DeleteApplication(&iam.DeleteApplicationRequest{
		ApplicationID: app.ID,
	})
	if err != nil {
		return err
	}

	return s.DeleteResource(ctx, app)
}
