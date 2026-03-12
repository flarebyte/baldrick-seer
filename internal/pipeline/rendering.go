package pipeline

import "context"

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func (DefaultReportRenderer) RenderReports(ctx context.Context, input RenderReportsInput) (RenderReportsOutput, error) {
	if err := checkContext(ctx, input.Command.ConfigPath); err != nil {
		return RenderReportsOutput{}, err
	}

	if input.Config.Config == nil {
		return RenderReportsOutput{}, NewRenderingFailure("rendering.config_missing", input.Command.ConfigPath, "reports could not be rendered", ErrRenderingFailed)
	}

	reportByName := make(map[string]ReportConfig, len(input.Config.Config.Reports))
	for _, report := range input.Config.Config.Reports {
		reportByName[report.Name] = report
	}

	orderedDefinitions := domain.CanonicalReportDefinitions(input.ReportDefinitions)
	renderedParts := make([]string, 0, len(orderedDefinitions))
	for _, reportDefinition := range orderedDefinitions {
		reportConfig, exists := reportByName[reportDefinition.Name]
		if !exists {
			return RenderReportsOutput{}, NewRenderingFailure("rendering.report_missing", input.Command.ConfigPath, "reports could not be rendered", ErrRenderingFailed)
		}

		rendered, err := renderReport(
			reportConfig,
			input.Config.Config,
			input.ScenarioResults,
			input.Config.Config.Aggregation,
			input.ScenarioWeights,
		)
		if err != nil {
			return RenderReportsOutput{}, NewRenderingFailure("rendering.report_failed", input.Command.ConfigPath, "reports could not be rendered", ErrRenderingFailed)
		}
		renderedParts = append(renderedParts, rendered)
	}

	return RenderReportsOutput{
		ReportDefinitions: orderedDefinitions,
		RenderedOutput:    strings.Join(renderedParts, "\n"),
	}, nil
}

func renderReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
) (string, error) {
	filteredScenarioResults := filterScenarioResultsForReport(report, scenarioResults)
	displayScenarioResults := filterScenarioAlternativesForReport(report, filteredScenarioResults)
	filteredFinalRanking, err := aggregateScenarioResults(aggregation, filteredScenarioResults)
	if err != nil {
		return "", err
	}
	filteredFinalRanking = filterFinalRankingForReport(report, filteredFinalRanking)
	filteredScenarioWeights := filterScenarioWeightsForReport(report, scenarioWeights)

	switch report.Format {
	case "markdown":
		return renderMarkdownReport(report, config, displayScenarioResults, filteredFinalRanking), nil
	case "json":
		return renderJSONReport(report, config, displayScenarioResults, filteredFinalRanking, aggregation, filteredScenarioWeights)
	case "csv":
		return renderCSVReport(report, config, displayScenarioResults, filteredFinalRanking)
	default:
		return "", ErrRenderingFailed
	}
}

func renderMarkdownReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
) string {
	includeScores := reportArgumentValue(report.Arguments, "include-scores", "true") == "true"
	topAlternatives := reportArgumentInt(report.Arguments, "top-alternatives")

	var builder strings.Builder
	builder.WriteString("# ")
	builder.WriteString(report.Title)
	builder.WriteString("\n\n")
	builder.WriteString("Problem: ")
	builder.WriteString(problemName(config))
	builder.WriteString("\n\n")

	builder.WriteString("## Scenarios\n")
	for _, scenarioResult := range scenarioResults {
		builder.WriteString("\n### ")
		builder.WriteString(scenarioResult.ScenarioName)
		builder.WriteString("\n")
		rows := limitRankedAlternatives(scenarioResult.RankedAlternatives, topAlternatives)
		for _, alternative := range rows {
			writeMarkdownAlternative(&builder, alternative, includeScores)
		}
	}

	builder.WriteString("\n## Final Ranking\n")
	if len(finalRanking.RankedAlternatives) == 0 {
		builder.WriteString("\n(none)\n")
		return builder.String()
	}

	for _, alternative := range limitRankedAlternatives(finalRanking.RankedAlternatives, topAlternatives) {
		writeMarkdownAlternative(&builder, alternative, includeScores)
	}

	return builder.String()
}

func renderJSONReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
) (string, error) {
	type jsonRankedAlternative struct {
		AlternativeName string   `json:"alternativeName"`
		Rank            *int     `json:"rank,omitempty"`
		Score           *float64 `json:"score,omitempty"`
		Excluded        bool     `json:"excluded,omitempty"`
	}
	type jsonScenarioResult struct {
		ScenarioName string                  `json:"scenarioName"`
		Ranking      []jsonRankedAlternative `json:"ranking"`
	}
	type jsonCriterionWeightEntry struct {
		CriterionName string  `json:"criterionName"`
		Weight        float64 `json:"weight"`
	}
	type jsonScenarioWeights struct {
		ScenarioName string                     `json:"scenarioName"`
		Weights      []jsonCriterionWeightEntry `json:"weights,omitempty"`
	}
	type jsonAggregation struct {
		Method string `json:"method"`
	}
	type jsonReport struct {
		ProblemName     string                  `json:"problemName"`
		ReportName      string                  `json:"reportName"`
		Format          string                  `json:"format"`
		Aggregation     jsonAggregation         `json:"aggregation"`
		ScenarioResults []jsonScenarioResult    `json:"scenarioResults"`
		ScenarioWeights []jsonScenarioWeights   `json:"scenarioWeights,omitempty"`
		FinalRanking    []jsonRankedAlternative `json:"finalRanking"`
	}

	payload := jsonReport{
		ProblemName: problemName(config),
		ReportName:  report.Name,
		Format:      report.Format,
	}
	if aggregation != nil {
		payload.Aggregation = jsonAggregation{Method: aggregation.Method}
	}
	for _, scenarioResult := range scenarioResults {
		jsonResult := jsonScenarioResult{
			ScenarioName: scenarioResult.ScenarioName,
		}
		for _, alternative := range scenarioResult.RankedAlternatives {
			entry := jsonRankedAlternative{
				AlternativeName: alternative.Name,
				Excluded:        alternative.Excluded,
			}
			if !alternative.Excluded {
				entry.Rank = intPointer(alternative.Rank)
				entry.Score = floatPointer(alternative.Score)
			}
			jsonResult.Ranking = append(jsonResult.Ranking, entry)
		}
		payload.ScenarioResults = append(payload.ScenarioResults, jsonResult)
	}

	if reportArgumentValue(report.Arguments, "include-weights", "false") == "true" {
		for _, scenarioWeight := range canonicalScenarioWeights(scenarioWeights) {
			entry := jsonScenarioWeights{ScenarioName: scenarioWeight.ScenarioName}
			for _, weight := range scenarioWeight.CriterionWeights {
				entry.Weights = append(entry.Weights, jsonCriterionWeightEntry(weight))
			}
			payload.ScenarioWeights = append(payload.ScenarioWeights, entry)
		}
	}

	for _, alternative := range finalRanking.RankedAlternatives {
		payload.FinalRanking = append(payload.FinalRanking, jsonRankedAlternative{
			AlternativeName: alternative.Name,
			Rank:            intPointer(alternative.Rank),
			Score:           floatPointer(alternative.Score),
		})
	}

	if reportArgumentValue(report.Arguments, "pretty", "false") == "true" {
		content, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return "", err
		}
		return string(content) + "\n", nil
	}

	content, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(content) + "\n", nil
}

func renderCSVReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
) (string, error) {
	columnsValue := reportArgumentValue(report.Arguments, "columns", "scenario,alternative,score,rank")
	columns := strings.Split(columnsValue, ",")
	includeHeader := reportArgumentValue(report.Arguments, "header", "true") == "true"
	includeCriterionRows := csvColumnsIncludeCriterionData(columns)
	valueLookup := buildEvaluationValueLookup(config, report)

	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)
	if includeHeader {
		if err := writer.Write(columns); err != nil {
			return "", err
		}
	}

	for _, scenarioResult := range scenarioResults {
		for _, alternative := range scenarioResult.RankedAlternatives {
			records := csvRecords(columns, scenarioResult.ScenarioName, alternative, valueLookup[scenarioResult.ScenarioName][alternative.Name], includeCriterionRows)
			for _, record := range records {
				if err := writer.Write(record); err != nil {
					return "", err
				}
			}
		}
	}
	for _, alternative := range finalRanking.RankedAlternatives {
		record := csvRecord(columns, "overall", alternative, criterionValueRecord{})
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}
	writer.Flush()
	return buffer.String(), writer.Error()
}

func csvColumnsIncludeCriterionData(columns []string) bool {
	for _, column := range columns {
		if column == "criterion" || column == "value" {
			return true
		}
	}
	return false
}

func csvRecords(
	columns []string,
	scenarioName string,
	alternative domain.RankedAlternative,
	values []criterionValueRecord,
	includeCriterionRows bool,
) [][]string {
	if !includeCriterionRows {
		return [][]string{csvRecord(columns, scenarioName, alternative, criterionValueRecord{})}
	}
	if len(values) == 0 {
		return [][]string{csvRecord(columns, scenarioName, alternative, criterionValueRecord{})}
	}

	records := make([][]string, 0, len(values))
	for _, value := range values {
		records = append(records, csvRecord(columns, scenarioName, alternative, value))
	}
	return records
}

func csvRecord(columns []string, scenarioName string, alternative domain.RankedAlternative, value criterionValueRecord) []string {
	record := make([]string, 0, len(columns))
	for _, column := range columns {
		switch column {
		case "scenario":
			record = append(record, scenarioName)
		case "alternative":
			record = append(record, alternative.Name)
		case "score":
			if alternative.Excluded {
				record = append(record, "")
			} else {
				record = append(record, formatScore(alternative.Score))
			}
		case "rank":
			if alternative.Excluded {
				record = append(record, "")
			} else {
				record = append(record, strconv.Itoa(alternative.Rank))
			}
		case "criterion":
			record = append(record, value.Name)
		case "value":
			record = append(record, value.Rendered)
		default:
			record = append(record, "")
		}
	}
	return record
}

func filterScenarioResultsForReport(report ReportConfig, scenarioResults []domain.ScenarioRankingResult) []domain.ScenarioRankingResult {
	if report.Focus == nil || len(report.Focus.ScenarioNames) == 0 {
		return domain.CanonicalScenarioResults(scenarioResults)
	}

	allowed := make(map[string]struct{}, len(report.Focus.ScenarioNames))
	for _, scenarioName := range report.Focus.ScenarioNames {
		allowed[scenarioName] = struct{}{}
	}

	var filtered []domain.ScenarioRankingResult
	for _, scenarioResult := range domain.CanonicalScenarioResults(scenarioResults) {
		if _, exists := allowed[scenarioResult.ScenarioName]; exists {
			filtered = append(filtered, scenarioResult)
		}
	}
	return filtered
}

func filterScenarioAlternativesForReport(report ReportConfig, scenarioResults []domain.ScenarioRankingResult) []domain.ScenarioRankingResult {
	if report.Focus == nil || len(report.Focus.AlternativeNames) == 0 {
		return domain.CanonicalScenarioResults(scenarioResults)
	}

	allowed := toAllowedNameSet(report.Focus.AlternativeNames)
	filtered := make([]domain.ScenarioRankingResult, 0, len(scenarioResults))
	for _, scenarioResult := range domain.CanonicalScenarioResults(scenarioResults) {
		alternatives := make([]domain.RankedAlternative, 0, len(scenarioResult.RankedAlternatives))
		for _, alternative := range domain.CanonicalRankedAlternatives(scenarioResult.RankedAlternatives) {
			if _, exists := allowed[alternative.Name]; exists {
				alternatives = append(alternatives, alternative)
			}
		}
		filtered = append(filtered, domain.ScenarioRankingResult{
			ScenarioName:       scenarioResult.ScenarioName,
			RankedAlternatives: alternatives,
		})
	}
	return domain.CanonicalScenarioResults(filtered)
}

