package pipeline

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestDefaultConfigLoaderWithCue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		configPath       string
		wantErr          error
		wantCategory     domain.FailureCategory
		wantCode         string
		wantMessage      string
		wantFields       []string
		wantConfigFields []string
	}{
		{
			name:       "valid minimal cue config",
			configPath: filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
			wantFields: []string{"config"},
			wantConfigFields: []string{
				"aggregation",
				"alternatives",
				"criteriaCatalog",
				"evaluations",
				"problem",
				"reports",
				"scenarios",
			},
		},
		{
			name:       "valid cue directory config",
			configPath: filepath.Join("..", "..", "testdata", "config_split"),
			wantFields: []string{"config"},
			wantConfigFields: []string{
				"aggregation",
				"alternatives",
				"criteriaCatalog",
				"evaluations",
				"problem",
				"reports",
				"scenarios",
			},
		},
		{
			name:       "valid cue config with report filepath",
			configPath: filepath.Join("..", "..", "testdata", "config", "valid_report_filepath.cue"),
			wantFields: []string{"config"},
			wantConfigFields: []string{
				"aggregation",
				"alternatives",
				"criteriaCatalog",
				"evaluations",
				"problem",
				"reports",
				"scenarios",
			},
		},
		{
			name:         "non concrete cue config",
			configPath:   filepath.Join("..", "..", "testdata", "config", "non_concrete.cue"),
			wantErr:      ErrConfigNotConcrete,
			wantCategory: domain.FailureCategoryInput,
			wantCode:     "config.not_concrete",
			wantMessage:  "config must evaluate to a concrete value",
		},
		{
			name:         "malformed cue config",
			configPath:   filepath.Join("..", "..", "testdata", "config", "malformed.cue"),
			wantErr:      ErrConfigLoadInvalid,
			wantCategory: domain.FailureCategoryInput,
			wantCode:     "config.load_invalid",
			wantMessage:  "config could not be loaded",
		},
		{
			name:         "missing file path",
			configPath:   filepath.Join("..", "..", "testdata", "config", "missing.cue"),
			wantErr:      ErrConfigPathDoesNotExist,
			wantCategory: domain.FailureCategoryInput,
			wantCode:     "config.not_found",
			wantMessage:  "config path does not exist",
		},
		{
			name:         "invalid file extension",
			configPath:   filepath.Join("..", "..", "testdata", "config", "not_cue.txt"),
			wantErr:      ErrConfigFileExtension,
			wantCategory: domain.FailureCategoryInput,
			wantCode:     "config.invalid_file_type",
			wantMessage:  "config file must have .cue extension",
		},
		{
			name:         "empty directory",
			configPath:   filepath.Join("..", "..", "testdata", "config_empty"),
			wantErr:      ErrConfigDirectoryEmpty,
			wantCategory: domain.FailureCategoryInput,
			wantCode:     "config.directory_empty",
			wantMessage:  "config directory does not contain any .cue files",
		},
	}

	loader := DefaultConfigLoader{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := loader.LoadConfig(context.Background(), LoadConfigInput{ConfigPath: tt.configPath})
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("LoadConfig() error = %v", err)
				}
				assertLoadedConfigSuccess(t, got, tt.configPath, tt.wantFields, tt.wantConfigFields)
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			assertLoaderFailure(t, err, tt.wantErr, tt.wantCategory, tt.wantCode, tt.wantMessage)
		})
	}
}

func TestDefaultConfigLoaderDecodesReportFilepath(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}
	got := mustLoadConfig(t, loader, filepath.Join("..", "..", "testdata", "config", "valid_report_filepath.cue"))

	if got.Config.Config == nil {
		t.Fatal("Config = nil, want value")
	}
	if len(got.Config.Config.Reports) != 2 {
		t.Fatalf("len(Reports) = %d, want 2", len(got.Config.Config.Reports))
	}

	if got, want := got.Config.Config.Reports[0].Filepath, "../artifacts/summary.md"; got != want {
		t.Fatalf("Reports[0].Filepath = %q, want %q", got, want)
	}
	if got, want := got.Config.Config.Reports[1].Filepath, "artifacts/summary.json"; got != want {
		t.Fatalf("Reports[1].Filepath = %q, want %q", got, want)
	}
}

func TestDefaultConfigLoaderDirectoryDeterminism(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}
	path := filepath.Join("..", "..", "testdata", "config_split")

	assertRepeatedDeepEqual(t, 1, func() (LoadConfigOutput, error) {
		return loader.LoadConfig(context.Background(), LoadConfigInput{ConfigPath: path})
	})
}

func TestDefaultConfigLoaderFileAndDirectoryEquivalence(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}

	fromFile := mustLoadConfig(t, loader, filepath.Join("..", "..", "testdata", "config", "minimal.cue"))
	fromDirectory := mustLoadConfig(t, loader, filepath.Join("..", "..", "testdata", "config_split"))

	if !reflect.DeepEqual(fromFile.Config.TopLevelFields, fromDirectory.Config.TopLevelFields) {
		t.Fatalf("TopLevelFields(file) = %#v, TopLevelFields(directory) = %#v", fromFile.Config.TopLevelFields, fromDirectory.Config.TopLevelFields)
	}

	if !reflect.DeepEqual(fromFile.Config.ConfigFields, fromDirectory.Config.ConfigFields) {
		t.Fatalf("ConfigFields(file) = %#v, ConfigFields(directory) = %#v", fromFile.Config.ConfigFields, fromDirectory.Config.ConfigFields)
	}

	if !reflect.DeepEqual(fromFile.Config.Config, fromDirectory.Config.Config) {
		t.Fatalf("Config(file) = %#v, Config(directory) = %#v", fromFile.Config.Config, fromDirectory.Config.Config)
	}
}
