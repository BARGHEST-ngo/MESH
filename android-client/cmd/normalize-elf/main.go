// Command normalize-elf zeros non-deterministic sections in an ELF shared
// library to achieve cross-platform reproducible builds.
//
// The NDK linker produces subtly different output depending on the host
// environment (Ubuntu vs Debian, different glibc versions, etc.).  The
// affected sections (.eh_frame, .eh_frame_hdr, .relro_padding) reside
// inside LOAD segments, so llvm-objcopy --remove-section silently fails
// to remove them.  This tool instead zeros their content in-place,
// preserving the binary layout while making the content deterministic.
//
// It also zeros the ELF entry point, which differs between environments
// for shared libraries (some linkers set it to the .text start, others
// set it to 0).
//
// Usage: go run ./cmd/normalize-elf <path-to-libgojni.so>
package main

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"os"
)

// sectionsToZero lists ELF sections whose content should be zeroed.
// These sections differ across build environments but are not required
// at runtime for a shared library on Android.
var sectionsToZero = []string{
	".eh_frame",     // DWARF call frame information — differs by compiler/linker
	".eh_frame_hdr", // Index into .eh_frame
	".relro_padding", // RELRO alignment padding — presence varies by linker
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

	modified := false

	// Zero the content of non-deterministic sections.
	for _, name := range sectionsToZero {
		sec := f.Section(name)
		if sec == nil {
			continue
		}

		// SHT_NOBITS sections (like .relro_padding when it has no file
		// content) have Size > 0 but occupy no space in the file.
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

		// Check if already zeroed.
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

	// Zero the ELF entry point.
	// Shared libraries don't need an entry point, but some linkers set it
	// to the start of .text while others set it to 0.
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

// zeroEntryPoint zeros the e_entry field in the ELF header.
// Returns true if the field was modified.
func zeroEntryPoint(data []byte, class elf.Class) bool {
	// The e_entry field is at offset 0x18 in both ELF32 and ELF64 headers.
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
