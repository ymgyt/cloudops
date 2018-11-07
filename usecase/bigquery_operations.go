package usecase

import "github.com/ymgyt/cloudops/core"

// NewBigqueryOps -
func NewBigqueryOps(ctx *core.Context, service core.BigqueryService) (*BigqueryOps, error) {
	return &BigqueryOps{ctx: ctx, service: service}, nil
}

// BigqueryOps -
type BigqueryOps struct {
	ctx     *core.Context
	service core.BigqueryService
}

// BQueryInput -
type BQueryInput struct {
	QueryConfig *core.BQueryInput
}

// BQueryOutput -
type BQueryOutput struct {
	Status *core.BQueryOutput
}

// Query -
func (bq *BigqueryOps) Query(in *BQueryInput) (*BQueryOutput, error) {
	out, err := bq.service.Query(in.QueryConfig)

	return &BQueryOutput{Status: out}, err
}
