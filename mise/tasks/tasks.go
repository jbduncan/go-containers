package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sync/errgroup"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(max(1, runtime.NumCPU()/2))

	switch os.Args[1] {
	case "depaware":
		doDepaware(ctx, group)
	case "eg":
		doEg(ctx, group)
	}

	os.Exit(must(group.Wait()))
}

func doDepaware(ctx context.Context, group *errgroup.Group) {
	depawareFileDirs := make(chan string)
	group.Go(func() error {
		defer close(depawareFileDirs)
		return filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}
			if filepath.Base(path) != "depaware.txt" {
				return nil
			}

			select {
			case depawareFileDirs <- "./" + filepath.Dir(path):
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	for dir := range depawareFileDirs {
		group.Go(func() error {
			fmt.Printf("Linting with depaware concurrently on directory %s...\n", dir)
			return cmd("depaware", "--check", dir).Run()
		})
	}
}

func doEg(ctx context.Context, group *errgroup.Group) {
	egTemplateFiles := make(chan string)
	group.Go(func() error {
		defer close(egTemplateFiles)
		matches, err := filepath.Glob("eg/*.template")
		if err != nil {
			return err
		}

		for _, match := range matches {
			select {
			case egTemplateFiles <- match:
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return nil
	})

	for file := range egTemplateFiles {
		group.Go(func() error {
			fmt.Printf("Linting with eg template %s concurrently...\n", file)
			c := exec.Command("eg", "-t", file, "./...")

			// On a match, eg prints the whole contents of the matching file
			// which is too noisy, so stop it from being printed.
			c.Stdout = io.Discard

			buf := new(strings.Builder)
			c.Stderr = buf
			err := c.Run()
			if err != nil {
				fmt.Printf("%s: %s: %v\n", file, strings.TrimRight(buf.String(), "\n"), err)
				return err
			}
			if buf.Len() > 0 {
				fmt.Printf("%s: %s\n", file, strings.TrimRight(buf.String(), "\n"))
				return fmt.Errorf("eg found a problem (see above)")
			}
			return nil
		})
	}
}

func cmd(name string, args ...string) *exec.Cmd {
	result := exec.Command(name, args...)
	result.Stdout = os.Stdout
	result.Stderr = os.Stderr
	return result
}

func must(err error) int {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		fmt.Printf("%v\n", err)
		return exitErr.ExitCode()
	}

	fmt.Printf("%v\n", err)
	return 1
}
