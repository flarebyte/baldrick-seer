package flyb

notes: [
  {
    name: "analysis.robustness"
    title: "Sensitivity and Robustness Analysis"
    labels: ["design", "implementation"]
    markdown: "Support testing how changes in criteria importance or scenario assumptions affect the final ranking so users can see whether a result is stable or fragile."
  },
  {
    name: "analysis.robustness.method"
    title: "Robustness Analysis"
    labels: ["analysis", "design", "method"]
    markdown: "Evaluate how stable a decision remains when assumptions, scenarios, or parameter ranges vary."
  },
  {
    name: "analysis.sensitivity"
    title: "Sensitivity Analysis"
    labels: ["analysis", "design", "method"]
    markdown: "Evaluate how changes in weights or inputs affect the ranking of alternatives."
  },
  {
    name: "cli.output.machine"
    title: "Structured Output for Automation"
    labels: ["design", "implementation"]
    markdown: "Provide machine-readable output such as JSON in addition to human-readable summaries."
  },
  {
    name: "cli.output.readability"
    title: "Readable CLI Output"
    labels: ["design", "implementation"]
    markdown: "Present results in a clear terminal-friendly format with summaries, tables, and scenario breakdowns."
  },
  {
    name: "criteria.pairwise.clarity"
    title: "Clear Representation of Pairwise Judgments"
    labels: ["design", "implementation"]
    markdown: "Represent pairwise comparisons explicitly with named criteria instead of positional matrices so humans and AI can validate and generate them."
  },
  {
    name: "criteria.semantic.consistency"
    title: "Consistent Criteria Interpretation"
    labels: ["design", "implementation"]
    markdown: "Keep each criterion semantically stable across scenarios even when its importance changes."
  },
  {
    name: "decision.explainability"
    title: "Explainable Results"
    labels: ["design", "implementation"]
    markdown: "Explain ranking outputs in terms of criteria influence and scenario differences."
  },
  {
    name: "decision.multi-criteria-ranking"
    title: "General Multi-Criteria Ranking"
    labels: ["design", "usecase"]
    markdown: "Rank competing alternatives using multiple evaluation criteria and scenario-based priorities."
  },
  {
    name: "decision.robust-choice-identification"
    title: "Robust Choice Identification"
    labels: ["design", "usecase"]
    markdown: "Identify alternatives that perform consistently well across different scenarios or changing assumptions."
  },
  {
    name: "decision.traceability"
    title: "Traceable Decision Process"
    labels: ["design", "implementation"]
    markdown: "Show the reasoning path from inputs to outputs, including scenario weights, criteria importance, and contribution of each factor."
  },
  {
    name: "example.hosting-choice"
    title: "Hosting Choice Example"
    labels: ["design", "example"]
    markdown: "Minimal scenario-based MCDA example that compares two hosting providers across lean-startup and regulated-growth scenarios using equal-average aggregation."
  },
  {
    name: "example.hosting-choice.lean-startup"
    title: "Lean Startup Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Scenario emphasizing cost and speed when budget pressure matters more than peak performance."
  },
  {
    name: "example.hosting-choice.regulated-growth"
    title: "Regulated Growth Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Scenario adding compliance as a dominant concern alongside cost and speed."
  },
  {
    name: "example.platform-selection"
    title: "Platform Selection Example"
    labels: ["design", "example"]
    markdown: "Scenario-based MCDA example that compares candidate platforms across startup, unicorn, and established-enterprise contexts using weighted-average aggregation."
  },
  {
    name: "example.platform-selection.established-enterprise"
    title: "Established Enterprise Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Mature-organization scenario where reliability and compliance dominate cost and future flexibility."
  },
  {
    name: "example.platform-selection.startup"
    title: "Startup Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Early-stage scenario where cost and time to market matter more than enterprise controls."
  },
  {
    name: "example.platform-selection.unicorn"
    title: "Unicorn Scenario"
    labels: ["design", "example", "flow"]
    markdown: "Hyper-growth scenario where scalability becomes dominant while cost still matters."
  },
  {
    name: "execution.reproducibility"
    title: "Reproducible Decision Runs"
    labels: ["design", "implementation"]
    markdown: "Running the same model with the same inputs should always produce identical results for auditing and comparison."
  },
  {
    name: "infrastructure.system-design-selection"
    title: "System Design Selection"
    labels: ["design", "usecase"]
    markdown: "Compare system designs where trade-offs exist between cost, scalability, and reliability."
  },
  {
    name: "input.format"
    title: "Human and AI Friendly Input Format"
    labels: ["design", "implementation"]
    markdown: "Use a semantic format such as JSON or YAML that is easy for humans and AI systems to read and generate."
  },
  {
    name: "mcda.ahp"
    title: "Analytic Hierarchy Process (AHP)"
    labels: ["design", "method"]
    markdown: "Derive criteria weights from pairwise comparisons and turn qualitative judgments into a consistent numerical weighting system."
  },
  {
    name: "mcda.electre"
    title: "ELECTRE Outranking Method"
    labels: ["design", "method"]
    markdown: "Use concordance and discordance reasoning to determine whether one alternative sufficiently outranks another."
  },
  {
    name: "mcda.promethee"
    title: "PROMETHEE"
    labels: ["design", "method"]
    markdown: "Compare alternatives pairwise with preference functions to produce a transparent ranking."
  },
  {
    name: "mcda.topsis"
    title: "TOPSIS"
    labels: ["design", "method"]
    markdown: "Rank alternatives by their distance from an ideal best and an ideal worst solution."
  },
  {
    name: "mcda.vikor"
    title: "VIKOR"
    labels: ["design", "method"]
    markdown: "Identify a compromise solution that balances group utility and individual regret."
  },
  {
    name: "model.documentation"
    title: "Model Documentation"
    labels: ["design", "implementation"]
    markdown: "Allow decision models to carry descriptions, notes, and justifications for comparisons and values."
  },
  {
    name: "model.incomplete.data"
    title: "Handling Incomplete Information"
    labels: ["design", "implementation"]
    markdown: "Detect missing comparisons or evaluation values and provide clear feedback, with any inferred values marked explicitly."
  },
  {
    name: "model.structure"
    title: "Decision Model Structure"
    labels: ["design", "implementation"]
    markdown: "Represent the decision problem with clear structures for criteria, alternatives, and scenarios that remain understandable to humans and AI."
  },
  {
    name: "model.validation"
    title: "Model Validation"
    labels: ["design", "implementation"]
    markdown: "Validate referenced criteria, pairwise comparison completeness, and alternative evaluation coverage before computation."
  },
  {
    name: "planning.lifecycle-decision"
    title: "Lifecycle Decision Support"
    labels: ["design", "usecase"]
    markdown: "Compare options that must remain effective throughout different stages of organizational or system development."
  },
  {
    name: "planning.long-term-option-evaluation"
    title: "Long-Term Option Evaluation"
    labels: ["design", "usecase"]
    markdown: "Evaluate alternatives that must perform well across multiple possible future environments."
  },
  {
    name: "policy.policy-option-analysis"
    title: "Policy Option Analysis"
    labels: ["design", "usecase"]
    markdown: "Support evaluation of policy alternatives where multiple criteria such as impact, feasibility, and cost must be considered."
  },
  {
    name: "product.feature-prioritization"
    title: "Product Feature Prioritization"
    labels: ["design", "usecase"]
    markdown: "Rank product features using multiple criteria such as user value, development effort, and strategic importance."
  },
  {
    name: "product.roadmap-planning"
    title: "Product Roadmap Planning"
    labels: ["design", "usecase"]
    markdown: "Evaluate product initiatives across different market or growth scenarios to support long-term planning."
  },
  {
    name: "scenario.aggregation.policy"
    title: "Scenario Aggregation Strategy"
    labels: ["design", "implementation"]
    markdown: "Define how multiple scenarios are combined into a final decision, such as equal averaging, weighted scenarios, or robustness-focused approaches."
  },
  {
    name: "scenario.constraints"
    title: "Constraint Enforcement"
    labels: ["design", "implementation"]
    markdown: "Allow scenarios to define hard requirements that can exclude alternatives before ranking."
  },
  {
    name: "scenario.isolation"
    title: "Scenario Isolation"
    labels: ["design", "implementation"]
    markdown: "Evaluate each scenario independently with its own priorities and candidate evaluations."
  },
  {
    name: "strategy.growth-scenario-evaluation"
    title: "Growth Scenario Evaluation"
    labels: ["design", "usecase"]
    markdown: "Assess strategic options under different business growth trajectories such as startup, rapid expansion, or mature operations."
  },
  {
    name: "strategy.investment-decision"
    title: "Strategic Investment Decision"
    labels: ["design", "usecase"]
    markdown: "Compare investment alternatives considering financial return, risk exposure, and long-term strategic impact."
  },
  {
    name: "system.extensibility.methods"
    title: "Extensible Decision Methods"
    labels: ["design", "implementation"]
    markdown: "Design the system so additional MCDA methods can be added later without redesigning the data model or CLI interface."
  },
  {
    name: "technology.architecture-choice"
    title: "Software Architecture Decision"
    labels: ["design", "usecase"]
    markdown: "Evaluate architectural approaches under different system growth conditions, performance requirements, and reliability expectations."
  },
  {
    name: "technology.infrastructure-strategy"
    title: "Infrastructure Strategy Planning"
    labels: ["design", "usecase"]
    markdown: "Assess infrastructure alternatives such as cloud providers or deployment models under varying operational scenarios."
  },
  {
    name: "technology.platform-selection"
    title: "Technology Platform Selection"
    labels: ["design", "usecase"]
    markdown: "Compare multiple technology platforms across operational scenarios such as startup, scale-up, and enterprise maturity."
  },
  {
    name: "ux.model.guidance"
    title: "Guidance for Model Creation"
    labels: ["design", "implementation"]
    markdown: "Provide prompts and guidance that help users define criteria, comparisons, and scenario descriptions with fewer modeling errors."
  },
  {
    name: "vendor.service-provider-comparison"
    title: "Service Provider Comparison"
    labels: ["design", "usecase"]
    markdown: "Compare service providers where priorities may change depending on scale, regulatory environment, or organizational maturity."
  },
  {
    name: "vendor.supplier-selection"
    title: "Supplier or Vendor Selection"
    labels: ["design", "usecase"]
    markdown: "Evaluate competing suppliers using multiple criteria such as cost, reliability, and service quality under different operating conditions."
  },
]
