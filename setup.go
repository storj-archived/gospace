package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

type Setup struct {
	Common
	Overwrite bool
}

func (cmd *Setup) Name() string { return "setup" }

func (cmd *Setup) Parse(args []string) error {
	set := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	set.BoolVar(&cmd.Overwrite, "overwrite", false, "overwrite existing and reclone everything")
	return set.Parse(args)
}

func (cmd *Setup) Exec() {
	if !cmd.Overwrite && !Exists(cmd.Path("src")) {
		fmt.Fprintf(os.Stderr, "src directory %q already setup", cmd.Path("src"))
		os.Exit(1)
	}

	// _ = os.RemoveAll(cmd.Path("pkg"))
	_ = os.RemoveAll(cmd.Path("bin"))
	_ = os.RemoveAll(cmd.Path("src"))

	repodir := cmd.RepoDir()

	fmt.Fprintf(os.Stdout, "# Cloning repository\n")
	git := exec.Command("git", "clone", cmd.Repo, repodir)
	git.Stdout, git.Stderr = os.Stdout, os.Stderr
	ErrFatal(git.Run(), "git clone failed")

	cmd.VendorModules()
	cmd.FlattenVendor()
}
