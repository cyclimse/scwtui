package scaleway

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type KapsuleCluster sdk.Cluster

func (c KapsuleCluster) Metadata() resource.Metadata {
	return resource.Metadata{
		Name:        c.Name,
		ID:          c.ID,
		ProjectID:   c.ProjectID,
		Status:      statusPtr(c.Status),
		Description: &c.Description,
		Tags:        c.Tags,
		Type:        resource.TypeKapsuleCluster,
		Locality:    resource.Region(c.Region),
	}
}

func (c KapsuleCluster) CockpitMetadata() resource.CockpitMetadata {
	return resource.CockpitMetadata{
		CanViewLogs:  true,
		ResourceName: c.Name,
		ResourceType: "kubernetes_cluster",
	}
}

func (c KapsuleCluster) Delete(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	_, err := api.DeleteCluster(&sdk.DeleteClusterRequest{
		ClusterID: c.ID,
		Region:    c.Region,
	})
	if err != nil {
		return err
	}

	return s.DeleteResource(ctx, c)
}