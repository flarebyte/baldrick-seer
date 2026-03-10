package app

import (
	"github.com/flarebyte/baldrick-seer/internal/domain"
	"github.com/flarebyte/baldrick-seer/internal/pipeline"
)

var validateRunner = pipeline.NewDefaultRunner()

func RunValidate(req domain.CommandRequest) (domain.CommandResult, error) {
	return validateRunner.RunValidate(req)
}
