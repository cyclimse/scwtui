package scaleway

import (
	"context"

	"github.com/cyclimse/scaleway-dangling/internal/resource"
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
		Tags:        nil,
		Type:        resource.TypeRegistryNamespace,
		Locality:    resource.Region(ns.Region),
	}
}

func (ns RegistryNamespace) Delete(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := registry.NewAPI(client)
	_, err := api.DeleteNamespace(&registry.DeleteNamespaceRequest{
		NamespaceID: ns.ID,
		Region:      ns.Region,
	})
	if err != nil {
		return err
	}

	return s.DeleteResource(ctx, ns)
}
