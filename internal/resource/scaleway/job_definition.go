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
		CreatedAt:   def.CreatedAt,
		Type:        resource.TypeJobDefinition,
		Locality:    resource.Region(def.Region),
	}
}

func (def JobDefinition) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs: false,
	}
}

func (def JobDefinition) Delete(ctx context.Context, index resource.Indexer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	err := api.DeleteJobDefinition(&sdk.DeleteJobDefinitionRequest{
		JobDefinitionID: def.ID,
		Region:          def.Region,
	})
	if err != nil {
		return err
	}

	return index.Deindex(ctx, def)
}

func (def JobDefinition) Actions() []resource.Action {
	return []resource.Action{
		{
			Name: "Start",
			Do: func(ctx context.Context, index resource.Indexer, client *scw.Client) error {
				api := sdk.NewAPI(client)
				r, err := api.StartJobDefinition(&sdk.StartJobDefinitionRequest{
					JobDefinitionID: def.ID,
					Region:          def.Region,
				})
				if err != nil {
					return err
				}

				startedRun := &JobRun{
					JobRun:        *r,
					JobDefinition: sdk.JobDefinition(def),
				}

				if err := index.Index(ctx, startedRun); err != nil {
					return err
				}

				go func() {
					_ = startedRun.pollUntilTerminated(ctx, index, client)
				}()

				return nil
			},
		},
	}
}
