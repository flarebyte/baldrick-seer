package app

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func assertAppResult(t *testing.T, result domain.CommandResult, wantCommandName domain.CommandName, wantConfigPath string) {
	t.Helper()

	if result.CommandName != wantCommandName {
		t.Fatalf("CommandName = %q, want %q", result.CommandName, wantCommandName)
	}

	if result.ValidatedModel == nil {
		t.Fatal("ValidatedModel = nil, want value")
	}

	if result.ValidatedModel.ConfigPath != wantConfigPath {
		t.Fatalf("ConfigPath = %q, want %q", result.ValidatedModel.ConfigPath, wantConfigPath)
	}
}

func TestAppFlows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		run             func(context.Context, domain.CommandRequest) (domain.CommandResult, error)
		wantCommandName domain.CommandName
	}{
		{
			name:            "validate",
			run:             RunValidate,
			wantCommandName: domain.CommandNameValidate,
		},
		{
			name:            "report generate",
			run:             RunReportGenerate,
			wantCommandName: domain.CommandNameReportGenerate,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.run(context.Background(), domain.CommandRequest{
				CommandName: tt.wantCommandName,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
			})
			if err != nil {
				t.Fatalf("run() error = %v", err)
			}

			assertAppResult(
				t,
				got,
				tt.wantCommandName,
				filepath.Clean(filepath.Join("..", "..", "testdata", "config", "minimal.cue")),
			)
		})
	}
}
