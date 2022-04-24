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
	"net"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/gonzojive/bazelgopackagesdriver/internal/cmdutil"
	"github.com/gonzojive/bazelgopackagesdriver/internal/configuration"
	"github.com/gonzojive/bazelgopackagesdriver/internal/idleness"
	"github.com/gonzojive/bazelgopackagesdriver/protocol"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "github.com/gonzojive/bazelgopackagesdriver/proto/gopackagesdriverpb"
)

const debugDisallowEmptyResponse = false // DO NOT SUBMIT

var (
	// It seems https://github.com/bazelbuild/bazel/issues/3115 isn't fixed when specifying
	// the aspect from the command line. Use this trick in the mean time.
	emptyResponse = &protocol.DriverResponse{
		NotHandled: false,
		Sizes:      types.SizesFor("gc", "amd64").(*types.StdSizes),
		Roots:      []string{},
		Packages:   []*protocol.FlatPackage{},
	}

	serverMode     = flag.String("mode", "normal", "Oneof 'server', 'client', 'normal'; If 'server,' start in gRPC server rather than in command line driver mode; if 'client', connects to the server at the given port.")
	grpcPort       = flag.Int("grpc_port", 50051, "The server port")
	lameduckPeriod = flag.Duration("lameduck_duration", time.Second*5, "Time to wait for the server to shut down gracefully if sent a signal.")
	idlePeriod     = flag.Duration("idle_duration", time.Minute*5, "If this period elapses between requests, the server shuts down.")
)

func main() {
	if err := run(); err != nil {
		glog.Exitf("error: %v", err)
	}
}

func run() error {
	flag.Parse()
	glog.Infof("bazelgopackagesdriver started with mode=%q, port=%d", *serverMode, *grpcPort)
	cmdutil.LogFlags()
	go func() {
		for _ = range time.Tick(time.Second * 2) {
			glog.Flush()
		}
	}()

	ctx, cancel := cmdutil.SignalCancelledContext(context.Background(), os.Interrupt)
	defer cancel()

	params := configuration.FromEnv()

	if *serverMode == "server" {
		if params.WorkspaceRoot == "" {
			return fmt.Errorf("must specify workspace root using env variable GOPACKAGESDRIVER_BUILD_WORKSPACE_DIRECTORY or BUILD_WORKSPACE_DIRECTORY")
		}
		glog.Infof("gRPC server starting up at %d", *grpcPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			return fmt.Errorf("failed to listen on --grpc_port: %w", err)
		}

		var serverOpts []grpc.ServerOption
		var shutdownServer func()

		if *idlePeriod > 0 {
			idlenessTracker := idleness.NewIdlenessMonitor(*idlePeriod, func() {
				if shutdownServer == nil {
					return
				}
				shutdownServer()
			})
			defer idlenessTracker.Close()
			serverOpts = append(serverOpts, grpc.UnaryInterceptor(idlenessTracker.UnaryInterceptor()))
		}

		s := grpc.NewServer(serverOpts...)
		shutdownServer = func() {
			glog.Infof("shutting down server due to idleness")
			s.GracefulStop()
		}
		server, err := startServer(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to start gopackagesdriver server: %w", err)
		}
		defer server.close()
		pb.RegisterGoPackagesDriverServiceServer(s, server)
		glog.Infof("gRPC server listening at %v", lis.Addr())

		return serveUntilCancelled(ctx, s, lis)
	}

	switch *serverMode {
	case "normal":
		response, err := runRegularMode(ctx, params)
		if err != nil {
			return err
		}
		if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
			fmt.Fprintf(os.Stderr, "unable to encode response: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("invalid --mode flag = %q", *serverMode)
	}
}

func serveUntilCancelled(ctx context.Context, s *grpc.Server, lis net.Listener) error {
	eg := &errgroup.Group{}
	fnDone := make(chan struct{})
	eg.Go(func() error {
		defer close(fnDone)
		if err := s.Serve(lis); err != nil {
			return fmt.Errorf("Serve() error: %w", err)
		}
		return nil
	})

	stopServerAsync := func() {
		eg.Go(func() error {
			glog.Infof("attempting to gracefully shut down gRPC server...")
			s.GracefulStop()
			return nil
		})
		eg.Go(func() error {
			timer := time.NewTimer(*lameduckPeriod)
			select {
			case <-fnDone:
				timer.Stop()
				return nil
			case <-timer.C:
				glog.Infof("forcing shut down gRPC server...")
				s.Stop()
				return nil
			}
		})
	}

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			stopServerAsync()
			return ctx.Err()
		case <-fnDone:
			return nil
		}
	})
	return eg.Wait()
}

