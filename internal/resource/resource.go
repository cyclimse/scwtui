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
	Locality Locality `json:"locality,omitempty"`
}

type CockpitMetadata struct {
	// CanViewLogs is true if the logs associated with the resource can be viewed Scaleway Cockpit.
	CanViewLogs bool

	// ResourceName is the name of the resource in Scaleway Cockpit.
	ResourceName string

	// ResourceID is the ID of the resource in Scaleway Cockpit.
	ResourceID string

	// ResourceType is the type of the resource in Scaleway Cockpit.
	ResourceType string
}

type Resource interface {
	// Metadata returns the metadatas of the resource.
	Metadata() Metadata

	// CockpitMetadata returns the metadatas of the resource for Scaleway Cockpit.
	CockpitMetadata() CockpitMetadata

	// Delete deletes the resource.
	Delete(ctx context.Context, s Storer, client *scw.Client) error
}

type Action struct {
	// Name is the name of the action.
	Name string

	// Do performs the action on the resource.
	Do func(ctx context.Context, s Storer, client *scw.Client) error
}

type Actionable interface {
	Resource

	// Actions returns the list of actions that can be performed on the resource.
	Actions() []Action
}
