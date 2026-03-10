package app

import (
	"path/filepath"
	"testing"
)

func TestLoadConfigExistingFile(t *testing.T) {
	t.Parallel()

	got, err := LoadConfig(ConfigRequest{
		Path: filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if want := filepath.Clean(filepath.Join("..", "..", "testdata", "config", "minimal.cue")); got.Path != want {
		t.Fatalf("Path = %q, want %q", got.Path, want)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	t.Parallel()

	_, err := LoadConfig(ConfigRequest{
		Path: filepath.Join("..", "..", "testdata", "config", "missing.cue"),
	})
	if err != ErrConfigPathDoesNotExist {
		t.Fatalf("error = %v, want %v", err, ErrConfigPathDoesNotExist)
	}
}

func TestLoadConfigDirectoryPath(t *testing.T) {
	t.Parallel()

	_, err := LoadConfig(ConfigRequest{
		Path: filepath.Join("..", "..", "testdata", "config"),
	})
	if err != ErrConfigPathIsDirectory {
		t.Fatalf("error = %v, want %v", err, ErrConfigPathIsDirectory)
	}
}
