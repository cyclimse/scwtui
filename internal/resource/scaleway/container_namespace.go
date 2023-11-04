package scaleway

import (
	"context"

	"github.com/cyclimse/scaleway-dangling/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type ContainerNamespace sdk.Namespace

func (ns ContainerNamespace) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          ns.ID,
		Name:        ns.Name,
		ProjectID:   ns.ProjectID,
		Status:      statusPtr(ns.Status),
		Description: ns.Description,
		Tags:        nil,
		Type:        resource.TypeContainerNamespace,
		Locality:    resource.Region(ns.Region),
	}
}

func (ns ContainerNamespace) Delete(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	_, err := api.DeleteNamespace(&sdk.DeleteNamespaceRequest{
		NamespaceID: ns.ID,
		Region:      ns.Region,
	})
	if err != nil {
		return err
	}

	return s.DeleteResource(ctx, ns)
}
