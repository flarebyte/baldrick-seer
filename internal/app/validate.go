package app

import (
	"context"

	"github.com/flarebyte/baldrick-seer/internal/domain"
	"github.com/flarebyte/baldrick-seer/internal/pipeline"
)

var validateRunner = pipeline.NewDefaultRunner()

func RunValidate(ctx context.Context, req domain.CommandRequest) (domain.CommandResult, error) {
	return validateRunner.RunValidate(ctx, req)
}
