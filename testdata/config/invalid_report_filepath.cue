package seer

config: {
	problem: {
		name: "invalid-report-filepath"
	}
	reports: [{
		name:      "summary"
		title:     "Summary"
		format:    "markdown"
		filepath:  "/tmp/summary.md"
		arguments: ["include-scenarios=all"]
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
