// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kylelemons/godebug/diff"
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

	before, err := ioutil.ReadFile(filepath.Join(repodir, "go.mod"))
	ErrFatal(err)

	cmd.Tidy()

	after, err := ioutil.ReadFile(filepath.Join(repodir, "go.mod"))
	ErrFatal(err)

	if string(before) != string(after) {
		diff, removed := difflines(string(before), string(after))
		fmt.Fprintln(os.Stderr, "go.mod is not tidy")
		fmt.Fprintln(os.Stderr, diff)
		if removed {
			os.Exit(1)
		}
	}
}

func difflines(a, b string) (patch string, removed bool) {
	alines, blines := strings.Split(a, "\n"), strings.Split(b, "\n")

	chunks := diff.DiffChunks(alines, blines)

	buf := new(bytes.Buffer)
	for _, c := range chunks {
		for _, line := range c.Added {
			fmt.Fprintf(buf, "+%s\n", line)
		}
		for _, line := range c.Deleted {
			fmt.Fprintf(buf, "-%s\n", line)
			removed = true
		}
	}

	return strings.TrimRight(buf.String(), "\n"), removed
}
