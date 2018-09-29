package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jawher/mow.cli"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/ymgyt/cloudops/backends"
	"github.com/ymgyt/cloudops/backends/filesystem"
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

	app.Spec = "[--log][--enc][--aws-region][--aws-access-key-id][--aws-secret-access-key][--aws-token]"

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
		ctx *core.Context

		fileSystem core.Backend
		s3Client   core.Backend
		confirmer  core.Confirmer

		fileOps usecase.FileOps
	)

	app.Before = func() {
		fail := func(err error) {
			fmt.Fprintln(os.Stderr, err)
			cli.Exit(1)
		}

		logger, err := core.NewLogger(*loggingLevel, *loggingEncode)
		if err != nil {
			fail(err)
		}
		validate := validator.New()
		ctx = core.NewContext(context.Background(), logger, validate)

		// backends
		fileSystem, err = filesystem.New(ctx)
		if err != nil {
			fail(err)
		}
		s3Client, err = s3.New(ctx, *awsRegion, *awsAccessKeyID, *awsSecretAccessKey, *awsToken)
		if err != nil {
			fail(err)
		}
		confirmer, err = backends.NewPromptConfirmer(os.Stdout, os.Stdin, "[yes/no]", []string{"yes", "y", "Y"})
		if err != nil {
			fail(err)
		}

		// usecase
		fileOps, err = usecase.NewFileOps(ctx, fileSystem, s3Client, confirmer)
		if err != nil {
			fail(err)
		}

	}

	app.Command("cp", "copy file(s) to/from remote datastorage", func(cmd *cli.Cmd) {

		cmd.Spec = "[--recursive[--regexp]][--dryrun][--yes] SRC DST"

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
