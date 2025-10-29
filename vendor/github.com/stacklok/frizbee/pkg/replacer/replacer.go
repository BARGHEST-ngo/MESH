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

// Package replacer provide common replacer implementation
package replacer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"golang.org/x/sync/errgroup"

	"github.com/stacklok/frizbee/internal/traverse"
	"github.com/stacklok/frizbee/pkg/interfaces"
	"github.com/stacklok/frizbee/pkg/replacer/actions"
	"github.com/stacklok/frizbee/pkg/replacer/image"
	"github.com/stacklok/frizbee/pkg/utils/config"
	"github.com/stacklok/frizbee/pkg/utils/ghrest"
)

// ReplaceResult holds a slice of all processed files along with a map of their modified content
type ReplaceResult struct {
	Processed []string
	Modified  map[string]string
}

// ListResult holds the result of the list methods
type ListResult struct {
	Processed []string
	Entities  []interfaces.EntityRef
}

// Replacer is an object with methods to replace references with digests
type Replacer struct {
	parser interfaces.Parser
	rest   interfaces.REST
	cfg    config.Config
}

// NewGitHubActionsReplacer creates a new replacer for GitHub actions
func NewGitHubActionsReplacer(cfg *config.Config) *Replacer {
	cfg = config.MergeUserConfig(cfg)

	return &Replacer{
		cfg:    *cfg,
		parser: actions.New(),
		rest:   ghrest.NewClient(""),
	}
}

// NewContainerImagesReplacer creates a new replacer for container images
func NewContainerImagesReplacer(cfg *config.Config) *Replacer {
	cfg = config.MergeUserConfig(cfg)

	return &Replacer{
		cfg:    *cfg,
		parser: image.New(),
		rest:   ghrest.NewClient(""),
	}
}

// WithGitHubClientFromToken creates an authenticated GitHub client from a token
func (r *Replacer) WithGitHubClientFromToken(token string) *Replacer {
	client := ghrest.NewClient(token)
	r.rest = client
	return r
}

// WithGitHubClient sets the GitHub client to use
func (r *Replacer) WithGitHubClient(client interfaces.REST) *Replacer {
	r.rest = client
	return r
}

// WithUserRegex sets a user-provided regex for the parser
func (r *Replacer) WithUserRegex(regex string) *Replacer {
	if r.parser != nil && regex != "" {
		r.parser.SetRegex(regex)
	}
	return r
}

// WithCacheDisabled disables caching
func (r *Replacer) WithCacheDisabled() *Replacer {
	r.parser.SetCache(nil)
	return r
}

// ParseString parses and returns the referenced entity pinned by its digest
func (r *Replacer) ParseString(ctx context.Context, entityRef string) (*interfaces.EntityRef, error) {
	return r.parser.Replace(ctx, entityRef, r.rest, r.cfg)
}

// ParsePath parses and replaces all entity references in the provided directory
func (r *Replacer) ParsePath(ctx context.Context, dir string) (*ReplaceResult, error) {
	return parsePathInFS(ctx, r.parser, r.rest, r.cfg, osfs.New(filepath.Dir(dir), osfs.WithBoundOS()), filepath.Base(dir))
}

// ParsePathInFS parses and replaces all entity references in the provided file system
func (r *Replacer) ParsePathInFS(ctx context.Context, bfs billy.Filesystem, base string) (*ReplaceResult, error) {
	return parsePathInFS(ctx, r.parser, r.rest, r.cfg, bfs, base)
}

// ParseFile parses and replaces all entity references in the provided file
func (r *Replacer) ParseFile(ctx context.Context, f io.Reader) (bool, string, error) {
	return parseAndReplaceReferencesInFile(ctx, f, r.parser, r.rest, r.cfg)
}

// ListPath lists all entity references in the provided directory
func (r *Replacer) ListPath(dir string) (*ListResult, error) {
	return listReferencesInFS(r.parser, osfs.New(filepath.Dir(dir), osfs.WithBoundOS()), filepath.Base(dir))
}

// ListPathInFS lists all entity references in the provided file system
func (r *Replacer) ListPathInFS(bfs billy.Filesystem, base string) (*ListResult, error) {
	return listReferencesInFS(r.parser, bfs, base)
}

// ListInFile lists all entities in the provided file
func (r *Replacer) ListInFile(f io.Reader) (*ListResult, error) {
	found, err := listReferencesInFile(f, r.parser)
	if err != nil {
		return nil, err
	}
	res := &ListResult{}
	res.Entities = found.ToSlice()

	// Sort the slice
	sort.Slice(res.Entities, func(i, j int) bool {
		return res.Entities[i].Name < res.Entities[j].Name
	})

	// All good
	return res, nil
}

