package workspace_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
	"github.com/gonzojive/bazelgopackagesdriver/internal/runfiles"
	"golang.org/x/tools/go/packages"
)

const (
	driverRunfilesPath   = "cmd/bazelgopackagesdriver/bazelgopackagesdriver_/bazelgopackagesdriver"
	bazelBinRunfilesPath = "external/build_bazel_bazel_5_1_1_binary/file/bazel_for_integration_testing"

	// typicalLoadMode matches what is commonly passed to packages.Load by
	// gopls. See
	// https://github.com/golang/tools/blob/master/internal/lsp/cache/snapshot.go.
	typicalLoadMode = packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypesSizes |
		packages.NeedModule
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: `
-- WORKSPACE --
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

# Bug in rules_go for 1.18: https://github.com/bazelbuild/rules_go/issues/3110
go_register_toolchains(version = "1.17")

-- BUILD.bazel --
load("@io_bazel_rules_go//go:def.bzl", "go_test")
go_test(
    name = "fail_fast_test",
    srcs = ["fail_fast_test.go"],
)
-- BUILD.bazel --
load("@io_bazel_rules_go//go:def.bzl", "go_test")
go_test(
    name = "example_test",
    srcs = ["example_test.go"],
	deps = [
		"//super/fun/thingy:thingy",
	],
)

-- example_test.go --
package test_fail_fast

import (
	"testing"

	_ "example.xyz/thingy"
)

func TestShouldFail(t *testing.T) {
	t.Fail()
}
-- super/fun/thingy/BUILD.bazel --
load("@io_bazel_rules_go//go:def.bzl", "go_library")
go_library(
    name = "thingy",
    srcs = ["thingy.go"],
	importpath = "example.xyz/thingy",
	visibility = ["//:__subpackages__"],
)

-- super/fun/thingy/thingy.go --
package thingy

func Foo() int { return 0 }
`,
	})
}

func TestFunctionalWorkspace(t *testing.T) {
	if err := bazel_testing.RunBazel("test", "//:example_test", "--test_runner_fail_fast"); err == nil {
		t.Fatal("got success; want failure")
	} else if bErr, ok := err.(*bazel_testing.StderrExitError); !ok {
		t.Fatalf("got %v; want StderrExitError", err)
	} else if code := bErr.Err.ExitCode(); code != 3 {
		t.Fatalf("got code %d; want code 3: %v", code, err)
	}

	logPath := filepath.FromSlash("bazel-testlogs/example_test/test.log")
	logData, err := ioutil.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(logData, []byte("TestShouldFail")) {
		t.Fatalf("test log does not contain 'TestShouldFail': %q", logData)
	}
}

func TestPackagesLoad(t *testing.T) {
	cmd := bazel_testing.BazelCmd()
	driver, err := runfiles.Runfile(driverRunfilesPath)
	if err != nil {
		t.Fatalf("failed to get path of driver: %v", err)
	}
	bazelBin, err := runfiles.Runfile(bazelBinRunfilesPath)
	if err != nil {
		t.Fatalf("failed to get path of bazel binary: %v", err)
	}
	cmd.Args[0] = bazelBin
	cmd.Path = bazelBin

	flags := cmd.Args[1:]
	workspaceDir, err := os.Getwd() // This works because bazel_testing's main function changes the working directory to the workspace dir.
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	cfgEnv := append([]string{}, cmd.Env...)
	cfgEnv = append(cfgEnv,
		fmt.Sprintf("GOPACKAGESDRIVER=%s", driver),
		fmt.Sprintf("GOPACKAGESDRIVER_BAZEL=%s", cmd.Path),
		fmt.Sprintf("GOPACKAGESDRIVER_BAZEL_FLAGS=%s", strings.Join(flags, " ")),
		fmt.Sprintf("GOPACKAGESDRIVER_BUILD_WORKSPACE_DIRECTORY=%s", workspaceDir),
	)
	t.Logf("got driver path %q", driver)
	t.Logf("got bazel: %q", cmd.Path)
	t.Logf("got bazel flags:\n  %s", strings.Join(flags, "\n  "))
	t.Logf("got workspace directory: %q", workspaceDir)

	os.Setenv("GOPACKAGESPRINTDRIVERERRORS", "1") // Rquired to displays stderr from the driver.
	packages, err := packages.Load(&packages.Config{
		Dir:  workspaceDir,
		Env:  cfgEnv,
		Mode: typicalLoadMode, //
	}, "file=example_test.go")

	if err != nil {
		t.Fatalf("failed to load packages: %v", err)
	}
	t.Errorf("got %v packages", packages)
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