func filterFinalRankingForReport(report ReportConfig, finalRanking domain.AggregatedRankingResult) domain.AggregatedRankingResult {
	if report.Focus == nil || len(report.Focus.AlternativeNames) == 0 {
		return domain.CanonicalAggregatedRankingResult(finalRanking)
	}

	allowed := toAllowedNameSet(report.Focus.AlternativeNames)

	var filtered []domain.RankedAlternative
	for _, alternative := range domain.CanonicalAggregatedRankingResult(finalRanking).RankedAlternatives {
		if _, exists := allowed[alternative.Name]; exists {
			filtered = append(filtered, alternative)
		}
	}
	for index := range filtered {
		filtered[index].Rank = index + 1
	}
	return domain.AggregatedRankingResult{RankedAlternatives: filtered}
}

func filterScenarioWeightsForReport(report ReportConfig, scenarioWeights []ScenarioCriterionWeights) []ScenarioCriterionWeights {
	if len(scenarioWeights) == 0 {
		return nil
	}

	allowedScenarios := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.ScenarioNames
	})
	allowedCriteria := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.CriterionNames
	})

	filtered := make([]ScenarioCriterionWeights, 0, len(scenarioWeights))
	for _, scenarioWeight := range canonicalScenarioWeights(scenarioWeights) {
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[scenarioWeight.ScenarioName]; !exists {
				continue
			}
		}

		weights := make([]CriterionWeight, 0, len(scenarioWeight.CriterionWeights))
		for _, weight := range canonicalCriterionWeights(scenarioWeight.CriterionWeights) {
			if len(allowedCriteria) > 0 {
				if _, exists := allowedCriteria[weight.CriterionName]; !exists {
					continue
				}
			}
			weights = append(weights, weight)
		}

		filtered = append(filtered, ScenarioCriterionWeights{
			ScenarioName:     scenarioWeight.ScenarioName,
			CriterionWeights: weights,
		})
	}

	return canonicalScenarioWeights(filtered)
}

type criterionValueRecord struct {
	Name     string
	Rendered string
}

func buildEvaluationValueLookup(config *ExecutionConfig, report ReportConfig) map[string]map[string][]criterionValueRecord {
	if config == nil || len(config.Evaluations) == 0 {
		return nil
	}

	allowedScenarios := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.ScenarioNames
	})
	allowedAlternatives := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.AlternativeNames
	})
	allowedCriteria := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.CriterionNames
	})

	output := make(map[string]map[string][]criterionValueRecord, len(config.Evaluations))
	for _, evaluation := range canonicalEvaluations(config.Evaluations) {
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[evaluation.ScenarioName]; !exists {
				continue
			}
		}

		alternatives := make(map[string][]criterionValueRecord, len(evaluation.Evaluations))
		for _, alternative := range canonicalAlternativeEvaluations(evaluation.Evaluations) {
			if len(allowedAlternatives) > 0 {
				if _, exists := allowedAlternatives[alternative.AlternativeName]; !exists {
					continue
				}
			}
			values := make([]criterionValueRecord, 0, len(alternative.Values))
			for _, criterionName := range domain.CanonicalNames(valueNames(alternative.Values)) {
				if len(allowedCriteria) > 0 {
					if _, exists := allowedCriteria[criterionName]; !exists {
						continue
					}
				}
				values = append(values, criterionValueRecord{
					Name:     criterionName,
					Rendered: renderCriterionValue(alternative.Values[criterionName].Value),
				})
			}
			alternatives[alternative.AlternativeName] = values
		}
		output[evaluation.ScenarioName] = alternatives
	}

	return output
}

func canonicalEvaluations(input []EvaluationConfig) []EvaluationConfig {
	if len(input) == 0 {
		return nil
	}

	output := append([]EvaluationConfig(nil), input...)
	for i := range output {
		output[i].Evaluations = canonicalAlternativeEvaluations(output[i].Evaluations)
	}
	for i := 1; i < len(output); i++ {
		current := output[i]
		j := i - 1
		for ; j >= 0 && output[j].ScenarioName > current.ScenarioName; j-- {
			output[j+1] = output[j]
		}
		output[j+1] = current
	}
	return output
}

