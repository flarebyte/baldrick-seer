# Implementation Considerations

Implementation guidance and method references for the CLI and model.

## Engineering conventions

Implementation rules intended to keep the codebase readable, testable, and deterministic.

### Code structure

#### Deterministic Output Ordering (v1)

Guarantee deterministic ordering in generated outputs so repeated runs produce stable markdown, JSON, and CSV artifacts.

#### Guard Clauses and Early Returns (v1)

Prefer early returns and guard clauses for error handling so failure paths stay short, obvious, and easy to test.

#### I/O and Core Logic Separation (v1)

Separate filesystem, terminal, and config-loading I/O from core decision logic so the computation pipeline stays testable and deterministic.

#### Named Predicates over Boolean Soup (v1)

Replace tangled boolean expressions with named predicates so validation and ranking rules read like domain logic instead of control noise.

#### Small Single-Purpose Functions (v1)

Keep functions small and single-purpose so validation, weighting, ranking, and rendering logic remain easy to understand and reuse.

#### Tiny Structs over Long Parameter Lists (v1)

Use small focused structs to carry grouped inputs instead of long parameter lists that are brittle and hard to read.

## Implementation stack

Primary languages, libraries, and tools chosen for the first release.

### Runtime and tooling

#### Cobra Command and Argument Parsing (v1)

Use Cobra for CLI command structure and argument parsing so command behavior and report argument handling follow one consistent model.

#### Go CLI Implementation (v1)

Implement the production CLI in Go so the tool remains fast, portable, and straightforward to distribute.

#### CUE as Configuration Source of Truth (v1)

Use CUE as the configuration source of truth so schema, defaults, validation, and concrete config evaluation live in one place.

#### Bun and TypeScript for E2E Tests (v1)

Implement end-to-end tests in TypeScript with Bun so CLI scenarios can be expressed tersely while staying fast to run in CI.

## Modeling guidance

Recommendations about the decision-model shape and validation.

### Model structure

#### Clear Representation of Pairwise Judgments (v1)

Represent pairwise comparisons explicitly with named criteria instead of positional matrices so humans and AI can validate and generate them.

#### Consistent Criteria Interpretation (v1)

Keep each criterion semantically stable across scenarios even when its importance changes.

#### Human and AI Friendly Input Format (v1)

Use a semantic format such as CUE that remains readable for humans and AI systems while supporting strong validation.

#### Model Documentation (v1)

Allow decision models to carry descriptions, notes, and justifications for comparisons and values.

#### Handling Incomplete Information (v1)

Detect missing comparisons or evaluation values early and return actionable diagnostics with both precise paths and readable named locations.

#### Decision Model Structure (v1)

Represent the decision problem with clear structures for criteria, alternatives, and scenarios that remain understandable to humans and AI.

#### Model Validation (v1)

Validate referenced criteria, pairwise comparison completeness, and alternative evaluation coverage before computation.

#### Scenario Aggregation Strategy (v1)

Define how multiple scenarios are combined into a final decision, starting with practical v1 aggregation approaches such as equal or weighted averaging.

#### Constraint Enforcement (v1)

Allow scenarios to define hard requirements that can exclude alternatives before ranking.

#### Scenario Isolation (v1)

Evaluate each scenario independently with its own priorities and candidate evaluations.

#### Extensible Decision Methods (v2)

Generalize the pipeline so additional MCDA methods can be added later without redesigning the data model or CLI interface.

## Referenced methods

Algorithms and analysis techniques explicitly named in the design.

### Algorithms

#### Robustness Analysis

Evaluate how stable a decision remains when assumptions, scenarios, or parameter ranges vary.

#### Sensitivity Analysis

Evaluate how changes in weights or inputs affect the ranking of alternatives.

#### Analytic Hierarchy Process (AHP)

Derive criteria weights from pairwise comparisons and turn qualitative judgments into a consistent numerical weighting system.

#### ELECTRE Outranking Method

Use concordance and discordance reasoning to determine whether one alternative sufficiently outranks another.

#### Multi-Criteria Decision Analysis (MCDA)

Evaluate alternatives against multiple criteria instead of reducing the decision to a single input dimension.

#### PROMETHEE

Compare alternatives pairwise with preference functions to produce a transparent ranking.

#### TOPSIS

Rank alternatives by their distance from an ideal best and an ideal worst solution.

