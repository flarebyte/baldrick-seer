package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
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

func executeAndCapture(args []string) (int, string, string) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	exitCode := Execute(args, stdout, stderr)
	return exitCode, stdout.String(), stderr.String()
}

func assertExecuteRepeated(t *testing.T, args []string, iterations int) {
	t.Helper()

	wantExit, wantStdout, wantStderr := executeAndCapture(args)

	for i := 0; i < iterations; i++ {
		gotExit, gotStdout, gotStderr := executeAndCapture(args)
		if gotExit != wantExit {
			t.Fatalf("iteration %d exitCode = %d, want %d", i, gotExit, wantExit)
		}
		if gotStdout != wantStdout {
			t.Fatalf("iteration %d stdout drifted", i)
		}
		if gotStderr != wantStderr {
			t.Fatalf("iteration %d stderr drifted", i)
		}
	}
}

func assertPresentedCommandError(t *testing.T, err error, wantExitCode int, wantStdout string, wantStderr string) {
	t.Helper()

	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}

	presentation := domain.PresentError(err)
	if got, want := int(presentation.ExitCode), wantExitCode; got != want {
		t.Fatalf("exitCode = %d, want %d", got, want)
	}
	if got, want := presentation.Stderr, wantStderr; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}

func assertRootCommandFailure(
	t *testing.T,
	deps dependencies,
	args []string,
	wantExitCode int,
	wantStdout string,
	wantStderr string,
) {
	t.Helper()

	cmd := newRootCmd(deps)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	err := cmd.Execute()
	if got, want := stdout.String(), wantStdout; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
	assertPresentedCommandError(t, err, wantExitCode, wantStdout, wantStderr)
}

func validateFailureDeps(err error) dependencies {
	return dependencies{
		runValidate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
			return domain.CommandResult{}, err
		},
		runReportGenerate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
			panic("runReportGenerate should not be called")
		},
	}
}

func reportGenerateFailureDeps(err error) dependencies {
	return dependencies{
		runValidate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
			panic("runValidate should not be called")
		},
		runReportGenerate: func(context.Context, domain.CommandRequest) (domain.CommandResult, error) {
			return domain.CommandResult{}, err
		},
	}
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

func TestRootCommandExposesBuildVersion(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd(dependencies{})

	if cmd.Version == "" {
		t.Fatal("Version = empty, want build metadata string")
	}
}

func TestHelpIncludesRepositoryAndExamples(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		wantSnippets []string
	}{
		{
			name: "root help",
			args: []string{"--help"},
			wantSnippets: []string{
				repositoryURL,
				"seer validate --config testdata/config/minimal.cue",
				"seer report generate --config testdata/config/valid_report.cue",
			},
		},
		{
			name: "validate help",
			args: []string{"validate", "--help"},
			wantSnippets: []string{
				repositoryURL,
				"Path to a .cue file or a directory containing a CUE package",
				"seer validate --config testdata/config_split",
			},
		},
		{
			name: "report generate help",
			args: []string{"report", "generate", "--help"},
			wantSnippets: []string{
				repositoryURL,
				"AHP weighting",
				"seer report generate --config testdata/config/valid_report.cue",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			exitCode, stdout, stderr := executeAndCapture(tt.args)
			if exitCode != 0 {
				t.Fatalf("exitCode = %d, want 0", exitCode)
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}
			for _, snippet := range tt.wantSnippets {
				if !strings.Contains(stdout, snippet) {
					t.Fatalf("stdout did not contain %q\nstdout=%s", snippet, stdout)
				}
			}
		})
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

	assertExecuteRepeated(t, []string{"validate", "--config", testConfigPath()}, 1)
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
			deps: validateFailureDeps(domain.NewFailure(
				domain.FailureCategoryValidation,
				[]domain.Diagnostic{
					domain.NewDiagnostic(domain.DiagnosticSeverityError, "validation.failed", "config", domain.DiagnosticLocation{}, "validation failed"),
				},
				nil,
			)),
			args:         []string{"validate", "--config", testConfigPath()},
			wantExitCode: 1,
			wantStderr:   "status: error\nmessage: validation failed\n",
		},
		{
			name: "report generate execution failure",
			deps: reportGenerateFailureDeps(domain.NewFailure(
				domain.FailureCategoryExecution,
				[]domain.Diagnostic{
					domain.NewDiagnostic(domain.DiagnosticSeverityError, "execution.failed", "config", domain.DiagnosticLocation{}, "command failed"),
				},
				nil,
			)),
			args:         []string{"report", "generate", "--config", testConfigPath()},
			wantExitCode: 1,
			wantStderr:   "status: error\nmessage: command failed\n",
		},
		{
			name: "validate cancellation failure",
			deps: validateFailureDeps(domain.NewFailureWithMessage(
				domain.FailureCategoryExecution,
				"command canceled",
				nil,
				nil,
			)),
			args:         []string{"validate", "--config", testConfigPath()},
			wantExitCode: 1,
			wantStderr:   "status: error\nmessage: command canceled\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertRootCommandFailure(t, tt.deps, tt.args, tt.wantExitCode, tt.wantStdout, tt.wantStderr)
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

			exitCode, _, _ := executeAndCapture(tt.args)
			if exitCode != tt.wantExitCode {
				t.Fatalf("exitCode = %d, want %d", exitCode, tt.wantExitCode)
			}
			assertExecuteRepeated(t, tt.args, 1)
		})
	}
}
