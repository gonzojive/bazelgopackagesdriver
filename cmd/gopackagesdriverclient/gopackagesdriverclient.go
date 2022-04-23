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
	"flag"
	"fmt"
	"go/types"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/gonzojive/bazelgopackagesdriver/internal/cmdutil"
	"github.com/gonzojive/bazelgopackagesdriver/internal/configuration"
	"github.com/gonzojive/bazelgopackagesdriver/protocol"
	ctxio "github.com/jbenet/go-context/io"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/johnsiilver/golib/development/portpicker"

	pb "github.com/gonzojive/bazelgopackagesdriver/proto/gopackagesdriverpb"
)

const (
	serverStartupDeadline = time.Second * 10
)

var (
	serverAddr    = os.Getenv("GOPACKAGESDRIVER_SERVER_ADDR") // If missing, looks for a config file
	emptyResponse = &protocol.DriverResponse{
		NotHandled: false,
		Sizes:      types.SizesFor("gc", "amd64").(*types.StdSizes),
		Roots:      []string{},
		Packages:   []*protocol.FlatPackage{},
	}
)

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

func run() (*protocol.DriverResponse, error) {
	flag.Parse()
	ctx, cancel := cmdutil.SignalCancelledContext(context.Background(), os.Interrupt)
	defer cancel()

	// Set up a connection to the server.
	c, closer, err := findOrStartServer(ctx, configuration.FromEnv())
	if closer != nil {
		defer closer.Close()
	}
	if err != nil {
		return nil, err
	}
	glog.Infof("connected to gRPC server... processing stdin")

	queries := flag.Args()

	request, err := protocol.ReadDriverRequest(ctxio.NewReader(ctx, os.Stdin))
	if err != nil {
		return emptyResponse, fmt.Errorf("unable to read request: %w", err)
	}
	glog.Infof("read driver request %v with queries %v", request, queries)

	// Contact the server and print out its response.
	respProto, err := c.LoadPackages(ctx, &pb.LoadPackagesRequest{
		Queries:   queries,
		LoadMode:  uint64(request.Mode),
		EnvParams: configuration.FromEnv().Proto(),
	})
	if err != nil {
		_, healthCheckErr := c.CheckStatus(ctx, &pb.CheckStatusRequest{})
		return emptyResponse, fmt.Errorf("error in server call to LoadPackages: %w; server check rpc code = %v", err, status.Code(healthCheckErr))
	}
	resp := &protocol.DriverResponse{}
	if err := json.Unmarshal([]byte(respProto.GetRawJson()), resp); err != nil {
		return emptyResponse, fmt.Errorf("error unmarshalling: %w", err)
	}
	return resp, nil
}

func connectToServer(ctx context.Context, addr string) (pb.GoPackagesDriverServiceClient, io.Closer, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect to %q: %w", addr, err)
	}
	c := pb.NewGoPackagesDriverServiceClient(conn)

	checkStatusCtx, cancel := context.WithDeadline(ctx, time.Now().Add(serverStartupDeadline))
	defer cancel()
	for {
		resp, err := c.CheckStatus(checkStatusCtx, &pb.CheckStatusRequest{})
		if status.Code(err) == codes.Unavailable {
			// TODO sleep until cancelled
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if err != nil {
			return nil, nil, fmt.Errorf("error checking status: %w", err)
		}
		glog.Infof("server is up with status: %s", resp.GetDebugMessage())
		break
	}

	return c, conn, nil
}

func findOrStartServer(ctx context.Context, params *configuration.Params) (pb.GoPackagesDriverServiceClient, io.Closer, error) {
	if serverAddr != "" {
		return connectToServer(ctx, serverAddr)
	}
	processInfoFiles, err := configuration.FindProcessInfoFiles(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("error finding active server: %w", err)
	}
	for _, f := range processInfoFiles {
		glog.Infof("connecting to already-running driver server at %q", f.Contents.GetGrpcAddress())
		client, closer, err := connectToServer(ctx, f.Contents.GetGrpcAddress())
		if err != nil {
			if closer != nil {
				defer closer.Close()
			}
			glog.Errorf("error connecting to already-running driver server: %v", err)
			continue
		}
		return client, closer, nil
	}
	return startServer(ctx)
}

func startServer(ctx context.Context) (pb.GoPackagesDriverServiceClient, io.Closer, error) {
	serverPort, err := portpicker.TCP(portpicker.Local4())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find a port for starting the gRPC server: %w", err)
	}

	cmd := exec.Command("nohup", "bazelgopackagesdriver",
		"--mode", "server",
		"--grpc_port", strconv.Itoa(serverPort),
	)
	if configDir, err := configuration.WorkspaceConfigDir(ctx, configuration.FromEnv()); err == nil {
		cmd.Args = append(cmd.Args, "-log_dir", configDir)
	}

	cmd.Env = append(cmd.Env, os.Environ()...)

	// Try to detach from the process so it remains after this process ends.
	// https://stackoverflow.com/questions/23031752/start-a-process-in-go-and-detach-from-it
	glog.Infof("starting gopackagesdriver gRPC server...")
	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("command to start server failed: %w", err)
	}
	if err := cmd.Process.Release(); err != nil {
		return nil, nil, fmt.Errorf("error detaching from server process: %w", err)
	}
	glog.Infof("connecting to server port %d", serverPort)
	// TODO(reddaly): Detach child process from parent.
	return connectToServer(ctx, fmt.Sprintf("localhost:%d", serverPort))
}
