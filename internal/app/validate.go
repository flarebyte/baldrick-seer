package app

import "github.com/flarebyte/baldrick-seer/internal/domain"

func RunValidate(req domain.CommandRequest) (domain.CommandResult, error) {
	config, err := LoadConfig(ConfigRequest{Path: req.ConfigPath})
	if err != nil {
		return domain.CommandResult{}, err
	}

	return domain.CommandResult{
		CommandName: domain.CommandNameValidate,
		ValidatedModel: &domain.ValidatedModelSummary{
			ConfigPath: config.Path,
		},
	}, nil
}
