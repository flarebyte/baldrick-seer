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
		_, _ = io.WriteString(stderr, renderFailure(err))
		return int(domain.ExitCodeForError(err))
	}

	return int(domain.ExitCodeSuccess)
}

func renderValidateSuccess() string {
	return "status: ok\ncommand: validate\nmessage: validate stub ok\n"
}

func renderReportGenerateSuccess(result domain.CommandResult) string {
	return result.RenderedOutput
}

func renderFailure(err error) string {
	return "status: error\nmessage: " + failureMessage(err) + "\n"
}

func failureMessage(err error) string {
	if failure := domain.AsCommandFailure(err); failure != nil {
		if len(failure.Diagnostics) > 0 && failure.Diagnostics[0].Message != "" {
			return failure.Diagnostics[0].Message
		}
	}

	return "command failed"
}
