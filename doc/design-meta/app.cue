package flyb

source: "baldrick-seer.design"
name: "baldrick-seer"
modules: ["design"]

#notesByName: {
  "analysis.robustness": {
    name: "analysis.robustness"
    title: "Sensitivity and Robustness Analysis"
    labels: ["design", "implementation"]
    markdown: "Support testing how changes in criteria importance or scenario assumptions affect the final ranking so users can see whether a result is stable or fragile."
  }
  "analysis.robustness.method": {
    name: "analysis.robustness.method"
    title: "Robustness Analysis"
    labels: ["analysis", "design", "method"]
    markdown: "Evaluate how stable a decision remains when assumptions, scenarios, or parameter ranges vary."
  }
  "analysis.sensitivity": {
    name: "analysis.sensitivity"
    title: "Sensitivity Analysis"
    labels: ["analysis", "design", "method"]
    markdown: "Evaluate how changes in weights or inputs affect the ranking of alternatives."
  }
  "cli.output.machine": {
    name: "cli.output.machine"
    title: "Structured Output for Automation"
    labels: ["design", "implementation"]
    markdown: "Provide machine-readable output such as JSON in addition to human-readable summaries."
  }
  "cli.output.readability": {
    name: "cli.output.readability"
    title: "Readable CLI Output"
    labels: ["design", "implementation"]
    markdown: "Present results in a clear terminal-friendly format with summaries, tables, and scenario breakdowns."
  }
  "call.validation.input-config": {
    name: "call.validation.input-config"
    title: "Validate Input Config Call"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Top-level CLI call flow for validating an input configuration file before any decision analysis runs."
  }
  "call.validation.input-config.parse-args": {
    name: "call.validation.input-config.parse-args"
    title: "Parse Validation Arguments"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Parse CLI arguments for the validate command, including the config path and output flags."
  }
  "call.validation.input-config.load-cue-config": {
    name: "call.validation.input-config.load-cue-config"
    title: "Load CUE Config"
    labels: ["call", "cue", "design", "flow", "implementation"]
    markdown: "Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value."
  }
  "call.validation.input-config.validate-model": {
    name: "call.validation.input-config.validate-model"
    title: "Validate Config Model"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data."
  }
  "criteria.pairwise.clarity": {
    name: "criteria.pairwise.clarity"
    title: "Clear Representation of Pairwise Judgments"
    labels: ["design", "implementation"]
    markdown: "Represent pairwise comparisons explicitly with named criteria instead of positional matrices so humans and AI can validate and generate them."
  }
  "criteria.semantic.consistency": {
    name: "criteria.semantic.consistency"
    title: "Consistent Criteria Interpretation"
    labels: ["design", "implementation"]
    markdown: "Keep each criterion semantically stable across scenarios even when its importance changes."
  }
  "decision.explainability": {
    name: "decision.explainability"
    title: "Explainable Results"
    labels: ["design", "implementation"]
    markdown: "Explain ranking outputs in terms of criteria influence and scenario differences."
  }
  "decision.multi-criteria-ranking": {
    name: "decision.multi-criteria-ranking"
    title: "General Multi-Criteria Ranking"
    labels: ["design", "usecase"]
    markdown: "Rank competing alternatives using multiple evaluation criteria and scenario-based priorities."
  }
  "decision.robust-choice-identification": {
    name: "decision.robust-choice-identification"
    title: "Robust Choice Identification"
    labels: ["design", "usecase"]
    markdown: "Identify alternatives that perform consistently well across different scenarios or changing assumptions."
  }
  "decision.traceability": {
    name: "decision.traceability"
    title: "Traceable Decision Process"
    labels: ["design", "implementation"]
    markdown: "Show the reasoning path from inputs to outputs, including scenario weights, criteria importance, and contribution of each factor."
  }
  "example.hosting-choice": {
    name: "example.hosting-choice"
    title: "Hosting Choice Example"
    labels: ["design", "example"]
    markdown: "Minimal scenario-based MCDA example that compares two hosting providers across lean-startup and regulated-growth scenarios using equal-average aggregation."
  }
  "example.hosting-choice.lean-startup": {
    name: "example.hosting-choice.lean-startup"
    title: "Lean Startup Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Scenario emphasizing cost and speed when budget pressure matters more than peak performance."
  }
  "example.hosting-choice.regulated-growth": {
    name: "example.hosting-choice.regulated-growth"
    title: "Regulated Growth Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Scenario adding compliance as a dominant concern alongside cost and speed."
  }
  "example.platform-selection": {
    name: "example.platform-selection"
    title: "Platform Selection Example"
    labels: ["design", "example"]
    markdown: "Scenario-based MCDA example that compares candidate platforms across startup, unicorn, and established-enterprise contexts using weighted-average aggregation."
  }
  "example.platform-selection.established-enterprise": {
    name: "example.platform-selection.established-enterprise"
    title: "Established Enterprise Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Mature-organization scenario where reliability and compliance dominate cost and future flexibility."
  }
  "example.platform-selection.startup": {
    name: "example.platform-selection.startup"
    title: "Startup Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Early-stage scenario where cost and time to market matter more than enterprise controls."
  }
  "example.platform-selection.unicorn": {
    name: "example.platform-selection.unicorn"
    title: "Unicorn Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Hyper-growth scenario where scalability becomes dominant while cost still matters."
  }
  "execution.reproducibility": {
    name: "execution.reproducibility"
    title: "Reproducible Decision Runs"
    labels: ["design", "implementation"]
    markdown: "Running the same model with the same inputs should always produce identical results for auditing and comparison."
  }
  "infrastructure.system-design-selection": {
    name: "infrastructure.system-design-selection"
    title: "System Design Selection"
    labels: ["design", "usecase"]
    markdown: "Compare system designs where trade-offs exist between cost, scalability, and reliability."
  }
  "input.format": {
    name: "input.format"
    title: "Human and AI Friendly Input Format"
    labels: ["design", "implementation"]
    markdown: "Use a semantic format such as JSON or YAML that is easy for humans and AI systems to read and generate."
  }
  "mcda.ahp": {
    name: "mcda.ahp"
    title: "Analytic Hierarchy Process (AHP)"
    labels: ["design", "method"]
    markdown: "Derive criteria weights from pairwise comparisons and turn qualitative judgments into a consistent numerical weighting system."
  }
  "mcda.electre": {
    name: "mcda.electre"
    title: "ELECTRE Outranking Method"
    labels: ["design", "method"]
    markdown: "Use concordance and discordance reasoning to determine whether one alternative sufficiently outranks another."
  }
  "mcda.promethee": {
    name: "mcda.promethee"
    title: "PROMETHEE"
    labels: ["design", "method"]
    markdown: "Compare alternatives pairwise with preference functions to produce a transparent ranking."
  }
  "mcda.topsis": {
    name: "mcda.topsis"
    title: "TOPSIS"
    labels: ["design", "method"]
    markdown: "Rank alternatives by their distance from an ideal best and an ideal worst solution."
  }
  "mcda.vikor": {
    name: "mcda.vikor"
    title: "VIKOR"
    labels: ["design", "method"]
    markdown: "Identify a compromise solution that balances group utility and individual regret."
  }
  "model.documentation": {
    name: "model.documentation"
    title: "Model Documentation"
    labels: ["design", "implementation"]
    markdown: "Allow decision models to carry descriptions, notes, and justifications for comparisons and values."
  }
  "model.incomplete.data": {
    name: "model.incomplete.data"
    title: "Handling Incomplete Information"
    labels: ["design", "implementation"]
    markdown: "Detect missing comparisons or evaluation values and provide clear feedback, with any inferred values marked explicitly."
  }
  "model.structure": {
    name: "model.structure"
    title: "Decision Model Structure"
    labels: ["design", "implementation"]
    markdown: "Represent the decision problem with clear structures for criteria, alternatives, and scenarios that remain understandable to humans and AI."
  }
  "model.validation": {
    name: "model.validation"
    title: "Model Validation"
    labels: ["design", "implementation"]
    markdown: "Validate referenced criteria, pairwise comparison completeness, and alternative evaluation coverage before computation."
  }
  "planning.lifecycle-decision": {
    name: "planning.lifecycle-decision"
    title: "Lifecycle Decision Support"
    labels: ["design", "usecase"]
    markdown: "Compare options that must remain effective throughout different stages of organizational or system development."
  }
  "planning.long-term-option-evaluation": {
    name: "planning.long-term-option-evaluation"
    title: "Long-Term Option Evaluation"
    labels: ["design", "usecase"]
    markdown: "Evaluate alternatives that must perform well across multiple possible future environments."
  }
  "policy.policy-option-analysis": {
    name: "policy.policy-option-analysis"
    title: "Policy Option Analysis"
    labels: ["design", "usecase"]
    markdown: "Support evaluation of policy alternatives where multiple criteria such as impact, feasibility, and cost must be considered."
  }
  "product.feature-prioritization": {
    name: "product.feature-prioritization"
    title: "Product Feature Prioritization"
    labels: ["design", "usecase"]
    markdown: "Rank product features using multiple criteria such as user value, development effort, and strategic importance."
  }
  "product.roadmap-planning": {
    name: "product.roadmap-planning"
    title: "Product Roadmap Planning"
    labels: ["design", "usecase"]
    markdown: "Evaluate product initiatives across different market or growth scenarios to support long-term planning."
  }
  "scenario.aggregation.policy": {
    name: "scenario.aggregation.policy"
    title: "Scenario Aggregation Strategy"
    labels: ["design", "implementation"]
    markdown: "Define how multiple scenarios are combined into a final decision, such as equal averaging, weighted scenarios, or robustness-focused approaches."
  }
  "scenario.constraints": {
    name: "scenario.constraints"
    title: "Constraint Enforcement"
    labels: ["design", "implementation"]
    markdown: "Allow scenarios to define hard requirements that can exclude alternatives before ranking."
  }
  "scenario.isolation": {
    name: "scenario.isolation"
    title: "Scenario Isolation"
    labels: ["design", "implementation"]
    markdown: "Evaluate each scenario independently with its own priorities and candidate evaluations."
  }
  "strategy.growth-scenario-evaluation": {
    name: "strategy.growth-scenario-evaluation"
    title: "Growth Scenario Evaluation"
    labels: ["design", "usecase"]
    markdown: "Assess strategic options under different business growth trajectories such as startup, rapid expansion, or mature operations."
  }
  "strategy.investment-decision": {
    name: "strategy.investment-decision"
    title: "Strategic Investment Decision"
    labels: ["design", "usecase"]
    markdown: "Compare investment alternatives considering financial return, risk exposure, and long-term strategic impact."
  }
  "system.extensibility.methods": {
    name: "system.extensibility.methods"
    title: "Extensible Decision Methods"
    labels: ["design", "implementation"]
    markdown: "Design the system so additional MCDA methods can be added later without redesigning the data model or CLI interface."
  }
  "technology.architecture-choice": {
    name: "technology.architecture-choice"
    title: "Software Architecture Decision"
    labels: ["design", "usecase"]
    markdown: "Evaluate architectural approaches under different system growth conditions, performance requirements, and reliability expectations."
  }
  "technology.infrastructure-strategy": {
    name: "technology.infrastructure-strategy"
    title: "Infrastructure Strategy Planning"
    labels: ["design", "usecase"]
    markdown: "Assess infrastructure alternatives such as cloud providers or deployment models under varying operational scenarios."
  }
  "technology.platform-selection": {
    name: "technology.platform-selection"
    title: "Technology Platform Selection"
    labels: ["design", "usecase"]
    markdown: "Compare multiple technology platforms across operational scenarios such as startup, scale-up, and enterprise maturity."
  }
  "ux.model.guidance": {
    name: "ux.model.guidance"
    title: "Guidance for Model Creation"
    labels: ["design", "implementation"]
    markdown: "Provide prompts and guidance that help users define criteria, comparisons, and scenario descriptions with fewer modeling errors."
  }
  "vendor.service-provider-comparison": {
    name: "vendor.service-provider-comparison"
    title: "Service Provider Comparison"
    labels: ["design", "usecase"]
    markdown: "Compare service providers where priorities may change depending on scale, regulatory environment, or organizational maturity."
  }
  "vendor.supplier-selection": {
    name: "vendor.supplier-selection"
    title: "Supplier or Vendor Selection"
    labels: ["design", "usecase"]
    markdown: "Evaluate competing suppliers using multiple criteria such as cost, reliability, and service quality under different operating conditions."
  }
}

