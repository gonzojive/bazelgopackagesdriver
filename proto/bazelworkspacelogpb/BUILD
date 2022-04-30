load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@build_stack_rules_proto//rules:proto_compiled_sources.bzl", "proto_compiled_sources")

# Regenerate compiled source with
# bazel run  //proto/bazelworkspacelogpb:compiled_sources.update
proto_compiled_sources(
    name = "compiled_sources",
    srcs = [
        "bazelworkspacelogpb.pb.go",
        #"@io_bazel//src/main/java/com/google/devtools/build/lib/bazel/debug:workspace_log.proto",
    ],
    args = [
        # See
        #"--go_out=/dev/null",
        #"--protoc-gen-go_opt",
        #"Msrc/main/java/com/google/devtools/build/lib/bazel/debug/workspace_log.proto=github.com/gonzojive/bazelgopackagesdriver/proto/bazelworkspacelogpb",
    ],
    output_mappings = [
        "bazelworkspacelogpb.pb.go=github.com/gonzojive/bazelgopackagesdriver/proto/bazelworkspacelogpb/workspace_log.pb.go",
    ],
    plugins = [
        "@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go",
    ],
    proto = "@io_bazel//src/main/java/com/google/devtools/build/lib/bazel/debug:workspace_log_proto",
    verbose = True,
    visibility = ["//proto:__subpackages__"],
)

go_library(
    name = "bazelworkspacelogpb",
    srcs = ["bazelworkspacelogpb.pb.go"],
    importpath = "github.com/gonzojive/bazelgopackagesdriver/proto/bazelworkspacelogpb",
    visibility = ["//visibility:public"],
    deps = [
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)
