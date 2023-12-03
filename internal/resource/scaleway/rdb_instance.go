package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type RdbInstance rdb.Instance

func (i RdbInstance) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          i.ID,
		Name:        i.Name,
		ProjectID:   i.ProjectID,
		Status:      statusPtr(i.Status),
		Description: nil,
		CreatedAt:   i.CreatedAt,
		Tags:        i.Tags,
		Type:        resource.TypeRdbInstance,
		Locality:    resource.Region(i.Region),
	}
}

func (i RdbInstance) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs:  true,
		ResourceID:   i.ID,
		ResourceType: "rdb_instance_postgresql",
	}
}

func (i RdbInstance) Delete(ctx context.Context, index resource.Indexer, client *scw.Client) error {
	api := rdb.NewAPI(client)
	_, err := api.DeleteInstance(&rdb.DeleteInstanceRequest{
		InstanceID: i.ID,
		Region:     i.Region,
	})
	if err != nil {
		return err
	}

	return index.Deindex(ctx, i)
}
