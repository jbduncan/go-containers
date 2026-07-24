package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	doMain := func() error {
		ctx, stop := signal.NotifyContext(context.Background(), interruptSignal())
		defer stop()

		switch os.Args[1] {
		case "depaware":
			return doDepaware(ctx)
		case "depaware-fix":
			return doDepawareFix(ctx)
		case "eg":
			return doEg(ctx)
		case "eg-fix":
			return doEgFix(ctx)
		default:
			return fmt.Errorf("invalid command: %s", os.Args[1])
		}
	}

	if err := doMain(); err != nil {
		exitCode := 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}

		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(exitCode)
	}
}

func interruptSignal() os.Signal {
	if runtime.GOOS == "windows" {
		return os.Kill // os.Interrupt is not implemented on Windows
	}
	return os.Interrupt
}

func doDepaware(baseCtx context.Context) error {
	depawareFileDirs := make(chan string)
	group, ctx := newErrorGroup(baseCtx)
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
				return cmd(ctx, "depaware", "--check", dir).Run()
			}
		})
	}

	return group.Wait()
}

func doDepawareFix(baseCtx context.Context) error {
	depawareFiles := make(chan string)
	group, ctx := newErrorGroup(baseCtx)
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
				return doDepawareFixFor(ctx, file)
			}
		})
	}

	return group.Wait()
}

func doDepawareFixFor(ctx context.Context, depawareFile string) error {
	dir := "./" + filepath.Dir(depawareFile)
	fmt.Printf("Fixing with depaware on directory %s...\n", dir)
	c := cmd(ctx, "depaware", dir)
	var b bytes.Buffer
	c.Stdout = &b
	if err := c.Run(); err != nil {
		return err
	}
	return os.WriteFile(depawareFile, b.Bytes(), 0o600)
}

func doEg(baseCtx context.Context) error {
	egTemplateFiles := make(chan string)
	group, ctx := newErrorGroup(baseCtx)
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
				return doEgFor(ctx, file)
			}
		})
	}

	return group.Wait()
}

func doEgFor(ctx context.Context, egTemplateFile string) error {
	fmt.Printf("Linting with eg template %s...\n", egTemplateFile)
	c := cmd(ctx, "eg", "-t", egTemplateFile, "./...")

	// On a match, eg prints the whole contents of the matching file
	// which is too noisy, so stop it from being printed.
	c.Stdout = io.Discard

	buf := new(strings.Builder)
	c.Stderr = buf
	if err := c.Run(); err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			egTemplateFile,
			strings.TrimRight(buf.String(), "\n"),
			err,
		)
	}
	if buf.Len() > 0 {
		return fmt.Errorf(
			"eg found a problem: %s: %s",
			egTemplateFile,
			strings.TrimRight(buf.String(), "\n"),
		)
	}
	return nil
}

func doEgFix(ctx context.Context) error {
	egTemplateFiles, err := filepath.Glob("eg/*.template")
	if err != nil {
		return err
	}
	for _, file := range egTemplateFiles {
		fmt.Printf("Fixing with eg template %s...\n", file)
		if err := cmd(ctx, "eg", "-t", file, "-w", "./...").Run(); err != nil {
			return err
		}
	}
	return nil
}

func newErrorGroup(ctx context.Context) (*errgroup.Group, context.Context) {
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(max(1, runtime.NumCPU()-1))
	return group, ctx
}

func cmd(ctx context.Context, name string, args ...string) *exec.Cmd {
	subprocess := exec.CommandContext(ctx, name, args...)
	subprocess.Stdout = os.Stdout
	subprocess.Stderr = os.Stderr
	subprocess.Cancel = func() error {
		return subprocess.Process.Signal(os.Interrupt)
	}
	// If Cancel above fails to shut down the subprocess within the wait delay
	// specified below, the subprocess will be killed automatically.
	subprocess.WaitDelay = 2 * time.Second
	return subprocess
}
