package pipeline

import (
	"encoding/json"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type jsonRankedAlternative struct {
	AlternativeName  string   `json:"alternativeName"`
	AlternativeTitle string   `json:"alternativeTitle,omitempty"`
	Rank             *int     `json:"rank,omitempty"`
	Score            *float64 `json:"score,omitempty"`
	Excluded         bool     `json:"excluded,omitempty"`
	ExclusionReason  string   `json:"exclusionReason,omitempty"`
}

type jsonScenarioResult struct {
	ScenarioName  string                  `json:"scenarioName"`
	ScenarioTitle string                  `json:"scenarioTitle,omitempty"`
	Ranking       []jsonRankedAlternative `json:"ranking"`
}

type jsonCriterionWeightEntry struct {
	CriterionName  string  `json:"criterionName"`
	CriterionTitle string  `json:"criterionTitle,omitempty"`
	Weight         float64 `json:"weight"`
}

type jsonScenarioWeights struct {
	ScenarioName  string                     `json:"scenarioName"`
	ScenarioTitle string                     `json:"scenarioTitle,omitempty"`
	Weights       []jsonCriterionWeightEntry `json:"weights,omitempty"`
}

type jsonAggregationWeight struct {
	ScenarioName  string  `json:"scenarioName"`
	ScenarioTitle string  `json:"scenarioTitle,omitempty"`
	Weight        float64 `json:"weight"`
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

type jsonPairwiseComparison struct {
	MoreImportantCriterionName  string  `json:"moreImportantCriterionName"`
	MoreImportantCriterionTitle string  `json:"moreImportantCriterionTitle,omitempty"`
	LessImportantCriterionName  string  `json:"lessImportantCriterionName"`
	LessImportantCriterionTitle string  `json:"lessImportantCriterionTitle,omitempty"`
	Strength                    float64 `json:"strength"`
}

type jsonScenarioPreferences struct {
	Method      string                   `json:"method,omitempty"`
	Scale       string                   `json:"scale,omitempty"`
	Comparisons []jsonPairwiseComparison `json:"comparisons,omitempty"`
}

type jsonScenarioConstraint struct {
	CriterionName  string `json:"criterionName"`
	CriterionTitle string `json:"criterionTitle,omitempty"`
	Operator       string `json:"operator"`
	Value          any    `json:"value"`
}

type jsonScenarioContext struct {
	Name           string                   `json:"name"`
	Title          string                   `json:"title,omitempty"`
	Description    string                   `json:"description,omitempty"`
	Narrative      string                   `json:"narrative,omitempty"`
	ActiveCriteria []string                 `json:"activeCriteria,omitempty"`
	Preferences    *jsonScenarioPreferences `json:"preferences,omitempty"`
	Constraints    []jsonScenarioConstraint `json:"constraints,omitempty"`
}

type jsonEvaluationValueContext struct {
	CriterionName  string `json:"criterionName"`
	CriterionTitle string `json:"criterionTitle,omitempty"`
	Value          any    `json:"value"`
	Rendered       string `json:"rendered,omitempty"`
}

type jsonAlternativeEvaluationContext struct {
	AlternativeName  string                       `json:"alternativeName"`
	AlternativeTitle string                       `json:"alternativeTitle,omitempty"`
	Description      string                       `json:"description,omitempty"`
	Values           []jsonEvaluationValueContext `json:"values,omitempty"`
}

type jsonEvaluationContext struct {
	ScenarioName  string                             `json:"scenarioName"`
	ScenarioTitle string                             `json:"scenarioTitle,omitempty"`
	Description   string                             `json:"description,omitempty"`
	Evaluations   []jsonAlternativeEvaluationContext `json:"evaluations,omitempty"`
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

type jsonRenderOptions struct {
	IncludeContext bool
	IncludeWeights bool
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
	options := resolveJSONRenderOptions(report)
	orderedResults := orderedScenarioResultsForMarkdown(report, config, scenarioResults)
	orderedFinalRanking := domain.CanonicalAggregatedRankingResult(finalRanking).RankedAlternatives

	payload := jsonReport{
		ProblemName: problemName(config),
		ReportName:  report.Name,
		Format:      report.Format,
	}
	if aggregation != nil {
		payload.Aggregation = jsonAggregation{Method: aggregation.Method}
		if options.IncludeWeights {
			for _, scenarioName := range orderedWeightNames(aggregationScenarioWeights) {
				payload.Aggregation.ScenarioWeights = append(payload.Aggregation.ScenarioWeights, jsonAggregationWeight{
					ScenarioName:  scenarioName,
					ScenarioTitle: scenarioTitleByName(config, scenarioName),
					Weight:        aggregationScenarioWeights[scenarioName],
				})
			}
		}
	}

	for _, scenarioResult := range orderedResults {
		jsonResult := jsonScenarioResult{
			ScenarioName:  scenarioResult.ScenarioName,
			ScenarioTitle: scenarioTitleByName(config, scenarioResult.ScenarioName),
		}
		for _, alternative := range scenarioResult.RankedAlternatives {
			jsonResult.Ranking = append(jsonResult.Ranking, buildJSONRankedAlternative(config, alternative))
		}
		payload.ScenarioResults = append(payload.ScenarioResults, jsonResult)
	}

	if options.IncludeWeights {
		for _, scenarioWeight := range canonicalScenarioWeights(scenarioWeights) {
			entry := jsonScenarioWeights{
				ScenarioName:  scenarioWeight.ScenarioName,
				ScenarioTitle: scenarioTitleByName(config, scenarioWeight.ScenarioName),
			}
			for _, weight := range canonicalCriterionWeights(scenarioWeight.CriterionWeights) {
				entry.Weights = append(entry.Weights, jsonCriterionWeightEntry{
					CriterionName:  weight.CriterionName,
					CriterionTitle: criterionTitleByName(config, weight.CriterionName),
					Weight:         weight.Weight,
				})
			}
			payload.ScenarioWeights = append(payload.ScenarioWeights, entry)
		}
	}

	if options.IncludeContext {
		payload.Problem = buildJSONProblemContext(config)
		payload.Report = buildJSONReportContext(report)
		payload.Alternatives = buildJSONAlternatives(report, config)
		payload.Criteria = buildJSONCriteria(report, config)
		payload.Scenarios = buildJSONScenarios(report, config)
		payload.Evaluations = buildJSONEvaluations(report, config)
	}

	for _, alternative := range orderedFinalRanking {
		payload.FinalRanking = append(payload.FinalRanking, buildJSONRankedAlternative(config, alternative))
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

func resolveJSONRenderOptions(report ReportConfig) jsonRenderOptions {
	options := jsonRenderOptions{
		IncludeContext: true,
		IncludeWeights: true,
	}
	if reportArgumentPresent(report.Arguments, "include-context") {
		options.IncludeContext = reportArgumentValue(report.Arguments, "include-context", "false") == "true"
	}
	if reportArgumentPresent(report.Arguments, "include-weights") {
		options.IncludeWeights = reportArgumentValue(report.Arguments, "include-weights", "false") == "true"
	}
	return options
}

func buildJSONRankedAlternative(config *ExecutionConfig, alternative domain.RankedAlternative) jsonRankedAlternative {
	entry := jsonRankedAlternative{
		AlternativeName:  alternative.Name,
		AlternativeTitle: alternativeTitleByName(config, alternative.Name),
		Excluded:         alternative.Excluded,
		ExclusionReason:  alternative.ExclusionReason,
	}
	if !alternative.Excluded {
		entry.Rank = intPointer(alternative.Rank)
		entry.Score = floatPointer(alternative.Score)
	}
	return entry
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

func buildJSONAlternatives(report ReportConfig, config *ExecutionConfig) []jsonAlternativeContext {
	alternatives := filteredAlternativesForMarkdown(report, config)
	output := make([]jsonAlternativeContext, 0, len(alternatives))
	for _, alternative := range alternatives {
		output = append(output, jsonAlternativeContext(alternative))
	}
	return output
}

func buildJSONCriteria(report ReportConfig, config *ExecutionConfig) []jsonCriterionContext {
	criteria := filteredCriteriaForMarkdown(report, config)
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

func buildJSONScenarios(report ReportConfig, config *ExecutionConfig) []jsonScenarioContext {
	scenarios := filteredScenariosForMarkdown(report, config)
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
		if scenario.Preferences != nil {
			entry.Preferences = &jsonScenarioPreferences{
				Method: scenario.Preferences.Method,
				Scale:  scenario.Preferences.Scale,
			}
			for _, comparison := range scenario.Preferences.Comparisons {
				entry.Preferences.Comparisons = append(entry.Preferences.Comparisons, jsonPairwiseComparison{
					MoreImportantCriterionName:  comparison.MoreImportantCriterionName,
					MoreImportantCriterionTitle: criterionTitleByName(config, comparison.MoreImportantCriterionName),
					LessImportantCriterionName:  comparison.LessImportantCriterionName,
					LessImportantCriterionTitle: criterionTitleByName(config, comparison.LessImportantCriterionName),
					Strength:                    comparison.Strength,
				})
			}
		}
		for _, constraint := range scenario.Constraints {
			entry.Constraints = append(entry.Constraints, jsonScenarioConstraint{
				CriterionName:  constraint.CriterionName,
				CriterionTitle: criterionTitleByName(config, constraint.CriterionName),
				Operator:       constraint.Operator,
				Value:          constraint.Value,
			})
		}
		output = append(output, entry)
	}
	return output
}

func buildJSONEvaluations(report ReportConfig, config *ExecutionConfig) []jsonEvaluationContext {
	if config == nil {
		return nil
	}

	criteria := filteredCriteriaForMarkdown(report, config)
	output := make([]jsonEvaluationContext, 0, len(config.Evaluations))
	for _, scenario := range filteredScenariosForMarkdown(report, config) {
		evaluation, exists := findScenarioEvaluation(report, config, scenario.Name)
		if !exists {
			continue
		}
		entry := jsonEvaluationContext{
			ScenarioName:  evaluation.ScenarioName,
			ScenarioTitle: scenario.Title,
			Description:   evaluation.Description,
		}
		for _, alternative := range orderedAlternativeEvaluations(report, config, evaluation.Evaluations) {
			alternativeEntry := jsonAlternativeEvaluationContext{
				AlternativeName:  alternative.AlternativeName,
				AlternativeTitle: alternativeTitleByName(config, alternative.AlternativeName),
				Description:      alternative.Description,
			}
			for _, value := range orderedCriterionValueRecords(criteria, alternative.Values) {
				criterionName := value.Name
				alternativeEntry.Values = append(alternativeEntry.Values, jsonEvaluationValueContext{
					CriterionName:  criterionName,
					CriterionTitle: criterionTitleByName(config, criterionName),
					Value:          alternative.Values[criterionName].Value,
					Rendered:       value.Rendered,
				})
			}
			entry.Evaluations = append(entry.Evaluations, alternativeEntry)
		}
		output = append(output, entry)
	}
	return output
}
