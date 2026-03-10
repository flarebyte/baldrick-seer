# Normative Specification

Authoritative v1 specification for the decision model, validation rules, and CLI execution behavior.

## Execution Behavior

Authoritative CLI behavior for validation-only runs and report generation.

### Method scope

#### Future Option: Rank with ELECTRE

Potential v2 branch where the validated model is ranked with ELECTRE instead of TOPSIS.

#### Future Option: TOPSIS with Sensitivity Analysis

Potential v2 branch where TOPSIS ranking is complemented by sensitivity analysis to assess robustness.

#### Analytic Hierarchy Process (AHP)

Derive criterion weights within a scenario from pairwise criterion comparisons and turn qualitative judgments into a consistent numerical weighting system.

#### Multi-Criteria Decision Analysis (MCDA)

Evaluate alternatives against multiple criteria instead of reducing the decision to a single input dimension. In this design, v1 applies MCDA through an AHP-derived weighting stage followed by TOPSIS ranking.

#### TOPSIS

Rank alternatives by their distance from an ideal best and an ideal worst solution.

#### Extensible Decision Methods (v2)

Extend the v1 AHP + TOPSIS pipeline so additional MCDA methods can be added later without replacing the overall CLI shape. Future methods such as ELECTRE, PROMETHEE, or VIKOR may require additional optional metadata or method-specific configuration beyond the v1 model.

### Report generation command

#### Generate Reports Call

Top-level CLI call flow for generating ranking reports from an input decision model. The command reuses the shared validation path and fails fast if the model is invalid.

#### Build AHP Inputs

Collect the validated full pairwise comparison set for each scenario into the normalized input structures needed for AHP computation of scenario-local criterion weights.

#### Build TOPSIS Inputs

Combine validated evaluations, criterion polarity, and AHP-derived scenario-local criterion weights into the decision matrices required by TOPSIS, after removing any alternatives excluded by scenario constraints.

#### Compute Criteria Weights with AHP

Transform pairwise criterion comparisons within each scenario into normalized scenario-local criterion weights using Analytic Hierarchy Process.

#### Parse Report Arguments

Parse CLI arguments for report generation, including the config path, requested report names, and output options.

#### Rank Alternatives with TOPSIS

Use the validated evaluations and scenario-local criterion weights derived with AHP to rank alternatives with TOPSIS among the alternatives that remain eligible after scenario-local constraint enforcement.

#### Render Requested Reports

Render the requested markdown, JSON, or CSV outputs only after validation succeeds and ranking results are computed. Invalid models do not reach report rendering.

#### Render CSV Report

Render flat tabular CSV output for spreadsheet analysis and data exchange.

#### Render JSON Report

Render machine-readable JSON ranking output for automation, downstream processing, and reproducibility only when validation succeeds. Scenario-level output should indicate alternatives excluded by constraints and omit scenario scores or ranks for them, while final aggregated rankings omit alternatives made ineligible by participating scenario constraints. If JSON output is requested and validation fails, the command may emit structured diagnostics as an error payload or via stderr, but that output is not a successful ranking report.

#### Render Markdown Report

Render narrative markdown output for human readers, including rankings, explanations, scenario summaries, and clear indication when an alternative was excluded by a scenario constraint. Final aggregated rankings include only alternatives that remain eligible across all participating scenarios.

#### Select Ranking Strategy

Select the ranking pipeline after computing scenario-local criterion weights with AHP. In v1, the design is built around an AHP + TOPSIS pipeline; v2 may add alternatives such as ELECTRE or TOPSIS followed by sensitivity analysis.

#### Select Requested Reports

Resolve which report definitions should run, applying any CLI filtering by report name or output target.

#### Reuse Shared Validation Flow

Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs. If validation fails, report generation stops immediately and no ranking report is produced.

### Validate command

#### Validate Input Config Call

Top-level CLI call flow for validating an input configuration file and returning validation results only, without scoring or report generation.

#### Load CUE Config

Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.

#### Parse Validation Arguments

Parse CLI arguments for the validate command, including the config path and output flags.

#### Validate Config Model

Run structural and graph validation on the loaded config and emit diagnostics with both machine paths and human-readable locations. For the `validate` command, this is the terminal result of the command.

## Model and Validation

Authoritative definitions for the input model, value semantics, aggregation, constraints, and validation rules.

### Model semantics

#### Clear Representation of Pairwise Judgments (v1)

Represent pairwise comparisons explicitly with named criteria and a single canonical direction, using one field for the more important criterion and one field for the less important criterion, so humans and AI can validate and generate exactly one comparison for each unordered criterion pair.

#### Document Ordinal Scales (v1)

Require ordinal criteria to document their scale with `scaleGuidance`, so each integer level has a clear ordered meaning before scoring.

#### Criterion Value Normalization (v1)

Define explicit v1 normalization rules for criterion values before ranking. Numbers are used directly in the decision matrix, ordinal values are validated as integers and then treated numerically, and boolean values are normalized to numeric form with `true = 1` and `false = 0`.

#### Supported Criterion Value Types (v1)

Support only three criterion value types in v1: number, ordinal, and boolean. Text criterion values are not part of the v1 model.

#### Human and AI Friendly Input Format (v1)

Use a semantic format such as CUE that remains readable for humans and AI systems while supporting strong validation.

#### Model Documentation (v1)

Allow decision models to carry descriptions, notes, and justifications for comparisons and values.

