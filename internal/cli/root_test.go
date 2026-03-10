package cli

import (
	"bytes"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/app"
)

func TestValidateCommand(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	stdout := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"validate"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got, want := stdout.String(), "validate: ok\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestReportGenerateCommand(t *testing.T) {
	t.Parallel()

	cmd := NewRootCmd()
	stdout := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"report", "generate"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got, want := stdout.String(), "report generate: ok\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestValidateCommandDelegatesToExecutor(t *testing.T) {
	t.Parallel()

	called := false
	cmd := newRootCmd(dependencies{
		runValidate: func(req app.ValidateRequest) (app.ValidateResponse, error) {
			called = true
			if req != (app.ValidateRequest{}) {
				t.Fatalf("request = %#v, want empty request", req)
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
	cmd.SetArgs([]string{"validate"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !called {
		t.Fatal("runValidate was not called")
	}

	if got, want := stdout.String(), "validate: ok\n"; got != want {
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
			if req != (app.ReportGenerateRequest{}) {
				t.Fatalf("request = %#v, want empty request", req)
			}

			return app.ReportGenerateResponse{Stdout: "report generate: ok\n"}, nil
		},
	})

	stdout := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"report", "generate"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !called {
		t.Fatal("runReportGenerate was not called")
	}

	if got, want := stdout.String(), "report generate: ok\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}
