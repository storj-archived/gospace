// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func HashFiles(files ...string) (string, error) {
	h := sha256.New()
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", err
		}

		h.Write(data)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func Zip(root string) ([]byte, error) {
	var b bytes.Buffer
	archive := zip.NewWriter(&b)

	archiveErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		// figure out path with forward slashes
		rpath, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		cpath := filepath.ToSlash(rpath)

		// create a zipped file
		zipFile, createErr := archive.Create(cpath)
		if createErr != nil {
			return createErr
		}

		// open source file
		dataFile, openErr := os.Open(path)
		if openErr != nil {
			return openErr
		}

		// copy data
		_, copyErr := io.Copy(zipFile, dataFile)
		if copyErr != nil && copyErr != io.EOF {
			return copyErr
		}

		// close source file
		if closeErr := dataFile.Close(); closeErr != nil {
			return closeErr
		}

		return nil
	})

	if archiveErr != nil {
		return nil, archiveErr
	}

	if err := archive.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func Unzip(data []byte, root string) error {
	b := bytes.NewReader(data)
	archive, err := zip.NewReader(b, b.Size())
	if err != nil {
		return err
	}

	for _, zipFile := range archive.File {
		if zipFile.FileInfo().IsDir() {
			continue
		}

		zipReader, openErr := zipFile.Open()
		if openErr != nil {
			return openErr
		}

		path := filepath.Join(root, filepath.FromSlash(zipFile.Name))
		mkdirErr := os.MkdirAll(filepath.Dir(path), 0700)
		if mkdirErr != nil {
			return mkdirErr
		}

		dataFile, createErr := os.Create(path)
		if createErr != nil {
			return createErr
		}

		_, copyErr := io.Copy(dataFile, zipReader)
		if copyErr != nil && copyErr != io.EOF {
			return copyErr
		}

		if closeErr := dataFile.Close(); closeErr != nil {
			return closeErr
		}
		if closeErr := zipReader.Close(); closeErr != nil {
			return closeErr
		}
	}

	return nil
}
