// Package configuration is used to configure the gopackagesdriver.
package configuration

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/golang/glog"
	"github.com/gonzojive/bazelgopackagesdriver/internal/cmdutil"
	"google.golang.org/protobuf/encoding/prototext"

	pb "github.com/gonzojive/bazelgopackagesdriver/proto/gopackagesdriverpb"
)

const (
	workspaceConfigDir    = ".bazelgopackagesdriver"
	processInfoPrefix     = "process-info-"
	configDirPermissions  = os.FileMode(0700) | os.ModeDir
	configFilePermissions = os.FileMode(0600)
)

// Params is a struct to hold the values of parameters used by the driver. These
// typically come from the
type Params struct {
	RulesGoRepositoryName string
	BazelBin              string
	BazelFlags            []string
	BazelQueryFlags       []string
	BazelQueryScope       string
	BazelBuildFlags       []string
	WorkspaceRoot         string
}

func (p *Params) Proto() *pb.EnvParams {
	return &pb.EnvParams{
		RulesGoRepositoryName: p.RulesGoRepositoryName,
		BazelBin:              p.BazelBin,
		BazelFlags:            p.BazelFlags,
		BazelQueryFlags:       p.BazelQueryFlags,
		BazelQueryScope:       p.BazelQueryScope,
		BazelBuildFlags:       p.BazelBuildFlags,
		WorkspaceRoot:         p.WorkspaceRoot,
	}
}

// FromEnv returns a *Params object based on os.Env values.
func FromEnv() *Params {
	// Workspace root if the binary is run using bazel run. See https://docs.bazel.build/versions/main/user-manual.html#run.
	workspaceRootFromBazelRun := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	// Let GOPACKAGESDRIVER_BUILD_WORKSPACE_DIRECTORY override BUILD_WORKSPACE_DIRECTORY... probably not necessary.
	workspaceRoot := cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_BUILD_WORKSPACE_DIRECTORY", workspaceRootFromBazelRun)

	return &Params{
		RulesGoRepositoryName: cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_RULES_GO_REPOSITORY_NAME", "@io_bazel_rules_go"),
		BazelBin:              cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_BAZEL", "bazel"),
		BazelFlags:            strings.Fields(os.Getenv("GOPACKAGESDRIVER_BAZEL_FLAGS")),
		BazelQueryFlags:       strings.Fields(os.Getenv("GOPACKAGESDRIVER_BAZEL_QUERY_FLAGS")),
		BazelQueryScope:       cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_BAZEL_QUERY_SCOPE", ""),
		BazelBuildFlags:       strings.Fields(os.Getenv("GOPACKAGESDRIVER_BAZEL_BUILD_FLAGS")),
		WorkspaceRoot:         workspaceRoot,
	}
}

// WriteProcessInfoToFile writes a .pbtxt (protobuf text) file to a file in the
// workspace with information about this process.
func WriteProcessInfoToFile(ctx context.Context, params *Params, grpcAddr string) (func() error, error) {
	info := &pb.ProcessInfo{
		Pid:         int64(os.Getpid()),
		GrpcAddress: grpcAddr,
	}
	contents, err := prototext.MarshalOptions{Multiline: true}.Marshal(info)
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(params.WorkspaceRoot, workspaceConfigDir)
	if err := mkDirIfNoExist(configDir, configDirPermissions); err != nil {
		return nil, fmt.Errorf("failed to make config directory: %w", err)
	}
	fileName := fmt.Sprintf("%s%d.pbtxt", processInfoPrefix, info.GetPid())
	fullPath := filepath.Join(configDir, fileName)
	if err := ioutil.WriteFile(fullPath, contents, configFilePermissions); err != nil {
		return nil, fmt.Errorf("error writing config file %q: %w", fullPath, err)
	}
	return func() error {
		if err := os.Remove(fullPath); err != nil {
			return fmt.Errorf("failed to remove ProcessInfo file %q: %w", fullPath, err)
		}
		return nil
	}, nil
}

// WorkspaceConfigDir returns the configuration directory in the workspace, if it exists.
func WorkspaceConfigDir(ctx context.Context, params *Params) (string, error) {
	configDir, err := filepath.Abs(filepath.Join(params.WorkspaceRoot, workspaceConfigDir))
	if err != nil {
		return "", fmt.Errorf("error getting absolute path of configuration dir in workspace: %w", err)
	}
	if _, err := os.Stat(configDir); err != nil {
		return configDir, err
	}
	return configDir, nil
}

func mkDirIfNoExist(dir string, perm os.FileMode) error {
	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.Mkdir(dir, perm)
	}
	if !fi.IsDir() {
		return fmt.Errorf("%q is not a directory", dir)
	}
	return nil
}

type ProcessInfoFile struct {
	Path     string
	Contents *pb.ProcessInfo
}

// FindProcessInfoFiles returns all of the process info files in the config
// directory. This function can be used to connect to an active server.
func FindProcessInfoFiles(ctx context.Context, params *Params) ([]*ProcessInfoFile, error) {
	configDir, err := WorkspaceConfigDir(ctx, params)
	if os.IsNotExist(err) {
		return nil, nil
	}
	paths, err := filepath.Glob(filepath.Join(configDir, fmt.Sprintf("%s*", processInfoPrefix)))
	if err != nil {
		return nil, fmt.Errorf("error searching for process info files: %w", err)
	}
	var out []*ProcessInfoFile
	for _, p := range paths {
		contents, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("error reading ProcessInfo file %q: %w", p, err)
		}
		infoProto := &pb.ProcessInfo{}
		if err := prototext.Unmarshal(contents, infoProto); err != nil {
			return nil, fmt.Errorf("error reading ProcessInfo file %q: %w", p, err)
		}
		infoObj := &ProcessInfoFile{
			Path:     p,
			Contents: infoProto,
		}
		if !infoObj.ProcessExists() {
			glog.Warningf("garbage collecting ProcessInfo file for expired process, %q", infoObj.Path)
			if err := os.Remove(infoObj.Path); err != nil {
				glog.Errorf("error removing old ProcessInfo file %q: %v", infoObj.Path, err)
			}
			continue
		}
		out = append(out, infoObj)
	}
	return out, nil
}

// ProcessExists reports if the process referenced by the ProcessInfo file is
// active.
//
// It does this using the approach in
// https://stackoverflow.com/questions/15204162/check-if-a-process-exists-in-go-way.
// As such, this may only work on UNIX-like systems.
func (f *ProcessInfoFile) ProcessExists() bool {
	proc, err := os.FindProcess(int(f.Contents.GetPid()))
	if err != nil {
		return false
	}

	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return false
	}
	return true
}
