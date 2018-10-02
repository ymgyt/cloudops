package usecase

import (
	"github.com/ymgyt/cloudops/core"
)

// NewGCPProjectsOps -
func NewGCPProjectsOps(ctx *core.Context, service core.GCPProjectService) (GCPProjectOps, error) {
	return &gcpProjectsOps{ctx: ctx, service: service}, nil
}

// GCPProjectOps -
type GCPProjectOps interface {
	List(*ListGCPProjectsInput) (*ListGCPProjectsOutput, error)
}

// ListGCPProjectsInput -
type ListGCPProjectsInput struct{}

// ListGCPProjectsOutput -
type ListGCPProjectsOutput struct {
	Projects []*core.GCPProject
}

type gcpProjectsOps struct {
	ctx     *core.Context
	service core.GCPProjectService
}

// List -
func (g *gcpProjectsOps) List(in *ListGCPProjectsInput) (*ListGCPProjectsOutput, error) {
	out, err := g.service.List(&core.ListGCPProjectsInput{})
	if err != nil {
		return nil, err
	}
	return &ListGCPProjectsOutput{Projects: out.Projects}, nil
}
