package scaleway

import (
	"context"
	"time"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type JobRun struct {
	sdk.JobRun    `json:"job_run"`
	JobDefinition sdk.JobDefinition `json:"job_definition"`
}

func (run JobRun) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:        run.ID,
		Name:      run.ID,
		ProjectID: run.JobDefinition.ProjectID,
		Status:    statusPtr(run.State),
		Type:      resource.TypeJobRun,
		Locality:  resource.Region(run.Region),
	}
}

func (run JobRun) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs:  true,
		ResourceType: "serverless_job",
		ResourceID:   run.ID,
	}
}

// Delete is a no-op because job the lifecycle of a job run is managed by the job itself.
func (run JobRun) Delete(_ context.Context, _ resource.Storer, _ *scw.Client) error {
	return nil
}

func (run JobRun) Actions() []resource.Action {
	return []resource.Action{
		{
			Name: "Retry",
			Do: func(ctx context.Context, s resource.Storer, client *scw.Client) error {
				api := sdk.NewAPI(client)
				r, err := api.StartJobDefinition(&sdk.StartJobDefinitionRequest{
					ID:     run.JobDefinition.ID,
					Region: run.JobDefinition.Region,
				})
				if err != nil {
					return err
				}

				retriedRun := &JobRun{
					JobRun:        *r,
					JobDefinition: run.JobDefinition,
				}

				err = s.Store(ctx, retriedRun)
				if err != nil {
					return err
				}

				go func() {
					_ = retriedRun.pollUntilTerminated(ctx, s, client)
				}()

				return nil
			},
		},
	}
}

func (run JobRun) pollUntilTerminated(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	for {
		r, err := api.GetJobRun(&sdk.GetJobRunRequest{
			ID:     run.ID,
			Region: run.Region,
		})
		if err != nil {
			return err
		}

		run.JobRun = *r

		if err = s.Store(ctx, run); err != nil {
			return err
		}

		if run.State == sdk.JobRunStateSucceeded ||
			run.State == sdk.JobRunStateFailed ||
			run.State == sdk.JobRunStateCanceled {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}
