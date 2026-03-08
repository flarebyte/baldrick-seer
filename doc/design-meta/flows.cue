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
            "call.reports.generate.shared-validation",
            "call.validation.input-config.load-cue-config",
            "call.validation.input-config.validate-model",
            "call.reports.generate.compute-ahp-weights",
            "call.reports.generate.rank-alternatives-topsis",
            "call.reports.generate.render-output",
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
    markdown: "Top-level CLI call flow for validating an input configuration file before any decision analysis runs."
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
    markdown: "Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data."
  },
  {
    name: "call.reports.generate"
    title: "Generate Reports Call"
    labels: ["call", "flow", "implementation"]
    markdown: "Top-level CLI call flow for generating reports from an input decision model."
  },
  {
    name: "call.reports.generate.parse-args"
    title: "Parse Report Arguments"
    labels: ["call", "flow", "implementation"]
    markdown: "Parse CLI arguments for report generation, including the config path, requested report names, and output options."
  },
  {
    name: "call.reports.generate.shared-validation"
    title: "Reuse Shared Validation Flow"
    labels: ["call", "flow", "implementation"]
    markdown: "Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs."
  },
  {
    name: "call.reports.generate.compute-ahp-weights"
    title: "Compute Criteria Weights with AHP"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Transform pairwise scenario preferences into normalized criteria weights using Analytic Hierarchy Process."
  },
  {
    name: "call.reports.generate.rank-alternatives-topsis"
    title: "Rank Alternatives with TOPSIS"
    labels: ["call", "flow", "implementation", "method"]
    markdown: "Use the validated evaluations and AHP-derived weights to rank alternatives with TOPSIS."
  },
  {
    name: "call.reports.generate.render-output"
    title: "Render Requested Reports"
    labels: ["call", "flow", "implementation"]
    markdown: "Render the requested markdown, JSON, or CSV reports from the computed ranking results."
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
    from: "call.reports.generate"
    to: "call.reports.generate.parse-args"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.parse-args"
    to: "call.reports.generate.shared-validation"
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: "call.reports.generate.shared-validation"
    to: "call.reports.generate.compute-ahp-weights"
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
    from: "call.reports.generate.compute-ahp-weights"
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
