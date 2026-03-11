package seer

config: {
	problem: {
		name: "invalid-evaluation"
	}
	reports: [{
		name:   "summary"
		title:  "Summary"
		format: "markdown"
	}]
	criteriaCatalog: [{
		name:      "cost"
		valueType: "number"
	}]
	alternatives: [{
		name: "option_a"
	}]
	scenarios: [{
		name: "baseline"
		activeCriteria: [{
			criterionName: "cost"
		}]
	}]
	evaluations: [{
		scenarioName: "baseline"
		evaluations: [{
			alternativeName: "option_a"
			values: {
				cost: {
					kind:  "boolean"
					value: true
				}
			}
		}]
	}]
	aggregation: {
		method: "equal_average"
	}
}
