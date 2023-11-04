package scaleway

import (
	"context"

	"github.com/cyclimse/scaleway-dangling/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Container struct {
	sdk.Container `json:"container"`
	Namespace     sdk.Namespace `json:"namespace"`
}

func (f Container) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          f.Container.ID,
		Name:        f.Container.Name,
		ProjectID:   f.Namespace.ProjectID,
		Status:      statusPtr(f.Container.Status),
		Description: f.Container.Description,
		Tags:        nil,
		Type:        resource.TypeContainer,
		Locality:    resource.Region(f.Container.Region),
	}
}

func (f Container) Delete(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	_, err := api.DeleteContainer(&sdk.DeleteContainerRequest{
		ContainerID: f.ID,
		Region:      f.Region,
	})
	if err != nil {
		return err
	}

	return s.DeleteResource(ctx, f)
}
