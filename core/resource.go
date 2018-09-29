package core

import (
	"io"
	"strings"
)

// Resource -
type Resource interface {
	Type() ResourceType
	URI() string
	Open() (io.ReadCloser, error)
}

// Resources -
type Resources []Resource

// ResourceType -
type ResourceType int

//go:generate stringer -type ResourceType resource.go
const (
	InvalidResource ResourceType = iota
	S3Resource
	LocalFileResource
)

// IsLocal -
func (rt ResourceType) IsLocal() bool {
	if rt == InvalidResource {
		return false
	}
	if rt == LocalFileResource {
		return true
	}
	return false
}

// IsRemote -
func (rt ResourceType) IsRemote() bool {
	if rt == InvalidResource {
		return false
	}
	return !rt.IsLocal()
}

// NOTE: not sure which layer should have resource type decision responsibility.
// InspectPath return corresponding resource type.
func InspectPath(path string) ResourceType {
	if path == "" {
		return InvalidResource
	}
	if strings.HasPrefix(path, "s3://") {
		return S3Resource
	}
	c := path[0]
	if c == '/' || c == '.' || c == '~' {
		return LocalFileResource
	}

	return LocalFileResource
}
