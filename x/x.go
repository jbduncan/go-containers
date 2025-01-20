package main

import (
	"bytes"
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

	astGrepBinary      = "ast-grep"
	egBinary           = "eg"
	goBinary           = "go"
	golangciLintBinary = "golangci-lint"
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
		doCheck()
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

func doAstGrepFix() {
	fmt.Printf("Fixing with %s...\n", astGrepBinary)
	mustRun(cmd(astGrepBinary, "scan", "--update-all"))
}

func doBuild() {
	mustRun(cmd("mise", "run", "build"))
}

func doCheck() {
	mustRun(cmd("mise", "run", "check"))
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

func doEgFix() {
	egTemplateFiles, err := filepath.Glob("eg/*.template")
	mustNotError(err)

	for _, file := range egTemplateFiles {
		fmt.Printf("Fixing with %s template %s...\n", egBinary, file)
		mustRun(cmd(
			egBinary,
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
	doMiseFmtFix()
	doAstGrepFix()
	doEgFix()
	doGolangciLintFix()
	doDepawareFix()
}

func doGoModTidy() {
	mustRun(cmd("mise", "go-mod-tidy"))
}

func doGoModVerify() {
	mustRun(cmd("mise", "go-mod-verify"))
}

func doGoModDownload() {
	fmt.Printf("Running '%s mod download'...\n", goBinary)
	mustRun(cmd(goBinary, "mod", "download"))
}

func doGolangciLintFix() {
	fmt.Printf("Fixing with %s...\n", golangciLintBinary)
	mustRun(cmd(golangciLintBinary, "run", "--fix"))
}

func doLint() {
	mustRun(cmd("mise", "run", "lint"))
}

func doMiseFmtFix() {
	mustRun(cmd("mise", "mise-fmt"))
}

func doTest() {
	mustRun(cmd("mise", "run", "test"))
}

func doUpdateVersions() {
	fmt.Println("Updating versions...")
	mustRun(cmd("mise", "up"))
	mustRun(cmd(goBinary, "get", "-u", "-t", "./..."))
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