reports: [
  {
    title: "Design Overview"
    filepath: "../design/overview.md"
    description: "High-level overview of the baldrick-seer decision-model design."
    sections: [
      {
        title: "Scope"
        description: "Primary use cases and example design anchors."
        sections: [
          {
            title: "Primary use cases"
            notes: [
              #notesByName["decision.multi-criteria-ranking"].name,
              #notesByName["decision.robust-choice-identification"].name,
              #notesByName["infrastructure.system-design-selection"].name,
              #notesByName["technology.infrastructure-strategy"].name,
              #notesByName["technology.platform-selection"].name,
            ]
          },
          {
            title: "Reference examples"
            notes: [
              #notesByName["example.hosting-choice"].name,
              #notesByName["example.platform-selection"].name,
            ]
          },
        ]
      },
      {
        title: "Method references"
        description: "Decision methods and analysis techniques mentioned by the design."
        sections: [
          {
            title: "Core methods"
            notes: [
              #notesByName["analysis.robustness.method"].name,
              #notesByName["analysis.sensitivity"].name,
              #notesByName["mcda.ahp"].name,
              #notesByName["mcda.electre"].name,
              #notesByName["mcda.promethee"].name,
              #notesByName["mcda.topsis"].name,
              #notesByName["mcda.vikor"].name,
            ]
          },
        ]
      },
    ]
  },
  {
    title: "Example Models"
    filepath: "../design/examples.md"
    description: "Scenario-based MCDA examples translated from the TypeScript design fixtures."
    sections: [
      {
        title: "Platform selection"
        description: "Three-scenario example for selecting a platform across company maturity stages."
        sections: [
          {
            title: "Model"
            notes: [
              #notesByName["example.platform-selection"].name,
            ]
          },
          {
            title: "Scenarios"
            notes: [
              #notesByName["example.platform-selection.established-enterprise"].name,
              #notesByName["example.platform-selection.startup"].name,
              #notesByName["example.platform-selection.unicorn"].name,
            ]
          },
        ]
      },
      {
        title: "Hosting choice"
        description: "Minimal example showing the same shape with fewer criteria and alternatives."
        sections: [
          {
            title: "Model"
            notes: [
              #notesByName["example.hosting-choice"].name,
            ]
          },
          {
            title: "Scenarios"
            notes: [
              #notesByName["example.hosting-choice.lean-startup"].name,
              #notesByName["example.hosting-choice.regulated-growth"].name,
            ]
          },
        ]
      },
    ]
  },
  {
    title: "Implementation Considerations"
    filepath: "../design/implementation.md"
    description: "Implementation guidance and method references for the CLI and model."
    sections: [
      {
        title: "Modeling guidance"
        description: "Recommendations about the decision-model shape and validation."
        sections: [
          {
            title: "Model structure"
            notes: [
              #notesByName["criteria.pairwise.clarity"].name,
              #notesByName["criteria.semantic.consistency"].name,
              #notesByName["input.format"].name,
              #notesByName["model.documentation"].name,
              #notesByName["model.incomplete.data"].name,
              #notesByName["model.structure"].name,
              #notesByName["model.validation"].name,
              #notesByName["scenario.aggregation.policy"].name,
              #notesByName["scenario.constraints"].name,
              #notesByName["scenario.isolation"].name,
              #notesByName["system.extensibility.methods"].name,
            ]
          },
        ]
      },
      {
        title: "User experience and output"
        description: "Guidance for readable, reproducible, and automatable execution."
        sections: [
          {
            title: "CLI and explainability"
            notes: [
              #notesByName["analysis.robustness"].name,
              #notesByName["cli.output.machine"].name,
              #notesByName["cli.output.readability"].name,
              #notesByName["decision.explainability"].name,
              #notesByName["decision.traceability"].name,
              #notesByName["execution.reproducibility"].name,
              #notesByName["ux.model.guidance"].name,
            ]
          },
        ]
      },
      {
        title: "Validation call flow"
        description: "Early CLI execution path for reading and validating an input config file."
        sections: [
          {
            title: "Input config validation"
            notes: [
              #notesByName["call.validation.input-config"].name,
              #notesByName["call.validation.input-config.parse-args"].name,
              #notesByName["call.validation.input-config.load-cue-config"].name,
              #notesByName["call.validation.input-config.validate-model"].name,
              #notesByName["model.validation"].name,
              #notesByName["model.incomplete.data"].name,
            ]
          },
        ]
      },
      {
        title: "Referenced methods"
        description: "Algorithms and analysis techniques explicitly named in the design."
        sections: [
          {
            title: "Algorithms"
            notes: [
              #notesByName["analysis.robustness.method"].name,
              #notesByName["analysis.sensitivity"].name,
              #notesByName["mcda.ahp"].name,
              #notesByName["mcda.electre"].name,
              #notesByName["mcda.promethee"].name,
              #notesByName["mcda.topsis"].name,
              #notesByName["mcda.vikor"].name,
            ]
          },
        ]
      },
    ]
  },
  {
    title: "Execution Flows"
    filepath: "../design/flows.md"
    description: "Call-oriented flows for CLI execution and validation."
    sections: [
      {
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
              #notesByName["call.validation.input-config"].name,
              #notesByName["call.validation.input-config.parse-args"].name,
              #notesByName["call.validation.input-config.load-cue-config"].name,
              #notesByName["call.validation.input-config.validate-model"].name,
              #notesByName["input.format"].name,
              #notesByName["model.validation"].name,
              #notesByName["model.incomplete.data"].name,
            ]
          },
        ]
      },
    ]
  },
  {
    title: "Use Cases"
    filepath: "../design/use-cases.md"
    description: "Use-case catalog translated from the TypeScript design source."
    sections: [
      {
        title: "Decision and planning"
        description: "Core decision support and longer-horizon planning problems."
        sections: [
          {
            title: "Decision support"
            notes: [
              #notesByName["decision.multi-criteria-ranking"].name,
              #notesByName["decision.robust-choice-identification"].name,
              #notesByName["planning.lifecycle-decision"].name,
              #notesByName["planning.long-term-option-evaluation"].name,
              #notesByName["policy.policy-option-analysis"].name,
              #notesByName["strategy.growth-scenario-evaluation"].name,
              #notesByName["strategy.investment-decision"].name,
            ]
          },
        ]
      },
      {
        title: "Technology and product"
        description: "Product, platform, and architecture comparisons."
        sections: [
          {
            title: "Technology and product choices"
            notes: [
              #notesByName["infrastructure.system-design-selection"].name,
              #notesByName["product.feature-prioritization"].name,
              #notesByName["product.roadmap-planning"].name,
              #notesByName["technology.architecture-choice"].name,
              #notesByName["technology.infrastructure-strategy"].name,
              #notesByName["technology.platform-selection"].name,
            ]
          },
        ]
      },
      {
        title: "Vendor comparisons"
        description: "Supplier and service-provider evaluation scenarios."
        sections: [
          {
            title: "Vendors"
            notes: [
              #notesByName["vendor.service-provider-comparison"].name,
              #notesByName["vendor.supplier-selection"].name,
            ]
          },
        ]
      },
    ]
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

notes: [
  #notesByName["analysis.robustness"],
  #notesByName["analysis.robustness.method"],
  #notesByName["analysis.sensitivity"],
  #notesByName["cli.output.machine"],
  #notesByName["cli.output.readability"],
  #notesByName["call.validation.input-config"],
  #notesByName["call.validation.input-config.parse-args"],
  #notesByName["call.validation.input-config.load-cue-config"],
  #notesByName["call.validation.input-config.validate-model"],
  #notesByName["criteria.pairwise.clarity"],
  #notesByName["criteria.semantic.consistency"],
  #notesByName["decision.explainability"],
  #notesByName["decision.multi-criteria-ranking"],
  #notesByName["decision.robust-choice-identification"],
  #notesByName["decision.traceability"],
  #notesByName["example.hosting-choice"],
  #notesByName["example.hosting-choice.lean-startup"],
  #notesByName["example.hosting-choice.regulated-growth"],
  #notesByName["example.platform-selection"],
  #notesByName["example.platform-selection.established-enterprise"],
  #notesByName["example.platform-selection.startup"],
  #notesByName["example.platform-selection.unicorn"],
  #notesByName["execution.reproducibility"],
  #notesByName["infrastructure.system-design-selection"],
  #notesByName["input.format"],
  #notesByName["mcda.ahp"],
  #notesByName["mcda.electre"],
  #notesByName["mcda.promethee"],
  #notesByName["mcda.topsis"],
  #notesByName["mcda.vikor"],
  #notesByName["model.documentation"],
  #notesByName["model.incomplete.data"],
  #notesByName["model.structure"],
  #notesByName["model.validation"],
  #notesByName["planning.lifecycle-decision"],
  #notesByName["planning.long-term-option-evaluation"],
  #notesByName["policy.policy-option-analysis"],
  #notesByName["product.feature-prioritization"],
  #notesByName["product.roadmap-planning"],
  #notesByName["scenario.aggregation.policy"],
  #notesByName["scenario.constraints"],
  #notesByName["scenario.isolation"],
  #notesByName["strategy.growth-scenario-evaluation"],
  #notesByName["strategy.investment-decision"],
  #notesByName["system.extensibility.methods"],
  #notesByName["technology.architecture-choice"],
  #notesByName["technology.infrastructure-strategy"],
  #notesByName["technology.platform-selection"],
  #notesByName["ux.model.guidance"],
  #notesByName["vendor.service-provider-comparison"],
  #notesByName["vendor.supplier-selection"],
]

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
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.validation.input-config.parse-args"].name
    to: #notesByName["call.validation.input-config.load-cue-config"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.validation.input-config.load-cue-config"].name
    to: #notesByName["call.validation.input-config.validate-model"].name
    label: "delegate_to"
    labels: ["delegate_to"]
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

graphIntegrityPolicy: {
  missingNode:              "error"
  orphanNode:               "warning"
  duplicateNoteName:        "error"
  unknownRelationshipLabel: "ignore"
  crossReportReference:     "allow"
}
