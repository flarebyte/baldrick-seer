# Glossary

Definitions of the main design terms, methods, and modeling concepts used by baldrick-seer.

## Decision methods

Core MCDA methods and analysis terms referenced by the design.

### Methods

#### Robustness Analysis

Evaluate how stable a decision remains when assumptions, scenarios, or parameter ranges vary.

#### Sensitivity Analysis

Evaluate how changes in weights or inputs affect the ranking of alternatives.

#### Criterion Value Normalization (v1)

Define explicit v1 normalization rules for criterion values before ranking. Numbers are used directly in the decision matrix, ordinal values are validated as integers and then treated numerically, and boolean values are normalized to numeric form with `true = 1` and `false = 0`.

#### Analytic Hierarchy Process (AHP)

Derive criterion weights within a scenario from pairwise criterion comparisons and turn qualitative judgments into a consistent numerical weighting system.

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

#### Boolean Criterion Scoring (v1)

Normalize boolean criterion values before scoring by mapping `true` to `1` and `false` to `0`. Criterion polarity determines whether `true` or `false` is preferred in the ranking.

#### Numeric Criterion Scoring (v1)

Treat numeric criterion values as measurable quantities used directly in the decision matrix. Criterion polarity determines whether higher or lower values are preferred during normalization and ranking.

#### Ordinal Criterion Scoring (v1)

Treat ordinal criterion values as ordered integer levels used numerically in the decision matrix. Higher integers represent a higher level of the criterion, polarity determines desirability, and ordinal criteria should include `scaleGuidance`.

## Modeling terms

Important concepts used to describe the input model and its validation rules.

### Model concepts

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

#### Handling Incomplete Information (v1)

Detect missing pairwise comparisons required for full AHP coverage or missing evaluation values early and return actionable diagnostics with both precise paths and readable named locations.

#### Decision Model Structure (v1)

Represent the decision problem with clear structures for criteria, alternatives, and scenarios that remain understandable to humans and AI.

#### Model Validation (v1)

Validate referenced criteria, exact full pairwise comparison coverage for each AHP scenario, supported v1 value types, integer ordinal values, ordinal scale documentation, boolean true-or-false values, compatible constraint operator/value combinations, and alternative evaluation coverage before computation.

#### Scenario Aggregation Strategy (v1)

Define how multiple scenarios are combined through cross-scenario aggregation into a final decision, starting with practical v1 approaches such as equal averaging or weighted averaging with explicit scenario aggregation weights defined in the aggregation configuration as the single source of truth.

#### Constraint Semantics (v1)

Keep the `ScenarioConstraint` shape as `criterionName`, `operator`, and `value`, but interpret it by criterion type. Number criteria accept numeric values with `<=`, `>=`, `=`, or `!=` for threshold-style rules. Ordinal criteria accept integer values within the defined scale with `<=`, `>=`, `=`, or `!=`, following the criterion's ordering. Boolean criteria accept only `=` or `!=` with `true` or `false`; comparison operators such as `<=` and `>=` are invalid.

#### Constraint Enforcement (v1)

Allow scenarios to define hard requirements that can exclude alternatives before ranking, using constraint operators and values that remain compatible with each referenced criterion type.

#### Scenario Isolation (v1)

Evaluate each scenario independently with its own priorities and candidate evaluations.

