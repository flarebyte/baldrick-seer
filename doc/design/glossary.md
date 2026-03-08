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

#### Sensitivity and Robustness Analysis

Support testing how changes in criteria importance or scenario assumptions affect the final ranking so users can see whether a result is stable or fragile.

#### Structured Output for Automation

Provide machine-readable output such as JSON in addition to human-readable summaries.

#### Readable CLI Output

Present results in a clear terminal-friendly format with summaries, tables, and scenario breakdowns.

#### Explainable Results

Explain ranking outputs in terms of criteria influence and scenario differences.

#### Traceable Decision Process

Show the reasoning path from inputs to outputs, including scenario weights, criteria importance, and contribution of each factor.

#### Reproducible Decision Runs

Running the same model with the same inputs should always produce identical results for auditing and comparison.

#### Extensible Decision Methods

Design the system so additional MCDA methods can be added later without redesigning the data model or CLI interface.

## Modeling terms

Important concepts used to describe the input model and its validation rules.

### Model concepts

#### Clear Representation of Pairwise Judgments

Represent pairwise comparisons explicitly with named criteria instead of positional matrices so humans and AI can validate and generate them.

#### Human and AI Friendly Input Format

Use a semantic format such as JSON or YAML that is easy for humans and AI systems to read and generate.

#### Model Documentation

Allow decision models to carry descriptions, notes, and justifications for comparisons and values.

#### Handling Incomplete Information

Detect missing comparisons or evaluation values and provide clear feedback, with any inferred values marked explicitly.

#### Decision Model Structure

Represent the decision problem with clear structures for criteria, alternatives, and scenarios that remain understandable to humans and AI.

#### Model Validation

Validate referenced criteria, pairwise comparison completeness, and alternative evaluation coverage before computation.

#### Scenario Aggregation Strategy

Define how multiple scenarios are combined into a final decision, such as equal averaging, weighted scenarios, or robustness-focused approaches.

#### Constraint Enforcement

Allow scenarios to define hard requirements that can exclude alternatives before ranking.

#### Scenario Isolation

Evaluate each scenario independently with its own priorities and candidate evaluations.

