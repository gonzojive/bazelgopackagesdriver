package workspace_ext_dep_test

// TODO(reddaly): Use bazel_external_dependency_archive to supply an external
// archive. Per
// https://github.com/bazelbuild/bazel-integration-testing/blob/165440b2dbda885f8d1ccb8d0f417e6cf8c54f17/tools/import.bzl
//
// Caching of archive data is described more here:
// https://github.com/bazelbuild/bazel-integration-testing/blob/165440b2dbda885f8d1ccb8d0f417e6cf8c54f17/java/build/bazel/tests/integration/RepositoryCache.java

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: `
-- WORKSPACE --
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715f",
    urls = [
		"https://bazel.build/docs/sandboxing.zip",
        #"https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
        #"https://github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.18")

-- BUILD.bazel --
load("@io_bazel_rules_go//go:def.bzl", "go_test")
go_test(
    name = "fail_fast_test",
    srcs = ["fail_fast_test.go"],
)
-- fail_fast_test.go --
package test_fail_fast
import "testing"
func TestShouldFail(t *testing.T) {
	t.Fail()
}
func TestShouldNotRun(t *testing.T) {
	t.Fail()
}
`,
	})
}

func TestSimple(t *testing.T) {
	if err := bazel_testing.RunBazel("test", "//:fail_fast_test", "--test_runner_fail_fast"); err == nil {
		t.Fatal("got success; want failure")
	} else if bErr, ok := err.(*bazel_testing.StderrExitError); !ok {
		t.Fatalf("got %v; want StderrExitError", err)
	} else if code := bErr.Err.ExitCode(); code != 3 {
		t.Fatalf("got code %d; want code 3: %v", code, err)
	}

	logPath := filepath.FromSlash("bazel-testlogs/fail_fast_test/test.log")
	logData, err := ioutil.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(logData, []byte("TestShouldFail")) {
		t.Fatalf("test log does not contain 'TestShouldFail': %q", logData)
	}

	if bytes.Contains(logData, []byte("TestShouldNotRun")) {
		t.Fatalf("test log contains 'TestShouldNotRun' but should not: %q", logData)
	}
}
