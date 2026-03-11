package seer

config: {
	problem: {
		name: "minimal"
	}
	reports: [{
		name:   "summary"
		title:  "Summary"
		format: "markdown"
	}]
	criteriaCatalog: [{
		name:      "cost"
		polarity:  "cost"
		valueType: "number"
	}]
}
