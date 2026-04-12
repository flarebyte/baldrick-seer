package pipeline

import (
	"strconv"
	"strings"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func renderMarkdownReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
) string {
	options := resolveMarkdownRenderOptions(report)
	return renderMarkdownRichReport(report, config, scenarioResults, finalRanking, aggregation, scenarioWeights, aggregationScenarioWeights, options)
}

type markdownRenderOptions struct {
	Detail                 string
	IncludeContext         bool
	IncludeWeights         bool
	IncludeAltDescriptions bool
	IncludeEvaluationNotes bool
	IncludeTradeoffs       bool
}

func resolveMarkdownRenderOptions(report ReportConfig) markdownRenderOptions {
	detail := reportArgumentValue(report.Arguments, "detail", "standard")
	options := markdownRenderOptions{
		Detail:                 detail,
		IncludeContext:         true,
		IncludeWeights:         false,
		IncludeAltDescriptions: true,
		IncludeEvaluationNotes: true,
		IncludeTradeoffs:       true,
	}

	switch detail {
	case "full":
		options.IncludeWeights = true
	case "standard":
		options.IncludeWeights = true
	case "brief":
		options.IncludeWeights = false
		options.IncludeTradeoffs = false
	default:
		options.IncludeWeights = true
	}

	if reportArgumentPresent(report.Arguments, "explain") && reportArgumentValue(report.Arguments, "explain", "true") == "false" {
		options.IncludeWeights = false
		options.IncludeTradeoffs = false
	}

	if reportArgumentPresent(report.Arguments, "include-context") {
		options.IncludeContext = reportArgumentValue(report.Arguments, "include-context", "false") == "true"
	}
	if reportArgumentPresent(report.Arguments, "include-weights") {
		options.IncludeWeights = reportArgumentValue(report.Arguments, "include-weights", "false") == "true"
	}
	if reportArgumentPresent(report.Arguments, "include-alternative-descriptions") {
		options.IncludeAltDescriptions = reportArgumentValue(report.Arguments, "include-alternative-descriptions", "false") == "true"
	}
	if reportArgumentPresent(report.Arguments, "include-evaluation-notes") {
		options.IncludeEvaluationNotes = reportArgumentValue(report.Arguments, "include-evaluation-notes", "false") == "true"
	}
	if reportArgumentPresent(report.Arguments, "include-tradeoffs") {
		options.IncludeTradeoffs = reportArgumentValue(report.Arguments, "include-tradeoffs", "false") == "true"
	}

	return options
}

func renderMarkdownRichReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
	options markdownRenderOptions,
) string {
	includeScores, topAlternatives := markdownReportSettings(report)
	orderedResults := orderedScenarioResultsForMarkdown(report, config, scenarioResults)
	orderedFinalRanking := limitRankedAlternatives(finalRanking.RankedAlternatives, topAlternatives)
	alternatives := filteredAlternativesForMarkdown(report, config)
	scenarios := filteredScenariosForMarkdown(report, config)
	criteria := filteredCriteriaForMarkdown(report, config)
	var builder strings.Builder

	writeMarkdownReportTitle(&builder, markdownDocumentTitle(report, config))

	if options.IncludeContext {
		writeMarkdownProblemSection(&builder, config)
		writeMarkdownAlternativesSection(&builder, alternatives, options.IncludeAltDescriptions)
		writeMarkdownScenariosSection(&builder, scenarios)
		writeMarkdownDecisionDriversSection(&builder, report, config, criteria, scenarioWeights, options.IncludeWeights)
	}

	writeMarkdownScenarioRankingsSection(&builder, report, config, orderedResults, criteria, includeScores, topAlternatives, options.IncludeEvaluationNotes)
	writeMarkdownFinalRankingSection(&builder, config, orderedFinalRanking, includeScores)

	if options.IncludeTradeoffs || options.IncludeWeights {
		writeMarkdownNotesAndTradeoffs(&builder, report, config, aggregation, scenarioWeights, aggregationScenarioWeights, orderedResults, options)
	}

	return builder.String()
}

func markdownReportSettings(report ReportConfig) (bool, int) {
	includeScores := reportArgumentValue(report.Arguments, "include-scores", "true") == "true"
	topAlternatives := reportArgumentInt(report.Arguments, "top-alternatives")
	return includeScores, topAlternatives
}

