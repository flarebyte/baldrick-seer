package flyb

relationships: [
  {
    from: #notesByName["analysis.robustness"].name
    to: #notesByName["analysis.robustness.method"].name
    label: "uses_method"
  },
  {
    from: #notesByName["analysis.robustness"].name
    to: #notesByName["analysis.sensitivity"].name
    label: "uses_method"
  },
  {
    from: #notesByName["cli.output.machine"].name
    to: #notesByName["decision.traceability"].name
    label: "supports"
  },
  {
    from: #notesByName["cli.output.readability"].name
    to: #notesByName["decision.explainability"].name
    label: "supports"
  },
  {
    from: #notesByName["call.validation.input-config"].name
    to: #notesByName["call.validation.input-config.parse-args"].name
    label: "delegate_to"
  },
  {
    from: #notesByName["call.validation.input-config.parse-args"].name
    to: #notesByName["call.validation.input-config.load-cue-config"].name
    label: "delegate_to"
  },
  {
    from: #notesByName["call.validation.input-config.load-cue-config"].name
    to: #notesByName["call.validation.input-config.validate-model"].name
    label: "delegate_to"
  },
  {
    from: #notesByName["call.validation.input-config.load-cue-config"].name
    to: #notesByName["input.format"].name
    label: "depends_on"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model"].name
    to: #notesByName["model.validation"].name
    label: "implements"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model"].name
    to: #notesByName["model.incomplete.data"].name
    label: "checks_for"
  },
  {
    from: #notesByName["criteria.pairwise.clarity"].name
    to: #notesByName["mcda.ahp"].name
    label: "documents_method"
  },
  {
    from: #notesByName["criteria.semantic.consistency"].name
    to: #notesByName["model.structure"].name
    label: "refines"
  },
  {
    from: #notesByName["decision.explainability"].name
    to: #notesByName["decision.traceability"].name
    label: "reinforces"
  },
  {
    from: #notesByName["example.hosting-choice"].name
    to: #notesByName["decision.multi-criteria-ranking"].name
    label: "addresses_usecase"
  },
  {
    from: #notesByName["example.hosting-choice"].name
    to: #notesByName["example.hosting-choice.lean-startup"].name
    label: "contains_scenario"
  },
  {
    from: #notesByName["example.hosting-choice"].name
    to: #notesByName["example.hosting-choice.regulated-growth"].name
    label: "contains_scenario"
  },
  {
    from: #notesByName["example.hosting-choice"].name
    to: #notesByName["mcda.ahp"].name
    label: "uses_method"
  },
  {
    from: #notesByName["example.hosting-choice"].name
    to: #notesByName["technology.infrastructure-strategy"].name
    label: "addresses_usecase"
  },
  {
    from: #notesByName["example.hosting-choice"].name
    to: #notesByName["vendor.service-provider-comparison"].name
    label: "addresses_usecase"
  },
  {
    from: #notesByName["example.platform-selection"].name
    to: #notesByName["decision.multi-criteria-ranking"].name
    label: "addresses_usecase"
  },
  {
    from: #notesByName["example.platform-selection"].name
    to: #notesByName["decision.robust-choice-identification"].name
    label: "addresses_usecase"
  },
  {
    from: #notesByName["example.platform-selection"].name
    to: #notesByName["example.platform-selection.established-enterprise"].name
    label: "contains_scenario"
  },
  {
    from: #notesByName["example.platform-selection"].name
    to: #notesByName["example.platform-selection.startup"].name
    label: "contains_scenario"
  },
  {
    from: #notesByName["example.platform-selection"].name
    to: #notesByName["example.platform-selection.unicorn"].name
    label: "contains_scenario"
  },
  {
    from: #notesByName["example.platform-selection"].name
    to: #notesByName["infrastructure.system-design-selection"].name
    label: "addresses_usecase"
  },
  {
    from: #notesByName["example.platform-selection"].name
    to: #notesByName["mcda.ahp"].name
    label: "uses_method"
  },
  {
    from: #notesByName["example.platform-selection"].name
    to: #notesByName["strategy.growth-scenario-evaluation"].name
    label: "addresses_usecase"
  },
  {
    from: #notesByName["example.platform-selection"].name
    to: #notesByName["technology.platform-selection"].name
    label: "addresses_usecase"
  },
  {
    from: #notesByName["execution.reproducibility"].name
    to: #notesByName["model.validation"].name
    label: "depends_on"
  },
  {
    from: #notesByName["input.format"].name
    to: #notesByName["model.structure"].name
    label: "supports"
  },
  {
    from: #notesByName["model.validation"].name
    to: #notesByName["model.incomplete.data"].name
    label: "reinforces"
  },
  {
    from: #notesByName["model.documentation"].name
    to: #notesByName["model.structure"].name
    label: "documents"
  },
  {
    from: #notesByName["planning.lifecycle-decision"].name
    to: #notesByName["planning.long-term-option-evaluation"].name
    label: "relates_to"
  },
  {
    from: #notesByName["planning.long-term-option-evaluation"].name
    to: #notesByName["decision.robust-choice-identification"].name
    label: "supports"
  },
  {
    from: #notesByName["policy.policy-option-analysis"].name
    to: #notesByName["decision.multi-criteria-ranking"].name
    label: "supports"
  },
  {
    from: #notesByName["product.feature-prioritization"].name
    to: #notesByName["decision.multi-criteria-ranking"].name
    label: "supports"
  },
  {
    from: #notesByName["product.roadmap-planning"].name
    to: #notesByName["planning.long-term-option-evaluation"].name
    label: "supports"
  },
  {
    from: #notesByName["scenario.aggregation.policy"].name
    to: #notesByName["analysis.robustness.method"].name
    label: "references_analysis"
  },
  {
    from: #notesByName["scenario.constraints"].name
    to: #notesByName["model.validation"].name
    label: "depends_on"
  },
  {
    from: #notesByName["scenario.isolation"].name
    to: #notesByName["example.hosting-choice"].name
    label: "illustrated_by"
  },
  {
    from: #notesByName["scenario.isolation"].name
    to: #notesByName["example.platform-selection"].name
    label: "illustrated_by"
  },
  {
    from: #notesByName["strategy.investment-decision"].name
    to: #notesByName["decision.multi-criteria-ranking"].name
    label: "supports"
  },
  {
    from: #notesByName["system.extensibility.methods"].name
    to: #notesByName["mcda.electre"].name
    label: "enables"
  },
  {
    from: #notesByName["system.extensibility.methods"].name
    to: #notesByName["mcda.promethee"].name
    label: "enables"
  },
  {
    from: #notesByName["system.extensibility.methods"].name
    to: #notesByName["mcda.topsis"].name
    label: "enables"
  },
  {
    from: #notesByName["system.extensibility.methods"].name
    to: #notesByName["mcda.vikor"].name
    label: "enables"
  },
  {
    from: #notesByName["technology.architecture-choice"].name
    to: #notesByName["infrastructure.system-design-selection"].name
    label: "relates_to"
  },
  {
    from: #notesByName["ux.model.guidance"].name
    to: #notesByName["model.structure"].name
    label: "supports"
  },
  {
    from: #notesByName["vendor.supplier-selection"].name
    to: #notesByName["vendor.service-provider-comparison"].name
    label: "relates_to"
  },
]
