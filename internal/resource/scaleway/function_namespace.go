package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FunctionNamespace sdk.Namespace

func (ns FunctionNamespace) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          ns.ID,
		Name:        ns.Name,
		ProjectID:   ns.ProjectID,
		Status:      statusPtr(ns.Status),
		Description: ns.Description,
		CreatedAt:   nil,
		Tags:        nil,
		Type:        resource.TypeFunctionNamespace,
		Locality:    resource.Region(ns.Region),
	}
}

func (ns FunctionNamespace) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs: false,
	}
}

func (ns FunctionNamespace) Delete(ctx context.Context, index resource.Indexer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	_, err := api.DeleteNamespace(&sdk.DeleteNamespaceRequest{
		NamespaceID: ns.ID,
		Region:      ns.Region,
	})
	if err != nil {
		return err
	}

	return index.Deindex(ctx, ns)
}
