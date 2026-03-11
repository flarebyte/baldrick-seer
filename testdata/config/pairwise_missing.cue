package seer

config: {
	problem: {
		name: "pairwise-missing"
	}
	reports: [{
		name:   "summary"
		title:  "Summary"
		format: "markdown"
	}]
	criteriaCatalog: [
		{name: "cost", valueType: "number"},
		{name: "speed", valueType: "number"},
		{name: "reliability", valueType: "number"},
	]
	alternatives: [{
		name: "option_a"
	}]
	scenarios: [{
		name: "baseline"
		activeCriteria: [
			{criterionName: "cost"},
			{criterionName: "speed"},
			{criterionName: "reliability"},
		]
		preferences: {
			method: "ahp_pairwise"
			scale:  "saaty_1_9"
			comparisons: [
				{
					moreImportantCriterionName: "cost"
					lessImportantCriterionName: "speed"
					strength:                   3
				},
				{
					moreImportantCriterionName: "cost"
					lessImportantCriterionName: "reliability"
					strength:                   5
				},
			]
		}
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
				reliability: {
					kind:  "number"
					value: 2
				}
				speed: {
					kind:  "number"
					value: 3
				}
			}
		}]
	}]
	aggregation: {
		method: "equal_average"
	}
}
