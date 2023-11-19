package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type JobDefinition sdk.JobDefinition

func (def JobDefinition) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          def.ID,
		Name:        def.Name,
		ProjectID:   def.ProjectID,
		Description: &def.Description,
		Type:        resource.TypeJobDefinition,
		Locality:    resource.Region(def.Region),
	}
}

func (def JobDefinition) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs: false,
	}
}

func (def JobDefinition) Trigger(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	run, err := api.StartJobDefinition(&sdk.StartJobDefinitionRequest{
		ID:     def.ID,
		Region: def.Region,
	})
	if err != nil {
		return err
	}

	err = s.Store(ctx, JobRun{
		JobRun:    *run,
		ProjectID: def.ProjectID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (def JobDefinition) Delete(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	err := api.DeleteJobDefinition(&sdk.DeleteJobDefinitionRequest{
		ID:     def.ID,
		Region: def.Region,
	})
	if err != nil {
		return err
	}

	return s.DeleteResource(ctx, def)
}
