package flyb

source: "baldrick-seer.flow"
name: "baldrick-seer-flows"
modules: ["design", "flow"]

reports: [{
  title: "Execution Flows"
  filepath: "../design/flows.md"
  description: "Call-oriented CLI graphs. Normative semantics are defined in the specification."
  sections: [
    {
      title: "Validation flows"
      description: "Graph view for validating an input config file. Refer to the normative specification for validation rules."
      sections: [
        {
          title: "Input config validation graph"
          description: "Text graph for the validate-config call chain."
          arguments: [
            "graph-subject-label=call",
            "graph-edge-label=delegate_to",
            "graph-start-node=call.validation.input-config",
            "graph-renderer=markdown-text",
            "cycle-policy=disallow",
          ]
        },
      ]
    },
    {
      title: "Report generation flows"
      description: "Graph view for generating reports from a validated input config. Refer to the normative specification for scoring and failure behavior."
      sections: [
        {
          title: "Report generation graph"
          description: "Text graph for the report-generation call chain, reusing the shared validation path."
          arguments: [
            "graph-subject-label=call",
            "graph-edge-label=delegate_to",
            "graph-start-node=call.reports.generate",
            "graph-renderer=markdown-text",
            "cycle-policy=disallow",
          ]
        },
      ]
    },
  ]
}]

notes: [
  {
    name: "call.validation.input-config"
    title: "Validate Input Config Call"
    labels: ["call", "flow", "implementation"]
    markdown: "Top-level validate-command call chain. See the normative specification for authoritative command semantics."
  },
  {
    name: "call.validation.input-config.parse-args"
    title: "Parse Validation Arguments"
    labels: ["call", "flow", "implementation"]
    markdown: "Parse CLI arguments for the validate command, including the config path and output flags."
  },
  {
    name: "call.validation.input-config.load-cue-config"
    title: "Load CUE Config"
    labels: ["call", "cue", "flow", "implementation"]
    markdown: "Load and evaluate the CUE configuration package before validation."
  },
  {
    name: "call.validation.input-config.validate-model"
    title: "Validate Config Model"
    labels: ["call", "flow", "implementation"]
    markdown: "Run the shared validation stage and emit diagnostics. See the normative specification for exact validation behavior."
  },
  {
    name: "call.validation.input-config.validate-model.check-structure"
    title: "Check Config Structure"
    labels: ["call", "flow", "implementation", "validation"]
    markdown: "Check that the loaded config matches the expected top-level shape, required sections, and field types after CUE evaluation."
  },
  {
    name: "call.validation.input-config.validate-model.check-references"
    title: "Check Named References"
    labels: ["call", "flow", "implementation", "validation"]
    markdown: "Check that all named references resolve, including criteria names, scenario names, alternative names, and report focus selectors."
  },
  {
    name: "call.validation.input-config.validate-model.check-pairwise-comparisons"
    title: "Check Pairwise Comparisons"
    labels: ["call", "flow", "implementation", "validation"]
    markdown: "Validate AHP pairwise-comparison coverage and canonical comparison structure."
  },
  {
    name: "call.validation.input-config.validate-model.check-evaluation-coverage"
    title: "Check Evaluation Coverage"
    labels: ["call", "flow", "implementation", "validation"]
    markdown: "Validate evaluation coverage and supported value forms for active criteria."
  },
  {
    name: "call.validation.input-config.validate-model.check-constraints"
    title: "Check Scenario Constraints"
    labels: ["call", "flow", "implementation", "validation"]
    markdown: "Validate constraint operator and value compatibility."
  },
  {
    name: "call.validation.input-config.validate-model.check-report-definitions"
    title: "Check Report Definitions"
    labels: ["call", "flow", "implementation", "validation"]
    markdown: "Validate report definitions, focus selectors, and report arguments."
  },
  {
    name: "call.reports.generate"
    title: "Generate Reports Call"
    labels: ["call", "flow", "implementation"]
    markdown: "Top-level report-generation call chain. See the normative specification for failure handling and output semantics."
  },
  {
    name: "call.reports.generate.parse-args"
    title: "Parse Report Arguments"
    labels: ["call", "flow", "implementation"]
    markdown: "Parse CLI arguments for report generation, including the config path, requested report names, and output options."
  },
  {
    name: "call.reports.generate.select-reports"
    title: "Select Requested Reports"
    labels: ["call", "flow", "implementation"]
    markdown: "Resolve which report definitions should run after CLI filtering."
  },
  {
    name: "call.reports.generate.shared-validation"
    title: "Reuse Shared Validation Flow"
    labels: ["call", "flow", "implementation"]
    markdown: "Reuse the shared validation stage before scoring."
  },
  {
    name: "call.reports.generate.build-ahp-inputs"
    title: "Build AHP Inputs"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Prepare validated pairwise comparisons for AHP weight computation."
  },
  {
    name: "call.reports.generate.compute-ahp-weights"
    title: "Compute Criteria Weights with AHP"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Compute scenario-local criterion weights from validated pairwise comparisons."
  },
  {
    name: "call.reports.generate.select-ranking-strategy"
    title: "Select Ranking Strategy"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Select the ranking branch after AHP weighting."
  },
  {
    name: "call.reports.generate.build-topsis-inputs"
    title: "Build TOPSIS Inputs"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Assemble TOPSIS decision matrices from validated evaluations, polarity, and AHP-derived weights."
  },
  {
    name: "call.reports.generate.rank-alternatives-topsis"
    title: "Rank Alternatives with TOPSIS"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Rank alternatives with TOPSIS using validated evaluations and scenario-local criterion weights."
  },
  {
    name: "call.reports.generate.future-rank-electre"
    title: "Future Option: Rank with ELECTRE"
    labels: ["call", "flow", "future", "method"]
    markdown: "Potential future branch for ELECTRE-based ranking."
  },
  {
    name: "call.reports.generate.future-rank-topsis-sensitivity"
    title: "Future Option: TOPSIS with Sensitivity Analysis"
    labels: ["call", "flow", "future", "method"]
    markdown: "Potential future branch that complements TOPSIS with sensitivity analysis."
  },
  {
    name: "call.reports.generate.render-output"
    title: "Render Requested Reports"
    labels: ["call", "flow", "implementation"]
    markdown: "Render requested outputs after validation succeeds and ranking results exist."
  },
  {
    name: "call.reports.generate.render-output.render-markdown"
    title: "Render Markdown Report"
    labels: ["call", "flow", "implementation"]
    markdown: "Render narrative markdown output for human readers, including rankings, explanations, and scenario summaries."
  },
  {
    name: "call.reports.generate.render-output.render-json"
    title: "Render JSON Report"
    labels: ["call", "flow", "implementation"]
    markdown: "Render JSON ranking output after successful validation and scoring."
  },
  {
    name: "call.reports.generate.render-output.render-csv"
    title: "Render CSV Report"
    labels: ["call", "flow", "implementation"]
    markdown: "Render flat tabular CSV output for spreadsheet analysis and data exchange."
  },
]

