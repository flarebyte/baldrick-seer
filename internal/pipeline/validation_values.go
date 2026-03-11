package pipeline

import (
	"fmt"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func validateAlternativeEvaluationValues(
	evaluationIndex int,
	alternativeIndex int,
	scenarioInfo scenarioValidationInfo,
	alternative AlternativeEvaluationConfig,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
	scenarioName string,
) []domain.Diagnostic {
	var diagnostics []domain.Diagnostic

	for _, criterionName := range scenarioInfo.ActiveCriterionNames {
		if _, exists := alternative.Values[criterionName]; !exists {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.missing_evaluation_value",
				fmt.Sprintf("config.evaluations[%d].evaluations[%d].values", evaluationIndex, alternativeIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("missing value for active criterion in scenario %s: %s", scenarioName, criterionName),
			))
		}
	}

	valueNames := make([]string, 0, len(alternative.Values))
	for criterionName := range alternative.Values {
		valueNames = append(valueNames, criterionName)
	}

	for _, criterionName := range domain.CanonicalNames(valueNames) {
		value := alternative.Values[criterionName]
		valuePath := fmt.Sprintf("config.evaluations[%d].evaluations[%d].values.%s", evaluationIndex, alternativeIndex, criterionName)

		if !hasName(criteriaNames, criterionName) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_evaluation_criterion",
				valuePath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown criterion name in evaluation values: %s", criterionName),
			))
			continue
		}

		if !hasName(scenarioInfo.ActiveCriterionNames, criterionName) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.inactive_evaluation_criterion",
				valuePath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("criterion value is not active in scenario %s: %s", scenarioName, criterionName),
			))
			continue
		}

		criterion := criteriaByName[criterionName]
		if !isSupportedCriterionValueKind(value.Kind) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unsupported_evaluation_value_kind",
				valuePath+".kind",
				domain.DiagnosticLocation{},
				fmt.Sprintf("unsupported evaluation value kind for criterion %s: %s", criterionName, value.Kind),
			))
			continue
		}

		if value.Kind != criterion.ValueType {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.evaluation_value_kind_mismatch",
				valuePath+".kind",
				domain.DiagnosticLocation{},
				fmt.Sprintf("evaluation value kind mismatch for criterion %s: want %s, got %s", criterionName, criterion.ValueType, value.Kind),
			))
			continue
		}

		if !isValidCriterionValue(value.Kind, value.Value) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				validationValueSemanticCode(value.Kind),
				valuePath+".value",
				domain.DiagnosticLocation{},
				validationValueSemanticMessage(criterionName, value.Kind),
			))
		}
	}

	return diagnostics
}

func isSupportedCriterionValueKind(kind string) bool {
	return isSupportedCriterionValueType(kind)
}

func isValidCriterionValue(kind string, value any) bool {
	switch kind {
	case "number":
		return isNumericValue(value)
	case "ordinal":
		return isIntegerValue(value)
	case "boolean":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}

func isNumericValue(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return true
	default:
		return false
	}
}

func isIntegerValue(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return true
	default:
		return false
	}
}

func validationValueSemanticCode(kind string) string {
	switch kind {
	case "number":
		return "validation.invalid_number_value"
	case "ordinal":
		return "validation.invalid_ordinal_value"
	case "boolean":
		return "validation.invalid_boolean_value"
	default:
		return "validation.invalid_value"
	}
}

func validationValueSemanticMessage(criterionName string, kind string) string {
	switch kind {
	case "number":
		return fmt.Sprintf("number criterion value must be numeric: %s", criterionName)
	case "ordinal":
		return fmt.Sprintf("ordinal criterion value must be an integer: %s", criterionName)
	case "boolean":
		return fmt.Sprintf("boolean criterion value must be true or false: %s", criterionName)
	default:
		return fmt.Sprintf("invalid criterion value: %s", criterionName)
	}
}

func isSupportedConstraintOperator(valueType string, operator string) bool {
	switch valueType {
	case "number", "ordinal":
		return operator == "<=" || operator == ">=" || operator == "=" || operator == "!="
	case "boolean":
		return operator == "=" || operator == "!="
	default:
		return false
	}
}

func validationConstraintValueCode(valueType string) string {
	switch valueType {
	case "number":
		return "validation.invalid_constraint_number_value"
	case "ordinal":
		return "validation.invalid_constraint_ordinal_value"
	case "boolean":
		return "validation.invalid_constraint_boolean_value"
	default:
		return "validation.invalid_constraint_value"
	}
}

func validationConstraintValueMessage(valueType string, criterionName string) string {
	switch valueType {
	case "number":
		return fmt.Sprintf("constraint value must be numeric for criterion: %s", criterionName)
	case "ordinal":
		return fmt.Sprintf("constraint value must be an integer for criterion: %s", criterionName)
	case "boolean":
		return fmt.Sprintf("constraint value must be true or false for criterion: %s", criterionName)
	default:
		return fmt.Sprintf("invalid constraint value for criterion: %s", criterionName)
	}
}
