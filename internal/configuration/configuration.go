// Package configuration is used to configure the gopackagesdriver.
package configuration

import (
	"os"
	"strings"

	"github.com/gonzojive/bazelgopackagesdriver/internal/cmdutil"
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

// FromEnv returns a *Params object based on os.Env values.
func FromEnv() *Params {
	return &Params{
		RulesGoRepositoryName: cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_RULES_GO_REPOSITORY_NAME", "@io_bazel_rules_go"),
		BazelBin:              cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_BAZEL", "bazel"),
		BazelFlags:            strings.Fields(os.Getenv("GOPACKAGESDRIVER_BAZEL_FLAGS")),
		BazelQueryFlags:       strings.Fields(os.Getenv("GOPACKAGESDRIVER_BAZEL_QUERY_FLAGS")),
		BazelQueryScope:       cmdutil.LookupEnvOrDefault("GOPACKAGESDRIVER_BAZEL_QUERY_SCOPE", ""),
		BazelBuildFlags:       strings.Fields(os.Getenv("GOPACKAGESDRIVER_BAZEL_BUILD_FLAGS")),
		WorkspaceRoot:         os.Getenv("BUILD_WORKSPACE_DIRECTORY"),
	}
}
