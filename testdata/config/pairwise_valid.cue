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
			comparisons: [
				{
					moreImportantCriterionName: "cost"
					lessImportantCriterionName: "speed"
				},
				{
					moreImportantCriterionName: "cost"
					lessImportantCriterionName: "reliability"
				},
				{
					moreImportantCriterionName: "speed"
					lessImportantCriterionName: "reliability"
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
