package core

import (
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1beta1"
)

// GCPProjectService -
type GCPProjectService interface {
	List(*ListGCPProjectsInput) (*ListGCPProjectsOutput, error)
}

// ListGCPProjectsInput -
type ListGCPProjectsInput struct {
}

// ListGCPProjectsOutput -
type ListGCPProjectsOutput struct {
	Projects []*GCPProject
}

// GCPProject -
type GCPProject cloudresourcemanager.Project
