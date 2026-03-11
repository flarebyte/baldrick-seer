package pipeline

import "context"

import "github.com/flarebyte/baldrick-seer/internal/domain"

type DefaultModelValidator struct{}

type scenarioValidationInfo struct {
	Index                int
	ActiveCriterionNames []string
}

func (DefaultModelValidator) ValidateModel(ctx context.Context, input ValidateModelInput) (ValidateModelOutput, error) {
	if err := checkContext(ctx, input.Command.ConfigPath); err != nil {
		return ValidateModelOutput{}, err
	}

	diagnostics := validateLoadedConfig(input.Config)
	if len(diagnostics) > 0 {
		return ValidateModelOutput{}, NewValidationDiagnosticsFailure(diagnostics, ErrValidationFailed)
	}

	reportDefinitions := make([]domain.ReportDefinition, 0, len(input.Config.Config.Reports))
	for _, report := range input.Config.Config.Reports {
		reportDefinitions = append(reportDefinitions, domain.ReportDefinition{
			Name:   report.Name,
			Title:  report.Title,
			Format: report.Format,
		})
	}

	return ValidateModelOutput{
		ValidatedModel: domain.CanonicalValidatedModelSummary(domain.ValidatedModelSummary{
			ConfigPath:        input.Config.Path,
			CriterionCount:    len(input.Config.Config.CriteriaCatalog),
			AlternativeCount:  len(input.Config.Config.Alternatives),
			ScenarioCount:     len(input.Config.Config.Scenarios),
			ReportDefinitions: reportDefinitions,
		}),
		ReportDefinitions: domain.CanonicalReportDefinitions(reportDefinitions),
	}, nil
}

func validateLoadedConfig(config LoadedConfig) []domain.Diagnostic {
	var diagnostics []domain.Diagnostic

	diagnostics = append(diagnostics, validateRequiredSections(config)...)
	if config.Config == nil {
		return domain.CanonicalDiagnostics(diagnostics)
	}

	criteriaNames := collectUniqueNames(
		&diagnostics,
		"config.criteriaCatalog",
		config.Config.CriteriaCatalog,
		func(item CriterionConfig) string { return item.Name },
		"validation.duplicate_criterion_name",
		"duplicate criterion name: %s",
	)
	alternativeNames := collectUniqueNames(
		&diagnostics,
		"config.alternatives",
		config.Config.Alternatives,
		func(item AlternativeConfig) string { return item.Name },
		"validation.duplicate_alternative_name",
		"duplicate alternative name: %s",
	)
	scenarioNames := collectUniqueNames(
		&diagnostics,
		"config.scenarios",
		config.Config.Scenarios,
		func(item ScenarioConfig) string { return item.Name },
		"validation.duplicate_scenario_name",
		"duplicate scenario name: %s",
	)
	collectUniqueNames(
		&diagnostics,
		"config.reports",
		config.Config.Reports,
		func(item ReportConfig) string { return item.Name },
		"validation.duplicate_report_name",
		"duplicate report name: %s",
	)

	criteriaByName := validateCriteriaCatalog(&diagnostics, config.Config.CriteriaCatalog)
	scenarioInfos := validateScenarios(&diagnostics, config.Config.Scenarios, criteriaNames, criteriaByName)
	validateEvaluations(&diagnostics, config.Config.Evaluations, scenarioInfos, scenarioNames, alternativeNames, criteriaNames, criteriaByName)
	validateReports(&diagnostics, config.Config.Reports, scenarioNames, alternativeNames, criteriaNames)
	validateAggregation(&diagnostics, config.Config.Aggregation, scenarioNames)

	return domain.CanonicalDiagnostics(diagnostics)
}
