package cli

import (
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

			assertExecuteRepeated(t, tt.args, 12)
		})
	}
}
