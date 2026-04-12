package seer

config: {
	problem: {
		name: "valid-report-filepath"
	}
	reports: [
		{
			name:      "summary-markdown"
			title:     "Summary Markdown"
			format:    "markdown"
			filepath:  "../artifacts/summary.md"
			arguments: ["include-scenarios=all", "top-alternatives=2", "include-scores=true"]
		},
		{
			name:      "summary-json"
			title:     "Summary JSON"
			format:    "json"
			filepath:  "artifacts/summary.json"
			arguments: ["include-evidence=false", "pretty=true"]
		},
	]
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
