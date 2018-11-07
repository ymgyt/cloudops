package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jawher/mow.cli"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1beta1"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/ymgyt/cloudops/backends"
	"github.com/ymgyt/cloudops/backends/filesystem"
	"github.com/ymgyt/cloudops/backends/gcp"
	"github.com/ymgyt/cloudops/backends/s3"
	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/usecase"
)

var (
	// injected build time
	version string
)

func main() {
	app := cli.App("cloudops", "utility tool for ops to make time to write more code")
	app.Version("version", version)

	app.Spec = "[--log][--enc][--aws-region][--aws-access-key-id][--aws-secret-access-key][--aws-token][--google-application-credentials]"

	var (
		loggingLevel  = app.StringOpt("log logging", "info", "logging level(debug,info,warn,error)")
		loggingEncode = app.StringOpt("enc encode", "color", "logging encode(json,console,color)")
		awsRegion     = app.String(cli.StringOpt{
			Name: "aws-region", Value: "ap-northeast-1", Desc: "aws region", EnvVar: "AWS_REGION"})
		awsAccessKeyID = app.String(cli.StringOpt{
			Name: "aws-access-key-id", Value: "", Desc: "aws access key id", EnvVar: "AWS_ACCESS_KEY_ID"})
		awsSecretAccessKey = app.String(cli.StringOpt{
			Name: "aws-secret-access-key", Value: "", Desc: "aws secret access key", EnvVar: "AWS_SECRET_ACCESS_KEY"})
		awsToken = app.String(cli.StringOpt{
			Name: "aws-token", Value: "", Desc: "aws token", EnvVar: "AWS_TOKEN"})
		googleApplicationCredentials = app.String(cli.StringOpt{
			Name: "google-application-credentials", Value: "", Desc: "gcp service account credential file path", EnvVar: "GOOGLE_APPLICATION_CREDENTIALS"})
		ctx          *core.Context
		googleClient *http.Client

		fileSystem core.Backend
		s3Client   core.Backend
		confirmer  core.Confirmer
		gcpProject core.GCPProjectService

		fileOps       usecase.FileOps
		gcpProjectOps usecase.GCPProjectOps
		err           error
	)

	app.Before = func() {

		logger, err := core.NewLogger(*loggingLevel, *loggingEncode)
		if err != nil {
			fail(err)
		}
		validate := validator.New()
		ctx = core.NewContext(context.Background(), logger, validate)

		// backends
		if fileSystem, err = filesystem.New(ctx); err != nil {
			fail(err)
		}
		if s3Client, err = s3.New(ctx, *awsRegion, *awsAccessKeyID, *awsSecretAccessKey, *awsToken); err != nil {
			fail(err)
		}
		if confirmer, err = backends.NewPromptConfirmer(os.Stdout, os.Stdin, "[yes/no]", []string{"yes", "y", "Y"}); err != nil {
			fail(err)
		}
		if googleClient, err = newGoogleClient(ctx.Ctx, *googleApplicationCredentials); err != nil {
			fail(err)
		}
		if gcpProject, err = gcp.NewGCPProjectService(ctx, googleClient); err != nil {
			fail(err)
		}

		// usecase
		if fileOps, err = usecase.NewFileOps(ctx, fileSystem, s3Client, confirmer); err != nil {
			fail(err)
		}
		if gcpProjectOps, err = usecase.NewGCPProjectsOps(ctx, gcpProject); err != nil {
			fail(err)
		}

	}

	app.Command("cp", "copy file(s) to/from remote datastorage", func(cmd *cli.Cmd) {

		cmd.Spec = "[--recursive[--regexp]][--dryrun][--yes][--remove] SRC DST"

		var (
			recursive  = cmd.BoolOpt("R recursive", false, "copy recursively")
			dryrun     = cmd.BoolOpt("dryrun", false, "no create/update/delete operation")
			createDir  = cmd.BoolOpt("create-dir", false, "create directory if not exists")
			skipPrompt = cmd.BoolOpt("y yes", false, "skip prompt message")
			remove     = cmd.BoolOpt("remove", false, "remove after copy(like mv)")
			regexp     = cmd.StringOpt("r regexp", "", "target files go regexp pattern")
			src        = cmd.StringArg("SRC", "", "source file to copy")
			dest       = cmd.StringArg("DST", "", "destination to copy")
		)

		cmd.Action = func() {
			copy := &CopyCommand{
				ctx:        ctx,
				fileOps:    fileOps,
				dryrun:     *dryrun,
				recursive:  *recursive,
				createDir:  *createDir,
				skipPrompt: *skipPrompt,
				remove:     *remove,
				regexp:     *regexp,
				src:        *src,
				dest:       *dest,
			}
			copy.Run()
		}
	})

	app.Command("project", "manage gcp project resources", func(project *cli.Cmd) {
		project.Command("list", "list projects", func(list *cli.Cmd) {
			list.Action = func() {
				(&ProjectListCommand{ctx: ctx, projectsOps: gcpProjectOps}).Run()
			}
		})
	})

	app.Command("bigquery bq", "gcp bigquery operations", func(bq *cli.Cmd) {
		bq.Spec = "[--project-id]"

		var (
			projectID = bq.String(cli.StringOpt{Name: "project-id", Desc: "gcp projectID", EnvVar: "GCP_PROJECT_ID"})
			credJSON  = readFile(*googleApplicationCredentials)
			bqSrv     core.BigqueryService
			bqOps     *usecase.BigqueryOps
		)

		bq.Before = func() {
			if bqSrv, err = gcp.NewBigqueryService(ctx, *projectID, credJSON); err != nil {
				fail(err)
			}
			if bqOps, err = usecase.NewBigqueryOps(ctx, bqSrv); err != nil {
				fail(err)
			}
		}

		bq.Command("query", "run query sql", func(query *cli.Cmd) {
			query.Spec = "[--dry-run][--max-bytes][--dest-dataset --dest-table --write-disposition[--create-disposition]] QUERY"

			var (
				queryArg          = query.StringArg("QUERY", "", "query")
				dryrun            = query.BoolOpt("dry-run", false, "dry run")
				maxBytes          = query.IntOpt("max-bytes", 0, "limit to which query read table data")
				destDataset       = query.StringOpt("dest-dataset", "", "dataset to which query result saved")
				destTable         = query.StringOpt("dest-table", "", "table to which query result saved")
				createDisposition = query.StringOpt("create-disposition", "never", "specify the circumstances under which destination table will be created. ifneeded and never are supported")
				writeDisposition  = query.StringOpt("write-disposition", "empty", "specify how existing data in a destination table is treated. empty, append, and truncate are supported")
			)

			query.Action = func() {
				(&BQueryCommand{
					ctx:               ctx,
					bqOps:             bqOps,
					query:             *queryArg,
					dryrun:            *dryrun,
					maxBytesBilled:    int64(*maxBytes),
					destDatasetID:     *destDataset,
					destTableID:       *destTable,
					createDisposition: *createDisposition,
					writeDisposition:  *writeDisposition,
				}).Run()
			}
		})
	})

	errCh := make(chan error)
	go func() {
		errCh <- app.Run(os.Args)
	}()
	sigCh := watchSignal()
	for {
		select {
		case sig := <-sigCh:
			ctx.Log.Info("main", zap.String("signal", sig.String()))
			switch sig {
			case syscall.SIGINT:
				ctx.Cancel()
			}
		case err := <-errCh:
			if err == nil {
				os.Exit(0)
			}
			ctx.Log.Error("main", zap.Error(err))
			os.Exit(2)
		}
	}
}

func watchSignal() chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	return ch
}

func readFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fail(err)
	}
	return b
}

func newGoogleClient(ctx context.Context, path string) (*http.Client, error) {
	if path == "" {
		return nil, nil
	}
	data := readFile(path)
	crd, err := google.CredentialsFromJSON(ctx, data, cloudresourcemanager.CloudPlatformReadOnlyScope)
	if err != nil {
		return nil, err
	}
	tokenSource := crd.TokenSource
	return oauth2.NewClient(ctx, tokenSource), nil
}

func fail(err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r)
		}
	}()
	fmt.Fprintln(os.Stderr, err)
	cli.Exit(1)
}
