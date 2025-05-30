package main

import (
	"bytes"
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
	doMain := func() error {
		switch os.Args[1] {
		case "depaware":
			return doDepaware()
		case "depaware-fix":
			return doDepawareFix()
		case "eg":
			return doEg()
		case "eg-fix":
			return doEgFix()
		case "go-fix-diff":
			return goFixDiff()
		default:
			return fmt.Errorf("invalid command: %s", os.Args[1])
		}
	}

	os.Exit(toExitCode(doMain()))
}

func doDepaware() error {
	depawareFileDirs := make(chan string)
	group, ctx := newErrorGroup()
	group.Go(func() error {
		defer close(depawareFileDirs)
		return filepath.WalkDir(
			".",
			func(path string, d os.DirEntry, err error) error {
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
				case <-ctx.Done():
					return ctx.Err()
				case depawareFileDirs <- "./" + filepath.Dir(path):
				}
				return nil
			},
		)
	})

	for dir := range depawareFileDirs {
		group.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				fmt.Printf("Linting with depaware on directory %s...\n", dir)
				return cmd("depaware", "--check", dir).Run()
			}
		})
	}

	return group.Wait()
}

func doDepawareFix() error {
	depawareFiles := make(chan string)
	group, ctx := newErrorGroup()
	group.Go(func() error {
		defer close(depawareFiles)
		return filepath.WalkDir(
			".",
			func(path string, d os.DirEntry, err error) error {
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
				case <-ctx.Done():
					return ctx.Err()
				case depawareFiles <- "./" + path:
				}
				return nil
			},
		)
	})

	for file := range depawareFiles {
		group.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return doDepawareFixFor(file)
			}
		})
	}

	return group.Wait()
}

func doDepawareFixFor(depawareFile string) error {
	dir := "./" + filepath.Dir(depawareFile)
	fmt.Printf("Fixing with depaware on directory %s...\n", dir)
	c := cmd("depaware", dir)
	var b bytes.Buffer
	c.Stdout = &b
	if err := c.Run(); err != nil {
		return err
	}
	return os.WriteFile(depawareFile, b.Bytes(), 0o600)
}

func doEg() error {
	egTemplateFiles := make(chan string)
	group, ctx := newErrorGroup()
	group.Go(func() error {
		defer close(egTemplateFiles)
		matches, err := filepath.Glob("eg/*.template")
		if err != nil {
			return err
		}

		for _, match := range matches {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case egTemplateFiles <- match:
				continue
			}
		}

		return nil
	})

	for file := range egTemplateFiles {
		group.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return doEgFor(file)
			}
		})
	}

	return group.Wait()
}

func doEgFor(egTemplateFile string) error {
	fmt.Printf("Linting with eg template %s...\n", egTemplateFile)
	c := exec.Command("eg", "-t", egTemplateFile, "./...")

	// On a match, eg prints the whole contents of the matching file
	// which is too noisy, so stop it from being printed.
	c.Stdout = io.Discard

	buf := new(strings.Builder)
	c.Stderr = buf
	if err := c.Run(); err != nil {
		fmt.Printf(
			"%s: %s: %v\n",
			egTemplateFile,
			strings.TrimRight(buf.String(), "\n"),
			err,
		)
		return err
	}
	if buf.Len() > 0 {
		fmt.Printf(
			"%s: %s\n",
			egTemplateFile,
			strings.TrimRight(buf.String(), "\n"),
		)
		return fmt.Errorf("eg found a problem (see above)")
	}
	return nil
}

func doEgFix() error {
	egTemplateFiles, err := filepath.Glob("eg/*.template")
	if err != nil {
		return err
	}
	for _, file := range egTemplateFiles {
		fmt.Printf("Fixing with eg template %s...\n", file)
		if err := cmd("eg", "-t", file, "-w", "./...").Run(); err != nil {
			return err
		}
	}
	return nil
}

func goFixDiff() error {
	goFileDirs := make(map[string]struct{})
	if err := filepath.WalkDir(
		".",
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}

			if filepath.Ext(d.Name()) != ".go" {
				return nil
			}

			goFileDirs["./"+filepath.Dir(path)] = struct{}{}
			return nil
		},
	); err != nil {
		return err
	}

	goVersion, err := minGoVersionForProject()
	if err != nil {
		return fmt.Errorf("min project Go version not found: %w", err)
	}
	fmt.Printf("Min project Go version of %q found\n", goVersion)

	group, ctx := newErrorGroup()
	for dir := range goFileDirs {
		group.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return doGoFixFor(dir, goVersion)
			}
		})
	}

	return group.Wait()
}

func minGoVersionForProject() (string, error) {
	s := new(strings.Builder)
	c := cmd("go", "list", "-f", "{{.GoVersion}}", "-m")
	c.Stdout = s
	if err := c.Run(); err != nil {
		return "", err
	}
	goVersion := strings.TrimRight(s.String(), "\n")
	return goVersion, nil
}

func doGoFixFor(dir string, goVersion string) error {
	fmt.Printf(
		"Linting with 'go tool fix -diff' on directory %s...\n",
		dir,
	)
	c := cmd(
		"go",
		"tool",
		"fix",
		"-diff",
		fmt.Sprintf("-go=go%s", goVersion),
		dir,
	)
	buf := new(strings.Builder)
	c.Stderr = buf
	if err := c.Run(); err != nil {
		return err
	}
	if buf.Len() > 0 {
		return fmt.Errorf(
			"'go tool fix -diff' found a problem (see above)",
		)
	}
	return nil
}

func newErrorGroup() (*errgroup.Group, context.Context) {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(max(1, runtime.NumCPU()-1))
	return group, ctx
}

func cmd(name string, args ...string) *exec.Cmd {
	result := exec.Command(name, args...)
	result.Stdout = os.Stdout
	result.Stderr = os.Stderr
	return result
}

func toExitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		return exitErr.ExitCode()
	}

	fmt.Printf("%v\n", err)
	return 1
}
