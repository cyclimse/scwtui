package testhelpers

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type MockResource struct {
	MetadataValue        resource.Metadata        `json:"-"`
	CockpitMetadataValue resource.CockpitMetadata `json:"-"`
}

func (f *MockResource) Metadata() resource.Metadata {
	return f.MetadataValue
}

func (f *MockResource) CockpitMetadata() resource.CockpitMetadata {
	return f.CockpitMetadataValue
}

func (f *MockResource) Delete(_ context.Context, _ resource.Storer, _ *scw.Client) error {
	return nil
}
