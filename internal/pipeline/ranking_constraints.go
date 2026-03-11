package pipeline

func alternativeViolatesConstraints(
	scenario ScenarioConfig,
	alternative AlternativeEvaluationConfig,
	criteriaByName map[string]CriterionConfig,
) (bool, error) {
	for _, constraint := range scenario.Constraints {
		criterion, exists := criteriaByName[constraint.CriterionName]
		if !exists {
			return false, ErrRankingFailed
		}
		value, exists := alternative.Values[constraint.CriterionName]
		if !exists {
			return false, ErrRankingFailed
		}
		matches, err := constraintMatches(criterion, constraint, value)
		if err != nil {
			return false, err
		}
		if !matches {
			return true, nil
		}
	}

	return false, nil
}

func constraintMatches(criterion CriterionConfig, constraint ConstraintConfig, value CriterionValue) (bool, error) {
	switch criterion.ValueType {
	case "number", "ordinal":
		left, err := numericValue(value.Value)
		if err != nil {
			return false, err
		}
		right, err := numericValue(constraint.Value)
		if err != nil {
			return false, err
		}
		switch constraint.Operator {
		case "<=":
			return left <= right, nil
		case ">=":
			return left >= right, nil
		case "=":
			return left == right, nil
		case "!=":
			return left != right, nil
		default:
			return false, ErrRankingFailed
		}
	case "boolean":
		left, ok := value.Value.(bool)
		if !ok {
			return false, ErrRankingFailed
		}
		right, ok := constraint.Value.(bool)
		if !ok {
			return false, ErrRankingFailed
		}
		switch constraint.Operator {
		case "=":
			return left == right, nil
		case "!=":
			return left != right, nil
		default:
			return false, ErrRankingFailed
		}
	default:
		return false, ErrRankingFailed
	}
}
