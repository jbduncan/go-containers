package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/sync/errgroup"
)

const (
	usageExitCode              = 2
	cdToRootProjectDirExitCode = 3
	otherErrorExitCode         = 4
)

func main() {
	cdToRootProjectDir()
	var (
		build bool
		check bool
		lint  bool
		test  bool
	)

	flag.BoolVar(&build, "build", false, "Builds the library source code")
	flag.BoolVar(&check, "check", false, "Runs all checks on the library source code")
	flag.BoolVar(&lint, "lint", false, "Lints the library source code")
	flag.BoolVar(&test, "test", false, "Tests the library source code")
	flag.Parse()

	switch {
	case build:
		doBuild()
	case check:
		doBuild()
		doLint()
		doTest()
	case lint:
		doLint()
	case test:
		doTest()
	}

	// TODO: Implement 'fix' from Makefile
	// TODO: Implement 'update_versions' from Makefile
	// TODO: Implement 'check' from Makefile

	flag.Usage()
	os.Exit(usageExitCode)
}

func cdToRootProjectDir() {
	must(os.Chdir(".."), cdToRootProjectDirExitCode)
}

func doBuild() {
	mustRun(cmd("go", "build", "./..."))
}

func doTest() {
	mustRun(cmd("go", "test", "-shuffle=on", "-race", "./..."))
}

func doLint() {
	fmt.Println("Linting 'go mod tidy' results...")
	mustRun(cmd("go", "mod", "tidy"))
	mustRun(cmd("git", "diff", "--exit-code", "--", "go.mod", "go.sum"), func() {
		fmt.Println("'go mod tidy' changed files")
	})

	fmt.Println("Linting with 'go mod verify'...")
	mustRun(cmd("go", "mod", "verify"))

	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(max(1, runtime.NumCPU()/2))

	runGolangciLint(group)
	runNilaway(group)
	runEg(ctx, group)
	runDepaware(ctx, group)

	err := group.Wait()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			must(exitErr, exitErr.ExitCode())
		} else {
			must(err, otherErrorExitCode)
		}
	}
}

func runGolangciLint(group *errgroup.Group) {
	group.Go(func() error {
		fmt.Println("Linting with golangci-lint concurrently...")
		return cmd(
			"go",
			"run",
			"github.com/golangci/golangci-lint/cmd/golangci-lint",
			"run",
		).Run()
	})
}

func runNilaway(group *errgroup.Group) {
	group.Go(func() error {
		fmt.Println("Linting with nilaway concurrently...")
		return cmd(
			"go",
			"run",
			"go.uber.org/nilaway/cmd/nilaway",
			"-include-pkgs",
			"github.com/jbduncan/go-containers",
			"./...",
		).Run()
	})
}

func runEg(ctx context.Context, group *errgroup.Group) {
	egTemplateFiles := make(chan []string, 1)
	group.Go(func() error {
		defer close(egTemplateFiles)
		matches, err := filepath.Glob("eg/*.template")
		if err != nil {
			return err
		}

		select {
		case egTemplateFiles <- matches:
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	})

	for files := range egTemplateFiles {
		for _, file := range files {
			file := file
			group.Go(func() error {
				fmt.Printf("Linting with eg concurrently with template %s...\n", file)
				return cmd(
					"go",
					"run",
					"golang.org/x/tools/cmd/eg",
					"-t",
					file,
					"./...",
				).Run()
			})
		}
	}
}

func runDepaware(ctx context.Context, group *errgroup.Group) {
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
		dir := dir
		group.Go(func() error {
			fmt.Printf("Linting with depaware concurrently on directory %s...\n", dir)
			return cmd(
				"go",
				"run",
				"github.com/tailscale/depaware",
				"--check",
				dir,
			).Run()
		})
	}
}

func cmd(name string, args ...string) *exec.Cmd {
	result := exec.Command(name, args...)
	result.Stdout = os.Stdout
	result.Stderr = os.Stderr
	return result
}

func must(err error, exitCode int) {
	if err != nil {
		fmt.Printf("err: %v\n", err)
		// A deep exit is allowed here because this is just a script.
		os.Exit(exitCode) //nolint:revive
	}
}

func mustRun(cmd *exec.Cmd, onError ...func()) {
	if err := cmd.Run(); err != nil {
		for _, f := range onError {
			f()
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			must(exitErr, exitErr.ExitCode())
		} else {
			must(err, otherErrorExitCode)
		}
	}
}