func markdownDocumentTitle(report ReportConfig, config *ExecutionConfig) string {
	if report.Title != "" {
		return report.Title
	}
	if config != nil && config.Problem != nil {
		if config.Problem.Title != "" {
			return config.Problem.Title
		}
		if config.Problem.Name != "" {
			return config.Problem.Name
		}
	}
	return ""
}

func writeMarkdownReportTitle(builder *strings.Builder, title string) {
	builder.WriteString("# ")
	builder.WriteString(title)
	builder.WriteString("\n")
}

func writeMarkdownProblemSection(builder *strings.Builder, config *ExecutionConfig) {
	if config == nil || config.Problem == nil {
		return
	}

	problem := config.Problem
	hasContent := problem.Title != "" || problem.Name != "" || problem.Goal != "" || problem.Description != "" || len(problem.Notes) > 0
	if !hasContent {
		return
	}

	builder.WriteString("\n## Problem\n")
	if title := firstNonEmpty(problem.Title, problem.Name); title != "" {
		builder.WriteString("\n### Title\n")
		builder.WriteString(title)
		builder.WriteString("\n")
	}
	if problem.Goal != "" {
		builder.WriteString("\n### Goal\n")
		builder.WriteString(problem.Goal)
		builder.WriteString("\n")
	}
	if problem.Description != "" {
		builder.WriteString("\n### Description\n")
		builder.WriteString(problem.Description)
		builder.WriteString("\n")
	}
	if len(problem.Notes) > 0 {
		builder.WriteString("\n### Notes\n")
		for _, note := range problem.Notes {
			builder.WriteString("- ")
			builder.WriteString(note)
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownAlternativesSection(builder *strings.Builder, alternatives []AlternativeConfig, includeDescriptions bool) {
	if len(alternatives) == 0 {
		return
	}

	builder.WriteString("\n## Alternatives\n")
	for _, alternative := range alternatives {
		builder.WriteString("\n### ")
		builder.WriteString(alternativeLabel(alternative))
		builder.WriteString("\n")
		if includeDescriptions && alternative.Description != "" {
			builder.WriteString(alternative.Description)
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownScenariosSection(builder *strings.Builder, scenarios []ScenarioConfig) {
	if len(scenarios) == 0 {
		return
	}

	sectionTitle := "## Scenario\n"
	if len(scenarios) > 1 {
		sectionTitle = "## Scenarios\n"
	}
	builder.WriteString("\n")
	builder.WriteString(sectionTitle)
	for _, scenario := range scenarios {
		builder.WriteString("\n### ")
		builder.WriteString(scenarioLabel(scenario))
		builder.WriteString("\n")
		if scenario.Description != "" {
			builder.WriteString(scenario.Description)
			builder.WriteString("\n")
		}
		if scenario.Narrative != "" {
			if scenario.Description != "" {
				builder.WriteString("\n")
			}
			builder.WriteString(scenario.Narrative)
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownDecisionDriversSection(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	criteria []CriterionConfig,
	scenarioWeights []ScenarioCriterionWeights,
	includeWeights bool,
) {
	hasCriteria := len(criteria) > 0
	comparisons := markdownPreferenceComparisons(report, config)
	weights := filterScenarioWeightsForReport(report, scenarioWeights)
	if !hasCriteria && len(comparisons) == 0 && (!includeWeights || len(weights) == 0) {
		return
	}

	builder.WriteString("\n## Decision Drivers\n")
	if hasCriteria {
		builder.WriteString("\n### Criteria\n")
		for _, criterion := range criteria {
			builder.WriteString("\n#### ")
			builder.WriteString(criterionLabel(criterion))
			builder.WriteString("\n")
			if criterion.Description != "" {
				builder.WriteString(criterion.Description)
				builder.WriteString("\n")
			}
		}
	}

	if len(comparisons) > 0 {
		builder.WriteString("\n### Preference Justifications\n")
		for _, comparison := range comparisons {
			builder.WriteString("- ")
			builder.WriteString(comparison)
			builder.WriteString("\n")
		}
	}

	if includeWeights && len(weights) > 0 {
		writeMarkdownScenarioWeightList(builder, "\n### Criteria Weights\n", config, weights)
	}
}

func writeMarkdownScenarioRankingsSection(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	criteria []CriterionConfig,
	includeScores bool,
	topAlternatives int,
	includeEvaluationNotes bool,
) {
	sectionTitle := "## Scenario Ranking\n"
	if len(scenarioResults) > 1 {
		sectionTitle = "## Scenario Rankings\n"
	}
	builder.WriteString("\n")
	builder.WriteString(sectionTitle)
	if len(scenarioResults) == 0 {
		builder.WriteString("\n(none)\n")
		return
	}

	for _, scenarioResult := range scenarioResults {
		builder.WriteString("\n### ")
		builder.WriteString(scenarioLabelByName(config, scenarioResult.ScenarioName))
		builder.WriteString("\n")
		rows := limitRankedAlternatives(scenarioResult.RankedAlternatives, topAlternatives)
		for _, alternative := range rows {
			writeMarkdownRankedAlternative(builder, config, alternative, includeScores)
		}
		if includeEvaluationNotes {
			writeMarkdownScenarioEvaluationNotes(builder, report, config, scenarioResult.ScenarioName, criteria)
		}
	}
}

func writeMarkdownFinalRankingSection(
	builder *strings.Builder,
	config *ExecutionConfig,
	rankedAlternatives []domain.RankedAlternative,
	includeScores bool,
) {
	builder.WriteString("\n## Final Ranking\n")
	if len(rankedAlternatives) == 0 {
		builder.WriteString("\n(none)\n")
		return
	}
	builder.WriteString("\n")
	for _, alternative := range rankedAlternatives {
		writeMarkdownRankedAlternative(builder, config, alternative, includeScores)
	}
}

func writeMarkdownNotesAndTradeoffs(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
	scenarioResults []domain.ScenarioRankingResult,
	options markdownRenderOptions,
) {
	builder.WriteString("\n## Notes and Tradeoffs\n")
	builder.WriteString("\n- Aggregation method: ")
	if aggregation != nil {
		builder.WriteString(aggregation.Method)
	}
	builder.WriteString("\n")

	if options.IncludeWeights {
		writeMarkdownAggregationWeightsSection(builder, config, aggregationScenarioWeights)
		writeMarkdownTradeoffCriteriaWeights(builder, report, config, scenarioWeights)
	}

	if options.IncludeTradeoffs {
		writeMarkdownExclusionSection(builder, config, scenarioResults)
	}
}

func writeMarkdownAggregationWeightsSection(builder *strings.Builder, config *ExecutionConfig, aggregationScenarioWeights map[string]float64) {
	if len(aggregationScenarioWeights) == 0 {
		return
	}
	builder.WriteString("\n### Scenario Weights\n")
	for _, scenarioName := range orderedWeightNames(aggregationScenarioWeights) {
		builder.WriteString("- ")
		builder.WriteString(scenarioLabelByName(config, scenarioName))
		builder.WriteString(": ")
		builder.WriteString(formatScore(aggregationScenarioWeights[scenarioName]))
		builder.WriteString("\n")
	}
}

func writeMarkdownTradeoffCriteriaWeights(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	scenarioWeights []ScenarioCriterionWeights,
) {
	weights := filterScenarioWeightsForReport(report, scenarioWeights)
	if len(weights) == 0 {
		return
	}
	writeMarkdownScenarioWeightList(builder, "\n### Scenario Criteria Weights\n", config, weights)
}

func writeMarkdownExclusionSection(builder *strings.Builder, config *ExecutionConfig, scenarioResults []domain.ScenarioRankingResult) {
	found := false
	for _, scenarioResult := range scenarioResults {
		for _, alternative := range scenarioResult.RankedAlternatives {
			if !alternative.Excluded {
				continue
			}
			if !found {
				builder.WriteString("\n### Exclusions\n")
				found = true
			}
			builder.WriteString("- ")
			builder.WriteString(scenarioLabelByName(config, scenarioResult.ScenarioName))
			builder.WriteString(": ")
			builder.WriteString(alternativeLabelByName(config, alternative.Name))
			if alternative.ExclusionReason != "" {
				builder.WriteString(" (")
				builder.WriteString(alternative.ExclusionReason)
				builder.WriteString(")")
			}
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownScenarioEvaluationNotes(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	scenarioName string,
	criteria []CriterionConfig,
) {
	evaluation, exists := findScenarioEvaluation(report, config, scenarioName)
	if !exists {
		return
	}

	builder.WriteString("\n#### Evaluation Notes\n")
	if evaluation.Description != "" {
		builder.WriteString(evaluation.Description)
		builder.WriteString("\n")
	}

	for _, alternative := range orderedAlternativeEvaluations(report, config, evaluation.Evaluations) {
		builder.WriteString("\n##### ")
		builder.WriteString(alternativeLabelByName(config, alternative.AlternativeName))
		builder.WriteString("\n")
		if alternative.Description != "" {
			builder.WriteString(alternative.Description)
			builder.WriteString("\n")
		}

		orderedValues := orderedCriterionValueRecords(criteria, alternative.Values)
		if len(orderedValues) == 0 {
			continue
		}

		builder.WriteString("\nScores:\n")
		for _, value := range orderedValues {
			builder.WriteString("- ")
			builder.WriteString(value.Label)
			builder.WriteString(": ")
			builder.WriteString(value.Rendered)
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownRankedAlternative(builder *strings.Builder, config *ExecutionConfig, alternative domain.RankedAlternative, includeScores bool) {
	if alternative.Excluded {
		builder.WriteString("- ")
		builder.WriteString(alternativeLabelByName(config, alternative.Name))
		builder.WriteString(": excluded")
		if alternative.ExclusionReason != "" {
			builder.WriteString(" (")
			builder.WriteString(alternative.ExclusionReason)
			builder.WriteString(")")
		}
		builder.WriteString("\n")
		return
	}

	builder.WriteString(strconv.Itoa(alternative.Rank))
	builder.WriteString(". ")
	builder.WriteString(alternativeLabelByName(config, alternative.Name))
	if includeScores {
		builder.WriteString(" (")
		builder.WriteString(formatScore(alternative.Score))
		builder.WriteString(")")
	}
	builder.WriteString("\n")
}

func writeMarkdownNamedWeights(builder *strings.Builder, config *ExecutionConfig, weights []CriterionWeight) {
	for index, weight := range weights {
		if index > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(criterionLabelByName(config, weight.CriterionName))
		builder.WriteString("=")
		builder.WriteString(formatScore(weight.Weight))
	}
}

func writeMarkdownScenarioWeightList(
	builder *strings.Builder,
	heading string,
	config *ExecutionConfig,
	scenarioWeights []ScenarioCriterionWeights,
) {
	builder.WriteString(heading)
	for _, scenarioWeight := range scenarioWeights {
		builder.WriteString("- ")
		builder.WriteString(scenarioLabelByName(config, scenarioWeight.ScenarioName))
		builder.WriteString(": ")
		writeMarkdownNamedWeights(builder, config, canonicalCriterionWeights(scenarioWeight.CriterionWeights))
		builder.WriteString("\n")
	}
}

func markdownPreferenceComparisons(report ReportConfig, config *ExecutionConfig) []string {
	if config == nil {
		return nil
	}

	allowedScenarios, allowedCriteria := focusedScenarioAndCriterionNames(report)

	var output []string
	for _, scenario := range config.Scenarios {
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[scenario.Name]; !exists {
				continue
			}
		}
		if scenario.Preferences == nil || len(scenario.Preferences.Comparisons) == 0 {
			continue
		}
		for _, comparison := range scenario.Preferences.Comparisons {
			if len(allowedCriteria) > 0 {
				if _, exists := allowedCriteria[comparison.MoreImportantCriterionName]; !exists {
					continue
				}
				if _, exists := allowedCriteria[comparison.LessImportantCriterionName]; !exists {
					continue
				}
			}
			var builder strings.Builder
			builder.WriteString(scenarioLabel(scenario))
			builder.WriteString(": ")
			builder.WriteString(criterionLabelByName(config, comparison.MoreImportantCriterionName))
			builder.WriteString(" over ")
			builder.WriteString(criterionLabelByName(config, comparison.LessImportantCriterionName))
			builder.WriteString(" (strength ")
			builder.WriteString(formatScore(comparison.Strength))
			builder.WriteString(")")
			output = append(output, builder.String())
		}
	}
	return output
}

func focusedScenarioAndCriterionNames(report ReportConfig) (map[string]struct{}, map[string]struct{}) {
	return allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
			return focus.ScenarioNames
		}),
		allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
			return focus.CriterionNames
		})
}

func orderedScenarioResultsForMarkdown(report ReportConfig, config *ExecutionConfig, scenarioResults []domain.ScenarioRankingResult) []domain.ScenarioRankingResult {
	resultByName := make(map[string]domain.ScenarioRankingResult, len(scenarioResults))
	for _, result := range scenarioResults {
		resultByName[result.ScenarioName] = result
	}

	ordered := make([]domain.ScenarioRankingResult, 0, len(scenarioResults))
	for _, scenario := range filteredScenariosForMarkdown(report, config) {
		result, exists := resultByName[scenario.Name]
		if !exists {
			continue
		}
		ordered = append(ordered, result)
		delete(resultByName, scenario.Name)
	}

	for _, result := range domain.CanonicalScenarioResults(scenarioResults) {
		if _, exists := resultByName[result.ScenarioName]; !exists {
			continue
		}
		ordered = append(ordered, result)
	}

	return ordered
}

func filteredAlternativesForMarkdown(report ReportConfig, config *ExecutionConfig) []AlternativeConfig {
	if config == nil {
		return nil
	}
	allowedAlternatives := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.AlternativeNames
	})

	output := make([]AlternativeConfig, 0, len(config.Alternatives))
	for _, alternative := range config.Alternatives {
		if len(allowedAlternatives) > 0 {
			if _, exists := allowedAlternatives[alternative.Name]; !exists {
				continue
			}
		}
		output = append(output, alternative)
	}
	return output
}

func filteredScenariosForMarkdown(report ReportConfig, config *ExecutionConfig) []ScenarioConfig {
	if config == nil {
		return nil
	}
	allowedScenarios := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.ScenarioNames
	})

	output := make([]ScenarioConfig, 0, len(config.Scenarios))
	for _, scenario := range config.Scenarios {
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[scenario.Name]; !exists {
				continue
			}
		}
		output = append(output, scenario)
	}
	return output
}

func filteredCriteriaForMarkdown(report ReportConfig, config *ExecutionConfig) []CriterionConfig {
	if config == nil {
		return nil
	}
	allowedCriteria := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.CriterionNames
	})

	output := make([]CriterionConfig, 0, len(config.CriteriaCatalog))
	for _, criterion := range config.CriteriaCatalog {
		if len(allowedCriteria) > 0 {
			if _, exists := allowedCriteria[criterion.Name]; !exists {
				continue
			}
		}
		output = append(output, criterion)
	}
	return output
}

func findScenarioEvaluation(report ReportConfig, config *ExecutionConfig, scenarioName string) (EvaluationConfig, bool) {
	if config == nil {
		return EvaluationConfig{}, false
	}
	allowedScenarios := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.ScenarioNames
	})
	for _, evaluation := range config.Evaluations {
		if evaluation.ScenarioName != scenarioName {
			continue
		}
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[evaluation.ScenarioName]; !exists {
				return EvaluationConfig{}, false
			}
		}
		return evaluation, true
	}
	return EvaluationConfig{}, false
}

