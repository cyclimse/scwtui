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
		Name:      run.JobDefinition.Name,
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
func (run JobRun) Delete(_ context.Context, _ resource.Indexer, _ *scw.Client) error {
	return nil
}

func (run JobRun) Actions() []resource.Action {
	runAgain := resource.Action{
		Name: "Start new run",
		Do: func(ctx context.Context, index resource.Indexer, client *scw.Client) error {
			api := sdk.NewAPI(client)
			r, err := api.StartJobDefinition(&sdk.StartJobDefinitionRequest{
				JobDefinitionID: run.JobDefinition.ID,
				Region:          run.JobDefinition.Region,
			})
			if err != nil {
				return err
			}

			retriedRun := &JobRun{
				JobRun:        *r,
				JobDefinition: run.JobDefinition,
			}

			if err := index.Index(ctx, retriedRun); err != nil {
				return err
			}

			go func() {
				_ = retriedRun.pollUntilTerminated(ctx, index, client)
			}()

			return nil
		},
	}

	if run.State == sdk.JobRunStateFailed {
		runAgain.Name = "Retry"
	}

	actions := []resource.Action{runAgain}

	if run.State == sdk.JobRunStateQueued || run.State == sdk.JobRunStateRunning {
		actions = append(actions, resource.Action{
			Name: "Cancel",
			Do: func(ctx context.Context, _ resource.Indexer, client *scw.Client) error {
				api := sdk.NewAPI(client)
				_, err := api.StopJobRun(&sdk.StopJobRunRequest{
					JobRunID: run.ID,
					Region:   run.Region,
				})
				return err
			},
		})
	}

	return actions
}

func (run JobRun) pollUntilTerminated(ctx context.Context, index resource.Indexer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	for {
		r, err := api.GetJobRun(&sdk.GetJobRunRequest{
			JobRunID: run.ID,
			Region:   run.Region,
		})
		if err != nil {
			return err
		}

		run.JobRun = *r
		if err := index.Index(ctx, run); err != nil {
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
