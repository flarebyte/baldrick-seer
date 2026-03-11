package pipeline

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestDefaultConfigLoaderWithCue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		configPath    string
		wantErr       error
		wantCategory  domain.FailureCategory
		wantCode      string
		wantMessage   string
		wantFields    []string
		wantEvaluated string
	}{
		{
			name:          "valid minimal cue config",
			configPath:    filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
			wantFields:    []string{"config"},
			wantEvaluated: "{\n\tconfig: {\n\t\tname: \"minimal\"\n\t}\n}",
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
			name:         "directory path",
			configPath:   filepath.Join("..", "..", "testdata", "config"),
			wantErr:      ErrConfigPathIsDirectory,
			wantCategory: domain.FailureCategoryInput,
			wantCode:     "config.is_directory",
			wantMessage:  "config path is a directory",
		},
	}

	loader := DefaultConfigLoader{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := loader.LoadConfig(LoadConfigInput{ConfigPath: tt.configPath})
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("LoadConfig() error = %v", err)
				}

				if got.Config.Path != filepath.Clean(tt.configPath) {
					t.Fatalf("Path = %q, want %q", got.Config.Path, filepath.Clean(tt.configPath))
				}

				if !reflect.DeepEqual(got.Config.TopLevelFields, tt.wantFields) {
					t.Fatalf("TopLevelFields = %#v, want %#v", got.Config.TopLevelFields, tt.wantFields)
				}

				if got.Config.Evaluated != tt.wantEvaluated {
					t.Fatalf("Evaluated = %q, want %q", got.Config.Evaluated, tt.wantEvaluated)
				}

				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			failure := domain.AsCommandFailure(err)
			if failure == nil {
				t.Fatal("AsCommandFailure(err) = nil, want value")
			}

			if failure.Category != tt.wantCategory {
				t.Fatalf("Category = %q, want %q", failure.Category, tt.wantCategory)
			}

			if failure.Diagnostics[0].Code != tt.wantCode {
				t.Fatalf("Code = %q, want %q", failure.Diagnostics[0].Code, tt.wantCode)
			}

			if failure.Diagnostics[0].Message != tt.wantMessage {
				t.Fatalf("Message = %q, want %q", failure.Diagnostics[0].Message, tt.wantMessage)
			}
		})
	}
}
