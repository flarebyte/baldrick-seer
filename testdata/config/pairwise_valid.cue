package seer

config: {
	problem: {
		name: "pairwise-valid"
	}
	reports: [{
		name:   "summary"
		title:  "Summary"
		format: "markdown"
	}]
	criteriaCatalog: [
		{name: "cost", polarity: "cost", valueType: "number"},
		{name: "speed", polarity: "benefit", valueType: "number"},
		{name: "reliability", polarity: "benefit", valueType: "number"},
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
				{
					moreImportantCriterionName: "speed"
					lessImportantCriterionName: "reliability"
					strength:                   2
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
