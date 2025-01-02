package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
		build          bool
		check          bool
		fix            bool
		lint           bool
		test           bool
		updateVersions bool
	)

	flag.BoolVar(&build, "build", false, "Builds the library source code")
	flag.BoolVar(&check, "check", false, "Runs all checks on the library source code")
	flag.BoolVar(&fix, "fix", false, "Fixes as many lint problems as possible on the library source code")
	flag.BoolVar(&lint, "lint", false, "Lints the library source code")
	flag.BoolVar(&test, "test", false, "Tests the library source code")
	flag.BoolVar(&updateVersions, "update-versions", false, "Updates the versions of all dependencies")
	flag.Parse()

	switch {
	case build:
		doBuild()
	case check:
		doBuild()
		doLint()
		doTest()
	case fix:
		doFix()
	case lint:
		doLint()
	case test:
		doTest()
	case updateVersions:
		doUpdateVersions()
	default:
		flag.Usage()
		os.Exit(usageExitCode)
	}
}

func cdToRootProjectDir() {
	must(os.Chdir(".."), cdToRootProjectDirExitCode)
}

func doBuild() {
	fmt.Println("Building...")
	mustRun(cmd("go", "build", "./..."))
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

func doDepawareFix() {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(max(1, runtime.NumCPU()/2))

	depawareFiles := make(chan string)
	group.Go(func() error {
		defer close(depawareFiles)
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
			case depawareFiles <- "./" + path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	for file := range depawareFiles {
		group.Go(func() error {
			dir := "./" + filepath.Dir(file)
			fmt.Printf("Fixing with depaware concurrently on directory %s...\n", dir)
			c := cmd("depaware", dir)
			var b bytes.Buffer
			c.Stdout = &b
			if err := c.Run(); err != nil {
				return err
			}
			return os.WriteFile(file, b.Bytes(), 0o600)
		})
	}

	mustNotError(group.Wait())
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
				return errors.New("eg found a problem (see above)")
			}
			return nil
		})
	}
}

func doEgFix() {
	egTemplateFiles, err := filepath.Glob("eg/*.template")
	mustNotError(err)

	for _, file := range egTemplateFiles {
		fmt.Printf("Fixing with eg template %s...\n", file)
		mustRun(cmd(
			"eg",
			"-t",
			file,
			"-w",
			"./...",
		))
	}
}

func doFix() {
	doGoModTidy()
	doGoModDownload()
	doEgFix()
	doGolangciLintFix()
	doDepawareFix()
}

func doGoModTidy() {
	fmt.Println("Running 'go mod tidy'...")
	mustRun(cmd("go", "mod", "tidy"))
}

func doGoModVerify() {
	fmt.Println("Running 'go mod verify'...")
	mustRun(cmd("go", "mod", "verify"))
}

func doGoModDownload() {
	fmt.Println("Running 'go mod download'...")
	mustRun(cmd("go", "mod", "download"))
}

func doGolangciLint(group *errgroup.Group) {
	group.Go(func() error {
		fmt.Println("Linting with golangci-lint concurrently...")
		return cmd("golangci-lint", "run").Run()
	})
}

func doGolangciLintFix() {
	fmt.Println("Fixing with golangci-lint...")
	mustRun(cmd("golangci-lint", "run", "--fix"))
}

func doLint() {
	doGoModTidy()
	mustRun(cmd("git", "diff", "--exit-code", "--", "go.mod", "go.sum"), func() {
		fmt.Println("'go mod tidy' changed files")
	})
	doGoModVerify()

	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(max(1, runtime.NumCPU()/2))

	doGolangciLint(group)
	doNilaway(group)
	doEg(ctx, group)
	doDepaware(ctx, group)

	mustNotError(group.Wait())
}

func doNilaway(group *errgroup.Group) {
	group.Go(func() error {
		fmt.Println("Linting with nilaway concurrently...")
		return cmd(
			"nilaway",
			"-include-pkgs",
			"github.com/jbduncan/go-containers",
			"./...",
		).Run()
	})
}

func doTest() {
	fmt.Println("Testing...")
	mustRun(cmd("go", "test", "-shuffle=on", "-race", "./..."))
}

func doUpdateVersions() {
	fmt.Println("Updating versions...")
	mustRun(cmd("mise", "up", "--bump"))
	mustRun(cmd("go", "get", "-u", "-t", "./..."))
	doGoModTidy()
	doGoModVerify()
	doGoModDownload()
	doDepawareFix()
}

func cmd(name string, args ...string) *exec.Cmd {
	result := exec.Command(name, args...)
	result.Stdout = os.Stdout
	result.Stderr = os.Stderr
	return result
}

func mustRun(cmd *exec.Cmd, onError ...func()) {
	if err := cmd.Run(); err != nil {
		for _, f := range onError {
			f()
		}
		mustNotError(err)
	}
}

func mustNotError(err error) {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		// Report the whole error
		must(err, exitErr.ExitCode())
	} else {
		must(err, otherErrorExitCode)
	}
}

func must(err error, exitCode int) {
	if err != nil {
		fmt.Printf("err: %v\n", err)
		// A deep exit is allowed here because this is just a script.
		os.Exit(exitCode) //nolint:revive
	}
}
