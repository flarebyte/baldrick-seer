package pipeline

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

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

func csvSchemaDescriptions() map[string]string {
	schema := make(map[string]string, len(csvColumnDefinitions))
	for _, definition := range csvColumnDefinitions {
		schema[definition.Name] = definition.Description
	}
	return schema
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
		case "excluded":
			if alternative.Excluded {
				record = append(record, "true")
			} else {
				record = append(record, "false")
			}
		case "exclusion_reason":
			record = append(record, alternative.ExclusionReason)
		default:
			record = append(record, "")
		}
	}
	return record
}
