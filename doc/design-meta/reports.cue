package flyb

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
              "decision.multi-criteria-ranking",
              "decision.robust-choice-identification",
              "infrastructure.system-design-selection",
              "technology.infrastructure-strategy",
              "technology.platform-selection",
            ]
          },
          {
            title: "Reference examples"
            notes: [
              "example.hosting-choice",
              "example.platform-selection",
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
              "analysis.robustness.method",
              "analysis.sensitivity",
              "mcda.ahp",
              "mcda.electre",
              "mcda.promethee",
              "mcda.topsis",
              "mcda.vikor",
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
              "example.platform-selection",
            ]
          },
          {
            title: "Scenarios"
            notes: [
              "example.platform-selection.established-enterprise",
              "example.platform-selection.startup",
              "example.platform-selection.unicorn",
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
              "example.hosting-choice",
            ]
          },
          {
            title: "Scenarios"
            notes: [
              "example.hosting-choice.lean-startup",
              "example.hosting-choice.regulated-growth",
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
              "criteria.pairwise.clarity",
              "criteria.semantic.consistency",
              "input.format",
              "model.documentation",
              "model.incomplete.data",
              "model.structure",
              "model.validation",
              "scenario.aggregation.policy",
              "scenario.constraints",
              "scenario.isolation",
              "system.extensibility.methods",
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
              "analysis.robustness",
              "cli.output.machine",
              "cli.output.readability",
              "decision.explainability",
              "decision.traceability",
              "execution.reproducibility",
              "ux.model.guidance",
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
              "analysis.robustness.method",
              "analysis.sensitivity",
              "mcda.ahp",
              "mcda.electre",
              "mcda.promethee",
              "mcda.topsis",
              "mcda.vikor",
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
              "decision.multi-criteria-ranking",
              "decision.robust-choice-identification",
              "planning.lifecycle-decision",
              "planning.long-term-option-evaluation",
              "policy.policy-option-analysis",
              "strategy.growth-scenario-evaluation",
              "strategy.investment-decision",
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
              "infrastructure.system-design-selection",
              "product.feature-prioritization",
              "product.roadmap-planning",
              "technology.architecture-choice",
              "technology.infrastructure-strategy",
              "technology.platform-selection",
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
              "vendor.service-provider-comparison",
              "vendor.supplier-selection",
            ]
          },
        ]
      },
    ]
  },
]
