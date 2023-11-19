package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type JobRun struct {
	sdk.JobRun `json:"job_run"`
	ProjectID  string `json:"project_id"`
}

func (run JobRun) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:        run.ID,
		Name:      run.ID,
		ProjectID: run.ProjectID,
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
