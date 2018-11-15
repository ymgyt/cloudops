package main

import (
	"fmt"
	"os"

	"github.com/ymgyt/cloudops/backends/filesystem"

	"github.com/ymgyt/cloudops/core"
)

// DiskUsageCommand -
type DiskUsageCommand struct {
	ctx *core.Context
	fs  *filesystem.FileSystem

	root string
}

// Run -
func (cmd *DiskUsageCommand) Run() {
	out, err := cmd.fs.DiskUsage(&filesystem.DiskUsageInput{
		Root: cmd.root,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// fmt.Printf("%s: %d(%s)\n", cmd.root, out.SizeSum, formatBytes(out.SizeSum))
	out.Root.Dump(os.Stdout)
}
