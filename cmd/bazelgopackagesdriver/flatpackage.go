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
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/gonzojive/bazelgopackagesdriver/protocol"
)

type ResolvePkgFunc func(importPath string) *FlatPackage

// Copy and pasted from golang.org/x/tools/go/packages
type FlatPackagesError struct {
	Pos  string // "file:line:col" or "file:line" or "" or "-"
	Msg  string
	Kind FlatPackagesErrorKind
}

type FlatPackagesErrorKind = protocol.FlatPackagesError

func (err FlatPackagesError) Error() string {
	pos := err.Pos
	if pos == "" {
		pos = "-" // like token.Position{}.String()
	}
	return pos + ": " + err.Msg
}

// FlatPackage is the JSON form of Package
// It drops all the type and syntax fields, and transforms the Imports
type FlatPackage = protocol.FlatPackage

type (
	PackageFunc      func(pkg *FlatPackage)
	PathResolverFunc func(path string) string
)

func resolvePathsInPlace(prf PathResolverFunc, paths []string) {
	for i, path := range paths {
		paths[i] = prf(path)
	}
}

func WalkFlatPackagesFromJSON(jsonFile string, onPkg PackageFunc) error {
	f, err := os.Open(jsonFile)
	if err != nil {
		return fmt.Errorf("unable to open package JSON file: %w", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	for decoder.More() {
		pkg := &FlatPackage{}
		if err := decoder.Decode(&pkg); err != nil {
			return fmt.Errorf("unable to decode package in %s: %w", f.Name(), err)
		}
		onPkg(pkg)
	}
	return nil
}

func ResolvePaths(fp *FlatPackage, prf PathResolverFunc) error {
	resolvePathsInPlace(prf, fp.CompiledGoFiles)
	resolvePathsInPlace(prf, fp.GoFiles)
	resolvePathsInPlace(prf, fp.OtherFiles)
	fp.ExportFile = prf(fp.ExportFile)
	return nil
}

// FilterFilesForBuildTags filters the source files given the current build
// tags.
func FilterFilesForBuildTags(fp *FlatPackage) {
	fp.GoFiles = filterSourceFilesForTags(fp.GoFiles)
	fp.CompiledGoFiles = filterSourceFilesForTags(fp.CompiledGoFiles)
}

func IsStdlib(fp *FlatPackage) bool {
	return fp.Standard
}

func ResolveImports(fp *FlatPackage, resolve ResolvePkgFunc) {
	// Stdlib packages are already complete import wise
	if fp.IsStdlib() {
		return
	}

	fset := token.NewFileSet()

	for _, file := range fp.CompiledGoFiles {
		f, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
		if err != nil {
			continue
		}
		// If the name is not provided, fetch it from the sources
		if fp.Name == "" {
			fp.Name = f.Name.Name
		}
		for _, rawImport := range f.Imports {
			imp, err := strconv.Unquote(rawImport.Path.Value)
			if err != nil {
				continue
			}
			// We don't handle CGo for now
			if imp == "C" {
				continue
			}
			if _, ok := fp.Imports[imp]; ok {
				continue
			}
			if pkg := resolve(imp); pkg != nil {
				if fp.Imports == nil {
					fp.Imports = map[string]string{}
				}
				fp.Imports[imp] = pkg.ID
			}
		}
	}
}

func IsRoot(fp *FlatPackage) bool {
	return strings.HasPrefix(fp.ID, "//")
}
