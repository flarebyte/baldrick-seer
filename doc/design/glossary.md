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

Derive criteria weights from pairwise comparisons and turn qualitative judgments into a consistent numerical weighting system.

#### ELECTRE Outranking Method

Use concordance and discordance reasoning to determine whether one alternative sufficiently outranks another.

#### PROMETHEE

Compare alternatives pairwise with preference functions to produce a transparent ranking.

#### TOPSIS

Rank alternatives by their distance from an ideal best and an ideal worst solution.

#### VIKOR

Identify a compromise solution that balances group utility and individual regret.

## Execution concepts

Terms used in the CLI and report-generation design.

### CLI and output

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

#### Extensible Decision Methods (v2)

Generalize the pipeline so additional MCDA methods can be added later without redesigning the data model or CLI interface.

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

Detect missing comparisons or evaluation values early and return actionable diagnostics instead of attempting opaque implicit recovery.

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