#### Decision Model Structure (v1)

Represent the decision problem with clear structures for criteria, alternatives, and scenarios that remain understandable to humans and AI.

#### Report Argument Validation (v1)

Keep `ReportDefinition.arguments` as `string[]` in `key=value` form so the model stays extensible, but validate it strictly in v1. Only documented arguments are accepted, unknown keys are errors, some keys may be shared across formats while others are format-specific, incompatible format-specific keys must be rejected, values must match the argument definition, and duplicate keys are invalid unless the spec explicitly allows them.

#### Scenario Aggregation Strategy (v1)

Define how multiple scenarios are combined through cross-scenario aggregation into a final decision, starting with practical v1 approaches such as equal averaging or weighted averaging with explicit scenario aggregation weights defined in the aggregation configuration as the single source of truth. In v1, if an alternative is excluded by constraints in any scenario that participates in the aggregation, that alternative is ineligible for the final ranking and must not appear in the final aggregated list.

#### Constraint Semantics (v1)

Keep the `ScenarioConstraint` shape as `criterionName`, `operator`, and `value`, but interpret it by criterion type. Number criteria accept numeric values with `<=`, `>=`, `=`, or `!=` for threshold-style rules. Ordinal criteria accept integer values within the defined scale with `<=`, `>=`, `=`, or `!=`, following the criterion's ordering. Boolean criteria accept only `=` or `!=` with `true` or `false`; comparison operators such as `<=` and `>=` are invalid.

#### Constraint Enforcement (v1)

Allow scenarios to define hard requirements that are enforced within each scenario before ranking, using constraint operators and values that remain compatible with each referenced criterion type. Alternatives that violate a scenario constraint are excluded from that scenario's scoring and ranking, receive no scenario score or rank there, and should be reported as excluded due to constraints.

#### Scenario Isolation (v1)

Evaluate each scenario independently with its own priorities and candidate evaluations.

#### Boolean Criterion Scoring (v1)

Normalize boolean criterion values before scoring by mapping `true` to `1` and `false` to `0`. Criterion polarity determines whether `true` or `false` is preferred in the ranking.

#### Numeric Criterion Scoring (v1)

Treat numeric criterion values as measurable quantities used directly in the decision matrix. Criterion polarity determines whether higher or lower values are preferred during normalization and ranking.

#### Ordinal Criterion Scoring (v1)

Treat ordinal criterion values as ordered integer levels used numerically in the decision matrix. Higher integers represent a higher level of the criterion, polarity determines desirability, and ordinal criteria should include `scaleGuidance`.

### Validation rules

#### Check Scenario Constraints

Check that each scenario constraint uses an operator and value compatible with the referenced criterion type: number criteria allow numeric values with `<=`, `>=`, `=`, or `!=`; ordinal criteria allow integer values with `<=`, `>=`, `=`, or `!=`; boolean criteria allow only `=` or `!=` with `true` or `false`. Invalid operator/type combinations must raise a validation error. Constraint enforcement itself happens during scenario-local scoring, where violating alternatives are excluded before ranking.

#### Check Evaluation Coverage

Check that evaluations reference known scenarios and alternatives and provide supported v1 criterion values for each scenario's active criteria: measurable numbers, integer ordinals, or booleans with only `true` and `false` values.

#### Check Pairwise Comparisons

Check that each scenario using AHP provides pairwise comparisons only between known active criteria, never compares a criterion with itself, and includes exactly one canonical comparison for every unordered pair of distinct active criteria. Reject duplicate comparisons, inverse duplicates, or any missing pair.

#### Check Report Definitions

Check that report definitions use supported formats, valid focus selectors, and strictly validated report arguments. In v1 every report argument must use `key=value`, unknown arguments are validation errors, argument names must be allowed globally or for the selected format, format-specific arguments must match the report format, invalid values must be rejected, and duplicate keys are invalid unless explicitly defined otherwise.

#### Handling Incomplete Information (v1)

Detect missing pairwise comparisons required for full AHP coverage or missing evaluation values early and return actionable diagnostics with both precise paths and readable named locations.

#### Model Validation (v1)

Validate referenced criteria, exact full pairwise comparison coverage for each AHP scenario, supported v1 value types, integer ordinal values, ordinal scale documentation, boolean true-or-false values, compatible constraint operator/value combinations, and alternative evaluation coverage before computation.

## Scope

Primary use cases and example design anchors.

### Primary use cases

#### General Multi-Criteria Ranking (v1)

Run a scenario-based CLI evaluation that ranks named alternatives from a structured config and emits decision reports for humans and tools.

#### Robust Choice Identification (v2)

Identify alternatives that remain strong across scenarios and under changing assumptions, likely using robustness or sensitivity-oriented post-analysis.

#### System Design Selection (v1)

Compare system designs where trade-offs exist between cost, scalability, and reliability.

#### Infrastructure Strategy Planning (v1)

Assess infrastructure alternatives such as cloud providers or deployment models under varying operational scenarios.

#### Technology Platform Selection (v1)

Compare multiple technology platforms across operational scenarios such as startup, scale-up, and enterprise maturity.

### Reference examples

#### Hosting Choice Example

Minimal scenario-based MCDA example that compares two hosting providers across lean-startup and regulated-growth scenarios using equal-average aggregation.

#### Platform Selection Example

Scenario-based MCDA example that compares candidate platforms across startup, unicorn, and established-enterprise contexts using weighted-average aggregation.

