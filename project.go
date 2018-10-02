package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	cli "github.com/jawher/mow.cli"
	"go.uber.org/zap"

	"github.com/ymgyt/cloudops/core"
	"github.com/ymgyt/cloudops/usecase"
)

// ProjectListCommand -
type ProjectListCommand struct {
	ctx         *core.Context
	projectsOps usecase.GCPProjectOps
	out         io.Writer
}

// Run -
func (cmd *ProjectListCommand) Run() {
	cmd.printStart()
	if err := cmd.run(); err != nil {
		cmd.ctx.Log.Error("project list", zap.Error(err))
		cli.Exit(2)
	}
}

func (cmd *ProjectListCommand) run() error {
	cmd.setDefault()
	out, err := cmd.projectsOps.List(&usecase.ListGCPProjectsInput{})
	if err != nil {
		return err
	}
	return cmd.print(out)
}

func (cmd *ProjectListCommand) print(out *usecase.ListGCPProjectsOutput) error {
	w := new(tabwriter.Writer)
	w.Init(cmd.out, 0, 0, 3, ' ', 0)
	header := strings.Join([]string{"Name", "ProjectID", "ProjectNumber", "State"}, "\t")
	fmt.Fprintln(w, header)

	for _, pj := range out.Projects {
		fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\t%v", pj.Name, pj.ProjectId, pj.ProjectNumber, pj.LifecycleState))
	}
	return w.Flush()
}
func (cmd *ProjectListCommand) printStart() {
	cmd.ctx.Log.Info("project list")
}

func (cmd *ProjectListCommand) setDefault() {
	if cmd.out == nil {
		cmd.out = os.Stdout
	}
}
