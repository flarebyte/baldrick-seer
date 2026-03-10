package flyb

source: "baldrick-seer.flow"
name: "baldrick-seer-flows"
modules: ["design", "flow"]

reports: [{
  title: "Execution Flows"
  filepath: "../design/flows.md"
  description: "Call-oriented flows for CLI validation and report generation."
  sections: [
    {
      title: "Validation flows"
      description: "Graph view for validating an input config file."
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
        {
          title: "Input config validation notes"
          notes: [
            "call.validation.input-config",
            "call.validation.input-config.parse-args",
            "call.validation.input-config.load-cue-config",
            "call.validation.input-config.validate-model",
            "call.validation.input-config.validate-model.check-structure",
            "call.validation.input-config.validate-model.check-references",
            "call.validation.input-config.validate-model.check-pairwise-comparisons",
            "call.validation.input-config.validate-model.check-evaluation-coverage",
            "call.validation.input-config.validate-model.check-constraints",
            "call.validation.input-config.validate-model.check-report-definitions",
          ]
        },
      ]
    },
    {
      title: "Report generation flows"
      description: "Graph view for generating reports from a validated input config."
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
        {
          title: "Report generation notes"
          notes: [
            "call.reports.generate",
            "call.reports.generate.parse-args",
            "call.reports.generate.select-reports",
            "call.reports.generate.shared-validation",
            "call.validation.input-config.load-cue-config",
            "call.validation.input-config.validate-model",
            "call.reports.generate.build-ahp-inputs",
            "call.reports.generate.compute-ahp-weights",
            "call.reports.generate.select-ranking-strategy",
            "call.reports.generate.build-topsis-inputs",
            "call.reports.generate.rank-alternatives-topsis",
            "call.reports.generate.future-rank-electre",
            "call.reports.generate.future-rank-topsis-sensitivity",
            "call.reports.generate.render-output",
            "call.reports.generate.render-output.render-markdown",
            "call.reports.generate.render-output.render-json",
            "call.reports.generate.render-output.render-csv",
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
    markdown: "Top-level CLI call flow for validating an input configuration file and returning validation results only, without scoring or report generation."
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
    markdown: "Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value."
  },
  {
    name: "call.validation.input-config.validate-model"
    title: "Validate Config Model"
    labels: ["call", "flow", "implementation"]
    markdown: "Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data. For the `validate` command, this is the terminal result of the command."
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
    markdown: "Check that each scenario using AHP provides pairwise comparisons only between known active criteria, never compares a criterion with itself, and includes exactly one canonical comparison for every unordered pair of distinct active criteria. Reject duplicate comparisons, inverse duplicates, or any missing pair."
  },
  {
    name: "call.validation.input-config.validate-model.check-evaluation-coverage"
    title: "Check Evaluation Coverage"
    labels: ["call", "flow", "implementation", "validation"]
    markdown: "Check that evaluations reference known scenarios and alternatives and provide supported v1 criterion values for each scenario's active criteria: measurable numbers, integer ordinals, or booleans with only `true` and `false` values."
  },
  {
    name: "call.validation.input-config.validate-model.check-constraints"
    title: "Check Scenario Constraints"
    labels: ["call", "flow", "implementation", "validation"]
    markdown: "Check that each scenario constraint uses an operator and value compatible with the referenced criterion type: number criteria allow numeric values with `<=`, `>=`, `=`, or `!=`; ordinal criteria allow integer values with `<=`, `>=`, `=`, or `!=`; boolean criteria allow only `=` or `!=` with `true` or `false`. Invalid operator/type combinations must raise a validation error."
  },
  {
    name: "call.validation.input-config.validate-model.check-report-definitions"
    title: "Check Report Definitions"
    labels: ["call", "flow", "implementation", "validation"]
    markdown: "Check that report definitions use supported formats, valid focus selectors, and strictly validated report arguments. In v1 every report argument must use `key=value`, unknown arguments are validation errors, argument names must be allowed globally or for the selected format, format-specific arguments must match the report format, invalid values must be rejected, and duplicate keys are invalid unless explicitly defined otherwise."
  },
  {
    name: "call.reports.generate"
    title: "Generate Reports Call"
    labels: ["call", "flow", "implementation"]
    markdown: "Top-level CLI call flow for generating ranking reports from an input decision model. The command reuses the shared validation path and fails fast if the model is invalid."
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
    markdown: "Resolve which report definitions should run, applying any CLI filtering by report name or output target."
  },
  {
    name: "call.reports.generate.shared-validation"
    title: "Reuse Shared Validation Flow"
    labels: ["call", "flow", "implementation"]
    markdown: "Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs. If validation fails, report generation stops immediately and no ranking report is produced."
  },
  {
    name: "call.reports.generate.build-ahp-inputs"
    title: "Build AHP Inputs"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Collect the validated full pairwise comparison set for each scenario into the normalized input structures needed for AHP computation of scenario-local criterion weights."
  },
  {
    name: "call.reports.generate.compute-ahp-weights"
    title: "Compute Criteria Weights with AHP"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Transform pairwise criterion comparisons within each scenario into normalized scenario-local criterion weights using Analytic Hierarchy Process."
  },
  {
    name: "call.reports.generate.select-ranking-strategy"
    title: "Select Ranking Strategy"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Select the ranking pipeline after computing scenario-local criterion weights with AHP. The current default path is TOPSIS, while v2 may add ELECTRE or TOPSIS followed by sensitivity analysis."
  },
  {
    name: "call.reports.generate.build-topsis-inputs"
    title: "Build TOPSIS Inputs"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Combine validated evaluations, criterion polarity, and AHP-derived scenario-local criterion weights into the decision matrices required by TOPSIS."
  },
  {
    name: "call.reports.generate.rank-alternatives-topsis"
    title: "Rank Alternatives with TOPSIS"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Use the validated evaluations and scenario-local criterion weights derived with AHP to rank alternatives with TOPSIS."
  },
  {
    name: "call.reports.generate.future-rank-electre"
    title: "Future Option: Rank with ELECTRE"
    labels: ["call", "flow", "future", "method"]
    markdown: "Potential v2 branch where the validated model is ranked with ELECTRE instead of TOPSIS."
  },
  {
    name: "call.reports.generate.future-rank-topsis-sensitivity"
    title: "Future Option: TOPSIS with Sensitivity Analysis"
    labels: ["call", "flow", "future", "method"]
    markdown: "Potential v2 branch where TOPSIS ranking is complemented by sensitivity analysis to assess robustness."
  },
  {
    name: "call.reports.generate.render-output"
    title: "Render Requested Reports"
    labels: ["call", "flow", "implementation"]
    markdown: "Render the requested markdown, JSON, or CSV outputs only after validation succeeds and ranking results are computed. Invalid models do not reach report rendering."
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
    markdown: "Render machine-readable JSON ranking output for automation, downstream processing, and reproducibility only when validation succeeds. If JSON output is requested and validation fails, the command may emit structured diagnostics as an error payload or via stderr, but that output is not a successful ranking report."
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
