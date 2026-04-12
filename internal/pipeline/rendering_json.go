package pipeline

import (
	"encoding/json"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type jsonRankedAlternative struct {
	AlternativeName string   `json:"alternativeName"`
	Rank            *int     `json:"rank,omitempty"`
	Score           *float64 `json:"score,omitempty"`
	Excluded        bool     `json:"excluded,omitempty"`
	ExclusionReason string   `json:"exclusionReason,omitempty"`
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

type jsonAggregationWeight struct {
	ScenarioName string  `json:"scenarioName"`
	Weight       float64 `json:"weight"`
}

type jsonAggregation struct {
	Method          string                  `json:"method"`
	ScenarioWeights []jsonAggregationWeight `json:"scenarioWeights,omitempty"`
}

type jsonProblemContext struct {
	Name        string   `json:"name"`
	Title       string   `json:"title,omitempty"`
	Goal        string   `json:"goal,omitempty"`
	Description string   `json:"description,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	Notes       []string `json:"notes,omitempty"`
}

type jsonReportContext struct {
	Name      string   `json:"name"`
	Title     string   `json:"title,omitempty"`
	Format    string   `json:"format"`
	Arguments []string `json:"arguments,omitempty"`
}

type jsonAlternativeContext struct {
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type jsonCriterionContext struct {
	Name          string `json:"name"`
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Polarity      string `json:"polarity,omitempty"`
	ValueType     string `json:"valueType,omitempty"`
	ScaleGuidance []any  `json:"scaleGuidance,omitempty"`
}

type jsonScenarioContext struct {
	Name           string   `json:"name"`
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	Narrative      string   `json:"narrative,omitempty"`
	ActiveCriteria []string `json:"activeCriteria,omitempty"`
}

type jsonAlternativeEvaluationContext struct {
	AlternativeName string `json:"alternativeName"`
	Description     string `json:"description,omitempty"`
}

type jsonEvaluationContext struct {
	ScenarioName string                             `json:"scenarioName"`
	Description  string                             `json:"description,omitempty"`
	Evaluations  []jsonAlternativeEvaluationContext `json:"evaluations,omitempty"`
}

type jsonReport struct {
	ProblemName     string                   `json:"problemName"`
	ReportName      string                   `json:"reportName"`
	Format          string                   `json:"format"`
	Aggregation     jsonAggregation          `json:"aggregation"`
	ScenarioResults []jsonScenarioResult     `json:"scenarioResults"`
	ScenarioWeights []jsonScenarioWeights    `json:"scenarioWeights,omitempty"`
	FinalRanking    []jsonRankedAlternative  `json:"finalRanking"`
	Problem         *jsonProblemContext      `json:"problem,omitempty"`
	Report          *jsonReportContext       `json:"report,omitempty"`
	Alternatives    []jsonAlternativeContext `json:"alternatives,omitempty"`
	Criteria        []jsonCriterionContext   `json:"criteria,omitempty"`
	Scenarios       []jsonScenarioContext    `json:"scenarios,omitempty"`
	Evaluations     []jsonEvaluationContext  `json:"evaluations,omitempty"`
}

func renderJSONReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
) (string, error) {
	payload := jsonReport{
		ProblemName: problemName(config),
		ReportName:  report.Name,
		Format:      report.Format,
	}
	if aggregation != nil {
		payload.Aggregation = jsonAggregation{Method: aggregation.Method}
		for _, scenarioName := range orderedWeightNames(aggregationScenarioWeights) {
			payload.Aggregation.ScenarioWeights = append(payload.Aggregation.ScenarioWeights, jsonAggregationWeight{
				ScenarioName: scenarioName,
				Weight:       aggregationScenarioWeights[scenarioName],
			})
		}
	}
	for _, scenarioResult := range scenarioResults {
		jsonResult := jsonScenarioResult{
			ScenarioName: scenarioResult.ScenarioName,
		}
		for _, alternative := range scenarioResult.RankedAlternatives {
			entry := jsonRankedAlternative{
				AlternativeName: alternative.Name,
				Excluded:        alternative.Excluded,
				ExclusionReason: alternative.ExclusionReason,
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

	if reportArgumentValue(report.Arguments, "include-context", "false") == "true" {
		payload.Problem = buildJSONProblemContext(config)
		payload.Report = buildJSONReportContext(report)
		payload.Alternatives = buildJSONAlternatives(config)
		payload.Criteria = buildJSONCriteria(config)
		payload.Scenarios = buildJSONScenarios(config)
		payload.Evaluations = buildJSONEvaluations(config)
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

func buildJSONProblemContext(config *ExecutionConfig) *jsonProblemContext {
	if config == nil || config.Problem == nil {
		return nil
	}
	return &jsonProblemContext{
		Name:        config.Problem.Name,
		Title:       config.Problem.Title,
		Goal:        config.Problem.Goal,
		Description: config.Problem.Description,
		Owner:       config.Problem.Owner,
		Notes:       append([]string(nil), config.Problem.Notes...),
	}
}

func buildJSONReportContext(report ReportConfig) *jsonReportContext {
	return &jsonReportContext{
		Name:      report.Name,
		Title:     report.Title,
		Format:    report.Format,
		Arguments: append([]string(nil), report.Arguments...),
	}
}

func buildJSONAlternatives(config *ExecutionConfig) []jsonAlternativeContext {
	if config == nil {
		return nil
	}
	alternatives := canonicalAlternatives(config.Alternatives)
	output := make([]jsonAlternativeContext, 0, len(alternatives))
	for _, alternative := range alternatives {
		output = append(output, jsonAlternativeContext(alternative))
	}
	return output
}

func buildJSONCriteria(config *ExecutionConfig) []jsonCriterionContext {
	if config == nil {
		return nil
	}
	criteria := canonicalCriteria(config.CriteriaCatalog)
	output := make([]jsonCriterionContext, 0, len(criteria))
	for _, criterion := range criteria {
		output = append(output, jsonCriterionContext{
			Name:          criterion.Name,
			Title:         criterion.Title,
			Description:   criterion.Description,
			Polarity:      criterion.Polarity,
			ValueType:     criterion.ValueType,
			ScaleGuidance: append([]any(nil), criterion.ScaleGuidance...),
		})
	}
	return output
}

func buildJSONScenarios(config *ExecutionConfig) []jsonScenarioContext {
	if config == nil {
		return nil
	}
	scenarios := canonicalScenarios(config.Scenarios)
	output := make([]jsonScenarioContext, 0, len(scenarios))
	for _, scenario := range scenarios {
		entry := jsonScenarioContext{
			Name:        scenario.Name,
			Title:       scenario.Title,
			Description: scenario.Description,
			Narrative:   scenario.Narrative,
		}
		for _, criterion := range scenario.ActiveCriteria {
			entry.ActiveCriteria = append(entry.ActiveCriteria, criterion.CriterionName)
		}
		output = append(output, entry)
	}
	return output
}

func buildJSONEvaluations(config *ExecutionConfig) []jsonEvaluationContext {
	if config == nil {
		return nil
	}
	evaluations := canonicalEvaluations(config.Evaluations)
	output := make([]jsonEvaluationContext, 0, len(evaluations))
	for _, evaluation := range evaluations {
		entry := jsonEvaluationContext{
			ScenarioName: evaluation.ScenarioName,
			Description:  evaluation.Description,
		}
		for _, alternative := range canonicalAlternativeEvaluations(evaluation.Evaluations) {
			entry.Evaluations = append(entry.Evaluations, jsonAlternativeEvaluationContext{
				AlternativeName: alternative.AlternativeName,
				Description:     alternative.Description,
			})
		}
		output = append(output, entry)
	}
	return output
}

func canonicalCriteria(criteria []CriterionConfig) []CriterionConfig {
	output := append([]CriterionConfig(nil), criteria...)
	for index := 0; index < len(output); index++ {
		for next := index + 1; next < len(output); next++ {
			if output[next].Name < output[index].Name {
				output[index], output[next] = output[next], output[index]
			}
		}
	}
	return output
}
