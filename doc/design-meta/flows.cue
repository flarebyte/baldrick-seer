package flyb

source: "baldrick-seer.flow"
name: "baldrick-seer-flows"
modules: ["design", "flow"]

reports: [{
  title: "Execution Flows"
  filepath: "../design/flows.md"
  description: "Call-oriented flows for CLI execution and validation."
  sections: [{
    title: "Validation flows"
    description: "Graph view for validating an input config file."
    sections: [
      {
        title: "Input config validation graph"
        description: "Mermaid graph for the validate-config call chain."
        arguments: [
          "graph-subject-label=call",
          "graph-edge-label=delegate_to",
          "graph-start-node=call.validation.input-config",
          "graph-renderer=mermaid",
          "mermaid-direction=TD",
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
  }]
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
      name: "mermaid-direction"
      valueType: "enum"
      scopes: ["h3-section", "note"]
      allowedValues: ["TD", "LR"]
      defaultValue: "TD"
    },
  ]
}
