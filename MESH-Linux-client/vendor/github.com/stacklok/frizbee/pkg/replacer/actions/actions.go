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

// Package actions provides utilities to work with GitHub Actions.
package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/google/go-github/v66/github"

	"github.com/stacklok/frizbee/pkg/interfaces"
	"github.com/stacklok/frizbee/pkg/replacer/image"
	"github.com/stacklok/frizbee/pkg/utils/config"
	"github.com/stacklok/frizbee/pkg/utils/store"
)

const (
	prefixUses   = "uses: "
	prefixDocker = "docker://"
	// GitHubActionsRegex is regular expression pattern to match GitHub Actions usage
	GitHubActionsRegex = `uses:\s*[^\s]+/[^\s]+@[^\s]+|uses:\s*docker://[^\s]+:[^\s]+`
	// ReferenceType is the type of the reference
	ReferenceType = "action"
)

var (
	// ErrInvalidAction is returned when parsing the action fails.
	ErrInvalidAction = errors.New("invalid action")
	// ErrInvalidActionReference is returned when parsing the action reference fails.
	ErrInvalidActionReference = errors.New("action reference is not a tag nor branch")
)

// Parser is a struct to replace action references with digests
type Parser struct {
	regex string
	cache store.RefCacher
}

// New creates a new Parser
func New() *Parser {
	return &Parser{
		regex: GitHubActionsRegex,
		cache: store.NewRefCacher(),
	}
}

// SetCache returns the regular expression pattern to match GitHub Actions usage
func (p *Parser) SetCache(cache store.RefCacher) {
	p.cache = cache
}

// SetRegex returns the regular expression pattern to match GitHub Actions usage
func (p *Parser) SetRegex(regex string) {
	p.regex = regex
}

// GetRegex returns the regular expression pattern to match GitHub Actions usage
func (p *Parser) GetRegex() string {
	return p.regex
}

// Replace replaces the action reference with the digest
func (p *Parser) Replace(
	ctx context.Context,
	matchedLine string,
	restIf interfaces.REST,
	cfg config.Config,
) (*interfaces.EntityRef, error) {
	var err error
	var actionRef *interfaces.EntityRef
	hasUsesPrefix := false

	// Trim the uses prefix
	if strings.HasPrefix(matchedLine, prefixUses) {
		matchedLine = strings.TrimPrefix(matchedLine, prefixUses)
		hasUsesPrefix = true
	}
	// Determine if the action reference has a docker prefix
	if strings.HasPrefix(matchedLine, prefixDocker) {
		actionRef, err = p.replaceDocker(ctx, matchedLine, restIf, cfg)
	} else {
		actionRef, err = p.replaceAction(ctx, matchedLine, restIf, cfg)
	}
	if err != nil {
		return nil, err
	}

	// Add back the uses prefix
	if hasUsesPrefix {
		actionRef.Prefix = fmt.Sprintf("%s%s", prefixUses, actionRef.Prefix)
	}

	// Return the new action reference
	return actionRef, nil
}

