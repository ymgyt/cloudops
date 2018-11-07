package main

import (
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/bigquery"

	cli "github.com/jawher/mow.cli"
	"go.uber.org/zap"

	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/usecase"
)

// BQueryCommand -
type BQueryCommand struct {
	ctx   *core.Context
	bqOps *usecase.BigqueryOps
	out   io.Writer

	query             string
	dryrun            bool
	maxBytesBilled    int64
	destDatasetID     string
	destTableID       string
	createDisposition string
	writeDisposition  string
}

// Run -
func (cmd *BQueryCommand) Run() {
	cmd.setDefault()

	out, err := cmd.bqOps.Query(&usecase.BQueryInput{
		QueryConfig: &core.BQueryInput{
			Query:             cmd.query,
			Dryrun:            cmd.dryrun,
			MaxBytesBilled:    cmd.maxBytesBilled,
			DestDatasetID:     cmd.destDatasetID,
			DestTableID:       cmd.destTableID,
			Handler:           cmd.printHandler,
			CreateDisposition: cmd.toBQCreateDisposition(),
			WriteDisposition:  cmd.toBQWriteDisposition(),
		},
	})
	if err != nil {
		cmd.ctx.Log.Error("bigquery query", zap.Error(err))
		cli.Exit(2)
	}

	// dryrunのときだけ有効..?
	if cmd.dryrun {
		b := out.Status.ProcessedBytes
		fmt.Fprintf(cmd.out, "ProcessedBytes: %s\n", formatBytes(b))
	}
}

func (cmd *BQueryCommand) printHandler(vs []bigquery.Value) {
	fmt.Println(vs)
}

func (cmd *BQueryCommand) toBQWriteDisposition() bigquery.TableWriteDisposition {
	var wd bigquery.TableWriteDisposition
	switch cmd.writeDisposition {
	case "append":
		wd = bigquery.WriteAppend
	case "empty":
		wd = bigquery.WriteEmpty
	case "truncate":
		wd = bigquery.WriteTruncate
	default:
		wd = bigquery.WriteEmpty
	}
	return wd
}

func (cmd *BQueryCommand) toBQCreateDisposition() bigquery.TableCreateDisposition {
	var cd bigquery.TableCreateDisposition
	switch cmd.createDisposition {
	case "ifneeded":
		cd = bigquery.CreateIfNeeded
	case "never":
		cd = bigquery.CreateNever
	default:
		cd = bigquery.CreateNever
	}
	return cd
}

func (cmd *BQueryCommand) setDefault() {
	if cmd.out == nil {
		cmd.out = os.Stdout
	}
}

func formatBytes(n int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}

	var u int64 = 1
	for i := 1; i <= len(units); i++ {
		u *= 1024
		if n/u == 0 || i == len(units) {
			return fmt.Sprintf("%d %s", n, units[i-1])
		}
		n = n / u
	}
	return ""
}
