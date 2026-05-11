// Copyright (c) BARGHEST
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tsDir := filepath.Join(root, "tailscale")

	info, err := os.Stat(tsDir)
	if err != nil || !info.IsDir() {
		log.Fatalf("tailscale submodule missing: %s", tsDir)
	}

	log.Println("tailscale submodule present")
	log.Println("No patching required")
}

// replaceOutsideComments applies string replacements to Go source code,
// skipping comments so that documentation and notes are preserved unchanged.
func replaceOutsideComments(src []byte, replacements map[string]string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Collect comment byte ranges, sorted by position.
	type span struct{ start, end int }
	var comments []span
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			comments = append(comments, span{
				start: fset.Position(c.Pos()).Offset,
				end:   fset.Position(c.End()).Offset,
			})
		}
	}
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].start < comments[j].start
	})

	// replaceSegment applies all replacements to a code segment.
	replaceSegment := func(segment string) string {
		for old, repl := range replacements {
			segment = strings.ReplaceAll(segment, old, repl)
		}
		return segment
	}

	// Build result: apply replacements only to non-comment segments.
	var result []byte
	pos := 0
	for _, c := range comments {
		result = append(result, replaceSegment(string(src[pos:c.start]))...)
		result = append(result, src[c.start:c.end]...)
		pos = c.end
	}
	result = append(result, replaceSegment(string(src[pos:]))...)

	return result, nil
}

// extractModuleZip extracts a Go module zip to destDir, stripping the
// "module@version/" prefix that all entries share per the module zip spec.
func extractModuleZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("opening zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		// Strip the "module@version/" prefix from each entry.
		_, relPath, ok := strings.Cut(f.Name, "/")
		if !ok || relPath == "" {
			continue
		}
		destPath := filepath.Join(destDir, relPath)

		// Verify the path doesn't escape the destination (zip slip protection).
		if !strings.HasPrefix(filepath.Clean(destPath), filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("zip entry %q would escape destination directory", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", destPath, err)
			}
			continue
		}

		if err := extractZipFile(f, destPath); err != nil {
			return err
		}
	}

	return nil
}

// extractZipFile extracts a single file from a zip archive to destPath,
// creating parent directories as needed.
func extractZipFile(f *zip.File, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("creating parent directory for %s: %w", destPath, err)
	}

	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening zip entry %s: %w", f.Name, err)
	}
	defer rc.Close()

	//nolint:gosec // G306 -- non-sensitive source file, must stay 0644 to match upstream tree
	outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("creating %s: %w", destPath, err)
	}
	defer outFile.Close()

	// Limit extraction size to guard against decompression bombs (gosec G110).
	// 100 MB is well above any single file we expect from the Tailscale module zip.
	const maxFileSize = 100 << 20 // 100 MB
	if _, err := io.Copy(outFile, io.LimitReader(rc, maxFileSize)); err != nil {
		return fmt.Errorf("writing %s: %w", destPath, err)
	}

	return nil
}

// rewriteHelpStrings applies substring replacements to the *contents* of
// string literals assigned to ShortHelp, LongHelp, or ShortUsage fields in
// Go source files under dir. The replacement is scoped to those three
// field names so that import paths, code identifiers, and unrelated strings
// (e.g. "tailscale.com/...", "tailscaled") are not touched.
func rewriteHelpStrings(dir string, replacements map[string]string) error {
	helpFields := map[string]bool{
		"ShortHelp":  true,
		"LongHelp":   true,
		"ShortUsage": true,
	}

	type edit struct {
		start, end int
		repl       string
	}

	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		var edits []edit
		ast.Inspect(f, func(n ast.Node) bool {
			kv, ok := n.(*ast.KeyValueExpr)
			if !ok {
				return true
			}
			id, ok := kv.Key.(*ast.Ident)
			if !ok || !helpFields[id.Name] {
				return true
			}
			ast.Inspect(kv.Value, func(x ast.Node) bool {
				lit, ok := x.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return true
				}
				start := fset.Position(lit.Pos()).Offset
				end := fset.Position(lit.End()).Offset
				original := string(src[start:end])
				if len(original) < 2 {
					return true
				}
				quote := original[0]
				inner := original[1 : len(original)-1]
				replaced := inner
				for old, new_ := range replacements {
					replaced = strings.ReplaceAll(replaced, old, new_)
				}
				if replaced != inner {
					edits = append(edits, edit{start, end, string(quote) + replaced + string(quote)})
				}
				return true
			})
			return true
		})

		if len(edits) == 0 {
			return nil
		}
		sort.Slice(edits, func(i, j int) bool { return edits[i].start < edits[j].start })
		var buf bytes.Buffer
		prev := 0
		for _, e := range edits {
			buf.Write(src[prev:e.start])
			buf.WriteString(e.repl)
			prev = e.end
		}
		buf.Write(src[prev:])
		//nolint:gosec // G306 -- non-sensitive source file, must stay 0644 to match upstream tree
		if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		log.Printf("Rebranded help strings in %s", path)
		return nil
	})
}
