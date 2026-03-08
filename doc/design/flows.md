# Execution Flows

Call-oriented flows for CLI validation and report generation.

## Report generation flows

Graph view for generating reports from a validated input config.

### Report generation graph

Text graph for the report-generation call chain, reusing the shared validation path.

- <a id="graph-node-call-reports-generate"></a> Generate Reports Call: Top-level CLI call flow for generating reports from an input decision model.
  - <a id="graph-node-call-reports-generate-parse-args"></a> Parse Report Arguments: Parse CLI arguments for report generation, including the config path, requested report names, and output options.
    - <a id="graph-node-call-reports-generate-select-reports"></a> Select Requested Reports: Resolve which report definitions should run, applying any CLI filtering by report name or output target.
      - <a id="graph-node-call-reports-generate-shared-validation"></a> Reuse Shared Validation Flow: Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs.
        - <a id="graph-node-call-reports-generate-build-ahp-inputs"></a> Build AHP Inputs: Collect scenario pairwise comparisons into the normalized input structures needed for AHP weight computation.
          - <a id="graph-node-call-reports-generate-compute-ahp-weights"></a> Compute Criteria Weights with AHP: Transform pairwise scenario preferences into normalized criteria weights using Analytic Hierarchy Process.
            - <a id="graph-node-call-reports-generate-select-ranking-strategy"></a> Select Ranking Strategy: Select the ranking pipeline after AHP weighting. The current default path is TOPSIS, while v2 may add ELECTRE or TOPSIS followed by sensitivity analysis.
              - <a id="graph-node-call-reports-generate-build-topsis-inputs"></a> Build TOPSIS Inputs: Combine validated evaluations, criterion polarity, and AHP-derived weights into the decision matrices required by TOPSIS.
                - <a id="graph-node-call-reports-generate-rank-alternatives-topsis"></a> Rank Alternatives with TOPSIS: Use the validated evaluations and AHP-derived weights to rank alternatives with TOPSIS.
                  - <a id="graph-node-call-reports-generate-render-output"></a> Render Requested Reports: Render the requested markdown, JSON, or CSV reports from the computed ranking results.
                    - <a id="graph-node-call-reports-generate-render-output-render-csv"></a> Render CSV Report: Render flat tabular CSV output for spreadsheet analysis and data exchange.
                    - <a id="graph-node-call-reports-generate-render-output-render-json"></a> Render JSON Report: Render machine-readable JSON output for automation, downstream processing, and reproducibility.
                    - <a id="graph-node-call-reports-generate-render-output-render-markdown"></a> Render Markdown Report: Render narrative markdown output for human readers, including rankings, explanations, and scenario summaries.
              - <a id="graph-node-call-reports-generate-future-rank-electre"></a> Future Option: Rank with ELECTRE: Potential v2 branch where the validated model is ranked with ELECTRE instead of TOPSIS.
              - <a id="graph-node-call-reports-generate-future-rank-topsis-sensitivity"></a> Future Option: TOPSIS with Sensitivity Analysis: Potential v2 branch where TOPSIS ranking is complemented by sensitivity analysis to assess robustness.

### Report generation notes

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

#### Analytic Hierarchy Process (AHP)

Derive criteria weights from pairwise comparisons and turn qualitative judgments into a consistent numerical weighting system.

#### TOPSIS

Rank alternatives by their distance from an ideal best and an ideal worst solution.

## Validation flows

Graph view for validating an input config file.

### Input config validation graph

Text graph for the validate-config call chain.

- <a id="graph-node-call-validation-input-config"></a> Validate Input Config Call: Top-level CLI call flow for validating an input configuration file before any decision analysis runs.
  - <a id="graph-node-call-validation-input-config-parse-args"></a> Parse Validation Arguments: Parse CLI arguments for the validate command, including the config path and output flags.
    - <a id="graph-node-call-validation-input-config-load-cue-config"></a> Load CUE Config: Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.
      - <a id="graph-node-call-validation-input-config-validate-model"></a> Validate Config Model: Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data.
        - <a id="graph-node-call-validation-input-config-validate-model-check-structure"></a> Check Config Structure: Check that the loaded config matches the expected top-level shape, required sections, and field types after CUE evaluation.
          - <a id="graph-node-call-validation-input-config-validate-model-check-references"></a> Check Named References: Check that all named references resolve, including criteria names, scenario names, alternative names, and report focus selectors.
            - <a id="graph-node-call-validation-input-config-validate-model-check-pairwise-comparisons"></a> Check Pairwise Comparisons: Check that pairwise comparisons are valid for each scenario, with known criteria, no self-comparisons, and sufficient coverage for AHP weighting.
              - <a id="graph-node-call-validation-input-config-validate-model-check-evaluation-coverage"></a> Check Evaluation Coverage: Check that evaluations reference known scenarios and alternatives and provide the values required by each scenario's active criteria.
                - <a id="graph-node-call-validation-input-config-validate-model-check-constraints"></a> Check Scenario Constraints: Check that scenario constraints target known criteria and use operators and values that are compatible with the referenced criterion types.
                  - <a id="graph-node-call-validation-input-config-validate-model-check-report-definitions"></a> Check Report Definitions: Check that report definitions use supported formats, valid focus selectors, and well-formed argument lists for later Cobra-style parsing.

### Input config validation notes

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

#### Clear Representation of Pairwise Judgments (v1)

Represent pairwise comparisons explicitly with named criteria instead of positional matrices so humans and AI can validate and generate them.

#### Human and AI Friendly Input Format (v1)

Use a semantic format such as CUE that remains readable for humans and AI systems while supporting strong validation.

#### Handling Incomplete Information (v1)

Detect missing comparisons or evaluation values early and return actionable diagnostics instead of attempting opaque implicit recovery.

#### Model Validation (v1)

Validate referenced criteria, pairwise comparison completeness, and alternative evaluation coverage before computation.

#### Constraint Enforcement (v1)

Allow scenarios to define hard requirements that can exclude alternatives before ranking.

