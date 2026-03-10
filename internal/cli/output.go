package cli

import (
	"errors"
	"io"

	"github.com/flarebyte/baldrick-seer/internal/app"
)

func Execute(args []string, stdout io.Writer, stderr io.Writer) int {
	cmd := NewRootCmd()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		_, _ = io.WriteString(stderr, renderFailure(err))
		return 1
	}

	return 0
}

func renderValidateSuccess() string {
	return "status: ok\ncommand: validate\nmessage: validate stub ok\n"
}

func renderReportGenerateSuccess() string {
	return "status: ok\ncommand: report generate\nmessage: report generate stub ok\n"
}

func renderFailure(err error) string {
	return "status: error\nmessage: " + failureMessage(err) + "\n"
}

func failureMessage(err error) string {
	switch {
	case errors.Is(err, app.ErrConfigPathRequired):
		return "config flag is required"
	case errors.Is(err, app.ErrConfigPathDoesNotExist):
		return "config path does not exist"
	case errors.Is(err, app.ErrConfigPathIsDirectory):
		return "config path is a directory"
	default:
		return "command failed"
	}
}
