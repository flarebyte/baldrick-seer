package seer

config: {
	problem: {
		name: "large-invalid"
	}
	reports: [
		{
			name:      "summary-markdown"
			title:     "Summary Markdown"
			format:    "markdown"
			arguments: ["include-scenarios=all", "top-alternatives=3", "include-scores=true"]
		},
		{
			name:      "summary-json"
			title:     "Summary JSON"
			format:    "json"
			arguments: ["include-evidence=false", "pretty=true"]
		},
	]
	criteriaCatalog: [
		{
			name:      "cost"
			polarity:  "cost"
			valueType: "number"
		},
		{
			name:          "quality"
			polarity:      "benefit"
			valueType:     "ordinal"
			scaleGuidance: [1, 2, 3, 4, 5]
		},
		{
			name:      "approved"
			polarity:  "benefit"
			valueType: "boolean"
		},
		{
			name:      "speed"
			polarity:  "benefit"
			valueType: "number"
		},
	]
	alternatives: [
		{name: "alpha"},
		{name: "beta"},
		{name: "gamma"},
	]
	scenarios: [
		{
			name: "baseline"
			activeCriteria: [
				{criterionName: "approved"},
				{criterionName: "cost"},
				{criterionName: "quality"},
			]
			preferences: {
				method: "ahp_pairwise"
				scale:  "saaty_1_9"
				comparisons: [
					{
						moreImportantCriterionName: "cost"
						lessImportantCriterionName: "approved"
						strength:                   2
					},
					{
						moreImportantCriterionName: "quality"
						lessImportantCriterionName: "approved"
						strength:                   3
					},
					{
						moreImportantCriterionName: "quality"
						lessImportantCriterionName: "cost"
						strength:                   2
					},
				]
			}
			constraints: [{
				criterionName: "approved"
				operator:      "="
				value:         true
			}]
		},
		{
			name: "expansion"
			activeCriteria: [
				{criterionName: "cost"},
				{criterionName: "quality"},
				{criterionName: "speed"},
			]
			preferences: {
				method: "ahp_pairwise"
				scale:  "saaty_1_9"
				comparisons: [
					{
						moreImportantCriterionName: "speed"
						lessImportantCriterionName: "cost"
						strength:                   4
					},
					{
						moreImportantCriterionName: "quality"
						lessImportantCriterionName: "cost"
						strength:                   2
					},
					{
						moreImportantCriterionName: "speed"
						lessImportantCriterionName: "quality"
						strength:                   2
					},
				]
			}
		},
	]
	evaluations: [
		{
			scenarioName: "baseline"
			evaluations: [
				{
					alternativeName: "alpha"
					values: {
						approved: {
							kind:  "boolean"
							value: true
						}
						cost: {
							kind:  "number"
							value: 120
						}
						quality: {
							kind:  "ordinal"
							value: 5
						}
					}
				},
				{
					alternativeName: "beta"
					values: {
						approved: {
							kind:  "boolean"
							value: true
						}
						cost: {
							kind:  "number"
							value: 140
						}
						quality: {
							kind:  "ordinal"
							value: 4
						}
					}
				},
				{
					alternativeName: "gamma"
					values: {
						approved: {
							kind:  "boolean"
							value: false
						}
						cost: {
							kind:  "number"
							value: 100
						}
						quality: {
							kind:  "ordinal"
							value: 3
						}
					}
				},
			]
		},
		{
			scenarioName: "expansion"
			evaluations: [
				{
					alternativeName: "alpha"
					values: {
						cost: {
							kind:  "number"
							value: 150
						}
						quality: {
							kind:  "ordinal"
							value: 5
						}
						speed: {
							kind:  "number"
							value: 70
						}
					}
				},
				{
					alternativeName: "beta"
					values: {
						cost: {
							kind:  "number"
							value: 170
						}
						quality: {
							kind:  "ordinal"
							value: 4
						}
						speed: {
							kind:  "number"
							value: 80
						}
						unknown: {
							kind:  "number"
							value: 1
						}
					}
				},
				{
					alternativeName: "gamma"
					values: {
						cost: {
							kind:  "number"
							value: 190
						}
						quality: {
							kind:  "ordinal"
							value: 3
						}
						speed: {
							kind:  "number"
							value: 65
						}
					}
				},
			]
		},
	]
	aggregation: {
		method: "equal_average"
	}
}
