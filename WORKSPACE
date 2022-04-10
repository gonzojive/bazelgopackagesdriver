# Bazel workspace file

workspace(
    name = "bazelgopackagesdriver",
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

### rules_go: https://github.com/bazelbuild/rules_go

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
    ],
)

git_repository(
    name = "rules_pkg",
    commit = "4f8f6ed027c07b92e4ee5a8b4b51d8673fcc60ee", # v0.7
    remote = "https://github.com/bazelbuild/rules_pkg.git",
)

git_repository(
    name = "com_google_protobuf",
    commit = "bc799d78f81115940eec953e2937245c70e3e6e4", # v3.20.0
    remote = "https://github.com/protocolbuffers/protobuf.git",
)


load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.18")

gazelle_dependencies()

### https://github.com/stackb/rules_proto related stuff

# Branch: master
# Commit: 7c95feba87ae269d09690fcebb18c77d8b8bcf6a
# Date: 2021-11-16 02:17:58 +0000 UTC
# URL: https://github.com/stackb/rules_proto/commit/7c95feba87ae269d09690fcebb18c77d8b8bcf6a
#
# V2 (#193)
# Size: 885598 (886 kB)
http_archive(
    name = "build_stack_rules_proto",
    sha256 = "1190c296a9f931343f70e58e5f6f9ee2331709be4e17001bb570e41237a6c497",
    strip_prefix = "rules_proto-7c95feba87ae269d09690fcebb18c77d8b8bcf6a",
    urls = ["https://github.com/stackb/rules_proto/archive/7c95feba87ae269d09690fcebb18c77d8b8bcf6a.tar.gz"],
)

register_toolchains("@build_stack_rules_proto//toolchain:standard")

load("@build_stack_rules_proto//deps:core_deps.bzl", "core_deps")

core_deps()

load("@build_stack_rules_proto//:go_deps.bzl", "gazelle_protobuf_extension_go_deps")

gazelle_protobuf_extension_go_deps()

load("@build_stack_rules_proto//deps:protobuf_core_deps.bzl", "protobuf_core_deps")

protobuf_core_deps()

load("@build_stack_rules_proto//deps:go_core_deps.bzl", "go_core_deps")

go_core_deps()
