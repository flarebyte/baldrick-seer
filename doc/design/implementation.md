# Implementation Considerations

Implementation guidance and method references for the CLI and model.

## Modeling guidance

Recommendations about the decision-model shape and validation.

### Model structure

#### Clear Representation of Pairwise Judgments

Represent pairwise comparisons explicitly with named criteria instead of positional matrices so humans and AI can validate and generate them.

#### Consistent Criteria Interpretation

Keep each criterion semantically stable across scenarios even when its importance changes.

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

#### Extensible Decision Methods

Design the system so additional MCDA methods can be added later without redesigning the data model or CLI interface.

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

Render machine-readable JSON output for automation, downstream processing, and reproducibility.

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

Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data.

#### Structured Output for Automation

Provide machine-readable output such as JSON in addition to human-readable summaries.

#### Readable CLI Output

Present results in a clear terminal-friendly format with summaries, tables, and scenario breakdowns.

#### Analytic Hierarchy Process (AHP)

Derive criteria weights from pairwise comparisons and turn qualitative judgments into a consistent numerical weighting system.

#### TOPSIS

Rank alternatives by their distance from an ideal best and an ideal worst solution.

## User experience and output

Guidance for readable, reproducible, and automatable execution.

### CLI and explainability

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

#### Guidance for Model Creation

Provide prompts and guidance that help users define criteria, comparisons, and scenario descriptions with fewer modeling errors.

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

Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data.

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

#### Clear Representation of Pairwise Judgments

Represent pairwise comparisons explicitly with named criteria instead of positional matrices so humans and AI can validate and generate them.

#### Handling Incomplete Information

Detect missing comparisons or evaluation values and provide clear feedback, with any inferred values marked explicitly.

#### Model Validation

Validate referenced criteria, pairwise comparison completeness, and alternative evaluation coverage before computation.

#### Constraint Enforcement

Allow scenarios to define hard requirements that can exclude alternatives before ranking.

