package resource

type Type int

//go:generate go run golang.org/x/tools/cmd/stringer -type=Type -linecomment -trimprefix=Type

const (
	TypeProject        Type = iota
	TypeIAMApplication      // IAM Application
	TypeCockpit
	TypeFunctionNamespace // Function Namespace
	TypeFunction
	TypeContainerNamespace // Container Namespace
	TypeContainer
	TypeRegistryNamespace // Registry Namespace
	TypeRdbInstance       // RDB Instance
	TypeKapsuleCluster    // Kapsule Cluster
	TypeInstance
	TypeJobDefinition // Job Definition
	TypeJobRun        // Job Run
	NumberOfResourceTypes
)
