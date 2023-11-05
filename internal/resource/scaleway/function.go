package scaleway

import (
	"context"
	"strings"

	"github.com/cyclimse/scwtui/internal/resource"
	sdk "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Function struct {
	sdk.Function `json:"function"`
	Namespace    sdk.Namespace `json:"namespace"`
}

func (f Function) Metadata() resource.Metadata {
	return resource.Metadata{
		ID:          f.Function.ID,
		Name:        f.Function.Name,
		ProjectID:   f.Namespace.ProjectID,
		Status:      statusPtr(f.Function.Status),
		Description: f.Function.Description,
		Tags:        nil,
		Type:        resource.TypeFunction,
		Locality:    resource.Region(f.Function.Region),
	}
}

func (f Function) CockpitMetadata() resource.CockpitMetadata {
	s := strings.TrimPrefix(f.DomainName, "https://")
	resourceName := strings.Split(s, ".")[0]
	return resource.CockpitMetadata{
		CanViewLogs:  true,
		ResourceName: resourceName,
		ResourceType: "serverless_function",
	}
}

func (f Function) Delete(ctx context.Context, s resource.Storer, client *scw.Client) error {
	api := sdk.NewAPI(client)
	_, err := api.DeleteFunction(&sdk.DeleteFunctionRequest{
		FunctionID: f.ID,
		Region:     f.Region,
	})
	if err != nil {
		return err
	}

	return s.DeleteResource(ctx, f)
}
