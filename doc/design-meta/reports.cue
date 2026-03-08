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
