// Command normalize-elf normalizes an ELF shared library for reproducible
// builds by zeroing non-deterministic content.
//
// It performs two categories of normalization:
//
//  1. Go build info: gomobile embeds a `replace` directive with the local
//     filesystem path in the Go build info stored in .rodata. This path
//     varies between build environments and breaks reproducible builds.
//     The tool finds the serialised build info strings (starting with
//     "path\tgobind" or "mod\tgobind") within the .rodata section and
//     overwrites them with null bytes.  The search is constrained to
//     .rodata to avoid corrupting other sections (e.g. .gopclntab, .text)
//     that could coincidentally contain the same byte pattern.
//
//  2. Linker-generated sections: the NDK linker produces subtly different
//     .relro_padding depending on the host environment.  This section
//     resides inside a LOAD segment, so llvm-objcopy --remove-section
//     cannot safely remove it (doing so re-lays out the LOAD segment and
//     corrupts Go's .gopclntab PC tables).  This tool zeros its content
//     in-place, preserving the binary layout.
//
//  3. ELF entry point: shared libraries have a non-deterministic e_entry
//     value that differs between environments.  Since the dynamic linker
//     uses DT_INIT/DT_INIT_ARRAY (not e_entry) for shared libraries,
//     zeroing it is safe.
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
//
// NOTE: .eh_frame and .eh_frame_hdr are intentionally excluded.
// They reside in LOAD segments and are processed by the Android dynamic
// linker at load time.  Zeroing them causes Go runtime stack overflows.
var sectionsToZero = []string{
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

	f, err := elf.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ELF %s: %v\n", path, err)
		os.Exit(1)
	}
	defer f.Close()

	modified := zeroBuildInfo(data, f)
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

// zeroBuildInfo finds and zeros Go build info blocks within the .rodata
// section of the binary.  The search is restricted to .rodata to prevent
// accidental corruption of other sections (.gopclntab, .text, etc.) that
// could contain the same byte pattern.
// Returns true if any modifications were made.
func zeroBuildInfo(data []byte, f *elf.File) bool {
	rodata := f.Section(".rodata")
	if rodata == nil {
		fmt.Fprintln(os.Stderr, "No .rodata section found — skipping build info zeroing")
		return false
	}

	start := rodata.Offset
	end := rodata.Offset + rodata.FileSize
	if end > uint64(len(data)) {
		fmt.Fprintf(os.Stderr, "Warning: .rodata extends beyond file — clamping\n")
		end = uint64(len(data))
	}

	modified := false

	for _, marker := range buildInfoMarkers {
		searchFrom := start
		for {
			region := data[searchFrom:end]
			idx := bytes.Index(region, marker)
			if idx == -1 {
				break
			}

			absIdx := searchFrom + uint64(idx)

			// Find the end of the build info block (terminated by \n\x00 or \x00).
			blockData := data[absIdx:end]
			blockEnd := bytes.Index(blockData, []byte("\n\x00"))
			if blockEnd == -1 {
				blockEnd = bytes.IndexByte(blockData, 0x00)
				if blockEnd == -1 {
					fmt.Fprintf(os.Stderr, "Could not find end of build info block at 0x%x — skipping\n", absIdx)
					searchFrom = absIdx + uint64(len(marker))
					continue
				}
			}
			blockEnd++ // include the terminator byte

			length := blockEnd
			fmt.Fprintf(os.Stderr, "Zeroing Go build info at offset 0x%x, length %d bytes (marker: %q)\n", absIdx, length, marker)

			for i := 0; i < blockEnd; i++ {
				data[absIdx+uint64(i)] = 0
			}
			modified = true
			searchFrom = absIdx + uint64(blockEnd)
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
