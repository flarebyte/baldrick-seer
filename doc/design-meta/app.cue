package flyb

source: "baldrick-seer.design"
name: "baldrick-seer"
modules: ["design"]

#notesByName: {
  "analysis.robustness": {
    name: "analysis.robustness"
    title: "Sensitivity and Robustness Analysis (v2)"
    labels: ["design", "implementation", "v2"]
    markdown: "Add post-ranking analysis that tests how changes in criteria importance or scenario assumptions affect the final result so users can judge stability."
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
    title: "Structured Output for Automation (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Provide machine-readable output such as JSON in addition to human-readable summaries."
  }
  "cli.output.readability": {
    name: "cli.output.readability"
    title: "Readable CLI Output (v1)"
    labels: ["design", "implementation", "v1"]
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
    markdown: "Run structural and graph validation on the loaded config and emit diagnostics with both machine paths and human-readable locations."
  }
  "call.validation.input-config.validate-model.check-structure": {
    name: "call.validation.input-config.validate-model.check-structure"
    title: "Check Config Structure"
    labels: ["call", "design", "flow", "implementation", "validation"]
    markdown: "Check that the loaded config matches the expected top-level shape, required sections, and field types after CUE evaluation."
  }
  "call.validation.input-config.validate-model.check-references": {
    name: "call.validation.input-config.validate-model.check-references"
    title: "Check Named References"
    labels: ["call", "design", "flow", "implementation", "validation"]
    markdown: "Check that all named references resolve, including criteria names, scenario names, alternative names, and report focus selectors."
  }
  "call.validation.input-config.validate-model.check-pairwise-comparisons": {
    name: "call.validation.input-config.validate-model.check-pairwise-comparisons"
    title: "Check Pairwise Comparisons"
    labels: ["call", "design", "flow", "implementation", "validation"]
    markdown: "Check that each scenario using AHP provides pairwise comparisons only between known active criteria, never compares a criterion with itself, and includes exactly one canonical comparison for every unordered pair of distinct active criteria. Reject duplicate comparisons, inverse duplicates, or any missing pair."
  }
  "call.validation.input-config.validate-model.check-evaluation-coverage": {
    name: "call.validation.input-config.validate-model.check-evaluation-coverage"
    title: "Check Evaluation Coverage"
    labels: ["call", "design", "flow", "implementation", "validation"]
    markdown: "Check that evaluations reference known scenarios and alternatives and provide supported v1 criterion values for each scenario's active criteria: numbers, integer ordinals, or booleans."
  }
  "call.validation.input-config.validate-model.check-constraints": {
    name: "call.validation.input-config.validate-model.check-constraints"
    title: "Check Scenario Constraints"
    labels: ["call", "design", "flow", "implementation", "validation"]
    markdown: "Check that scenario constraints target known criteria and use operators and values compatible with the referenced criterion types, including requiring equality-only constraints for boolean criteria."
  }
  "call.validation.input-config.validate-model.check-report-definitions": {
    name: "call.validation.input-config.validate-model.check-report-definitions"
    title: "Check Report Definitions"
    labels: ["call", "design", "flow", "implementation", "validation"]
    markdown: "Check that report definitions use supported formats, valid focus selectors, and well-formed argument lists for later Cobra-style parsing."
  }
  "call.reports.generate": {
    name: "call.reports.generate"
    title: "Generate Reports Call"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Top-level CLI call flow for generating reports from an input decision model."
  }
  "call.reports.generate.parse-args": {
    name: "call.reports.generate.parse-args"
    title: "Parse Report Arguments"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Parse CLI arguments for report generation, including the config path, requested report names, and output options."
  }
  "call.reports.generate.select-reports": {
    name: "call.reports.generate.select-reports"
    title: "Select Requested Reports"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Resolve which report definitions should run, applying any CLI filtering by report name or output target."
  }
  "call.reports.generate.shared-validation": {
    name: "call.reports.generate.shared-validation"
    title: "Reuse Shared Validation Flow"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs."
  }
  "call.reports.generate.build-ahp-inputs": {
    name: "call.reports.generate.build-ahp-inputs"
    title: "Build AHP Inputs"
    labels: ["call", "design", "flow", "implementation", "method"]
    markdown: "Collect the validated full pairwise comparison set for each scenario into the normalized input structures needed for AHP computation of scenario-local criterion weights."
  }
  "call.reports.generate.compute-ahp-weights": {
    name: "call.reports.generate.compute-ahp-weights"
    title: "Compute Criteria Weights with AHP"
    labels: ["call", "design", "flow", "implementation", "method"]
    markdown: "Transform pairwise criterion comparisons within each scenario into normalized scenario-local criterion weights using Analytic Hierarchy Process."
  }
  "call.reports.generate.select-ranking-strategy": {
    name: "call.reports.generate.select-ranking-strategy"
    title: "Select Ranking Strategy"
    labels: ["call", "design", "flow", "implementation", "method"]
    markdown: "Select the ranking pipeline after computing scenario-local criterion weights with AHP. The current default path is TOPSIS, while v2 may add ELECTRE or TOPSIS followed by sensitivity analysis."
  }
  "call.reports.generate.build-topsis-inputs": {
    name: "call.reports.generate.build-topsis-inputs"
    title: "Build TOPSIS Inputs"
    labels: ["call", "design", "flow", "implementation", "method"]
    markdown: "Combine validated evaluations, criterion polarity, and AHP-derived scenario-local criterion weights into the decision matrices required by TOPSIS."
  }
  "call.reports.generate.rank-alternatives-topsis": {
    name: "call.reports.generate.rank-alternatives-topsis"
    title: "Rank Alternatives with TOPSIS"
    labels: ["call", "design", "flow", "implementation", "method"]
    markdown: "Use the validated evaluations and scenario-local criterion weights derived with AHP to rank alternatives with TOPSIS."
  }
  "call.reports.generate.future-rank-electre": {
    name: "call.reports.generate.future-rank-electre"
    title: "Future Option: Rank with ELECTRE"
    labels: ["call", "design", "flow", "future", "method"]
    markdown: "Potential v2 branch where the validated model is ranked with ELECTRE instead of TOPSIS."
  }
  "call.reports.generate.future-rank-topsis-sensitivity": {
    name: "call.reports.generate.future-rank-topsis-sensitivity"
    title: "Future Option: TOPSIS with Sensitivity Analysis"
    labels: ["call", "design", "flow", "future", "method"]
    markdown: "Potential v2 branch where TOPSIS ranking is complemented by sensitivity analysis to assess robustness."
  }
  "call.reports.generate.render-output": {
    name: "call.reports.generate.render-output"
    title: "Render Requested Reports"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Render the requested markdown, JSON, or CSV reports from the computed ranking results."
  }
  "call.reports.generate.render-output.render-markdown": {
    name: "call.reports.generate.render-output.render-markdown"
    title: "Render Markdown Report"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Render narrative markdown output for human readers, including rankings, explanations, and scenario summaries."
  }
  "call.reports.generate.render-output.render-json": {
    name: "call.reports.generate.render-output.render-json"
    title: "Render JSON Report"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Render machine-readable JSON output for automation, downstream processing, and reproducibility, including structured diagnostics when validation fails."
  }
  "call.reports.generate.render-output.render-csv": {
    name: "call.reports.generate.render-output.render-csv"
    title: "Render CSV Report"
    labels: ["call", "design", "flow", "implementation"]
    markdown: "Render flat tabular CSV output for spreadsheet analysis and data exchange."
  }
  "criteria.pairwise.clarity": {
    name: "criteria.pairwise.clarity"
    title: "Clear Representation of Pairwise Judgments (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Represent pairwise comparisons explicitly with named criteria and a single canonical direction, using one field for the more important criterion and one field for the less important criterion, so humans and AI can validate and generate exactly one comparison for each unordered criterion pair."
  }
  "criteria.semantic.consistency": {
    name: "criteria.semantic.consistency"
    title: "Consistent Criteria Interpretation (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Keep each criterion semantically stable across scenarios even when its importance changes."
  }
  "criteria.value-types.v1": {
    name: "criteria.value-types.v1"
    title: "Supported Criterion Value Types (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Support only three criterion value types in v1: number, ordinal, and boolean. Text criterion values are not part of the v1 model."
  }
  "criteria.scale-guidance.ordinal": {
    name: "criteria.scale-guidance.ordinal"
    title: "Document Ordinal Scales (v1)"
    labels: ["design", "implementation", "validation", "v1"]
    markdown: "Require ordinal criteria to document their scale with `scaleGuidance`, so each integer level has a clear ordered meaning before scoring."
  }
  "decision.explainability": {
    name: "decision.explainability"
    title: "Explainable Results (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Explain ranking outputs in terms of criteria influence and scenario differences."
  }
  "decision.multi-criteria-ranking": {
    name: "decision.multi-criteria-ranking"
    title: "General Multi-Criteria Ranking (v1)"
    labels: ["design", "usecase", "v1"]
    markdown: "Run a scenario-based CLI evaluation that ranks named alternatives from a structured config and emits decision reports for humans and tools."
  }
  "decision.robust-choice-identification": {
    name: "decision.robust-choice-identification"
    title: "Robust Choice Identification (v2)"
    labels: ["design", "usecase", "v2"]
    markdown: "Identify alternatives that remain strong across scenarios and under changing assumptions, likely using robustness or sensitivity-oriented post-analysis."
  }
  "decision.traceability": {
    name: "decision.traceability"
    title: "Traceable Decision Process (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Show the reasoning path from inputs to outputs, including scenario-local criterion weights, scenario aggregation weights, and contribution of each factor."
  }
  "engineering.deterministic-ordering": {
    name: "engineering.deterministic-ordering"
    title: "Deterministic Output Ordering (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Guarantee deterministic ordering in generated outputs so repeated runs produce stable markdown, JSON, and CSV artifacts."
  }
  "engineering.guard-clauses": {
    name: "engineering.guard-clauses"
    title: "Guard Clauses and Early Returns (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Prefer early returns and guard clauses for error handling so failure paths stay short, obvious, and easy to test."
  }
  "engineering.small-functions": {
    name: "engineering.small-functions"
    title: "Small Single-Purpose Functions (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Keep functions small and single-purpose so validation, weighting, ranking, and rendering logic remain easy to understand and reuse."
  }
  "engineering.io-core-separation": {
    name: "engineering.io-core-separation"
    title: "I/O and Core Logic Separation (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Separate filesystem, terminal, and config-loading I/O from core decision logic so the computation pipeline stays testable and deterministic."
  }
  "engineering.tiny-structs": {
    name: "engineering.tiny-structs"
    title: "Tiny Structs over Long Parameter Lists (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Use small focused structs to carry grouped inputs instead of long parameter lists that are brittle and hard to read."
  }
  "engineering.named-predicates": {
    name: "engineering.named-predicates"
    title: "Named Predicates over Boolean Soup (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Replace tangled boolean expressions with named predicates so validation and ranking rules read like domain logic instead of control noise."
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
  "example.input-schema.ts": {
    name: "example.input-schema.ts"
    title: "TypeScript Input Model"
    labels: ["design", "example", "source", "typescript"]
    filepath: "examples/model.ts"
  }
  "example.input-platform-selection.ts": {
    name: "example.input-platform-selection.ts"
    title: "Platform Selection TypeScript Example"
    labels: ["design", "example", "source", "typescript"]
    filepath: "examples/mcda.ts"
  }
  "example.input-hosting-choice.ts": {
    name: "example.input-hosting-choice.ts"
    title: "Hosting Choice TypeScript Example"
    labels: ["design", "example", "source", "typescript"]
    filepath: "examples/minimum-mcda.ts"
  }
  "example.output-hosting-choice.markdown": {
    name: "example.output-hosting-choice.markdown"
    title: "Hosting Choice Markdown Output"
    labels: ["design", "example", "output", "markdown"]
    filepath: "examples/hosting-choice-summary.md"
  }
  "example.output-hosting-choice.json": {
    name: "example.output-hosting-choice.json"
    title: "Hosting Choice JSON Output"
    labels: ["design", "example", "output", "json"]
    filepath: "examples/hosting-choice-results.json"
  }
  "example.output-hosting-choice.csv": {
    name: "example.output-hosting-choice.csv"
    title: "Hosting Choice CSV Output"
    labels: ["design", "example", "output", "csv"]
    filepath: "examples/hosting-choice-scores.csv"
  }
  "example.output-validation-failure.markdown": {
    name: "example.output-validation-failure.markdown"
    title: "Validation Failure Markdown Output"
    labels: ["design", "example", "output", "markdown", "validation"]
    filepath: "examples/validation-failure.md"
  }
  "example.output-validation-failure.json": {
    name: "example.output-validation-failure.json"
    title: "Validation Failure JSON Output"
    labels: ["design", "example", "output", "json", "validation"]
    filepath: "examples/validation-failure.json"
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
    title: "Reproducible Decision Runs (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Running the same model with the same inputs should always produce identical results for auditing and comparison."
  }
  "infrastructure.system-design-selection": {
    name: "infrastructure.system-design-selection"
    title: "System Design Selection (v1)"
    labels: ["design", "usecase", "v1"]
    markdown: "Compare system designs where trade-offs exist between cost, scalability, and reliability."
  }
  "input.format": {
    name: "input.format"
    title: "Human and AI Friendly Input Format (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Use a semantic format such as CUE that remains readable for humans and AI systems while supporting strong validation."
  }
  "mcda.ahp": {
    name: "mcda.ahp"
    title: "Analytic Hierarchy Process (AHP)"
    labels: ["design", "method"]
    markdown: "Derive criterion weights within a scenario from pairwise criterion comparisons and turn qualitative judgments into a consistent numerical weighting system."
  }
  "mcda.general": {
    name: "mcda.general"
    title: "Multi-Criteria Decision Analysis (MCDA)"
    labels: ["design", "method"]
    markdown: "Evaluate alternatives against multiple criteria instead of reducing the decision to a single input dimension."
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
  "scoring.number-normalization.v1": {
    name: "scoring.number-normalization.v1"
    title: "Numeric Criterion Scoring (v1)"
    labels: ["design", "implementation", "method", "v1"]
    markdown: "Treat numeric criterion values as measurable quantities used directly in the decision matrix. Criterion polarity determines whether higher or lower values are preferred during normalization and ranking."
  }
  "scoring.ordinal-normalization.v1": {
    name: "scoring.ordinal-normalization.v1"
    title: "Ordinal Criterion Scoring (v1)"
    labels: ["design", "implementation", "method", "validation", "v1"]
    markdown: "Treat ordinal criterion values as ordered integer levels used numerically in the decision matrix. Higher integers represent a higher level of the criterion, polarity determines desirability, and ordinal criteria should include `scaleGuidance`."
  }
  "scoring.boolean-normalization.v1": {
    name: "scoring.boolean-normalization.v1"
    title: "Boolean Criterion Scoring (v1)"
    labels: ["design", "implementation", "method", "validation", "v1"]
    markdown: "Normalize boolean criterion values before scoring by mapping `true` to `1` and `false` to `0`. Criterion polarity determines whether `true` or `false` is preferred in the ranking."
  }
  "mcda.vikor": {
    name: "mcda.vikor"
    title: "VIKOR"
    labels: ["design", "method"]
    markdown: "Identify a compromise solution that balances group utility and individual regret."
  }
  "model.documentation": {
    name: "model.documentation"
    title: "Model Documentation (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Allow decision models to carry descriptions, notes, and justifications for comparisons and values."
  }
  "model.incomplete.data": {
    name: "model.incomplete.data"
    title: "Handling Incomplete Information (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Detect missing pairwise comparisons required for full AHP coverage or missing evaluation values early and return actionable diagnostics with both precise paths and readable named locations."
  }
  "model.structure": {
    name: "model.structure"
    title: "Decision Model Structure (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Represent the decision problem with clear structures for criteria, alternatives, and scenarios that remain understandable to humans and AI."
  }
  "model.validation": {
    name: "model.validation"
    title: "Model Validation (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Validate referenced criteria, exact full pairwise comparison coverage for each AHP scenario, supported v1 value types, ordinal scale documentation, boolean constraint operators, and alternative evaluation coverage before computation."
  }
  "planning.lifecycle-decision": {
    name: "planning.lifecycle-decision"
    title: "Lifecycle Decision Support (v2)"
    labels: ["design", "usecase", "v2"]
    markdown: "Compare options that must remain effective throughout different stages of organizational or system development."
  }
  "planning.long-term-option-evaluation": {
    name: "planning.long-term-option-evaluation"
    title: "Long-Term Option Evaluation (v2)"
    labels: ["design", "usecase", "v2"]
    markdown: "Evaluate alternatives that must perform well across multiple possible future environments."
  }
  "policy.policy-option-analysis": {
    name: "policy.policy-option-analysis"
    title: "Policy Option Analysis (v2)"
    labels: ["design", "usecase", "v2"]
    markdown: "Support evaluation of policy alternatives where multiple criteria such as impact, feasibility, and cost must be considered."
  }
  "product.feature-prioritization": {
    name: "product.feature-prioritization"
    title: "Product Feature Prioritization (v1)"
    labels: ["design", "usecase", "v1"]
    markdown: "Rank product features using multiple criteria such as user value, development effort, and strategic importance."
  }
  "product.roadmap-planning": {
    name: "product.roadmap-planning"
    title: "Product Roadmap Planning (v2)"
    labels: ["design", "usecase", "v2"]
    markdown: "Evaluate product initiatives across different market or growth scenarios to support long-term planning."
  }
  "scenario.aggregation.policy": {
    name: "scenario.aggregation.policy"
    title: "Scenario Aggregation Strategy (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Define how multiple scenarios are combined through cross-scenario aggregation into a final decision, starting with practical v1 approaches such as equal averaging or weighted averaging with explicit scenario aggregation weights defined in the aggregation configuration as the single source of truth."
  }
  "scenario.constraints": {
    name: "scenario.constraints"
    title: "Constraint Enforcement (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Allow scenarios to define hard requirements that can exclude alternatives before ranking."
  }
  "scenario.isolation": {
    name: "scenario.isolation"
    title: "Scenario Isolation (v1)"
    labels: ["design", "implementation", "v1"]
    markdown: "Evaluate each scenario independently with its own priorities and candidate evaluations."
  }
  "stack.cli.go": {
    name: "stack.cli.go"
    title: "Go CLI Implementation (v1)"
    labels: ["design", "implementation", "stack", "v1"]
    markdown: "Implement the production CLI in Go so the tool remains fast, portable, and straightforward to distribute."
  }
  "stack.cli.cobra": {
    name: "stack.cli.cobra"
    title: "Cobra Command and Argument Parsing (v1)"
    labels: ["design", "implementation", "stack", "v1"]
    markdown: "Use Cobra for CLI command structure and argument parsing so command behavior and report argument handling follow one consistent model."
  }
  "stack.config.cue": {
    name: "stack.config.cue"
    title: "CUE as Configuration Source of Truth (v1)"
    labels: ["design", "implementation", "stack", "v1"]
    markdown: "Use CUE as the configuration source of truth so schema, defaults, validation, and concrete config evaluation live in one place."
  }
  "strategy.growth-scenario-evaluation": {
    name: "strategy.growth-scenario-evaluation"
    title: "Growth Scenario Evaluation (v1)"
    labels: ["design", "usecase", "v1"]
    markdown: "Assess the same strategic or technical options under startup, scale-up, and mature-operation scenarios using one shared decision model."
  }
  "strategy.investment-decision": {
    name: "strategy.investment-decision"
    title: "Strategic Investment Decision (v2)"
    labels: ["design", "usecase", "v2"]
    markdown: "Compare investment alternatives considering financial return, risk exposure, and long-term strategic impact."
  }
  "system.extensibility.methods": {
    name: "system.extensibility.methods"
    title: "Extensible Decision Methods (v2)"
    labels: ["design", "implementation", "v2"]
    markdown: "Generalize the pipeline so additional MCDA methods can be added later without redesigning the data model or CLI interface."
  }
  "testing.e2e.bun-typescript": {
    name: "testing.e2e.bun-typescript"
    title: "Bun and TypeScript for E2E Tests (v1)"
    labels: ["design", "implementation", "testing", "v1"]
    markdown: "Implement end-to-end tests in TypeScript with Bun so CLI scenarios can be expressed tersely while staying fast to run in CI."
  }
  "technology.architecture-choice": {
    name: "technology.architecture-choice"
    title: "Software Architecture Decision (v1)"
    labels: ["design", "usecase", "v1"]
    markdown: "Evaluate architectural approaches under different system growth conditions, performance requirements, and reliability expectations."
  }
  "technology.infrastructure-strategy": {
    name: "technology.infrastructure-strategy"
    title: "Infrastructure Strategy Planning (v1)"
    labels: ["design", "usecase", "v1"]
    markdown: "Assess infrastructure alternatives such as cloud providers or deployment models under varying operational scenarios."
  }
  "technology.platform-selection": {
    name: "technology.platform-selection"
    title: "Technology Platform Selection (v1)"
    labels: ["design", "usecase", "v1"]
    markdown: "Compare multiple technology platforms across operational scenarios such as startup, scale-up, and enterprise maturity."
  }
  "ux.model.guidance": {
    name: "ux.model.guidance"
    title: "Guidance for Model Creation (v2)"
    labels: ["design", "implementation", "v2"]
    markdown: "Provide richer prompts and guidance that help users define criteria, comparisons, and scenario descriptions with fewer modeling errors."
  }
  "vendor.service-provider-comparison": {
    name: "vendor.service-provider-comparison"
    title: "Service Provider Comparison (v1)"
    labels: ["design", "usecase", "v1"]
    markdown: "Compare service providers where priorities may change depending on scale, regulatory environment, or organizational maturity."
  }
  "vendor.supplier-selection": {
    name: "vendor.supplier-selection"
    title: "Supplier or Vendor Selection (v1)"
    labels: ["design", "usecase", "v1"]
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
      {
        title: "TypeScript source examples"
        description: "Readable draft examples for the future CUE input model."
        sections: [
          {
            title: "Model types"
            notes: [
              #notesByName["example.input-schema.ts"].name,
            ]
          },
          {
            title: "Scenario-based examples"
            notes: [
              #notesByName["example.input-platform-selection.ts"].name,
              #notesByName["example.input-hosting-choice.ts"].name,
            ]
          },
        ]
      },
      {
        title: "Example outputs"
        description: "Illustrative outputs for the hosting-choice example across the main v1 report formats."
        sections: [
          {
            title: "Successful reports"
            notes: [
              #notesByName["example.output-hosting-choice.markdown"].name,
              #notesByName["example.output-hosting-choice.json"].name,
              #notesByName["example.output-hosting-choice.csv"].name,
            ]
          },
          {
            title: "Validation failures"
            notes: [
              #notesByName["example.output-validation-failure.markdown"].name,
              #notesByName["example.output-validation-failure.json"].name,
            ]
          },
        ]
      },
    ]
  },
  {
    title: "Glossary"
    filepath: "../design/glossary.md"
    description: "Definitions of the main design terms, methods, and modeling concepts used by baldrick-seer."
    sections: [
      {
        title: "Decision methods"
        description: "Core MCDA methods and analysis terms referenced by the design."
        sections: [
          {
            title: "Methods"
            notes: [
              #notesByName["mcda.general"].name,
              #notesByName["analysis.robustness.method"].name,
              #notesByName["analysis.sensitivity"].name,
              #notesByName["mcda.ahp"].name,
              #notesByName["mcda.electre"].name,
              #notesByName["mcda.promethee"].name,
              #notesByName["mcda.topsis"].name,
              #notesByName["scoring.number-normalization.v1"].name,
              #notesByName["scoring.ordinal-normalization.v1"].name,
              #notesByName["scoring.boolean-normalization.v1"].name,
              #notesByName["mcda.vikor"].name,
            ]
          },
        ]
      },
      {
        title: "Modeling terms"
        description: "Important concepts used to describe the input model and its validation rules."
        sections: [
          {
            title: "Model concepts"
            notes: [
              #notesByName["criteria.pairwise.clarity"].name,
              #notesByName["criteria.value-types.v1"].name,
              #notesByName["criteria.scale-guidance.ordinal"].name,
              #notesByName["input.format"].name,
              #notesByName["model.documentation"].name,
              #notesByName["model.incomplete.data"].name,
              #notesByName["model.structure"].name,
              #notesByName["model.validation"].name,
              #notesByName["scenario.aggregation.policy"].name,
              #notesByName["scenario.constraints"].name,
              #notesByName["scenario.isolation"].name,
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
              #notesByName["criteria.value-types.v1"].name,
              #notesByName["criteria.scale-guidance.ordinal"].name,
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
        title: "Engineering conventions"
        description: "Implementation rules intended to keep the codebase readable, testable, and deterministic."
        sections: [
          {
            title: "Code structure"
            notes: [
              #notesByName["engineering.deterministic-ordering"].name,
              #notesByName["engineering.guard-clauses"].name,
              #notesByName["engineering.small-functions"].name,
              #notesByName["engineering.io-core-separation"].name,
              #notesByName["engineering.tiny-structs"].name,
              #notesByName["engineering.named-predicates"].name,
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
        title: "Implementation stack"
        description: "Primary languages, libraries, and tools chosen for the first release."
        sections: [
          {
            title: "Runtime and tooling"
            notes: [
              #notesByName["stack.cli.go"].name,
              #notesByName["stack.cli.cobra"].name,
              #notesByName["stack.config.cue"].name,
              #notesByName["testing.e2e.bun-typescript"].name,
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
              #notesByName["call.validation.input-config.validate-model.check-structure"].name,
              #notesByName["call.validation.input-config.validate-model.check-references"].name,
              #notesByName["call.validation.input-config.validate-model.check-pairwise-comparisons"].name,
              #notesByName["call.validation.input-config.validate-model.check-evaluation-coverage"].name,
              #notesByName["call.validation.input-config.validate-model.check-constraints"].name,
              #notesByName["call.validation.input-config.validate-model.check-report-definitions"].name,
            ]
          },
        ]
      },
      {
        title: "Report generation flow"
        description: "CLI execution path for generating reports after the shared validation stage."
        sections: [
          {
            title: "Generate reports"
            notes: [
              #notesByName["call.reports.generate"].name,
              #notesByName["call.reports.generate.parse-args"].name,
              #notesByName["call.reports.generate.select-reports"].name,
              #notesByName["call.reports.generate.shared-validation"].name,
              #notesByName["call.validation.input-config.load-cue-config"].name,
              #notesByName["call.validation.input-config.validate-model"].name,
              #notesByName["call.reports.generate.build-ahp-inputs"].name,
              #notesByName["call.reports.generate.compute-ahp-weights"].name,
              #notesByName["call.reports.generate.select-ranking-strategy"].name,
              #notesByName["call.reports.generate.build-topsis-inputs"].name,
              #notesByName["call.reports.generate.rank-alternatives-topsis"].name,
              #notesByName["call.reports.generate.future-rank-electre"].name,
              #notesByName["call.reports.generate.future-rank-topsis-sensitivity"].name,
              #notesByName["call.reports.generate.render-output"].name,
              #notesByName["call.reports.generate.render-output.render-markdown"].name,
              #notesByName["call.reports.generate.render-output.render-json"].name,
              #notesByName["call.reports.generate.render-output.render-csv"].name,
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
              #notesByName["mcda.general"].name,
              #notesByName["mcda.ahp"].name,
              #notesByName["mcda.electre"].name,
              #notesByName["mcda.promethee"].name,
              #notesByName["mcda.topsis"].name,
              #notesByName["scoring.number-normalization.v1"].name,
              #notesByName["scoring.ordinal-normalization.v1"].name,
              #notesByName["scoring.boolean-normalization.v1"].name,
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
              #notesByName["call.validation.input-config"].name,
              #notesByName["call.validation.input-config.parse-args"].name,
              #notesByName["call.validation.input-config.load-cue-config"].name,
              #notesByName["call.validation.input-config.validate-model"].name,
              #notesByName["call.validation.input-config.validate-model.check-structure"].name,
              #notesByName["call.validation.input-config.validate-model.check-references"].name,
              #notesByName["call.validation.input-config.validate-model.check-pairwise-comparisons"].name,
              #notesByName["call.validation.input-config.validate-model.check-evaluation-coverage"].name,
              #notesByName["call.validation.input-config.validate-model.check-constraints"].name,
              #notesByName["call.validation.input-config.validate-model.check-report-definitions"].name,
              #notesByName["input.format"].name,
              #notesByName["model.validation"].name,
              #notesByName["model.incomplete.data"].name,
              #notesByName["criteria.pairwise.clarity"].name,
              #notesByName["scenario.constraints"].name,
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
              #notesByName["call.reports.generate"].name,
              #notesByName["call.reports.generate.parse-args"].name,
              #notesByName["call.reports.generate.select-reports"].name,
              #notesByName["call.reports.generate.shared-validation"].name,
              #notesByName["call.validation.input-config.load-cue-config"].name,
              #notesByName["call.validation.input-config.validate-model"].name,
              #notesByName["call.reports.generate.build-ahp-inputs"].name,
              #notesByName["call.reports.generate.compute-ahp-weights"].name,
              #notesByName["call.reports.generate.select-ranking-strategy"].name,
              #notesByName["call.reports.generate.build-topsis-inputs"].name,
              #notesByName["call.reports.generate.rank-alternatives-topsis"].name,
              #notesByName["call.reports.generate.future-rank-electre"].name,
              #notesByName["call.reports.generate.future-rank-topsis-sensitivity"].name,
              #notesByName["call.reports.generate.render-output"].name,
              #notesByName["call.reports.generate.render-output.render-markdown"].name,
              #notesByName["call.reports.generate.render-output.render-json"].name,
              #notesByName["call.reports.generate.render-output.render-csv"].name,
              #notesByName["mcda.ahp"].name,
              #notesByName["mcda.topsis"].name,
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
  #notesByName["call.validation.input-config.validate-model.check-structure"],
  #notesByName["call.validation.input-config.validate-model.check-references"],
  #notesByName["call.validation.input-config.validate-model.check-pairwise-comparisons"],
  #notesByName["call.validation.input-config.validate-model.check-evaluation-coverage"],
  #notesByName["call.validation.input-config.validate-model.check-constraints"],
  #notesByName["call.validation.input-config.validate-model.check-report-definitions"],
  #notesByName["call.reports.generate"],
  #notesByName["call.reports.generate.parse-args"],
  #notesByName["call.reports.generate.select-reports"],
  #notesByName["call.reports.generate.shared-validation"],
  #notesByName["call.reports.generate.build-ahp-inputs"],
  #notesByName["call.reports.generate.compute-ahp-weights"],
  #notesByName["call.reports.generate.select-ranking-strategy"],
  #notesByName["call.reports.generate.build-topsis-inputs"],
  #notesByName["call.reports.generate.rank-alternatives-topsis"],
  #notesByName["call.reports.generate.future-rank-electre"],
  #notesByName["call.reports.generate.future-rank-topsis-sensitivity"],
  #notesByName["call.reports.generate.render-output"],
  #notesByName["call.reports.generate.render-output.render-markdown"],
  #notesByName["call.reports.generate.render-output.render-json"],
  #notesByName["call.reports.generate.render-output.render-csv"],
  #notesByName["criteria.pairwise.clarity"],
  #notesByName["criteria.value-types.v1"],
  #notesByName["criteria.scale-guidance.ordinal"],
  #notesByName["criteria.semantic.consistency"],
  #notesByName["decision.explainability"],
  #notesByName["decision.multi-criteria-ranking"],
  #notesByName["decision.robust-choice-identification"],
  #notesByName["decision.traceability"],
  #notesByName["engineering.deterministic-ordering"],
  #notesByName["engineering.guard-clauses"],
  #notesByName["engineering.io-core-separation"],
  #notesByName["engineering.named-predicates"],
  #notesByName["engineering.small-functions"],
  #notesByName["engineering.tiny-structs"],
  #notesByName["example.hosting-choice"],
  #notesByName["example.hosting-choice.lean-startup"],
  #notesByName["example.hosting-choice.regulated-growth"],
  #notesByName["example.input-hosting-choice.ts"],
  #notesByName["example.output-hosting-choice.csv"],
  #notesByName["example.output-hosting-choice.json"],
  #notesByName["example.output-hosting-choice.markdown"],
  #notesByName["example.output-validation-failure.json"],
  #notesByName["example.output-validation-failure.markdown"],
  #notesByName["example.input-platform-selection.ts"],
  #notesByName["example.input-schema.ts"],
  #notesByName["example.platform-selection"],
  #notesByName["example.platform-selection.established-enterprise"],
  #notesByName["example.platform-selection.startup"],
  #notesByName["example.platform-selection.unicorn"],
  #notesByName["execution.reproducibility"],
  #notesByName["infrastructure.system-design-selection"],
  #notesByName["input.format"],
  #notesByName["mcda.ahp"],
  #notesByName["mcda.general"],
  #notesByName["mcda.electre"],
  #notesByName["mcda.promethee"],
  #notesByName["mcda.topsis"],
  #notesByName["scoring.number-normalization.v1"],
  #notesByName["scoring.ordinal-normalization.v1"],
  #notesByName["scoring.boolean-normalization.v1"],
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
  #notesByName["stack.cli.cobra"],
  #notesByName["stack.cli.go"],
  #notesByName["stack.config.cue"],
  #notesByName["strategy.growth-scenario-evaluation"],
  #notesByName["strategy.investment-decision"],
  #notesByName["system.extensibility.methods"],
  #notesByName["testing.e2e.bun-typescript"],
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
    from: #notesByName["call.validation.input-config.validate-model"].name
    to: #notesByName["call.validation.input-config.validate-model.check-structure"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-structure"].name
    to: #notesByName["call.validation.input-config.validate-model.check-references"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-references"].name
    to: #notesByName["call.validation.input-config.validate-model.check-pairwise-comparisons"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-pairwise-comparisons"].name
    to: #notesByName["call.validation.input-config.validate-model.check-evaluation-coverage"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-evaluation-coverage"].name
    to: #notesByName["call.validation.input-config.validate-model.check-constraints"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-constraints"].name
    to: #notesByName["call.validation.input-config.validate-model.check-report-definitions"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate"].name
    to: #notesByName["call.reports.generate.parse-args"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.parse-args"].name
    to: #notesByName["call.reports.generate.select-reports"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.select-reports"].name
    to: #notesByName["call.reports.generate.shared-validation"].name
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
    from: #notesByName["call.validation.input-config.validate-model.check-structure"].name
    to: #notesByName["model.validation"].name
    label: "implements"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-references"].name
    to: #notesByName["model.validation"].name
    label: "implements"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model"].name
    to: #notesByName["model.incomplete.data"].name
    label: "checks_for"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-pairwise-comparisons"].name
    to: #notesByName["criteria.pairwise.clarity"].name
    label: "implements"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-pairwise-comparisons"].name
    to: #notesByName["model.incomplete.data"].name
    label: "checks_for"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-evaluation-coverage"].name
    to: #notesByName["model.incomplete.data"].name
    label: "checks_for"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-evaluation-coverage"].name
    to: #notesByName["model.validation"].name
    label: "implements"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-constraints"].name
    to: #notesByName["scenario.constraints"].name
    label: "implements"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-constraints"].name
    to: #notesByName["model.validation"].name
    label: "implements"
  },
  {
    from: #notesByName["call.validation.input-config.validate-model.check-report-definitions"].name
    to: #notesByName["model.validation"].name
    label: "implements"
  },
  {
    from: #notesByName["call.reports.generate.shared-validation"].name
    to: #notesByName["call.validation.input-config.load-cue-config"].name
    label: "reuses"
  },
  {
    from: #notesByName["call.reports.generate.shared-validation"].name
    to: #notesByName["call.validation.input-config.validate-model"].name
    label: "reuses"
  },
  {
    from: #notesByName["call.reports.generate.compute-ahp-weights"].name
    to: #notesByName["mcda.ahp"].name
    label: "implements"
  },
  {
    from: #notesByName["call.reports.generate.shared-validation"].name
    to: #notesByName["call.reports.generate.build-ahp-inputs"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.build-ahp-inputs"].name
    to: #notesByName["call.reports.generate.compute-ahp-weights"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.build-ahp-inputs"].name
    to: #notesByName["mcda.ahp"].name
    label: "implements"
  },
  {
    from: #notesByName["call.reports.generate.compute-ahp-weights"].name
    to: #notesByName["call.reports.generate.select-ranking-strategy"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.select-ranking-strategy"].name
    to: #notesByName["call.reports.generate.build-topsis-inputs"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.select-ranking-strategy"].name
    to: #notesByName["call.reports.generate.future-rank-electre"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.select-ranking-strategy"].name
    to: #notesByName["call.reports.generate.future-rank-topsis-sensitivity"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.build-topsis-inputs"].name
    to: #notesByName["call.reports.generate.rank-alternatives-topsis"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.build-topsis-inputs"].name
    to: #notesByName["mcda.topsis"].name
    label: "implements"
  },
  {
    from: #notesByName["call.reports.generate.rank-alternatives-topsis"].name
    to: #notesByName["mcda.topsis"].name
    label: "implements"
  },
  {
    from: #notesByName["call.reports.generate.rank-alternatives-topsis"].name
    to: #notesByName["call.reports.generate.render-output"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.render-output"].name
    to: #notesByName["cli.output.machine"].name
    label: "supports"
  },
  {
    from: #notesByName["call.reports.generate.render-output"].name
    to: #notesByName["cli.output.readability"].name
    label: "supports"
  },
  {
    from: #notesByName["call.reports.generate.render-output"].name
    to: #notesByName["call.reports.generate.render-output.render-markdown"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.render-output"].name
    to: #notesByName["call.reports.generate.render-output.render-json"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.render-output"].name
    to: #notesByName["call.reports.generate.render-output.render-csv"].name
    label: "delegate_to"
    labels: ["delegate_to"]
  },
  {
    from: #notesByName["call.reports.generate.render-output.render-markdown"].name
    to: #notesByName["cli.output.readability"].name
    label: "supports"
  },
  {
    from: #notesByName["call.reports.generate.render-output.render-json"].name
    to: #notesByName["cli.output.machine"].name
    label: "supports"
  },
  {
    from: #notesByName["call.reports.generate.render-output.render-csv"].name
    to: #notesByName["cli.output.machine"].name
    label: "supports"
  },
  {
    from: #notesByName["mcda.general"].name
    to: #notesByName["decision.multi-criteria-ranking"].name
    label: "supports"
  },
  {
    from: #notesByName["call.reports.generate.future-rank-electre"].name
    to: #notesByName["mcda.electre"].name
    label: "implements"
  },
  {
    from: #notesByName["call.reports.generate.future-rank-topsis-sensitivity"].name
    to: #notesByName["analysis.sensitivity"].name
    label: "implements"
  },
  {
    from: #notesByName["criteria.pairwise.clarity"].name
    to: #notesByName["mcda.ahp"].name
    label: "documents_method"
  },
  {
    from: #notesByName["mcda.general"].name
    to: #notesByName["mcda.ahp"].name
    label: "includes_method"
  },
  {
    from: #notesByName["mcda.general"].name
    to: #notesByName["mcda.electre"].name
    label: "includes_method"
  },
  {
    from: #notesByName["mcda.general"].name
    to: #notesByName["mcda.promethee"].name
    label: "includes_method"
  },
  {
    from: #notesByName["mcda.general"].name
    to: #notesByName["mcda.topsis"].name
    label: "includes_method"
  },
  {
    from: #notesByName["mcda.general"].name
    to: #notesByName["mcda.vikor"].name
    label: "includes_method"
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
    from: #notesByName["engineering.deterministic-ordering"].name
    to: #notesByName["execution.reproducibility"].name
    label: "supports"
  },
  {
    from: #notesByName["engineering.guard-clauses"].name
    to: #notesByName["model.validation"].name
    label: "supports"
  },
  {
    from: #notesByName["engineering.small-functions"].name
    to: #notesByName["engineering.io-core-separation"].name
    label: "supports"
  },
  {
    from: #notesByName["engineering.io-core-separation"].name
    to: #notesByName["execution.reproducibility"].name
    label: "supports"
  },
  {
    from: #notesByName["engineering.tiny-structs"].name
    to: #notesByName["model.structure"].name
    label: "supports"
  },
  {
    from: #notesByName["engineering.named-predicates"].name
    to: #notesByName["model.validation"].name
    label: "supports"
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
    from: #notesByName["stack.cli.go"].name
    to: #notesByName["call.reports.generate"].name
    label: "implements"
  },
  {
    from: #notesByName["stack.cli.go"].name
    to: #notesByName["call.validation.input-config"].name
    label: "implements"
  },
  {
    from: #notesByName["stack.cli.cobra"].name
    to: #notesByName["call.validation.input-config.parse-args"].name
    label: "implements"
  },
  {
    from: #notesByName["stack.cli.cobra"].name
    to: #notesByName["call.reports.generate.parse-args"].name
    label: "implements"
  },
  {
    from: #notesByName["stack.config.cue"].name
    to: #notesByName["input.format"].name
    label: "implements"
  },
  {
    from: #notesByName["stack.config.cue"].name
    to: #notesByName["call.validation.input-config.load-cue-config"].name
    label: "implements"
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
    from: #notesByName["example.input-schema.ts"].name
    to:   #notesByName["model.structure"].name
    label: "documents"
  },
  {
    from: #notesByName["example.input-platform-selection.ts"].name
    to:   #notesByName["example.platform-selection"].name
    label: "documents"
  },
  {
    from: #notesByName["example.input-hosting-choice.ts"].name
    to:   #notesByName["example.hosting-choice"].name
    label: "documents"
  },
  {
    from: #notesByName["example.output-hosting-choice.markdown"].name
    to:   #notesByName["example.hosting-choice"].name
    label: "documents"
  },
  {
    from: #notesByName["example.output-hosting-choice.json"].name
    to:   #notesByName["example.hosting-choice"].name
    label: "documents"
  },
  {
    from: #notesByName["example.output-hosting-choice.csv"].name
    to:   #notesByName["example.hosting-choice"].name
    label: "documents"
  },
  {
    from: #notesByName["example.output-validation-failure.markdown"].name
    to:   #notesByName["call.validation.input-config.validate-model"].name
    label: "documents"
  },
  {
    from: #notesByName["example.output-validation-failure.json"].name
    to:   #notesByName["call.validation.input-config.validate-model"].name
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
    from: #notesByName["testing.e2e.bun-typescript"].name
    to: #notesByName["execution.reproducibility"].name
    label: "supports"
  },
  {
    from: #notesByName["testing.e2e.bun-typescript"].name
    to: #notesByName["call.reports.generate"].name
    label: "tests"
  },
  {
    from: #notesByName["testing.e2e.bun-typescript"].name
    to: #notesByName["call.validation.input-config"].name
    label: "tests"
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
