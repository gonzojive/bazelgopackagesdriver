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
	"log"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/gonzojive/bazelgopackagesdriver/internal/cmdutil"
	"github.com/gonzojive/bazelgopackagesdriver/protocol"
	ctxio "github.com/jbenet/go-context/io"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/gonzojive/bazelgopackagesdriver/proto/gopackagesdriverpb"
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
	serverAddr            = cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_SERVER_ADDR", "localhost:50051")
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

	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGoPackagesDriverServiceClient(conn)

	// Contact the server and print out its response.
	respProto, err := c.LoadPackages(ctx, &pb.LoadPackagesRequest{
		Queries:  queries,
		LoadMode: uint64(request.Mode),
	})
	if err != nil {
		return emptyResponse, fmt.Errorf("error: %w", err)
	}
	resp := &protocol.DriverResponse{}
	if err := json.Unmarshal([]byte(respProto.GetRawJson()), resp); err != nil {
		return emptyResponse, fmt.Errorf("error unmarshalling: %w", err)
	}
	return resp, nil
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
