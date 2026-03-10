package cli

import (
	"bytes"
	"testing"
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

	if got, want := stdout.String(), validateStubOutput; got != want {
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

	if got, want := stdout.String(), reportGenerateStubOutput; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}
