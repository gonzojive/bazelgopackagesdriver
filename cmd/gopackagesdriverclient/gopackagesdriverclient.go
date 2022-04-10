// Copyright 2021 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go/types"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/gonzojive/bazelgopackagesdriver/internal/cmdutil"
	"github.com/gonzojive/bazelgopackagesdriver/protocol"
	ctxio "github.com/jbenet/go-context/io"
)

var (
	// It seems https://github.com/bazelbuild/bazel/issues/3115 isn't fixed when specifying
	// the aspect from the command line. Use this trick in the mean time.
	rulesGoRepositoryName = cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_RULES_GO_REPOSITORY_NAME", "@io_bazel_rules_go")
	bazelBin              = cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_BAZEL", "bazel")
	bazelFlags            = strings.Fields(os.Getenv("GOPACKAGESDRIVER_BAZEL_FLAGS"))
	bazelQueryFlags       = strings.Fields(os.Getenv("GOPACKAGESDRIVER_BAZEL_QUERY_FLAGS"))
	bazelQueryScope       = cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_BAZEL_QUERY_SCOPE", "")
	bazelBuildFlags       = strings.Fields(os.Getenv("GOPACKAGESDRIVER_BAZEL_BUILD_FLAGS"))
	workspaceRoot         = os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	emptyResponse         = &protocol.DriverResponse{
		NotHandled: false,
		Sizes:      types.SizesFor("gc", "amd64").(*types.StdSizes),
		Roots:      []string{},
		Packages:   []*protocol.FlatPackage{},
	}
)

func run() (*protocol.DriverResponse, error) {
	ctx, cancel := cmdutil.SignalCancelledContext(context.Background(), os.Interrupt)
	defer cancel()

	queries := os.Args[1:]

	request, err := protocol.ReadDriverRequest(ctxio.NewReader(ctx, os.Stdin))
	if err != nil {
		return emptyResponse, fmt.Errorf("unable to read request: %w", err)
	}
	glog.Infof("read driver request %v with queries %v", request, queries)
	return nil, fmt.Errorf("not implemented")
}

func main() {
	response, err := run()
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "unable to encode response: %v", err)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		// gopls will check the packages driver exit code, and if there is an
		// error, it will fall back to go list. Obviously we don't want that,
		// so force a 0 exit code.
		os.Exit(0)
	}
}
