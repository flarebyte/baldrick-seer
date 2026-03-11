package seer

config: {
	problem: {
		name: "invalid-constraint"
	}
	reports: [{
		name:   "summary"
		title:  "Summary"
		format: "markdown"
	}]
	criteriaCatalog: [{
		name:      "approved"
		valueType: "boolean"
	}]
	alternatives: [{
		name: "option_a"
	}]
	scenarios: [{
		name: "baseline"
		activeCriteria: [{
			criterionName: "approved"
		}]
		constraints: [{
			criterionName: "approved"
			operator:      "<="
			value:         true
		}]
	}]
	evaluations: [{
		scenarioName: "baseline"
		evaluations: [{
			alternativeName: "option_a"
			values: {
				approved: {
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
