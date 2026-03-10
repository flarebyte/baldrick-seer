package app

import "github.com/flarebyte/baldrick-seer/internal/domain"

func RunReportGenerate(req domain.CommandRequest) (domain.CommandResult, error) {
	config, err := LoadConfig(ConfigRequest{Path: req.ConfigPath})
	if err != nil {
		return domain.CommandResult{}, err
	}

	return domain.CommandResult{
		CommandName: domain.CommandNameReportGenerate,
		ValidatedModel: &domain.ValidatedModelSummary{
			ConfigPath: config.Path,
		},
	}, nil
}
