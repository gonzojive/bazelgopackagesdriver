# Copyright 2018 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""Rules to be used in WORKSPACE files for building a cache of files.

To run bazel integration tests with no network access requires building a cache
directory.
"""

# TODO: have bazel_dependency_archive and bazel_git_dependency_archive rules
# to 1) register any Bazel target that produces an archive and 2) register git_repository
# dependencies. For bazel_dependency_archive, we would just need to compute the sha256 checksum.
# Implementing bazel_git_dependency_archive would mean calling "git archive" directly.

def _get_basename(url):
  return url[url.rindex("/") + 1:]

def _http_files_impl(repository_ctx):
  repository_ctx.file("cache_manifest.json", json.encode_indent(
      {
          "DownloadedFiles": [
              {
                  "ChecksumType": "sha256",
                  "Checksum": sha256,
                  "URLs": urls,
                  "RepoRelativePath": sha256,
              }
              for (sha256, urls) in repository_ctx.attr.srcs.items()
          ],
      },
      indent="  ",
  ))
  filegroup_name = "filegroup"  # was repository_ctx.name
  build_lines = [
      "filegroup(",
      "    name = \"%s\"," % filegroup_name,
      "    srcs = [",
      "      \"cache_manifest.json\",",
  ]

  for sha256 in repository_ctx.attr.srcs:
    urls = repository_ctx.attr.srcs[sha256]
    basename = sha256
    repository_ctx.download(
        url = urls,
        sha256 = sha256,
        output = basename,
    )
    build_lines.append("        \"%s\",  # %s" % (basename, urls[0]))

  build_lines += [
      "    ],",
      "    visibility = [\"//visibility:public\"],",
      ")",
  ]

  repository_ctx.file("BUILD", "\n".join(build_lines))

http_files = repository_rule(
    attrs = {
        "srcs": attr.string_list_dict(allow_empty = False),
    },
    implementation = _http_files_impl,
)
"""
Import external dependencies by their download URLs and sha256 checksums,
for Bazel 0.12.0 and above (below, it is a no-op).
```python
load(
    "@build_bazel_integration_testing//:bazel_integration_test.bzl", 
    "http_files",
)
http_files(
    name = "test_archive",
    srcs = {
        "90a8e1603eeca48e7e879f3afbc9560715322985f39a274f6f6070b43f9d06fe": [
            "http://repo1.maven.org/maven2/junit/junit/4.11/junit-4.11.jar",
            "http://maven.ibiblio.org/maven2/junit/junit/4.11/junit-4.11.jar",
        ],
    },
)
```
By specifying the filegroup `@test_archive` as, e.g., an `external_deps` of
`bazel_java_integration_test`, the file above (`junit-4.11.jar`) is made available without network
access and without changing the import logic.  For instance, the scratch WORKSPACE could contain
the following, without downloading actually taking place:
```
load("@bazel_tools//tools/build_defs/repo:java.bzl", "java_import_external")
java_import_external(
    name = "org_junit",
    licenses = ["restricted"],  # Eclipse Public License 1.0",
    jar_urls = [
        "http://repo1.maven.org/maven2/junit/junit/4.11/junit-4.11.jar",
    ],
    jar_sha256 = "90a8e1603eeca48e7e879f3afbc9560715322985f39a274f6f6070b43f9d06fe",
)
```
Note that external artifacts generated by the following rules cannot be imported using this means:
`maven_jar`, `git_repository`, `new_git_repository`.  This is due to their not relying on
sha256 checksums.  However, there is almost always a work around that does not use these
repository rules.  For instance, `maven_jar` can be replaced by `java_import_external` as above,
and if a Git repository is hosted on Github or another Git provider that allows for archiving, it is
probably better to use an http_archive rule:
# ```
# http_files(
#     name = "test_archive",
#     srcs = {
#         "25e8256d9e5b5e6a9e4196a9b17b1d6cf5faf2fdab247c00615a51250d8a7e10": [
#             "https://github.com/kubernetes/charts/archive/13f273aef2db18211ff0c5e6953461201be89c5d.tar.gz",
#         ],
#     },
# )
# ```
# and in the scratch WORKSPACE:
# ```
# load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
# http_archive(
#     name = "com_github_kubernetes_charts",
#     build_file_content = "...",
#     sha256 = "a3d6f8ee7db572fed090376b53561d76cdac41417b9373cbce2b8587543a8a0b",
#     strip_prefix = "charts-790457b26dcfdc4ad868ef64115ae0f10f0af0dd",
#     urls = ["https://github.com/kubernetes/charts/archive/790457b26dcfdc4ad868ef64115ae0f10f0af0dd.tar.gz"],
# )
# ```
"""
