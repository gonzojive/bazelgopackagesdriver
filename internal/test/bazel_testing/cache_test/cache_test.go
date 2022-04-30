package cache_test

import (
	"testing"

	"github.com/gonzojive/bazelgopackagesdriver/internal/runfiles"
	"github.com/gonzojive/bazelgopackagesdriver/internal/test/bazel_testing"
	"github.com/google/go-cmp/cmp"
)

func TestManifest(t *testing.T) {
	for _, tt := range []struct {
		manifest string
		want     *bazel_testing.CacheManifest
	}{
		{
			"external/my_filez/cache_manifest.json",
			&bazel_testing.CacheManifest{
				DownloadedFiles: []*bazel_testing.DownloadedFile{
					{
						ChecksumType: "sha256",
						Checksum:     "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715c",
						URLs: []string{
							"https://github.com/bazelbuild/bazel/releases/download/5.1.1/bazel-5.1.1-linux-x86_64",
						},
						RepoRelativePath: "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715c",
					},
				},
			},
		},
	} {
		t.Run(tt.manifest, func(t *testing.T) {
			manifestPath, err := runfiles.Runfile(tt.manifest)
			fatalIfErr(t, err, "must get manifest path")
			got, err := bazel_testing.LoadCacheManifest(manifestPath)
			fatalIfErr(t, err, "invalid manifest data")
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("unexpected manifest diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func fatalIfErr(t *testing.T, err error, message string) {
	if err == nil {
		return
	}
	t.Fatalf("%s: %v", message, err)
}
