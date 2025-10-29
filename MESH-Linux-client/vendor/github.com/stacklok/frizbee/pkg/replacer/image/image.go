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

// Package image provides utilities to work with container images.
package image

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	dockerparser "github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/stacklok/frizbee/internal/cli"
	"github.com/stacklok/frizbee/pkg/interfaces"
	"github.com/stacklok/frizbee/pkg/utils/config"
	"github.com/stacklok/frizbee/pkg/utils/store"
)

const (
	// ContainerImageRegex is regular expression pattern to match container image usage in YAML
	// nolint:lll
	ContainerImageRegex = `image\s*:\s*["']?([^\s"']+/[^\s"']+|[^\s"']+)(:[^\s"']+)?(@[^\s"']+)?["']?|FROM\s+(--platform=[^\s]+[^\s]*\s+)?([^\s]+(/[^\s]+)?(:[^\s]+)?(@[^\s]+)?)`
	prefixFROM          = "FROM "
	prefixImage         = "image: "
	// ReferenceType is the type of the reference
	ReferenceType = "container"
)

// Parser is a struct to replace container image references with digests
type Parser struct {
	regex string
	cache store.RefCacher
}

type unresolvedImage struct {
	imageRef string
	flags    []string
}

// New creates a new Parser
func New() *Parser {
	return &Parser{
		regex: ContainerImageRegex,
		cache: store.NewRefCacher(),
	}
}

// SetCache sets the cache to store the image references
func (p *Parser) SetCache(cache store.RefCacher) {
	p.cache = cache
}

// SetRegex sets the regular expression pattern to match container image usage
func (p *Parser) SetRegex(regex string) {
	p.regex = regex
}

// GetRegex returns the regular expression pattern to match container image usage
func (p *Parser) GetRegex() string {
	return p.regex
}

// Replace replaces the container image reference with the digest
func (p *Parser) Replace(
	ctx context.Context,
	matchedLine string,
	_ interfaces.REST,
	cfg config.Config,
) (*interfaces.EntityRef, error) {
	var imageRef string
	var extraArgs string

	// Trim the prefix
	hasFROMPrefix := false
	hasImagePrefix := false
	// Check if the image reference has the FROM prefix, i.e. Dockerfile
	if strings.HasPrefix(matchedLine, prefixFROM) {
		parsedFrom, err := getRefFromDockerfileFROM(matchedLine)
		if err != nil {
			return nil, err
		}

		// Check if the image reference should be excluded, i.e. scratch
		if shouldSkipImageRef(&cfg, parsedFrom.imageRef) {
			return nil, fmt.Errorf("image reference %s should be excluded - %w", matchedLine, interfaces.ErrReferenceSkipped)
		}

		imageRef = parsedFrom.imageRef
		extraArgs = strings.Join(parsedFrom.flags, " ")
		if extraArgs != "" {
			extraArgs += " "
		}

		hasFROMPrefix = true
	} else if strings.HasPrefix(matchedLine, prefixImage) {
		// Check if the image reference has the image prefix, i.e. Kubernetes or Docker Compose YAML
		imageRef = strings.TrimPrefix(matchedLine, prefixImage)
		// Check if the image reference should be excluded, i.e. scratch
		if shouldSkipImageRef(&cfg, imageRef) {
			return nil, fmt.Errorf("image reference %s should be excluded - %w", matchedLine, interfaces.ErrReferenceSkipped)
		}
		hasImagePrefix = true
	} else {
		imageRef = matchedLine
	}

	// Get the digest of the image reference
	imageRefWithDigest, err := GetImageDigestFromRef(ctx, imageRef, cfg.Platform, p.cache)
	if err != nil {
		return nil, err
	}

	// Add the prefix back
	if hasFROMPrefix {
		imageRefWithDigest.Prefix = fmt.Sprintf("%s%s%s", prefixFROM, extraArgs, imageRefWithDigest.Prefix)
	} else if hasImagePrefix {
		imageRefWithDigest.Prefix = fmt.Sprintf("%s%s", prefixImage, imageRefWithDigest.Prefix)
	}

	// Return the reference
	return imageRefWithDigest, nil
}