func parsePathInFS(
	ctx context.Context,
	parser interfaces.Parser,
	rest interfaces.REST,
	cfg config.Config,
	bfs billy.Filesystem,
	base string,
) (*ReplaceResult, error) {
	var eg errgroup.Group
	var mu sync.Mutex

	res := ReplaceResult{
		Processed: make([]string, 0),
		Modified:  make(map[string]string),
	}

	// Traverse all YAML/YML files in dir
	err := traverse.YamlDockerfiles(bfs, base, func(path string) error {
		eg.Go(func() error {
			file, err := bfs.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			// nolint:errcheck // ignore error
			defer file.Close()

			// Parse the content of the file and update the matching references
			modified, updatedFile, err := parseAndReplaceReferencesInFile(ctx, file, parser, rest, cfg)
			if err != nil {
				return fmt.Errorf("failed to modify references in %s: %w", path, err)
			}

			mu.Lock()
			// Store the file name to the processed batch
			res.Processed = append(res.Processed, path)
			// Store the updated file content if it was modified
			if modified {
				res.Modified[path] = updatedFile
			}
			mu.Unlock()

			// All good
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// All good
	return &res, nil
}

func listReferencesInFS(parser interfaces.Parser, bfs billy.Filesystem, base string) (*ListResult, error) {
	var eg errgroup.Group
	var mu sync.Mutex

	res := ListResult{
		Processed: make([]string, 0),
		Entities:  make([]interfaces.EntityRef, 0),
	}

	found := mapset.NewSet[interfaces.EntityRef]()

	// Traverse all related files
	err := traverse.YamlDockerfiles(bfs, base, func(path string) error {
		eg.Go(func() error {
			file, err := bfs.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer file.Close() // nolint:errcheck

			// Parse the content of the file and list the matching references
			foundRefs, err := listReferencesInFile(file, parser)
			if err != nil {
				return fmt.Errorf("failed to list references in %s: %w", path, err)
			}

			// Store the file name to the processed batch
			mu.Lock()
			res.Processed = append(res.Processed, path)
			found = found.Union(foundRefs)
			mu.Unlock()

			// All good
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	res.Entities = found.ToSlice()

	// Sort the slice
	sort.Slice(res.Entities, func(i, j int) bool {
		return res.Entities[i].Name < res.Entities[j].Name
	})

	// All good
	return &res, nil
}

func parseAndReplaceReferencesInFile(
	ctx context.Context,
	f io.Reader,
	parser interfaces.Parser,
	rest interfaces.REST,
	cfg config.Config,
) (bool, string, error) {
	var contentBuilder strings.Builder
	var ret *interfaces.EntityRef

	modified := false

	// Compile the regular expression
	re, err := regexp.Compile(parser.GetRegex())
	if err != nil {
		return false, "", err
	}

	// Read the file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip commented lines
		if strings.HasPrefix(strings.TrimLeft(line, " \t\n\r"), "#") {
			// Write the line to the content builder buffer
			contentBuilder.WriteString(line + "\n")
			continue
		}

		// See if we can match an entity reference in the line
		newLine := re.ReplaceAllStringFunc(line, func(matchedLine string) string {
			// Modify the reference in the line
			ret, err = parser.Replace(ctx, matchedLine, rest, cfg)
			if err != nil {
				// Return the original line as we don't want to update it in case something errored out
				return matchedLine
			}
			// Construct the new line, comments in dockerfiles are handled differently than yml files
			if strings.Contains(matchedLine, "FROM") {
				return fmt.Sprintf("%s%s:%s@%s", ret.Prefix, ret.Name, ret.Tag, ret.Ref)
			}
			return fmt.Sprintf("%s%s@%s # %s", ret.Prefix, ret.Name, ret.Ref, ret.Tag)
		})

		// Check if the line was modified and set the modified flag to true if it was
		if newLine != line {
			modified = true
		}

		// Write the line to the content builder buffer
		contentBuilder.WriteString(newLine + "\n")
	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		return false, "", err
	}

	// Return the workflow content
	return modified, contentBuilder.String(), nil
}

// listReferencesInFile takes the given file reader and returns a map of all references, action or images it finds
func listReferencesInFile(
	f io.Reader,
	parser interfaces.Parser,
) (mapset.Set[interfaces.EntityRef], error) {
	found := mapset.NewSet[interfaces.EntityRef]()

	// Compile the regular expression
	re, err := regexp.Compile(parser.GetRegex())
	if err != nil {
		return nil, err
	}

	// Read the file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip commented lines
		if strings.HasPrefix(strings.TrimLeft(line, " \t\n\r"), "#") {
			continue
		}

		// See if we can match an entity reference in the line
		foundEntries := re.FindAllString(line, -1)
		// nolint:gosimple
		if foundEntries != nil {
			for _, entry := range foundEntries {
				e, err := parser.ConvertToEntityRef(entry)
				if err != nil {
					continue
				}
				found.Add(*e)
			}
		}
	}

	// Check for errors during the scan
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Return the found references
	return found, nil
}
