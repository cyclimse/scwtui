package scaleway

import (
	"context"

	"github.com/cyclimse/scaleway-dangling/internal/resource"
	"github.com/cyclimse/scaleway-dangling/internal/resource/scaleway"
	sdk "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (d *ResourceDiscover) discoverFunctionsInRegion(ctx context.Context, region scw.Region) ([]resource.Resource, error) {
	api := sdk.NewAPI(d.client)

	resources := make([]resource.Resource, 0)

	for _, project := range d.projects {
		projectID := project.Metadata().ID

		nss, err := api.ListNamespaces(&sdk.ListNamespacesRequest{
			Region:    region,
			ProjectID: &projectID,
		}, scw.WithAllPages(), scw.WithContext(ctx))
		if handleRequestError(err) != nil {
			return nil, err
		}

		for _, ns := range nss.Namespaces {
			if ns == nil {
				continue
			}

			resources = append(resources, scaleway.FunctionNamespace(*ns))

			fs, err := api.ListFunctions(&sdk.ListFunctionsRequest{
				Region:      region,
				NamespaceID: ns.ID,
			}, scw.WithAllPages(), scw.WithContext(ctx))
			if handleRequestError(err) != nil {
				return nil, err
			}

			for _, f := range fs.Functions {
				if f == nil {
					continue
				}

				resources = append(resources, scaleway.Function{
					Function:  *f,
					Namespace: *ns,
				})
			}
		}
	}

	return resources, nil
}
