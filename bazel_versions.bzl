# Created according to 
# https://github.com/bazel-contrib/rules_bazel_integration_test#1-configure-your-workspace-to-use-rules_bazel_integration_test

CURRENT_BAZEL_VERSION = "5.1.0"

OTHER_BAZEL_VERSIONS = [
    "4.2.2",
    "6.0.0-pre.20220328.1",
]

SUPPORTED_BAZEL_VERSIONS = [
    CURRENT_BAZEL_VERSION,
] + OTHER_BAZEL_VERSIONS