relationships: [
  {
    from: "call.validation.input-config"
    to: "call.validation.input-config.parse-args"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.validation.input-config.parse-args"
    to: "call.validation.input-config.load-cue-config"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.validation.input-config.load-cue-config"
    to: "call.validation.input-config.validate-model"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.validation.input-config.validate-model"
    to: "call.validation.input-config.validate-model.check-structure"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.validation.input-config.validate-model.check-structure"
    to: "call.validation.input-config.validate-model.check-references"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.validation.input-config.validate-model.check-references"
    to: "call.validation.input-config.validate-model.check-pairwise-comparisons"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.validation.input-config.validate-model.check-pairwise-comparisons"
    to: "call.validation.input-config.validate-model.check-evaluation-coverage"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.validation.input-config.validate-model.check-evaluation-coverage"
    to: "call.validation.input-config.validate-model.check-constraints"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.validation.input-config.validate-model.check-constraints"
    to: "call.validation.input-config.validate-model.check-report-definitions"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate"
    to: "call.reports.generate.parse-args"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.parse-args"
    to: "call.reports.generate.select-reports"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.select-reports"
    to: "call.reports.generate.shared-validation"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.shared-validation"
    to: "call.reports.generate.build-ahp-inputs"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.shared-validation"
    to: "call.validation.input-config.load-cue-config"
    label: "reuses"
  },
  {
    from: "call.reports.generate.shared-validation"
    to: "call.validation.input-config.validate-model"
    label: "reuses"
  },
  {
    from: "call.reports.generate.build-ahp-inputs"
    to: "call.reports.generate.compute-ahp-weights"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.compute-ahp-weights"
    to: "call.reports.generate.select-ranking-strategy"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.select-ranking-strategy"
    to: "call.reports.generate.build-topsis-inputs"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.select-ranking-strategy"
    to: "call.reports.generate.future-rank-electre"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.select-ranking-strategy"
    to: "call.reports.generate.future-rank-topsis-sensitivity"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.build-topsis-inputs"
    to: "call.reports.generate.rank-alternatives-topsis"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.rank-alternatives-topsis"
    to: "call.reports.generate.render-output"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.render-output"
    to: "call.reports.generate.render-output.render-markdown"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.render-output"
    to: "call.reports.generate.render-output.render-json"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.render-output"
    to: "call.reports.generate.render-output.render-csv"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
]

argumentRegistry: {
  version: "1"
  arguments: [
    {
      name: "graph-subject-label"
      valueType: "string"
      scopes: ["h3-section"]
    },
    {
      name: "graph-edge-label"
      valueType: "string"
      scopes: ["h3-section"]
    },
    {
      name: "graph-start-node"
      valueType: "string"
      scopes: ["h3-section"]
    },
    {
      name: "graph-renderer"
      valueType: "enum"
      scopes: ["h3-section", "note"]
      allowedValues: ["markdown-text", "mermaid"]
      defaultValue: "markdown-text"
    },
    {
      name: "cycle-policy"
      valueType: "enum"
      scopes: ["h3-section"]
      allowedValues: ["allow", "disallow"]
      defaultValue: "allow"
    },
    {
      name: "mermaid-direction"
      valueType: "enum"
      scopes: ["h3-section", "note"]
      allowedValues: ["TD", "LR"]
      defaultValue: "TD"
    },
  ]
}
