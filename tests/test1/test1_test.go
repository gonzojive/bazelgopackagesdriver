package test1

import (
	"context"
	"fmt"
	"os"
	"testing"

	"golang.org/x/tools/go/packages"
)

type exampleWorkspace struct {
	files []fileEntry
}

type fileEntry struct {
	path     string
	contents string
}

var exampleProject1 = exampleWorkspace{
	files: []fileEntry{
		{
			path: "WORKSPACE",
			contents: `
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
    ],
)

`,
		},
	},
}

func TestPackagesLoad1(t *testing.T) {
	// tmpDir, err := ioutil.TempDir(t.TempDir(), "project-root")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	workingDir := "/home/red/git/bazelgopackagesdriver"

	env := []string{
		//"GOPACKAGESDRIVER=gopackagesdriverclient",
		"GOPACKAGESDRIVER=mydriver.sh",
		"GOPACKAGESPRINTDRIVERERRORS=1",
		"BUILD_WORKSPACE_DIRECTORY=/home/red/git/bazelgopackagesdriver",
	}
	env = append(env, os.Environ()...)

	cfg := &packages.Config{
		Context: context.Background(),
		Dir:     workingDir,
		Env:     env,
		//BuildFlags: ,
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypesSizes |
			packages.NeedModule,
		Fset:    nil,
		Overlay: nil,
		// ParseFile: func(*token.FileSet, string, []byte) (*ast.File, error) {
		// 	panic("go/packages must not be used to parse files")
		// },
		Logf: func(format string, args ...interface{}) {
			t.Logf("packages config: "+format, args)
		},
		Tests: true,
	}
	query := fmt.Sprintf("file=%s/tests/test1/test1_test.go", workingDir)

	pkgs, err := packages.Load(cfg, query)
	if err != nil {
		t.Fatalf("error with packages.Load: %v", err)
	}
	t.Errorf("got %d packages for query %q: %v", len(pkgs), query, pkgs)
}
