// Command normalize-elf normalizes an ELF shared library for reproducible
// builds by zeroing non-deterministic content.
//
// It performs two categories of normalization:
//
//  1. Go build info: gomobile embeds a `replace` directive with the local
//     filesystem path in the Go build info stored in .rodata. This path
//     varies between build environments and breaks reproducible builds.
//     The tool finds the serialised build info strings (starting with
//     "path\tgobind" or "mod\tgobind") and overwrites them with null bytes.
//
//  2. Linker-generated sections: the NDK linker produces subtly different
//     output depending on the host environment. The affected sections
//     (.eh_frame, .eh_frame_hdr, .relro_padding) reside inside LOAD
//     segments, so llvm-objcopy --remove-section silently fails to remove
//     them. This tool zeros their content in-place, preserving the binary
//     layout while making the content deterministic. It also zeros the ELF
//     entry point, which differs between environments for shared libraries.
//
// Usage: go run ./cmd/normalize-elf <path-to-libgojni.so>
package main

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"os"
)

// sectionsToZero lists ELF sections whose content should be zeroed.
// These sections differ across build environments but are not required
// at runtime for a shared library on Android.
var sectionsToZero = []string{
	".eh_frame",      // DWARF call frame information — differs by compiler/linker
	".eh_frame_hdr",  // Index into .eh_frame
	".relro_padding", // RELRO alignment padding — presence varies by linker
}

// buildInfoMarkers are byte patterns that identify Go build info blocks
// in .rodata. These contain replace directives with absolute filesystem
// paths that differ between build environments.
var buildInfoMarkers = [][]byte{
	[]byte("path\tgobind"),
	[]byte("mod\tgobind"),
}

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

	modified := zeroBuildInfo(data)

	f, err := elf.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ELF %s: %v\n", path, err)
		os.Exit(1)
	}
	defer f.Close()

	modified = zeroSections(data, f) || modified
	modified = zeroEntryPoint(data, f.Class) || modified

	if !modified {
		fmt.Fprintln(os.Stderr, "No modifications needed.")
		return
	}

	if err := os.WriteFile(path, data, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Done: %s\n", path)
}

// zeroBuildInfo finds and zeros Go build info blocks in the raw binary data.
// Returns true if any modifications were made.
func zeroBuildInfo(data []byte) bool {
	modified := false

	for _, marker := range buildInfoMarkers {
		for {
			idx := bytes.Index(data, marker)
			if idx == -1 {
				break
			}

			// Find the end of the build info block (terminated by \n\x00 or \x00).
			end := bytes.Index(data[idx:], []byte("\n\x00"))
			if end == -1 {
				end = bytes.IndexByte(data[idx:], 0x00)
				if end == -1 {
					fmt.Fprintf(os.Stderr, "Could not find end of build info block at 0x%x.\n", idx)
					os.Exit(1)
				}
				end += idx
			} else {
				end += idx
			}
			end++ // include the newline

			length := end - idx
			fmt.Fprintf(os.Stderr, "Zeroing Go build info at offset 0x%x, length %d bytes (marker: %q)\n", idx, length, marker)

			for i := idx; i < end; i++ {
				data[i] = 0
			}
			modified = true
		}
	}

	return modified
}

// zeroSections zeros the content of non-deterministic ELF sections.
// Returns true if any modifications were made.
func zeroSections(data []byte, f *elf.File) bool {
	modified := false

	for _, name := range sectionsToZero {
		sec := f.Section(name)
		if sec == nil {
			continue
		}

		if sec.Type == elf.SHT_NOBITS {
			fmt.Fprintf(os.Stderr, "Skipping %s (NOBITS, no file content)\n", name)
			continue
		}

		offset := sec.Offset
		size := sec.FileSize
		if offset+size > uint64(len(data)) {
			fmt.Fprintf(os.Stderr, "Warning: section %s (offset=0x%x, size=0x%x) extends beyond file — skipping\n",
				name, offset, size)
			continue
		}

		allZero := true
		for i := uint64(0); i < size; i++ {
			if data[offset+i] != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			fmt.Fprintf(os.Stderr, "Section %s already zeroed (0x%x bytes)\n", name, size)
			continue
		}

		fmt.Fprintf(os.Stderr, "Zeroing section %s at offset 0x%x, size 0x%x (%d bytes)\n",
			name, offset, size, size)
		for i := uint64(0); i < size; i++ {
			data[offset+i] = 0
		}
		modified = true
	}

	return modified
}

// zeroEntryPoint zeros the e_entry field in the ELF header.
// Returns true if the field was modified.
func zeroEntryPoint(data []byte, class elf.Class) bool {
	const entryOffset = 0x18

	switch class {
	case elf.ELFCLASS64:
		if len(data) < entryOffset+8 {
			return false
		}
		current := binary.LittleEndian.Uint64(data[entryOffset : entryOffset+8])
		if current == 0 {
			fmt.Fprintln(os.Stderr, "Entry point already zero")
			return false
		}
		fmt.Fprintf(os.Stderr, "Zeroing entry point (was 0x%x)\n", current)
		binary.LittleEndian.PutUint64(data[entryOffset:entryOffset+8], 0)
		return true

	case elf.ELFCLASS32:
		if len(data) < entryOffset+4 {
			return false
		}
		current := binary.LittleEndian.Uint32(data[entryOffset : entryOffset+4])
		if current == 0 {
			fmt.Fprintln(os.Stderr, "Entry point already zero")
			return false
		}
		fmt.Fprintf(os.Stderr, "Zeroing entry point (was 0x%x)\n", current)
		binary.LittleEndian.PutUint32(data[entryOffset:entryOffset+4], 0)
		return true

	default:
		fmt.Fprintf(os.Stderr, "Warning: unknown ELF class %v — skipping entry point\n", class)
		return false
	}
}
