package pipeline

import (
	"fmt"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func validateCriteriaCatalog(diagnostics *[]domain.Diagnostic, criteria []CriterionConfig) map[string]CriterionConfig {
	criteriaByName := make(map[string]CriterionConfig, len(criteria))
	for criterionIndex, criterion := range criteria {
		if _, exists := criteriaByName[criterion.Name]; !exists {
			criteriaByName[criterion.Name] = criterion
		}

		if !isSupportedCriterionValueType(criterion.ValueType) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unsupported_criterion_value_type",
				fmt.Sprintf("config.criteriaCatalog[%d].valueType", criterionIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unsupported criterion value type: %s", criterion.ValueType),
			))
		}

		if criterion.ValueType == "ordinal" && len(criterion.ScaleGuidance) == 0 {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.ordinal_scale_guidance_missing",
				fmt.Sprintf("config.criteriaCatalog[%d].scaleGuidance", criterionIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("ordinal criterion is missing scaleGuidance: %s", criterion.Name),
			))
		}
	}

	return criteriaByName
}

func validateScenarios(
	diagnostics *[]domain.Diagnostic,
	scenarios []ScenarioConfig,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
) map[string]scenarioValidationInfo {
	scenarioNameCounts := countScenarioNames(scenarios)
	scenarioInfos := make(map[string]scenarioValidationInfo, len(scenarios))

	for scenarioIndex, scenario := range scenarios {
		activeCriterionNames := make([]string, 0, len(scenario.ActiveCriteria))
		hasUnknownActiveCriteria := false
		for criterionIndex, ref := range scenario.ActiveCriteria {
			activeCriterionNames = append(activeCriterionNames, ref.CriterionName)
			if !hasName(criteriaNames, ref.CriterionName) {
				hasUnknownActiveCriteria = true
				*diagnostics = append(*diagnostics, domain.NewDiagnostic(
					domain.DiagnosticSeverityError,
					"validation.unknown_active_criterion",
					fmt.Sprintf("config.scenarios[%d].activeCriteria[%d].criterionName", scenarioIndex, criterionIndex),
					domain.DiagnosticLocation{},
					fmt.Sprintf("unknown criterion name in active criteria: %s", ref.CriterionName),
				))
			}
		}

		activeCriterionNames = domain.CanonicalNames(activeCriterionNames)
		if scenarioNameCounts[scenario.Name] == 1 && !hasUnknownActiveCriteria {
			scenarioInfos[scenario.Name] = scenarioValidationInfo{
				Index:                scenarioIndex,
				ActiveCriterionNames: activeCriterionNames,
			}
		}

		validateScenarioConstraints(diagnostics, scenarioIndex, scenario, activeCriterionNames, criteriaNames, criteriaByName)
		*diagnostics = append(*diagnostics, validateScenarioPairwiseComparisons(
			scenarioIndex,
			scenario,
			criteriaNames,
			activeCriterionNames,
		)...)
	}

	return scenarioInfos
}

func validateScenarioConstraints(
	diagnostics *[]domain.Diagnostic,
	scenarioIndex int,
	scenario ScenarioConfig,
	activeCriterionNames []string,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
) {
	for constraintIndex, constraint := range scenario.Constraints {
		if !hasName(criteriaNames, constraint.CriterionName) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_constraint_criterion",
				fmt.Sprintf("config.scenarios[%d].constraints[%d].criterionName", scenarioIndex, constraintIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown criterion name in constraints: %s", constraint.CriterionName),
			))
			continue
		}

		if !hasName(activeCriterionNames, constraint.CriterionName) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.inactive_constraint_criterion",
				fmt.Sprintf("config.scenarios[%d].constraints[%d].criterionName", scenarioIndex, constraintIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("constraint references criterion not active in scenario: %s", constraint.CriterionName),
			))
			continue
		}

		criterion := criteriaByName[constraint.CriterionName]
		if !isSupportedConstraintOperator(criterion.ValueType, constraint.Operator) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.invalid_constraint_operator",
				fmt.Sprintf("config.scenarios[%d].constraints[%d].operator", scenarioIndex, constraintIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("invalid constraint operator for %s criterion %s: %s", criterion.ValueType, constraint.CriterionName, constraint.Operator),
			))
		}

		if !isValidCriterionValue(criterion.ValueType, constraint.Value) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				validationConstraintValueCode(criterion.ValueType),
				fmt.Sprintf("config.scenarios[%d].constraints[%d].value", scenarioIndex, constraintIndex),
				domain.DiagnosticLocation{},
				validationConstraintValueMessage(criterion.ValueType, constraint.CriterionName),
			))
		}
	}
}

func validateEvaluations(
	diagnostics *[]domain.Diagnostic,
	evaluations []EvaluationConfig,
	scenarioInfos map[string]scenarioValidationInfo,
	scenarioNames []string,
	alternativeNames []string,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
) {
	seenEvaluationScenarios := make(map[string]struct{}, len(evaluations))
	for evaluationIndex, evaluation := range evaluations {
		if _, exists := seenEvaluationScenarios[evaluation.ScenarioName]; exists {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.duplicate_evaluation_scenario",
				fmt.Sprintf("config.evaluations[%d].scenarioName", evaluationIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("duplicate evaluation block for scenario: %s", evaluation.ScenarioName),
			))
		} else {
			seenEvaluationScenarios[evaluation.ScenarioName] = struct{}{}
		}

		scenarioInfo, hasScenario := scenarioInfos[evaluation.ScenarioName]
		if !hasName(scenarioNames, evaluation.ScenarioName) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_evaluation_scenario",
				fmt.Sprintf("config.evaluations[%d].scenarioName", evaluationIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown scenario name in evaluations: %s", evaluation.ScenarioName),
			))
		}

		validateScenarioEvaluationAlternatives(
			diagnostics,
			evaluationIndex,
			evaluation,
			hasScenario,
			scenarioInfo,
			alternativeNames,
			criteriaNames,
			criteriaByName,
		)
	}
}

func validateScenarioEvaluationAlternatives(
	diagnostics *[]domain.Diagnostic,
	evaluationIndex int,
	evaluation EvaluationConfig,
	hasScenario bool,
	scenarioInfo scenarioValidationInfo,
	alternativeNames []string,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
) {
	seenAlternatives := make(map[string]struct{}, len(evaluation.Evaluations))
	for alternativeIndex, alternative := range evaluation.Evaluations {
		if _, exists := seenAlternatives[alternative.AlternativeName]; exists {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.duplicate_evaluation_alternative",
				fmt.Sprintf("config.evaluations[%d].evaluations[%d].alternativeName", evaluationIndex, alternativeIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("duplicate alternative evaluation in scenario %s: %s", evaluation.ScenarioName, alternative.AlternativeName),
			))
		} else {
			seenAlternatives[alternative.AlternativeName] = struct{}{}
		}

		if !hasName(alternativeNames, alternative.AlternativeName) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_evaluation_alternative",
				fmt.Sprintf("config.evaluations[%d].evaluations[%d].alternativeName", evaluationIndex, alternativeIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown alternative name in evaluations: %s", alternative.AlternativeName),
			))
		}

		if hasScenario {
			*diagnostics = append(*diagnostics, validateAlternativeEvaluationValues(
				evaluationIndex,
				alternativeIndex,
				scenarioInfo,
				alternative,
				criteriaNames,
				criteriaByName,
				evaluation.ScenarioName,
			)...)
		}
	}
}
