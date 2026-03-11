package app

import (
	"context"

	"github.com/flarebyte/baldrick-seer/internal/domain"
	"github.com/flarebyte/baldrick-seer/internal/pipeline"
)

var reportGenerateRunner = pipeline.NewDefaultRunner()

func RunReportGenerate(ctx context.Context, req domain.CommandRequest) (domain.CommandResult, error) {
	return reportGenerateRunner.RunReportGenerate(ctx, req)
}