func (p *Parser) replaceAction(
	ctx context.Context,
	matchedLine string,
	restIf interfaces.REST,
	cfg config.Config,
) (*interfaces.EntityRef, error) {

	// If the value is a local path or should be excluded, skip it
	if isLocal(matchedLine) || shouldExclude(&cfg.GHActions, matchedLine) {
		return nil, fmt.Errorf("%w: %s", interfaces.ErrReferenceSkipped, matchedLine)
	}

	// Parse the action reference
	act, ref, err := ParseActionReference(matchedLine)
	if err != nil {
		return nil, fmt.Errorf("failed to parse action reference '%s': %w", matchedLine, err)
	}

	// Check if the parsed reference should be excluded
	if shouldExclude(&cfg.GHActions, act) {
		return nil, fmt.Errorf("%w: %s", interfaces.ErrReferenceSkipped, matchedLine)
	}
	var sum string

	// Check if we have a cache
	if p.cache != nil {
		// Check if we have a cached value
		if val, ok := p.cache.Load(matchedLine); ok {
			sum = val
		} else {
			// Get the checksum for the action reference
			sum, err = GetChecksum(ctx, cfg.GHActions, restIf, act, ref)
			if err != nil {
				return nil, fmt.Errorf("failed to get checksum for action '%s': %w", matchedLine, err)
			}
			// Store the checksum in the cache
			p.cache.Store(matchedLine, sum)
		}
	} else {
		// Get the checksum for the action reference
		sum, err = GetChecksum(ctx, cfg.GHActions, restIf, act, ref)
		if err != nil {
			return nil, fmt.Errorf("failed to get checksum for action '%s': %w", matchedLine, err)
		}
	}

	// Compare the digest with the reference and return the original reference if they already match
	if ref == sum {
		return nil, fmt.Errorf("image already referenced by digest: %s %w", matchedLine, interfaces.ErrReferenceSkipped)
	}

	return &interfaces.EntityRef{
		Name: act,
		Ref:  sum,
		Type: ReferenceType,
		Tag:  ref,
	}, nil
}

func (p *Parser) replaceDocker(
	ctx context.Context,
	matchedLine string,
	_ interfaces.REST,
	cfg config.Config,
) (*interfaces.EntityRef, error) {
	// Trim the docker prefix
	trimmedRef := strings.TrimPrefix(matchedLine, prefixDocker)

	// If the value is a local path or should be excluded, skip it
	if isLocal(trimmedRef) || shouldExclude(&cfg.GHActions, trimmedRef) {
		return nil, fmt.Errorf("%w: %s", interfaces.ErrReferenceSkipped, matchedLine)
	}

	// Get the digest of the docker:// image reference
	actionRef, err := image.GetImageDigestFromRef(ctx, trimmedRef, cfg.Platform, p.cache)
	if err != nil {
		return nil, err
	}

	// Check if the parsed reference should be excluded
	if shouldExclude(&cfg.GHActions, actionRef.Name) {
		return nil, fmt.Errorf("%w: %s", interfaces.ErrReferenceSkipped, matchedLine)
	}

	// Add back the docker prefix
	if strings.HasPrefix(matchedLine, prefixDocker) {
		actionRef.Prefix = fmt.Sprintf("%s%s", prefixDocker, actionRef.Prefix)
	}

	return actionRef, nil
}

// ConvertToEntityRef converts an action reference to an EntityRef
func (_ *Parser) ConvertToEntityRef(reference string) (*interfaces.EntityRef, error) {
	reference = strings.TrimPrefix(reference, prefixUses)
	refType := ReferenceType
	separator := "@"
	// Update the separator in case this is a docker reference with a digest
	if strings.Contains(reference, prefixDocker) {
		reference = strings.TrimPrefix(reference, prefixDocker)
		if !strings.Contains(reference, separator) && strings.Contains(reference, ":") {
			separator = ":"
		}
		refType = image.ReferenceType
	}
	frags := strings.Split(reference, separator)
	if len(frags) != 2 {
		return nil, fmt.Errorf("invalid action reference: %s", reference)
	}

	return &interfaces.EntityRef{
		Name: frags[0],
		Ref:  frags[1],
		Type: refType,
	}, nil
}

// isLocal returns true if the input is a local path.
func isLocal(input string) bool {
	return strings.HasPrefix(input, "./") || strings.HasPrefix(input, "../")
}

func shouldExclude(cfg *config.GHActions, input string) bool {
	for _, e := range cfg.Exclude {
		if e == input {
			return true
		}
	}
	return false
}

// ParseActionReference parses an action reference into action and reference.
func ParseActionReference(input string) (action string, reference string, err error) {
	frags := strings.Split(input, "@")
	if len(frags) != 2 {
		return "", "", fmt.Errorf("invalid action reference: %s", input)
	}

	return frags[0], frags[1], nil
}

