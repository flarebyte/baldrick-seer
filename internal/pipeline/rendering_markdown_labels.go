package pipeline

func alternativeLabelByName(config *ExecutionConfig, name string) string {
	for _, alternative := range filteredAlternativesForMarkdown(ReportConfig{}, config) {
		if alternative.Name == name {
			return alternativeLabel(alternative)
		}
	}
	return name
}

func alternativeTitleByName(config *ExecutionConfig, name string) string {
	if config == nil {
		return ""
	}
	for _, alternative := range config.Alternatives {
		if alternative.Name == name {
			return alternative.Title
		}
	}
	return ""
}

func criterionLabelByName(config *ExecutionConfig, name string) string {
	if config == nil {
		return name
	}
	for _, criterion := range config.CriteriaCatalog {
		if criterion.Name == name {
			return criterionLabel(criterion)
		}
	}
	return name
}

func scenarioLabelByName(config *ExecutionConfig, name string) string {
	if config == nil {
		return name
	}
	for _, scenario := range config.Scenarios {
		if scenario.Name == name {
			return scenarioLabel(scenario)
		}
	}
	return name
}

func scenarioTitleByName(config *ExecutionConfig, name string) string {
	if config == nil {
		return ""
	}
	for _, scenario := range config.Scenarios {
		if scenario.Name == name {
			return scenario.Title
		}
	}
	return ""
}

func criterionTitleByName(config *ExecutionConfig, name string) string {
	if config == nil {
		return ""
	}
	for _, criterion := range config.CriteriaCatalog {
		if criterion.Name == name {
			return criterion.Title
		}
	}
	return ""
}

func alternativeLabel(alternative AlternativeConfig) string {
	return firstNonEmpty(alternative.Title, alternative.Name)
}

func criterionLabel(criterion CriterionConfig) string {
	return firstNonEmpty(criterion.Title, criterion.Name)
}

func scenarioLabel(scenario ScenarioConfig) string {
	return firstNonEmpty(scenario.Title, scenario.Name)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
