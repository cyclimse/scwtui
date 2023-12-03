package scaleway

import (
	"context"
	"strings"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Container struct {
	sdk.Container `json:"container"`
	Namespace     sdk.Namespace `json:"namespace"`
}

func (c Container) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          c.Container.ID,
		Name:        c.Container.Name,
		ProjectID:   c.Namespace.ProjectID,
		Status:      statusPtr(c.Container.Status),
		Description: c.Container.Description,
		Tags:        nil,
		Type:        resource.TypeContainer,
		Locality:    resource.Region(c.Container.Region),
	}
}

func (c Container) CockpitMetadata() resource.CockpitMetadata {
	s := strings.TrimPrefix(c.DomainName, "https://")
	resourceName := strings.Split(s, ".")[0]
	return resource.CockpitMetadata{
		CanViewLogs:  true,
		ResourceName: resourceName,
		ResourceType: "serverless_container",
	}
}

func (c Container) Delete(ctx context.Context, index resource.Indexer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	_, err := api.DeleteContainer(&sdk.DeleteContainerRequest{
		ContainerID: c.ID,
		Region:      c.Region,
	})
	if err != nil {
		return err
	}

	return index.Deindex(ctx, c)
}
