package cli

import (
	"bytes"
	"path/filepath"
	"testing"
)

func largeCLIValidConfigPath() string {
	return filepath.Join("..", "..", "testdata", "config", "large_valid.cue")
}

func largeCLIInvalidConfigPath() string {
	return filepath.Join("..", "..", "testdata", "config", "large_invalid.cue")
}

func TestExecuteLargeFixtureStressDeterminism(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "validate success",
			args: []string{"validate", "--config", largeCLIValidConfigPath()},
		},
		{
			name: "report generate success",
			args: []string{"report", "generate", "--config", largeCLIValidConfigPath()},
		},
		{
			name: "validate failure",
			args: []string{"validate", "--config", largeCLIInvalidConfigPath()},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)
			wantExit := Execute(tt.args, stdout, stderr)
			wantStdout := stdout.String()
			wantStderr := stderr.String()

			for i := 0; i < 12; i++ {
				nextStdout := new(bytes.Buffer)
				nextStderr := new(bytes.Buffer)
				gotExit := Execute(tt.args, nextStdout, nextStderr)

				if gotExit != wantExit {
					t.Fatalf("iteration %d exitCode = %d, want %d", i, gotExit, wantExit)
				}
				if nextStdout.String() != wantStdout {
					t.Fatalf("iteration %d stdout drifted", i)
				}
				if nextStderr.String() != wantStderr {
					t.Fatalf("iteration %d stderr drifted", i)
				}
			}
		})
	}
}
