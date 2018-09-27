package main

import (
	"fmt"
	"os"
)

func Exists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil || !os.IsNotExist(err)
}

func Fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	panic("aborted")
}

func Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	panic("aborted")
}

func ErrFatal(err error, args ...interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, args...)
	panic("aborted")
}

func ErrFatalf(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format, args...)
	panic("aborted")
}

func Errf(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, format, args...)
}
