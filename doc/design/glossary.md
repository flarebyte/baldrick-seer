# Glossary

Concise term definitions used by the normative specification and examples.

## Decision methods

Short definitions of the named decision methods and analysis terms.

### Methods

#### Robustness Analysis

Evaluate how stable a decision remains when assumptions, scenarios, or parameter ranges vary.

#### Sensitivity Analysis

Evaluate how changes in weights or inputs affect the ranking of alternatives.

#### Analytic Hierarchy Process (AHP)

Derive criterion weights within a scenario from pairwise criterion comparisons and turn qualitative judgments into a consistent numerical weighting system.

#### ELECTRE Outranking Method

Use concordance and discordance reasoning to determine whether one alternative sufficiently outranks another.

#### Multi-Criteria Decision Analysis (MCDA)

Evaluate alternatives against multiple criteria instead of reducing the decision to a single input dimension. In this design, v1 applies MCDA through an AHP-derived weighting stage followed by TOPSIS ranking.

#### PROMETHEE

Compare alternatives pairwise with preference functions to produce a transparent ranking.

#### TOPSIS

Rank alternatives by their distance from an ideal best and an ideal worst solution.

#### VIKOR

Identify a compromise solution that balances group utility and individual regret.

## Modeling terms

Short definitions of recurring model terms. Normative rules are defined in the specification.

### Model concepts

#### Supported Criterion Value Types (v1)

Support only three criterion value types in v1: number, ordinal, and boolean. Text criterion values are not part of the v1 model.

#### Human and AI Friendly Input Format (v1)

Use a semantic format such as CUE that remains readable for humans and AI systems while supporting strong validation.

#### Report Argument Validation (v1)

Keep `ReportDefinition.arguments` as `string[]` in `key=value` form so the model stays extensible, but validate it strictly in v1. Only documented arguments are accepted, unknown keys are errors, some keys may be shared across formats while others are format-specific, incompatible format-specific keys must be rejected, values must match the argument definition, and duplicate keys are invalid unless the spec explicitly allows them.

#### Scenario Aggregation Strategy (v1)

Define how multiple scenarios are combined through cross-scenario aggregation into a final decision, starting with practical v1 approaches such as equal averaging or weighted averaging with explicit scenario aggregation weights defined in the aggregation configuration as the single source of truth. In v1, if an alternative is excluded by constraints in any scenario that participates in the aggregation, that alternative is ineligible for the final ranking and must not appear in the final aggregated list.

#### Scenario Isolation (v1)

Evaluate each scenario independently with its own priorities and candidate evaluations.

