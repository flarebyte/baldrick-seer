package pipeline

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestDefaultConfigLoader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		configPath string
	}{
		{
			name:       "single file",
			configPath: fixtureConfigPath(),
		},
		{
			name:       "directory package",
			configPath: filepath.Join("..", "..", "testdata", "config_split"),
		},
	}

	loader := DefaultConfigLoader{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mustLoadConfig(t, loader, tt.configPath)
			if got.Config.Path != filepath.Clean(tt.configPath) {
				t.Fatalf("ConfigPath = %q, want %q", got.Config.Path, filepath.Clean(tt.configPath))
			}
		})
	}
}

func TestDefaultConfigLoaderMissingFile(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}

	_, err := loader.LoadConfig(context.Background(), LoadConfigInput{
		ConfigPath: filepath.Join("..", "..", "testdata", "config", "missing.cue"),
	})
	assertLoaderFailure(t, err, ErrConfigPathDoesNotExist, domain.FailureCategoryInput, "config.not_found", "config path does not exist")
}

func TestDefaultConfigLoaderEmptyDirectoryPath(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}

	_, err := loader.LoadConfig(context.Background(), LoadConfigInput{
		ConfigPath: filepath.Join("..", "..", "testdata", "config_empty"),
	})
	assertLoaderFailure(t, err, ErrConfigDirectoryEmpty, domain.FailureCategoryInput, "config.directory_empty", "config directory does not contain any .cue files")
}