#### VIKOR

Identify a compromise solution that balances group utility and individual regret.

## Report generation flow

CLI execution path for generating reports after the shared validation stage.

### Generate reports

#### Generate Reports Call

Top-level CLI call flow for generating reports from an input decision model.

#### Build AHP Inputs

Collect scenario pairwise comparisons into the normalized input structures needed for AHP weight computation.

#### Build TOPSIS Inputs

Combine validated evaluations, criterion polarity, and AHP-derived weights into the decision matrices required by TOPSIS.

#### Compute Criteria Weights with AHP

Transform pairwise scenario preferences into normalized criteria weights using Analytic Hierarchy Process.

#### Future Option: Rank with ELECTRE

Potential v2 branch where the validated model is ranked with ELECTRE instead of TOPSIS.

#### Future Option: TOPSIS with Sensitivity Analysis

Potential v2 branch where TOPSIS ranking is complemented by sensitivity analysis to assess robustness.

#### Parse Report Arguments

Parse CLI arguments for report generation, including the config path, requested report names, and output options.

#### Rank Alternatives with TOPSIS

Use the validated evaluations and AHP-derived weights to rank alternatives with TOPSIS.

#### Render Requested Reports

Render the requested markdown, JSON, or CSV reports from the computed ranking results.

#### Render CSV Report

Render flat tabular CSV output for spreadsheet analysis and data exchange.

#### Render JSON Report

Render machine-readable JSON output for automation, downstream processing, and reproducibility, including structured diagnostics when validation fails.

#### Render Markdown Report

Render narrative markdown output for human readers, including rankings, explanations, and scenario summaries.

#### Select Ranking Strategy

Select the ranking pipeline after AHP weighting. The current default path is TOPSIS, while v2 may add ELECTRE or TOPSIS followed by sensitivity analysis.

#### Select Requested Reports

Resolve which report definitions should run, applying any CLI filtering by report name or output target.

#### Reuse Shared Validation Flow

Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs.

#### Load CUE Config

Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.

#### Validate Config Model

Run structural and graph validation on the loaded config and emit diagnostics with both machine paths and human-readable locations.

## User experience and output

Guidance for readable, reproducible, and automatable execution.

### CLI and explainability

#### Sensitivity and Robustness Analysis (v2)

Add post-ranking analysis that tests how changes in criteria importance or scenario assumptions affect the final result so users can judge stability.

#### Structured Output for Automation (v1)

Provide machine-readable output such as JSON in addition to human-readable summaries.

#### Readable CLI Output (v1)

Present results in a clear terminal-friendly format with summaries, tables, and scenario breakdowns.

#### Explainable Results (v1)

Explain ranking outputs in terms of criteria influence and scenario differences.

#### Traceable Decision Process (v1)

Show the reasoning path from inputs to outputs, including scenario weights, criteria importance, and contribution of each factor.

#### Reproducible Decision Runs (v1)

Running the same model with the same inputs should always produce identical results for auditing and comparison.

#### Guidance for Model Creation (v2)

Provide richer prompts and guidance that help users define criteria, comparisons, and scenario descriptions with fewer modeling errors.

## Validation call flow

Early CLI execution path for reading and validating an input config file.

### Input config validation

#### Validate Input Config Call

Top-level CLI call flow for validating an input configuration file before any decision analysis runs.

#### Load CUE Config

Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.

#### Parse Validation Arguments

Parse CLI arguments for the validate command, including the config path and output flags.

#### Validate Config Model

Run structural and graph validation on the loaded config and emit diagnostics with both machine paths and human-readable locations.

#### Check Scenario Constraints

Check that scenario constraints target known criteria and use operators and values that are compatible with the referenced criterion types.

#### Check Evaluation Coverage

Check that evaluations reference known scenarios and alternatives and provide the values required by each scenario's active criteria.

#### Check Pairwise Comparisons

Check that pairwise comparisons are valid for each scenario, with known criteria, no self-comparisons, and sufficient coverage for AHP weighting.

#### Check Named References

Check that all named references resolve, including criteria names, scenario names, alternative names, and report focus selectors.

#### Check Report Definitions

Check that report definitions use supported formats, valid focus selectors, and well-formed argument lists for later Cobra-style parsing.

#### Check Config Structure

Check that the loaded config matches the expected top-level shape, required sections, and field types after CUE evaluation.

