package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type RegistryNamespace registry.Namespace

func (ns RegistryNamespace) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          ns.ID,
		Name:        ns.Name,
		ProjectID:   ns.ProjectID,
		Status:      statusPtr(ns.Status),
		Description: &ns.Description,
		CreatedAt:   ns.CreatedAt,
		Tags:        nil,
		Type:        resource.TypeRegistryNamespace,
		Locality:    resource.Region(ns.Region),
	}
}

func (ns RegistryNamespace) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs: false,
	}
}

func (ns RegistryNamespace) Delete(ctx context.Context, index resource.Indexer, client *scw.Client) error {
	api := registry.NewAPI(client)
	_, err := api.DeleteNamespace(&registry.DeleteNamespaceRequest{
		NamespaceID: ns.ID,
		Region:      ns.Region,
	})
	if err != nil {
		return err
	}

	return index.Deindex(ctx, ns)
}
