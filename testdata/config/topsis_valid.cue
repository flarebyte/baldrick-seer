package seer

config: {
	problem: {
		name: "topsis-valid"
	}
	reports: [{
		name:   "summary"
		title:  "Summary"
		format: "markdown"
	}]
	criteriaCatalog: [
		{name: "cost", polarity: "cost", valueType: "number"},
		{name: "speed", polarity: "benefit", valueType: "number"},
	]
	alternatives: [
		{name: "alpha"},
		{name: "beta"},
	]
	scenarios: [{
		name: "baseline"
		activeCriteria: [
			{criterionName: "cost"},
			{criterionName: "speed"},
		]
		preferences: {
			method: "ahp_pairwise"
			scale:  "saaty_1_9"
			comparisons: [{
				moreImportantCriterionName: "cost"
				lessImportantCriterionName: "speed"
				strength:                   3
			}]
		}
	}]
	evaluations: [{
		scenarioName: "baseline"
		evaluations: [
			{
				alternativeName: "alpha"
				values: {
					cost: {
						kind:  "number"
						value: 100
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
						value: 180
					}
					speed: {
						kind:  "number"
						value: 90
					}
				}
			},
		]
	}]
	aggregation: {
		method: "equal_average"
	}
}