// GetChecksum returns the checksum for a given action and tag.
func GetChecksum(ctx context.Context, cfg config.GHActions, restIf interfaces.REST, action, ref string) (string, error) {
	owner, repo, err := parseActionFragments(action)
	if err != nil {
		return "", err
	}

	// Check if we're using a checksum
	if isChecksum(ref) {
		return ref, nil
	}

	res, err := getCheckSumForTag(ctx, restIf, owner, repo, ref)
	if err != nil {
		return "", fmt.Errorf("failed to get checksum for tag: %w", err)
	} else if res != "" {
		return res, nil
	}

	// check branch
	if excludeBranch(cfg.Filter.ExcludeBranches, ref) {
		// if a branch is excluded, we won't know if it's a valid reference
		// but that's OK - we just won't touch that reference
		return "", fmt.Errorf("%w: %s", interfaces.ErrReferenceSkipped, ref)
	}

	res, err = getCheckSumForBranch(ctx, restIf, owner, repo, ref)
	if err != nil {
		return "", fmt.Errorf("failed to get checksum for branch: %w", err)
	} else if res != "" {
		return res, nil
	}

	return "", ErrInvalidActionReference
}

func parseActionFragments(action string) (owner string, repo string, err error) {
	frags := strings.Split(action, "/")

	// if we have more than 2 fragments, we're probably dealing with
	// sub-actions, so we take the first two fragments as the owner and repo
	if len(frags) < 2 {
		return "", "", fmt.Errorf("%w: '%s' reference is incorrect", ErrInvalidAction, action)
	}

	return frags[0], frags[1], nil
}

// isChecksum returns true if the input is a checksum.
func isChecksum(ref string) bool {
	return len(ref) == 40
}

func getCheckSumForTag(ctx context.Context, restIf interfaces.REST, owner, repo, tag string) (string, error) {
	path, err := url.JoinPath("repos", owner, repo, "git", "refs", "tags", tag)
	if err != nil {
		return "", fmt.Errorf("failed to join path: %w", err)
	}

	sha, otype, err := doGetReference(ctx, restIf, path)
	if err != nil {
		return "", err
	}

	if otype == "commit" {
		return sha, nil
	}

	// assume otype == "tag"
	path, err = url.JoinPath("repos", owner, repo, "git", "tags", sha)
	if err != nil {
		return "", fmt.Errorf("failed to join path: %w", err)
	}

	sha, _, err = doGetReference(ctx, restIf, path)
	return sha, err
}

func getCheckSumForBranch(ctx context.Context, restIf interfaces.REST, owner, repo, branch string) (string, error) {
	path, err := url.JoinPath("repos", owner, repo, "git", "refs", "heads", branch)
	if err != nil {
		return "", fmt.Errorf("failed to join path: %w", err)
	}

	sha, _, err := doGetReference(ctx, restIf, path)
	return sha, err
}

func excludeBranch(excludes []string, branch string) bool {
	if len(excludes) == 0 {
		return false
	}
	if slices.Contains(excludes, "*") {
		return true
	}

	return slices.Contains(excludes, branch)
}

func doGetReference(ctx context.Context, restIf interfaces.REST, path string) (string, string, error) {
	req, err := restIf.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", "", fmt.Errorf("cannot create REST request: %w", err)
	}

	resp, err := restIf.Do(ctx, req)

	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	if err != nil && resp.StatusCode != http.StatusNotFound {
		return "", "", fmt.Errorf("failed to do API request: %w", err)
	} else if resp.StatusCode == http.StatusNotFound {
		// No error, but no tag found
		return "", "", nil
	}

	var t github.Reference
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil && strings.Contains(err.Error(), "cannot unmarshal array into Go value of type") {
		// This is a branch, not a tag
		return "", "", nil
	} else if err != nil {
		return "", "", fmt.Errorf("canont decode response: %w", err)
	}

	return t.GetObject().GetSHA(), t.GetObject().GetType(), nil
}