func runRegularMode(ctx context.Context, params *configuration.Params) (*protocol.DriverResponse, error) {
	queries := os.Args[1:]

	request, err := protocol.ReadDriverRequest(os.Stdin)
	if err != nil {
		return emptyResponse, fmt.Errorf("unable to read request: %w", err)
	}
	fmt.Fprintf(os.Stderr, "read driver request %v with queries %v", request, queries)

	bazel, err := NewBazel(ctx, params.BazelBin, params.WorkspaceRoot, params.BazelFlags)
	if err != nil {
		return emptyResponse, fmt.Errorf("unable to create bazel instance: %w", err)
	}

	bazelJsonBuilder, err := NewBazelJSONBuilder(params, bazel, queries...)
	if err != nil {
		return emptyResponse, fmt.Errorf("unable to build JSON files: %w", err)
	}

	jsonFiles, err := bazelJsonBuilder.Build(ctx, request.Mode)
	if err != nil {
		return emptyResponse, fmt.Errorf("unable to build JSON files: %w", err)
	}

	driver, err := NewJSONPackagesDriver(params.WorkspaceRoot, jsonFiles, bazelJsonBuilder.PathResolver())
	if err != nil {
		return emptyResponse, fmt.Errorf("unable to load JSON files: %w", err)
	}

	resp := driver.Match(queries...)
	if debugDisallowEmptyResponse && len(resp.Packages) == 0 {
		return nil, fmt.Errorf("got no packages")
	}
	return resp, nil
}

type server struct {
	pb.UnimplementedGoPackagesDriverServiceServer
	params   *configuration.Params
	rmConfig func() error
}

func startServer(ctx context.Context, params *configuration.Params) (*server, error) {
	s := &server{params: params}
	if _, err := NewBazel(ctx, s.params.BazelBin, s.params.WorkspaceRoot, s.params.BazelFlags); err != nil {
		return nil, fmt.Errorf("unable to create bazel instance: %w", err)
	}

	rmConfig, err := configuration.WriteProcessInfoToFile(ctx, params, fmt.Sprintf("localhost:%d", *grpcPort))
	if err != nil {
		return nil, err
	}
	s.rmConfig = rmConfig

	return s, nil
}

func (s *server) close() error {
	if s.rmConfig != nil {
		return s.rmConfig()
	}
	return nil
}

func (s *server) CheckStatus(ctx context.Context, req *pb.CheckStatusRequest) (*pb.CheckStatusResponse, error) {
	return &pb.CheckStatusResponse{
		DebugMessage: fmt.Sprintf("everything is up... %v", s.params),
	}, nil
}

func (s *server) LoadPackages(ctx context.Context, req *pb.LoadPackagesRequest) (*pb.LoadPackagesResponse, error) {
	request := &protocol.DriverRequest{
		Mode: protocol.LoadMode(req.LoadMode),
	}
	glog.Infof("read driver request %v with queries %v", request, req.GetQueries())

	bazel, err := NewBazel(ctx, s.params.BazelBin, s.params.WorkspaceRoot, s.params.BazelFlags)
	if err != nil {
		return nil, fmt.Errorf("unable to create bazel instance: %w", err)
	}

	bazelJsonBuilder, err := NewBazelJSONBuilder(s.params, bazel, req.GetQueries()...)
	if err != nil {
		return nil, fmt.Errorf("unable to build JSON files: %w", err)
	}

	jsonFiles, err := bazelJsonBuilder.Build(ctx, request.Mode)
	if err != nil {
		return nil, fmt.Errorf("unable to build JSON files: %w", err)
	}

	driver, err := NewJSONPackagesDriver(s.params.WorkspaceRoot, jsonFiles, bazelJsonBuilder.PathResolver())
	if err != nil {
		return nil, fmt.Errorf("unable to load JSON files: %w", err)
	}

	respObj := driver.Match(req.GetQueries()...)
	respJSONBytes, err := json.Marshal(respObj)
	if err != nil {
		return nil, fmt.Errorf("internal error encoding JSON: %w", err)
	}
	glog.Infof("responding with JSON:\n  %s", string(respJSONBytes))
	return &pb.LoadPackagesResponse{
		RawJson: string(respJSONBytes),
	}, nil
}
