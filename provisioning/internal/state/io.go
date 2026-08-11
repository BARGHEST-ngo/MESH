package state

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func loadJSON[T any](path string, v *T) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func saveJSON[T any](path string, v T) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return os.Rename(tmp, path)
}
