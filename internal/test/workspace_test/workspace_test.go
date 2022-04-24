package workspace_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
	"github.com/phst/runfiles"
	"golang.org/x/tools/go/packages"
)

const (
	driverRunfilesPath = "cmd/bazelgopackagesdriver/bazelgopackagesdriver"
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: `
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
	driver, err := runfiles.Path(driverRunfilesPath)
	if err != nil {
		t.Fatalf("failed to get path of driver: %v", err)
	}
	cfgEnv := []string{
		fmt.Sprintf("GOPACKAGESDRIVER=blah", driver),
	}
	packages, err := packages.Load(&packages.Config{
		Dir: cmd.Dir,
		Env: cfgEnv,
	}, "example.xyz/...")
	if err != nil {
		t.Fatalf("failed to load packages: %v", err)
	}
	t.Errorf("got %v packages", packages)
}
