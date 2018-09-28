// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Cache struct {
	Common
}

func (cmd *Cache) Name() string { return "cache" }

func (cmd *Cache) Parse(args []string) error {
	set := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	return set.Parse(args)
}

func (cmd *Cache) Exec() {
	fmt.Fprintf(os.Stderr, "# Caching\n")
	repodir := cmd.RepoDir()

	hash, err := HashFiles(
		filepath.Join(repodir, "go.mod"),
		filepath.Join(repodir, "go.sum"))
	ErrFatal(err)

	fmt.Println("HASH:", hash)

	cmd.VendorModules()
	zipdata, err := Zip(filepath.Join(repodir, "vendor"))
	ErrFatal(err)

	ioutil.WriteFile(filepath.Join(repodir, "vendor.zip"), zipdata, 0755)

	err = Unzip(zipdata, filepath.Join(repodir, "vendor2"))
	ErrFatal(err)
}

type Hash struct {
	Common
}

func (cmd *Hash) Name() string { return "hash" }

func (cmd *Hash) Parse(args []string) error {
	set := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	return set.Parse(args)
}

func (cmd *Hash) Exec() {
	repodir := cmd.RepoDir()

	hash, err := HashFiles(
		filepath.Join(repodir, "go.mod"),
		filepath.Join(repodir, "go.sum"))
	ErrFatal(err)

	fmt.Println(hash)
}

type ZipVendor struct {
	Common
	Destination string
}

func (cmd *ZipVendor) Name() string { return "zip-vendor" }

func (cmd *ZipVendor) Parse(args []string) error {
	set := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	err := set.Parse(args)
	if err != nil {
		return err
	}

	cmd.Destination = set.Arg(0)
	if cmd.Destination == "" {
		return errors.New("destination file unspecified")
	}

	return nil
}

func (cmd *ZipVendor) Exec() {
	fmt.Fprintf(os.Stderr, "# Zipping vendor\n")

	cmd.VendorModules()

	zipdata, err := Zip(filepath.Join(cmd.RepoDir(), "vendor"))
	ErrFatalf(err, "unable to zip vendor: %v", err)

	cmd.DeleteVendor()

	writeErr := ioutil.WriteFile(cmd.Destination, zipdata, 0644)
	ErrFatalf(writeErr, "unable to write zip: %v", writeErr)
}

type UnzipVendor struct {
	Common
	Source string
}

func (cmd *UnzipVendor) Name() string { return "unzip-vendor" }

func (cmd *UnzipVendor) Parse(args []string) error {
	set := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	err := set.Parse(args)
	if err != nil {
		return err
	}

	cmd.Source = set.Arg(0)
	if cmd.Source == "" {
		return errors.New("source file unspecified")
	}

	return nil
}

func (cmd *UnzipVendor) Exec() {
	fmt.Fprintf(os.Stderr, "# Unzipping vendor\n")
	zipdata, err := ioutil.ReadFile(cmd.Source)
	ErrFatalf(err, "unable to read zip: %v", err)

	cmd.DeleteVendor()

	unzipErr := Unzip(zipdata, filepath.Join(cmd.RepoDir(), "vendor"))
	ErrFatalf(unzipErr, "unable to unzip: %v", unzipErr)
}

type FlattenVendor struct {
	Common
	Source string
}

func (cmd *FlattenVendor) Name() string { return "flatten-vendor" }

func (cmd *FlattenVendor) Parse(args []string) error {
	set := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	return set.Parse(args)
}

func (cmd *FlattenVendor) Exec() {
	cmd.DeleteNonRepos()
	cmd.FlattenVendor()
}
