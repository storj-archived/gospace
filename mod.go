// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type IsTidy struct {
	Common
}

func (cmd *IsTidy) Name() string { return "istidy" }

func (cmd *IsTidy) Parse(args []string) error {
	set := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	return set.Parse(args)
}

func (cmd *IsTidy) Exec() {
	repodir := cmd.RepoDir()

	before, err := HashFiles(filepath.Join(repodir, "go.mod"))
	ErrFatal(err)

	cmd.Tidy()

	after, err := HashFiles(filepath.Join(repodir, "go.mod"))
	ErrFatal(err)

	if before != after {
		fmt.Fprintln(os.Stderr, "go.mod is not tidy")
		os.Exit(1)
	}
}
