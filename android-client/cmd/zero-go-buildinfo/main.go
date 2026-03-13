// Command zero-go-buildinfo zeros out Go module build info in an ELF shared
// library.
//
// gomobile embeds a `replace` directive with the local filesystem path in the
// Go build info stored in .rodata. This path varies between build environments
// and breaks reproducible builds. This tool finds the serialised build info
// string (which starts with "path\tgobind") and overwrites it with null bytes,
// producing a deterministic binary regardless of the build directory.
//
// Usage: go run ./cmd/zero-go-buildinfo <path-to-libgojni.so>
package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <elf-file>\n", os.Args[0])
		os.Exit(1)
	}

	path := os.Args[1]

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", path, err)
		os.Exit(1)
	}

	marker := []byte("path\tgobind")
	idx := bytes.Index(data, marker)
	if idx == -1 {
		fmt.Fprintln(os.Stderr, "No Go build info found — nothing to do.")
		return
	}

	// Find the end of the build info block (terminated by \n\x00 or \x00).
	end := bytes.Index(data[idx:], []byte("\n\x00"))
	if end == -1 {
		end = bytes.IndexByte(data[idx:], 0x00)
		if end == -1 {
			fmt.Fprintln(os.Stderr, "Could not find end of build info block.")
			os.Exit(1)
		}
		end += idx
	} else {
		end += idx
	}
	end++ // include the newline

	length := end - idx
	fmt.Fprintf(os.Stderr, "Zeroing Go build info at offset 0x%x, length %d bytes\n", idx, length)

	// Zero out the build info block in place.
	for i := idx; i < end; i++ {
		data[i] = 0
	}

	if err := os.WriteFile(path, data, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Done: %s\n", path)
}
