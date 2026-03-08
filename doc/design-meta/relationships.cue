package flyb

relationships: [
  {
    from: "analysis.robustness"
    to: "analysis.robustness.method"
    label: "uses_method"
  },
  {
    from: "analysis.robustness"
    to: "analysis.sensitivity"
    label: "uses_method"
  },
  {
    from: "cli.output.machine"
    to: "decision.traceability"
    label: "supports"
  },
  {
    from: "cli.output.readability"
    to: "decision.explainability"
    label: "supports"
  },
  {
    from: "criteria.pairwise.clarity"
    to: "mcda.ahp"
    label: "documents_method"
  },
  {
    from: "criteria.semantic.consistency"
    to: "model.structure"
    label: "refines"
  },
  {
    from: "decision.explainability"
    to: "decision.traceability"
    label: "reinforces"
  },
  {
    from: "example.hosting-choice"
    to: "decision.multi-criteria-ranking"
    label: "addresses_usecase"
  },
  {
    from: "example.hosting-choice"
    to: "example.hosting-choice.lean-startup"
    label: "contains_scenario"
  },
  {
    from: "example.hosting-choice"
    to: "example.hosting-choice.regulated-growth"
    label: "contains_scenario"
  },
  {
    from: "example.hosting-choice"
    to: "mcda.ahp"
    label: "uses_method"
  },
  {
    from: "example.hosting-choice"
    to: "technology.infrastructure-strategy"
    label: "addresses_usecase"
  },
  {
    from: "example.hosting-choice"
    to: "vendor.service-provider-comparison"
    label: "addresses_usecase"
  },
  {
    from: "example.platform-selection"
    to: "decision.multi-criteria-ranking"
    label: "addresses_usecase"
  },
  {
    from: "example.platform-selection"
    to: "decision.robust-choice-identification"
    label: "addresses_usecase"
  },
  {
    from: "example.platform-selection"
    to: "example.platform-selection.established-enterprise"
    label: "contains_scenario"
  },
  {
    from: "example.platform-selection"
    to: "example.platform-selection.startup"
    label: "contains_scenario"
  },
  {
    from: "example.platform-selection"
    to: "example.platform-selection.unicorn"
    label: "contains_scenario"
  },
  {
    from: "example.platform-selection"
    to: "infrastructure.system-design-selection"
    label: "addresses_usecase"
  },
  {
    from: "example.platform-selection"
    to: "mcda.ahp"
    label: "uses_method"
  },
  {
    from: "example.platform-selection"
    to: "strategy.growth-scenario-evaluation"
    label: "addresses_usecase"
  },
  {
    from: "example.platform-selection"
    to: "technology.platform-selection"
    label: "addresses_usecase"
  },
  {
    from: "execution.reproducibility"
    to: "model.validation"
    label: "depends_on"
  },
  {
    from: "input.format"
    to: "model.structure"
    label: "supports"
  },
  {
    from: "model.validation"
    to: "model.incomplete.data"
    label: "reinforces"
  },
  {
    from: "model.documentation"
    to: "model.structure"
    label: "documents"
  },
  {
    from: "planning.lifecycle-decision"
    to: "planning.long-term-option-evaluation"
    label: "relates_to"
  },
  {
    from: "planning.long-term-option-evaluation"
    to: "decision.robust-choice-identification"
    label: "supports"
  },
  {
    from: "policy.policy-option-analysis"
    to: "decision.multi-criteria-ranking"
    label: "supports"
  },
  {
    from: "product.feature-prioritization"
    to: "decision.multi-criteria-ranking"
    label: "supports"
  },
  {
    from: "product.roadmap-planning"
    to: "planning.long-term-option-evaluation"
    label: "supports"
  },
  {
    from: "scenario.aggregation.policy"
    to: "analysis.robustness.method"
    label: "references_analysis"
  },
  {
    from: "scenario.constraints"
    to: "model.validation"
    label: "depends_on"
  },
  {
    from: "scenario.isolation"
    to: "example.hosting-choice"
    label: "illustrated_by"
  },
  {
    from: "scenario.isolation"
    to: "example.platform-selection"
    label: "illustrated_by"
  },
  {
    from: "strategy.investment-decision"
    to: "decision.multi-criteria-ranking"
    label: "supports"
  },
  {
    from: "system.extensibility.methods"
    to: "mcda.electre"
    label: "enables"
  },
  {
    from: "system.extensibility.methods"
    to: "mcda.promethee"
    label: "enables"
  },
  {
    from: "system.extensibility.methods"
    to: "mcda.topsis"
    label: "enables"
  },
  {
    from: "system.extensibility.methods"
    to: "mcda.vikor"
    label: "enables"
  },
  {
    from: "technology.architecture-choice"
    to: "infrastructure.system-design-selection"
    label: "relates_to"
  },
  {
    from: "ux.model.guidance"
    to: "model.structure"
    label: "supports"
  },
  {
    from: "vendor.supplier-selection"
    to: "vendor.service-provider-comparison"
    label: "relates_to"
  },
]
