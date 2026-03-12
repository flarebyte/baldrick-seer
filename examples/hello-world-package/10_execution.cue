package seer

config: {
	reports: [{
		name:   "summary"
		title:  "Hello World Package Summary"
		format: "markdown"
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
