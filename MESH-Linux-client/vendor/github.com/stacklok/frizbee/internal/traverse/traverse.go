//
// Copyright 2024 Stacklok, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package traverse provides utilities to traverse directories.
package traverse

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
)

// GhwFunc is a function that gets called with each file in a GitHub Actions workflow
// directory. It receives the path to the file.
type GhwFunc func(path string) error

// FuncTraverse is a function that gets called with each file in a directory.
type FuncTraverse func(path string, info fs.FileInfo) error

// YamlDockerfiles traverses all yaml/yml in the given directory
// and calls the given function with each workflow.
func YamlDockerfiles(bfs billy.Filesystem, base string, fun GhwFunc) error {
	return Traverse(bfs, base, func(path string, info fs.FileInfo) error {
		if !isYAMLOrDockerfile(info) {
			return nil
		}

		if err := fun(path); err != nil {
			return fmt.Errorf("failed to process file %s: %w", path, err)
		}

		return nil
	})
}

// Traverse traverses the given directory and calls the given function with each file.
func Traverse(bfs billy.Filesystem, base string, fun FuncTraverse) error {
	return Walk(bfs, base, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		return fun(path, info)
	})
}

// isYAMLOrDockerfile returns true if the given file is a YAML or Dockerfile.
func isYAMLOrDockerfile(info fs.FileInfo) bool {
	// Skip if not a file
	if info.IsDir() {
		return false
	}

	// Filter out files that are not yml, yaml or dockerfiles
	if strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") ||
		strings.Contains(strings.ToLower(info.Name()), "dockerfile") {
		return true
	}

	return false
}

// walk recursively descends path, calling walkFn
// adapted from https://golang.org/src/path/filepath/path.go
func walk(bfs billy.Filesystem, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}

	names, err := readDirNames(bfs, path)
	err1 := walkFn(path, info, err)
	// If err != nil, walk can't walk into this directory.
	// err1 != nil means walkFn want walk to skip this directory or stop walking.
	// Therefore, if one of err and err1 isn't nil, walk will return.
	if err != nil || err1 != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err1
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := bfs.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walk(bfs, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

// Walk walks the file tree rooted at root, calling fn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by fn: see the WalkFunc documentation for
// details.
//
// The files are walked in lexical order, which makes the output deterministic
// but requires Walk to read an entire directory into memory before proceeding
// to walk that directory. Walk does not follow symbolic links.
//
// Function adapted from https://github.com/golang/go/blob/3b770f2ccb1fa6fecc22ea822a19447b10b70c5c/src/path/filepath/path.go#L500
func Walk(bfs billy.Filesystem, root string, walkFn filepath.WalkFunc) error {
	info, err := bfs.Lstat(root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = walk(bfs, root, info, walkFn)
	}

	if err == filepath.SkipDir {
		return nil
	}

	return err
}

func readDirNames(bfs billy.Filesystem, dir string) ([]string, error) {
	files, err := bfs.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, file := range files {
		names = append(names, file.Name())
	}

	return names, nil
}
