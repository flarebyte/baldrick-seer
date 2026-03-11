package cli

import (
	"io"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func Execute(args []string, stdout io.Writer, stderr io.Writer) int {
	cmd := NewRootCmd()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		presentation := domain.PresentError(err)
		_, _ = io.WriteString(stderr, presentation.Stderr)
		return int(presentation.ExitCode)
	}

	return int(domain.ExitCodeSuccess)
}

func renderValidateSuccess() string {
	return "status: ok\ncommand: validate\nmessage: validate stub ok\n"
}

func renderReportGenerateSuccess(result domain.CommandResult) string {
	return result.RenderedOutput
}
