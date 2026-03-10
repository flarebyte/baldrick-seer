package app

import (
	"github.com/flarebyte/baldrick-seer/internal/domain"
	"github.com/flarebyte/baldrick-seer/internal/pipeline"
)

var reportGenerateRunner = pipeline.NewDefaultRunner()

func RunReportGenerate(req domain.CommandRequest) (domain.CommandResult, error) {
	return reportGenerateRunner.RunReportGenerate(req)
}
