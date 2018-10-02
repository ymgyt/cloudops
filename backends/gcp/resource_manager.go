package gcp

import (
	"net/http"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1beta1"

	"github.com/ymgyt/cloudops/core"
)

// NewGCPProjectService -
func NewGCPProjectService(ctx *core.Context, client *http.Client) (core.GCPProjectService, error) {
	svc, err := cloudresourcemanager.New(client)
	if err != nil {
		return nil, core.WrapError(core.Internal, "", err)
	}
	return &gcpProjectService{ctx: ctx, service: svc}, nil
}

type gcpProjectService struct {
	ctx     *core.Context
	service *cloudresourcemanager.Service
}

// List -
func (s *gcpProjectService) List(in *core.ListGCPProjectsInput) (*core.ListGCPProjectsOutput, error) {
	call := s.service.Projects.List()
	if call == nil {
		return nil, core.NewError(core.Internal, "failed to call gcp project List API")
	}

	var projects []*core.GCPProject
	if err := call.Pages(s.ctx.Ctx, func(res *cloudresourcemanager.ListProjectsResponse) error {
		for _, pj := range res.Projects {
			projects = append(projects, (*core.GCPProject)(pj))
		}
		return nil
	}); err != nil {
		return nil, core.WrapError(core.Internal, "", err)
	}

	return &core.ListGCPProjectsOutput{
		Projects: projects,
	}, nil
}
