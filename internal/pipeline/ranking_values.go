package pipeline

func normalizeCriterionValue(criterion CriterionConfig, value CriterionValue) (float64, error) {
	switch criterion.ValueType {
	case "number":
		return numericValue(value.Value)
	case "ordinal":
		return numericValue(value.Value)
	case "boolean":
		booleanValue, ok := value.Value.(bool)
		if !ok {
			return 0, ErrRankingFailed
		}
		if booleanValue {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ErrRankingFailed
	}
}

func numericValue(value any) (float64, error) {
	switch typed := value.(type) {
	case int:
		return float64(typed), nil
	case int8:
		return float64(typed), nil
	case int16:
		return float64(typed), nil
	case int32:
		return float64(typed), nil
	case int64:
		return float64(typed), nil
	case float32:
		return float64(typed), nil
	case float64:
		return typed, nil
	default:
		return 0, ErrRankingFailed
	}
}
