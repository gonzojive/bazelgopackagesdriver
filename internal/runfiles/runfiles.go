// Package runfiles helps read runtime data files for programs built using
// Bazel.
//
// This library uses "github.com/bazelbuild/rules_go/go/tools/bazel" with some
// extra usability features.
package runfiles

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func mapSlice[T, R any](s []T, fn func(T) R) []R {
	var out []R
	for _, t := range s {
		out = append(out, fn(t))
	}
	return out
}

func filter[T any](s []T, fn func(T) bool) []T {
	var out []T
	for _, t := range s {
		if fn(t) {
			out = append(out, t)
		}
	}
	return out
}

func Runfile(filePath string) (string, error) {

	got, err := bazel.Runfile(filePath)
	if err == nil {
		return got, nil
	}
	if err != nil && strings.Contains(err.Error(), "could not locate file") {
		err = fmt.Errorf("%s: %w", err, os.ErrNotExist)
	}

	if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("unknown error locating runfile %q: %w", filePath, err)
	}
	base := path.Base(filePath)
	entries, err := bazel.ListRunfiles()
	if err != nil {
		return "", fmt.Errorf("could not locate %q: failed to list runfile entries: %w", filePath, err)
	}
	sameBaseEntries := filter(entries, func(entry bazel.RunfileEntry) bool {
		return path.Base(entry.ShortPath) == base
	})

	if len(sameBaseEntries) == 0 {
		return "", fmt.Errorf("could not locate %q among %d runfiles:\n  %s", filePath, len(entries),
			strings.Join(mapSlice(entries, func(e bazel.RunfileEntry) string { return e.ShortPath }), "\n  "))
	}
	return "", fmt.Errorf("could not locate %q among %d runfiles; suggested match(es):\n  %s", filePath, len(entries),
		strings.Join(mapSlice(sameBaseEntries, func(e bazel.RunfileEntry) string {
			return fmt.Sprintf("%q", e.ShortPath)
		}), "\n  "))

	// var sameBaseShortPaths []*bazel.RunfileEntry
	// for _, e := range entries {
	// 	if path.Base(e.ShortPath)
	// }
}