func valueNames(values map[string]CriterionValue) []string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	return names
}

func renderCriterionValue(value any) string {
	switch typed := value.(type) {
	case bool:
		if typed {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(typed)
	case int8, int16, int32, int64, float32, float64:
		number, err := numericValue(typed)
		if err != nil {
			return ""
		}
		return strconv.FormatFloat(number, 'f', -1, 64)
	default:
		return ""
	}
}

func allowedFocusedNames(focus *ReportFocus, selectNames func(*ReportFocus) []string) map[string]struct{} {
	if focus == nil {
		return nil
	}
	names := selectNames(focus)
	if len(names) == 0 {
		return nil
	}
	return toAllowedNameSet(names)
}

func toAllowedNameSet(names []string) map[string]struct{} {
	allowed := make(map[string]struct{}, len(names))
	for _, name := range names {
		allowed[name] = struct{}{}
	}
	return allowed
}

func reportArgumentValue(arguments []string, key string, fallback string) string {
	for _, argument := range arguments {
		parsedKey, parsedValue, ok := strings.Cut(argument, "=")
		if ok && parsedKey == key {
			return parsedValue
		}
	}
	return fallback
}

func reportArgumentInt(arguments []string, key string) int {
	value := reportArgumentValue(arguments, key, "")
	if value == "" {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}

func limitRankedAlternatives(input []domain.RankedAlternative, top int) []domain.RankedAlternative {
	if top <= 0 {
		return domain.CanonicalRankedAlternatives(input)
	}

	var ranked []domain.RankedAlternative
	var excluded []domain.RankedAlternative
	for _, alternative := range domain.CanonicalRankedAlternatives(input) {
		if alternative.Excluded {
			excluded = append(excluded, alternative)
			continue
		}
		if len(ranked) < top {
			ranked = append(ranked, alternative)
		}
	}
	return append(ranked, excluded...)
}

func formatScore(score float64) string {
	return strconv.FormatFloat(score, 'f', 6, 64)
}

func intPointer(value int) *int {
	return &value
}

func floatPointer(value float64) *float64 {
	return &value
}

func canonicalScenarioWeights(input []ScenarioCriterionWeights) []ScenarioCriterionWeights {
	if len(input) == 0 {
		return nil
	}

	output := append([]ScenarioCriterionWeights(nil), input...)
	for index := range output {
		output[index].CriterionWeights = canonicalCriterionWeights(output[index].CriterionWeights)
	}
	return sortScenarioWeightsByName(output)
}

func sortScenarioWeightsByName(input []ScenarioCriterionWeights) []ScenarioCriterionWeights {
	if len(input) == 0 {
		return nil
	}

	output := append([]ScenarioCriterionWeights(nil), input...)
	for i := 1; i < len(output); i++ {
		current := output[i]
		j := i - 1
		for ; j >= 0 && output[j].ScenarioName > current.ScenarioName; j-- {
			output[j+1] = output[j]
		}
		output[j+1] = current
	}
	return output
}

func problemName(config *ExecutionConfig) string {
	if config == nil || config.Problem == nil || config.Problem.Name == "" {
		return ""
	}
	return config.Problem.Name
}

func writeMarkdownAlternative(builder *strings.Builder, alternative domain.RankedAlternative, includeScores bool) {
	builder.WriteString("- ")
	if alternative.Excluded {
		builder.WriteString(alternative.Name)
		builder.WriteString(": excluded\n")
		return
	}

	builder.WriteString(strconv.Itoa(alternative.Rank))
	builder.WriteString(". ")
	builder.WriteString(alternative.Name)
	if includeScores {
		builder.WriteString(" (")
		builder.WriteString(formatScore(alternative.Score))
		builder.WriteString(")")
	}
	builder.WriteString("\n")
}
