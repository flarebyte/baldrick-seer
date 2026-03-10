package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func testConfigPath() string {
	return filepath.Join("..", "..", "testdata", "config", "minimal.cue")
}

func readGolden(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join("..", "..", "testdata", "golden", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}

	return string(content)
}

func TestCommandOutputGoldens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		stdoutGolden string
		stderrGolden string
	}{
		{
			name:         "validate success",
			args:         []string{"validate", "--config", testConfigPath()},
			wantExitCode: 0,
			stdoutGolden: "validate_success.stdout.golden",
		},
		{
			name:         "report generate success",
			args:         []string{"report", "generate", "--config", testConfigPath()},
			wantExitCode: 0,
			stdoutGolden: "report_generate_success.stdout.golden",
		},
		{
			name:         "validate missing config",
			args:         []string{"validate"},
			wantExitCode: 1,
			stderrGolden: "missing_config.stderr.golden",
		},
		{
			name:         "report generate missing config",
			args:         []string{"report", "generate"},
			wantExitCode: 1,
			stderrGolden: "missing_config.stderr.golden",
		},
		{
			name:         "validate missing file",
			args:         []string{"validate", "--config", filepath.Join("..", "..", "testdata", "config", "missing.cue")},
			wantExitCode: 1,
			stderrGolden: "missing_file.stderr.golden",
		},
		{
			name:         "report generate directory path",
			args:         []string{"report", "generate", "--config", filepath.Join("..", "..", "testdata", "config")},
			wantExitCode: 1,
			stderrGolden: "directory_path.stderr.golden",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)

			exitCode := Execute(tt.args, stdout, stderr)
			if exitCode != tt.wantExitCode {
				t.Fatalf("exitCode = %d, want %d", exitCode, tt.wantExitCode)
			}

			wantStdout := ""
			if tt.stdoutGolden != "" {
				wantStdout = readGolden(t, tt.stdoutGolden)
			}

			if got, want := stdout.String(), wantStdout; got != want {
				t.Fatalf("stdout = %q, want %q", got, want)
			}

			wantStderr := ""
			if tt.stderrGolden != "" {
				wantStderr = readGolden(t, tt.stderrGolden)
			}

			if got, want := stderr.String(), wantStderr; got != want {
				t.Fatalf("stderr = %q, want %q", got, want)
			}
		})
	}
}

func TestValidateCommandDelegatesToExecutor(t *testing.T) {
	t.Parallel()

	called := false
	cmd := newRootCmd(dependencies{
		runValidate: func(req domain.CommandRequest) (domain.CommandResult, error) {
			called = true
			if req.CommandName != domain.CommandNameValidate {
				t.Fatalf("CommandName = %q, want %q", req.CommandName, domain.CommandNameValidate)
			}
			if want := testConfigPath(); req.ConfigPath != want {
				t.Fatalf("ConfigPath = %q, want %q", req.ConfigPath, want)
			}

			return domain.CommandResult{CommandName: domain.CommandNameValidate}, nil
		},
		runReportGenerate: func(domain.CommandRequest) (domain.CommandResult, error) {
			t.Fatal("runReportGenerate should not be called")
			return domain.CommandResult{}, nil
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

	if got, want := stdout.String(), readGolden(t, "validate_success.stdout.golden"); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestReportGenerateCommandDelegatesToExecutor(t *testing.T) {
	t.Parallel()

	called := false
	cmd := newRootCmd(dependencies{
		runValidate: func(domain.CommandRequest) (domain.CommandResult, error) {
			t.Fatal("runValidate should not be called")
			return domain.CommandResult{}, nil
		},
		runReportGenerate: func(req domain.CommandRequest) (domain.CommandResult, error) {
			called = true
			if req.CommandName != domain.CommandNameReportGenerate {
				t.Fatalf("CommandName = %q, want %q", req.CommandName, domain.CommandNameReportGenerate)
			}
			if want := testConfigPath(); req.ConfigPath != want {
				t.Fatalf("ConfigPath = %q, want %q", req.ConfigPath, want)
			}

			return domain.CommandResult{CommandName: domain.CommandNameReportGenerate}, nil
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

	if got, want := stdout.String(), readGolden(t, "report_generate_success.stdout.golden"); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestExecuteIsDeterministic(t *testing.T) {
	t.Parallel()

	args := []string{"validate", "--config", testConfigPath()}

	stdout1 := new(bytes.Buffer)
	stderr1 := new(bytes.Buffer)
	exitCode1 := Execute(args, stdout1, stderr1)

	stdout2 := new(bytes.Buffer)
	stderr2 := new(bytes.Buffer)
	exitCode2 := Execute(args, stdout2, stderr2)

	if exitCode1 != exitCode2 {
		t.Fatalf("exitCode1 = %d, exitCode2 = %d", exitCode1, exitCode2)
	}

	if stdout1.String() != stdout2.String() {
		t.Fatalf("stdout1 = %q, stdout2 = %q", stdout1.String(), stdout2.String())
	}

	if stderr1.String() != stderr2.String() {
		t.Fatalf("stderr1 = %q, stderr2 = %q", stderr1.String(), stderr2.String())
	}
}
