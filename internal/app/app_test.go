package app

import (
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestRunValidate(t *testing.T) {
	t.Parallel()

	got, err := RunValidate(domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("RunValidate() error = %v", err)
	}

	if got.CommandName != domain.CommandNameValidate {
		t.Fatalf("CommandName = %q, want %q", got.CommandName, domain.CommandNameValidate)
	}

	if got.ValidatedModel == nil {
		t.Fatal("ValidatedModel = nil, want value")
	}

	if want := filepath.Clean(filepath.Join("..", "..", "testdata", "config", "minimal.cue")); got.ValidatedModel.ConfigPath != want {
		t.Fatalf("ConfigPath = %q, want %q", got.ValidatedModel.ConfigPath, want)
	}
}

func TestRunReportGenerate(t *testing.T) {
	t.Parallel()

	got, err := RunReportGenerate(domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	if got.CommandName != domain.CommandNameReportGenerate {
		t.Fatalf("CommandName = %q, want %q", got.CommandName, domain.CommandNameReportGenerate)
	}

	if got.ValidatedModel == nil {
		t.Fatal("ValidatedModel = nil, want value")
	}

	if want := filepath.Clean(filepath.Join("..", "..", "testdata", "config", "minimal.cue")); got.ValidatedModel.ConfigPath != want {
		t.Fatalf("ConfigPath = %q, want %q", got.ValidatedModel.ConfigPath, want)
	}
}