// ConvertToEntityRef converts a container image reference to an EntityRef
func (_ *Parser) ConvertToEntityRef(reference string) (*interfaces.EntityRef, error) {
	reference = strings.TrimPrefix(reference, prefixImage)
	reference = strings.TrimPrefix(reference, prefixFROM)
	var sep string
	var frags []string
	if strings.Contains(reference, "@") {
		sep = "@"
	} else if strings.Contains(reference, ":") {
		sep = ":"
	}

	if sep != "" {
		frags = strings.Split(reference, sep)
		if len(frags) != 2 {
			return nil, fmt.Errorf("invalid container reference: %s", reference)
		}
	} else {
		frags = []string{reference, "latest"}
	}

	return &interfaces.EntityRef{
		Name: frags[0],
		Ref:  frags[1],
		Type: ReferenceType,
	}, nil
}

// GetImageDigestFromRef returns the digest of a container image reference
// from a name.Reference.
func GetImageDigestFromRef(ctx context.Context, imageRef, platform string, cache store.RefCacher) (*interfaces.EntityRef, error) {
	// Parse the image reference
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, err
	}
	opts := []remote.Option{
		remote.WithContext(ctx),
		remote.WithUserAgent(cli.UserAgent),
		remote.WithAuthFromKeychain(authn.DefaultKeychain),
	}

	// Set the platform if provided
	if platform != "" {
		platformSplit := strings.Split(platform, "/")
		if len(platformSplit) != 2 {
			return nil, errors.New("platform must be in the format os/arch")
		}
		opts = append(opts, remote.WithPlatform(v1.Platform{
			OS:           platformSplit[0],
			Architecture: platformSplit[1],
		}))
	}

	// Get the digest of the image reference
	var digest string

	if cache != nil {
		if d, ok := cache.Load(imageRef); ok {
			digest = d
		} else {
			desc, err := remote.Get(ref, opts...)
			if err != nil {
				return nil, err
			}
			digest = desc.Digest.String()
			cache.Store(imageRef, digest)
		}
	} else {
		desc, err := remote.Get(ref, opts...)
		if err != nil {
			return nil, err
		}
		digest = desc.Digest.String()
	}

	// Compare the digest with the reference and return the original reference if they already match
	if digest == ref.Identifier() {
		return nil, fmt.Errorf("image already referenced by digest: %s %w", imageRef, interfaces.ErrReferenceSkipped)
	}

	return &interfaces.EntityRef{
		Name: ref.Context().Name(),
		Ref:  digest,
		Type: ReferenceType,
		Tag:  ref.Identifier(),
	}, nil
}

func shouldSkipImageRef(cfg *config.Config, ref string) bool {
	// Parse the image reference
	nameRef, err := name.ParseReference(ref)
	if err != nil {
		// we wouldn't know how to resolve this reference, so let's skip
		return true
	}

	imageName := getImageNameFromRef(nameRef)
	if slices.Contains(cfg.Images.ImageFilter.ExcludeImages, imageName) {
		return true
	}

	tag := nameRef.Identifier()
	return slices.Contains(cfg.Images.ImageFilter.ExcludeTags, tag)
}

// TODO(jakub): this is a bit of a hack, but I didn't find a better way to get just the name
func getImageNameFromRef(nameRef name.Reference) string {
	fullRepositoryName := nameRef.Context().Name()
	parts := strings.Split(fullRepositoryName, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}

	return ""
}
func getRefFromDockerfileFROM(line string) (unresolvedImage, error) {
	parseResult, err := dockerparser.Parse(strings.NewReader(line))
	if err != nil {
		return unresolvedImage{}, fmt.Errorf("failed to parse Dockerfile line: %w", err)
	}

	if len(parseResult.AST.Children) == 0 ||
		parseResult.AST.Children[0] == nil ||
		strings.ToUpper(parseResult.AST.Children[0].Value) != "FROM" {
		return unresolvedImage{}, errors.New("invalid Dockerfile line: the first parsed node is not FROM")
	}

	fromNode := parseResult.AST.Children[0]

	imgNode := parseResult.AST.Children[0].Next
	if imgNode == nil {
		return unresolvedImage{}, errors.New("invalid Dockerfile line: no image node found")
	}

	return unresolvedImage{
		imageRef: imgNode.Value,
		flags:    fromNode.Flags,
	}, nil
}
