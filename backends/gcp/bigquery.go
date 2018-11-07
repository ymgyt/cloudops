package gcp

import (
	"go.uber.org/zap"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/ymgyt/cloudops/core"
)

// NewBigqueryService -
func NewBigqueryService(ctx *core.Context, projectID string, credJSON []byte) (*Bigquery, error) {
	client, err := bigquery.NewClient(ctx.Ctx, projectID, option.WithCredentialsJSON(credJSON))
	if err != nil {
		return nil, err
	}

	bq := &Bigquery{
		ctx:    ctx,
		client: client,
	}

	return bq, nil
}

// Bigquery -
type Bigquery struct {
	ctx    *core.Context
	client *bigquery.Client
}

// Query -
func (b *Bigquery) Query(in *core.BQueryInput) (*core.BQueryOutput, error) {
	b.ctx.Log.Debug("bq",
		zap.String("query", in.Query), zap.Bool("dryRun", in.Dryrun),
		zap.Int64("maxBytes", in.MaxBytesBilled),
		zap.String("destDataset", in.DestDatasetID), zap.String("destTable", in.DestTableID),
		zap.String("writeDisposition", string(in.WriteDisposition)), zap.String("createDisposition", string(in.CreateDisposition)))

	q := b.client.Query(in.Query)

	q.DryRun = in.Dryrun
	q.MaxBytesBilled = in.MaxBytesBilled
	if in.DestDatasetID != "" {
		// TODO empty case handling
		q.QueryConfig.Dst = b.client.Dataset(in.DestDatasetID).Table(in.DestTableID)
	}
	q.CreateDisposition = in.CreateDisposition
	q.WriteDisposition = in.WriteDisposition

	job, err := q.Run(b.ctx.Ctx)
	if err != nil {
		return nil, core.WrapError(core.Internal, "query.Run()", err)
	}

	status := job.LastStatus()
	if status == nil {
		return nil, core.WrapError(core.Internal, "failed to job.LastStatus()", err)
	}

	out := core.BQueryOutput{ProcessedBytes: status.Statistics.TotalBytesProcessed}
	if in.Dryrun || in.Handler == nil {
		return &out, nil
	}

	itr, err := job.Read(b.ctx.Ctx)
	if err != nil {
		return nil, core.WrapError(core.Internal, "job.Read()", err)
	}

	for {
		var vs []bigquery.Value
		err := itr.Next(&vs)
		in.Handler(vs)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, core.WrapError(core.Internal, "itr.Next()", err)
		}
	}

	return &out, nil
}
