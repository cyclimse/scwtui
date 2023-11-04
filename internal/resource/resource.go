package resource

import (
	"context"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Metadata struct {
	// Name is the name of the resource.
	Name string `json:"name"`

	// ID is the unique identifier of the resource.
	ID string `json:"id"`

	// ProjectID is the unique identifier of the project the resource belongs to.
	ProjectID string `json:"project_id"`

	// Status is the status of the resource.
	// Maybe nil if not available.
	Status *Status `json:"status"`

	// Description is the description of the resource.
	// May be empty.
	Description *string `json:"description"`

	// Tags is the list of tags associated with the resource.
	Tags []string `json:"tags"`

	// Type is the type of the resource.
	Type Type `json:"type"`

	// Locality is the locality of the resource.
	Locality Locality `json:"locality"`
}

type Resource interface {
	// Metadata returns the metadatas of the resource.
	Metadata() Metadata

	// Delete deletes the resource.
	Delete(ctx context.Context, s Storer, client *scw.Client) error
}

type MonitoredResource interface {
	Resource

	// GetLogs returns the logs of the resource.
	Logs(ctx context.Context, s Storer, client *scw.Client) error
}
