package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
	sdk "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (d *ResourceDiscover) discoverJobsInRegion(ctx context.Context, region scw.Region) ([]resource.Resource, error) {
	api := sdk.NewAPI(d.client)

	jobs, err := api.ListJobDefinitions(&sdk.ListJobDefinitionsRequest{
		Region: region,
	}, scw.WithAllPages(), scw.WithContext(ctx))
	if handleRequestError(err) != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0, len(jobs.JobDefinitions))

	for _, jobDef := range jobs.JobDefinitions {
		if jobDef == nil {
			continue
		}

		resources = append(resources, scaleway.JobDefinition(*jobDef))

		jobRuns, err := api.ListJobRuns(&sdk.ListJobRunsRequest{
			Region: region,
			ID:     &jobDef.ID,
		}, scw.WithAllPages(), scw.WithContext(ctx))
		if handleRequestError(err) != nil {
			return nil, err
		}

		for _, jobRun := range jobRuns.JobRuns {
			if jobRun == nil {
				continue
			}

			resources = append(resources, scaleway.JobRun{
				JobRun:        *jobRun,
				JobDefinition: *jobDef,
			})
		}
	}

	return resources, nil
}
