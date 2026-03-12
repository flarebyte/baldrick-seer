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

type jsonReport struct {
	ProblemName     string                  `json:"problemName"`
	ReportName      string                  `json:"reportName"`
	Format          string                  `json:"format"`
	Aggregation     jsonAggregation         `json:"aggregation"`
	ScenarioResults []jsonScenarioResult    `json:"scenarioResults"`
	ScenarioWeights []jsonScenarioWeights   `json:"scenarioWeights,omitempty"`
	FinalRanking    []jsonRankedAlternative `json:"finalRanking"`
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
