package core

import "strings"

// Resource -
type Resource struct {
	Type ResourceType
}

// ResourceType -
type ResourceType int

//go:generate stringer -type ResourceType resource.go
const (
	InvalidResource ResourceType = iota
	RemoteResource
	LocalResource
)

// NewResource -
func NewResource(path string) (*Resource, error) {
	return &Resource{
		Type: inspectPath(path),
	}, nil
}

func inspectPath(path string) ResourceType {
	if path == "" {
		return InvalidResource
	}

	if strings.HasPrefix(path, "s3://") {
		return RemoteResource
	}

	c := path[0]
	if c == '/' || c == '.' {
		return LocalResource
	}

	return LocalResource
}
