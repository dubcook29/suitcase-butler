package model

import (
	"os"
	"strings"
)

func writeJSONFile(filename string, in []byte) error {
	filename = strings.TrimSuffix(filename, ".json")
	if err := os.WriteFile(filename+".json", in, 0644); err != nil {
		return err
	}
	return nil
}

func readJSONFile(filename string) ([]byte, error) {
	filename = strings.TrimSuffix(filename, ".json")
	return os.ReadFile(filename + ".json")
}
