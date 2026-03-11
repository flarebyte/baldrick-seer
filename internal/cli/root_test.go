package cli

import (
	"bytes"
	"context"
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
			name:         "report generate empty directory",
			args:         []string{"report", "generate", "--config", filepath.Join("..", "..", "testdata", "config_empty")},
			wantExitCode: 1,
			stderrGolden: "directory_path.stderr.golden",
		},
		{
			name:         "validate non concrete cue",
			args:         []string{"validate", "--config", filepath.Join("..", "..", "testdata", "config", "non_concrete.cue")},
			wantExitCode: 1,
			stderrGolden: "invalid_cue.stderr.golden",
		},
		{
			name:         "validate semantic validation failure",
			args:         []string{"validate", "--config", filepath.Join("..", "..", "testdata", "config", "invalid_reference.cue")},
			wantExitCode: 1,
			stderrGolden: "invalid_validation.stderr.golden",
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
		runValidate: func(ctx context.Context, req domain.CommandRequest) (domain.CommandResult, error) {
			called = true
			if ctx == nil {
				t.Fatal("context = nil, want value")
			}
			if req.CommandName != domain.CommandNameValidate {
				t.Fatalf("CommandName = %q, want %q", req.CommandName, domain.CommandNameValidate)
			}
			if want := testConfigPath(); req.ConfigPath != want {
				t.Fatalf("ConfigPath = %q, want %q", req.ConfigPath, want)
			}

			return domain.CommandResult{CommandName: domain.CommandNameValidate}, nil
		},
		runReportGenerate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
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
		runValidate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
			t.Fatal("runValidate should not be called")
			return domain.CommandResult{}, nil
		},
		runReportGenerate: func(ctx context.Context, req domain.CommandRequest) (domain.CommandResult, error) {
			called = true
			if ctx == nil {
				t.Fatal("context = nil, want value")
			}
			if req.CommandName != domain.CommandNameReportGenerate {
				t.Fatalf("CommandName = %q, want %q", req.CommandName, domain.CommandNameReportGenerate)
			}
			if want := testConfigPath(); req.ConfigPath != want {
				t.Fatalf("ConfigPath = %q, want %q", req.ConfigPath, want)
			}

			return domain.CommandResult{
				CommandName:    domain.CommandNameReportGenerate,
				RenderedOutput: readGolden(t, "report_generate_success.stdout.golden"),
			}, nil
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

func TestExecuteUsesCentralizedFailureEmission(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		deps         dependencies
		args         []string
		wantExitCode int
		wantStdout   string
		wantStderr   string
	}{
		{
			name: "validate validation failure",
			deps: dependencies{
				runValidate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
					return domain.CommandResult{}, domain.NewFailure(
						domain.FailureCategoryValidation,
						[]domain.Diagnostic{
							domain.NewDiagnostic(domain.DiagnosticSeverityError, "validation.failed", "config", domain.DiagnosticLocation{}, "validation failed"),
						},
						nil,
					)
				},
				runReportGenerate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
					t.Fatal("runReportGenerate should not be called")
					return domain.CommandResult{}, nil
				},
			},
			args:         []string{"validate", "--config", testConfigPath()},
			wantExitCode: 1,
			wantStderr:   "status: error\nmessage: validation failed\n",
		},
		{
			name: "report generate execution failure",
			deps: dependencies{
				runValidate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
					t.Fatal("runValidate should not be called")
					return domain.CommandResult{}, nil
				},
				runReportGenerate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
					return domain.CommandResult{}, domain.NewFailure(
						domain.FailureCategoryExecution,
						[]domain.Diagnostic{
							domain.NewDiagnostic(domain.DiagnosticSeverityError, "execution.failed", "config", domain.DiagnosticLocation{}, "command failed"),
						},
						nil,
					)
				},
			},
			args:         []string{"report", "generate", "--config", testConfigPath()},
			wantExitCode: 1,
			wantStderr:   "status: error\nmessage: command failed\n",
		},
		{
			name: "validate cancellation failure",
			deps: dependencies{
				runValidate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
					return domain.CommandResult{}, domain.NewFailureWithMessage(
						domain.FailureCategoryExecution,
						"command canceled",
						nil,
						nil,
					)
				},
				runReportGenerate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
					t.Fatal("runReportGenerate should not be called")
					return domain.CommandResult{}, nil
				},
			},
			args:         []string{"validate", "--config", testConfigPath()},
			wantExitCode: 1,
			wantStderr:   "status: error\nmessage: command canceled\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd := newRootCmd(tt.deps)
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("Execute() error = nil, want error")
			}

			presentation := domain.PresentError(err)
			if got, want := int(presentation.ExitCode), tt.wantExitCode; got != want {
				t.Fatalf("exitCode = %d, want %d", got, want)
			}
			if got, want := stdout.String(), tt.wantStdout; got != want {
				t.Fatalf("stdout = %q, want %q", got, want)
			}
			if got, want := presentation.Stderr, tt.wantStderr; got != want {
				t.Fatalf("stderr = %q, want %q", got, want)
			}
		})
	}
}

func TestExecuteRepeatedRunDeterminism(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
	}{
		{
			name:         "validate success",
			args:         []string{"validate", "--config", testConfigPath()},
			wantExitCode: 0,
		},
		{
			name:         "report generate success",
			args:         []string{"report", "generate", "--config", testConfigPath()},
			wantExitCode: 0,
		},
		{
			name:         "validate validation failure",
			args:         []string{"validate", "--config", filepath.Join("..", "..", "testdata", "config", "invalid_reference.cue")},
			wantExitCode: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout1 := new(bytes.Buffer)
			stderr1 := new(bytes.Buffer)
			exitCode1 := Execute(tt.args, stdout1, stderr1)

			stdout2 := new(bytes.Buffer)
			stderr2 := new(bytes.Buffer)
			exitCode2 := Execute(tt.args, stdout2, stderr2)

			if exitCode1 != tt.wantExitCode || exitCode2 != tt.wantExitCode {
				t.Fatalf("exit codes = (%d, %d), want (%d, %d)", exitCode1, exitCode2, tt.wantExitCode, tt.wantExitCode)
			}
			if stdout1.String() != stdout2.String() {
				t.Fatalf("stdout1 = %q, stdout2 = %q", stdout1.String(), stdout2.String())
			}
			if stderr1.String() != stderr2.String() {
				t.Fatalf("stderr1 = %q, stderr2 = %q", stderr1.String(), stderr2.String())
			}
		})
	}
}
