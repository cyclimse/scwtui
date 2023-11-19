package sqlite

import (
	"encoding/json"
	"fmt"

	"github.com/cyclimse/scwtui/internal/resource"
	"github.com/cyclimse/scwtui/internal/resource/demo"
	"github.com/cyclimse/scwtui/internal/resource/scaleway"
)

type ResourceUnmarshaler interface {
	// UnmarshalResource unmarshals a resource from a string.
	UnmarshalResource(resourceType resource.Type, resourceData string) (resource.Resource, error)
}

type ScalewayResourceUnmarshal struct{}

func fromString[T resource.Resource](s string) (resource.Resource, error) {
	var r T
	err := json.Unmarshal([]byte(s), &r)
	if err != nil {
		return nil, fmt.Errorf("store: failed to unmarshal resource: %w", err)
	}
	return r, nil
}

// UnmarshalResource implements ResourceUnmarshaler.
func (s ScalewayResourceUnmarshal) UnmarshalResource(resourceType resource.Type, resourceData string) (resource.Resource, error) {
	switch resourceType {
	case resource.TypeProject:
		return fromString[scaleway.Project](resourceData)
	case resource.TypeIAMApplication:
		return fromString[scaleway.IAMApplication](resourceData)
	case resource.TypeCockpit:
		return fromString[scaleway.Cockpit](resourceData)
	case resource.TypeFunctionNamespace:
		return fromString[scaleway.FunctionNamespace](resourceData)
	case resource.TypeFunction:
		return fromString[scaleway.Function](resourceData)
	case resource.TypeContainerNamespace:
		return fromString[scaleway.ContainerNamespace](resourceData)
	case resource.TypeContainer:
		return fromString[scaleway.Container](resourceData)
	case resource.TypeRegistryNamespace:
		return fromString[scaleway.RegistryNamespace](resourceData)
	case resource.TypeRdbInstance:
		return fromString[scaleway.RdbInstance](resourceData)
	case resource.TypeKapsuleCluster:
		return fromString[scaleway.KapsuleCluster](resourceData)
	case resource.TypeInstance:
		return fromString[scaleway.Instance](resourceData)
	case resource.TypeJobDefinition:
		return fromString[scaleway.JobDefinition](resourceData)
	case resource.TypeJobRun:
		return fromString[scaleway.JobRun](resourceData)
	default:
		return nil, fmt.Errorf("store: unknown resource type %s", resourceType)
	}
}

type DemoResourceUnmarshal struct{}

// UnmarshalResource implements ResourceUnmarshaler.
// The reason this has to be split for demo resources is that the demo resources can take on any type.
// Eg: a demo resource can be a scaleway.Project, a scaleway.Container, a scaleway.Function, etc.
func (d DemoResourceUnmarshal) UnmarshalResource(_ resource.Type, resourceData string) (resource.Resource, error) {
	return fromString[demo.Resource](resourceData)
}
