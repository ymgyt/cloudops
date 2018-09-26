package main

import (
	"fmt"
	"os"

	"github.com/jawher/mow.cli"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/ymgyt/cloudops/core"
)

var (
	// injected build time
	version string
)

func main() {
	app := cli.App("cloudops", "utility tool for ops to make time to write more code")
	app.Version("version", version)

	app.Spec = "[-le]"

	var (
		loggingLevel  = app.StringOpt("l log logging", "info", "logging level(debug,info,warn,error)")
		loggingEncode = app.StringOpt("e enc encode", "color", "logging encode(json,console,color)")
		ctx           *core.Context
	)

	app.Before = func() {
		logger, err := core.NewLogger(*loggingLevel, *loggingEncode)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			cli.Exit(1)
		}
		validate := validator.New()
		ctx = core.NewContext(logger, validate)
	}

	app.Command("cp", "copy file(s) to/from remote datastorage", func(cmd *cli.Cmd) {

		cmd.Spec = "[--recursive --dryrun] SRC DST"

		var (
			recursive = cmd.BoolOpt("R recursive", false, "copy recursively")
			dryrun    = cmd.BoolOpt("D dryrun", false, "no create/update/delete operation")
			src       = cmd.StringArg("SRC", "", "source file to copy")
			dest      = cmd.StringArg("DST", "", "destination to copy")
		)

		cmd.Action = func() {
			copy := &CopyCommand{
				ctx:       ctx,
				dryrun:    *dryrun,
				recursive: *recursive,
				src:       *src,
				dest:      *dest,
			}
			copy.Run()
		}
	})

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	} else {
		os.Exit(0)
	}
}
