package seer

config: {
	problem: {
		name: "invalid-report"
	}
	reports: [{
		name:      "summary"
		title:     "Summary"
		format:    "json"
		arguments: ["header=true"]
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
					kind:  "number"
					value: 1
				}
			}
		}]
	}]
	aggregation: {
		method: "equal_average"
	}
}
