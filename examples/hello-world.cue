package seer

config: {
	problem: {
		name: "hello-world"
	}
	reports: [{
		name:   "summary"
		title:  "Hello World Summary"
		format: "markdown"
	}]
	criteriaCatalog: [{
		name:      "cost"
		polarity:  "cost"
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
