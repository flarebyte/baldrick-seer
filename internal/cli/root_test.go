package cli

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/app"
)

func testConfigPath() string {
	return filepath.Join("..", "..", "testdata", "config", "minimal.cue")
}

func TestValidateCommand(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	stdout := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"validate", "--config", testConfigPath()})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got, want := stdout.String(), "status: ok\ncommand: validate\nmessage: validate stub ok\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestReportGenerateCommand(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	stdout := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"report", "generate", "--config", testConfigPath()})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got, want := stdout.String(), "status: ok\ncommand: report generate\nmessage: report generate stub ok\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestValidateCommandDelegatesToExecutor(t *testing.T) {
	t.Parallel()

	called := false
	cmd := newRootCmd(dependencies{
		runValidate: func(req app.ValidateRequest) (app.ValidateResponse, error) {
			called = true
			if want := testConfigPath(); req.ConfigPath != want {
				t.Fatalf("ConfigPath = %q, want %q", req.ConfigPath, want)
			}

			return app.ValidateResponse{Stdout: "validate: ok\n"}, nil
		},
		runReportGenerate: func(app.ReportGenerateRequest) (app.ReportGenerateResponse, error) {
			t.Fatal("runReportGenerate should not be called")
			return app.ReportGenerateResponse{}, nil
		},
	})

	stdout := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"validate", "--config", testConfigPath()})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !called {
		t.Fatal("runValidate was not called")
	}

	if got, want := stdout.String(), "status: ok\ncommand: validate\nmessage: validate stub ok\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestReportGenerateCommandDelegatesToExecutor(t *testing.T) {
	t.Parallel()

	called := false
	cmd := newRootCmd(dependencies{
		runValidate: func(app.ValidateRequest) (app.ValidateResponse, error) {
			t.Fatal("runValidate should not be called")
			return app.ValidateResponse{}, nil
		},
		runReportGenerate: func(req app.ReportGenerateRequest) (app.ReportGenerateResponse, error) {
			called = true
			if want := testConfigPath(); req.ConfigPath != want {
				t.Fatalf("ConfigPath = %q, want %q", req.ConfigPath, want)
			}

			return app.ReportGenerateResponse{Stdout: "report generate: ok\n"}, nil
		},
	})

	stdout := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"report", "generate", "--config", testConfigPath()})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !called {
		t.Fatal("runReportGenerate was not called")
	}

	if got, want := stdout.String(), "status: ok\ncommand: report generate\nmessage: report generate stub ok\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestValidateCommandRequiresConfig(t *testing.T) {
	t.Parallel()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	exitCode := Execute([]string{"validate"}, stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}

	if got, want := stdout.String(), ""; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}

	if got, want := stderr.String(), "status: error\nmessage: config flag is required\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}

func TestValidateCommandRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	exitCode := Execute([]string{"validate", "--config", filepath.Join("..", "..", "testdata", "config", "missing.cue")}, stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}

	if got, want := stdout.String(), ""; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}

	if got, want := stderr.String(), "status: error\nmessage: config path does not exist\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}

func TestReportGenerateCommandRequiresConfig(t *testing.T) {
	t.Parallel()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	exitCode := Execute([]string{"report", "generate"}, stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}

	if got, want := stdout.String(), ""; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}

	if got, want := stderr.String(), "status: error\nmessage: config flag is required\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}

func TestReportGenerateCommandRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	exitCode := Execute([]string{"report", "generate", "--config", filepath.Join("..", "..", "testdata", "config")}, stdout, stderr)
	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}

	if got, want := stdout.String(), ""; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}

	if got, want := stderr.String(), "status: error\nmessage: config path is a directory\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}