func orderedAlternativeEvaluations(report ReportConfig, config *ExecutionConfig, evaluations []AlternativeEvaluationConfig) []AlternativeEvaluationConfig {
	allowedAlternatives := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.AlternativeNames
	})
	evaluationByName := make(map[string]AlternativeEvaluationConfig, len(evaluations))
	for _, evaluation := range evaluations {
		if len(allowedAlternatives) > 0 {
			if _, exists := allowedAlternatives[evaluation.AlternativeName]; !exists {
				continue
			}
		}
		evaluationByName[evaluation.AlternativeName] = evaluation
	}

	ordered := make([]AlternativeEvaluationConfig, 0, len(evaluationByName))
	for _, alternative := range filteredAlternativesForMarkdown(report, config) {
		evaluation, exists := evaluationByName[alternative.Name]
		if !exists {
			continue
		}
		ordered = append(ordered, evaluation)
		delete(evaluationByName, alternative.Name)
	}

	for _, evaluation := range canonicalAlternativeEvaluations(evaluations) {
		if _, exists := evaluationByName[evaluation.AlternativeName]; !exists {
			continue
		}
		ordered = append(ordered, evaluation)
	}
	return ordered
}

type markdownCriterionValueRecord struct {
	Name     string
	Label    string
	Rendered string
}

func orderedCriterionValueRecords(criteria []CriterionConfig, values map[string]CriterionValue) []markdownCriterionValueRecord {
	if len(values) == 0 {
		return nil
	}

	ordered := make([]markdownCriterionValueRecord, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, criterion := range criteria {
		value, exists := values[criterion.Name]
		if !exists {
			continue
		}
		ordered = append(ordered, markdownCriterionValueRecord{
			Name:     criterion.Name,
			Label:    criterionLabel(criterion),
			Rendered: renderCriterionValue(value.Value),
		})
		seen[criterion.Name] = struct{}{}
	}

	for _, criterionName := range domain.CanonicalNames(valueNames(values)) {
		if _, exists := seen[criterionName]; exists {
			continue
		}
		ordered = append(ordered, markdownCriterionValueRecord{
			Name:     criterionName,
			Label:    criterionName,
			Rendered: renderCriterionValue(values[criterionName].Value),
		})
	}
	return ordered
}

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
