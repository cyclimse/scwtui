package demo

import (
	"context"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Resource struct {
	MetadataValue        resource.Metadata        `json:"metadata"`
	CockpitMetadataValue resource.CockpitMetadata `json:"cockpit_metadata"`
	RawLocality          string                   `json:"locality"`
}

func (r Resource) Metadata() resource.Metadata {
	v := r.MetadataValue
	switch r.RawLocality {
	case "global":
		v.Locality = resource.Global
	case "fr-par", "nl-ams", "pl-waw":
		v.Locality = resource.Region(scw.Region(r.RawLocality))
	default:
		v.Locality = resource.Zone(scw.Zone(r.RawLocality))
	}
	return v
}

func (r Resource) CockpitMetadata() resource.CockpitMetadata {
	return r.CockpitMetadataValue
}

func (r Resource) Delete(_ context.Context, _ resource.Indexer, _ *scw.Client) error {
	return nil
}
