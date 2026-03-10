# Glossary

Definitions of the main design terms, methods, and modeling concepts used by baldrick-seer.

## Decision methods

Core MCDA methods and analysis terms referenced by the design.

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

Evaluate alternatives against multiple criteria instead of reducing the decision to a single input dimension.

#### PROMETHEE

Compare alternatives pairwise with preference functions to produce a transparent ranking.

#### TOPSIS

Rank alternatives by their distance from an ideal best and an ideal worst solution.

#### VIKOR

Identify a compromise solution that balances group utility and individual regret.

## Modeling terms

Important concepts used to describe the input model and its validation rules.

### Model concepts

#### Clear Representation of Pairwise Judgments (v1)

Represent pairwise comparisons explicitly with named criteria instead of positional matrices so humans and AI can validate and generate them.

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

Define how multiple scenarios are combined through cross-scenario aggregation into a final decision, starting with practical v1 approaches such as equal averaging or weighted averaging with explicit scenario aggregation weights defined in the aggregation configuration as the single source of truth.

#### Constraint Enforcement (v1)

Allow scenarios to define hard requirements that can exclude alternatives before ranking.

#### Scenario Isolation (v1)

Evaluate each scenario independently with its own priorities and candidate evaluations.

