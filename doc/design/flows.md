# Execution Flows

Call-oriented flows for CLI validation and report generation.

## Report generation flows

Graph view for generating reports from a validated input config.

### Report generation graph

Text graph for the report-generation call chain, reusing the shared validation path.

- <a id="graph-node-call-reports-generate"></a> Generate Reports Call: Top-level CLI call flow for generating ranking reports from an input decision model. The command reuses the shared validation path and fails fast if the model is invalid.
  - <a id="graph-node-call-reports-generate-parse-args"></a> Parse Report Arguments: Parse CLI arguments for report generation, including the config path, requested report names, and output options.
    - <a id="graph-node-call-reports-generate-select-reports"></a> Select Requested Reports: Resolve which report definitions should run, applying any CLI filtering by report name or output target.
      - <a id="graph-node-call-reports-generate-shared-validation"></a> Reuse Shared Validation Flow: Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs. If validation fails, report generation stops immediately and no ranking report is produced.
        - <a id="graph-node-call-reports-generate-build-ahp-inputs"></a> Build AHP Inputs: Collect the validated full pairwise comparison set for each scenario into the normalized input structures needed for AHP computation of scenario-local criterion weights.
          - <a id="graph-node-call-reports-generate-compute-ahp-weights"></a> Compute Criteria Weights with AHP: Transform pairwise criterion comparisons within each scenario into normalized scenario-local criterion weights using Analytic Hierarchy Process.
            - <a id="graph-node-call-reports-generate-select-ranking-strategy"></a> Select Ranking Strategy: Select the ranking pipeline after computing scenario-local criterion weights with AHP. In v1, the design is built around an AHP + TOPSIS pipeline; v2 may add alternatives such as ELECTRE or TOPSIS followed by sensitivity analysis.
              - <a id="graph-node-call-reports-generate-build-topsis-inputs"></a> Build TOPSIS Inputs: Combine validated evaluations, criterion polarity, and AHP-derived scenario-local criterion weights into the decision matrices required by TOPSIS.
                - <a id="graph-node-call-reports-generate-rank-alternatives-topsis"></a> Rank Alternatives with TOPSIS: Use the validated evaluations and scenario-local criterion weights derived with AHP to rank alternatives with TOPSIS.
                  - <a id="graph-node-call-reports-generate-render-output"></a> Render Requested Reports: Render the requested markdown, JSON, or CSV outputs only after validation succeeds and ranking results are computed. Invalid models do not reach report rendering.
                    - <a id="graph-node-call-reports-generate-render-output-render-csv"></a> Render CSV Report: Render flat tabular CSV output for spreadsheet analysis and data exchange.
                    - <a id="graph-node-call-reports-generate-render-output-render-json"></a> Render JSON Report: Render machine-readable JSON ranking output for automation, downstream processing, and reproducibility only when validation succeeds. If JSON output is requested and validation fails, the command may emit structured diagnostics as an error payload or via stderr, but that output is not a successful ranking report.
                    - <a id="graph-node-call-reports-generate-render-output-render-markdown"></a> Render Markdown Report: Render narrative markdown output for human readers, including rankings, explanations, and scenario summaries.
              - <a id="graph-node-call-reports-generate-future-rank-electre"></a> Future Option: Rank with ELECTRE: Potential v2 branch where the validated model is ranked with ELECTRE instead of TOPSIS.
              - <a id="graph-node-call-reports-generate-future-rank-topsis-sensitivity"></a> Future Option: TOPSIS with Sensitivity Analysis: Potential v2 branch where TOPSIS ranking is complemented by sensitivity analysis to assess robustness.

### Report generation notes

#### Generate Reports Call

Top-level CLI call flow for generating ranking reports from an input decision model. The command reuses the shared validation path and fails fast if the model is invalid.

#### Build AHP Inputs

Collect the validated full pairwise comparison set for each scenario into the normalized input structures needed for AHP computation of scenario-local criterion weights.

#### Build TOPSIS Inputs

Combine validated evaluations, criterion polarity, and AHP-derived scenario-local criterion weights into the decision matrices required by TOPSIS.

#### Compute Criteria Weights with AHP

Transform pairwise criterion comparisons within each scenario into normalized scenario-local criterion weights using Analytic Hierarchy Process.

#### Future Option: Rank with ELECTRE

Potential v2 branch where the validated model is ranked with ELECTRE instead of TOPSIS.

#### Future Option: TOPSIS with Sensitivity Analysis

Potential v2 branch where TOPSIS ranking is complemented by sensitivity analysis to assess robustness.

#### Parse Report Arguments

Parse CLI arguments for report generation, including the config path, requested report names, and output options.

#### Rank Alternatives with TOPSIS

Use the validated evaluations and scenario-local criterion weights derived with AHP to rank alternatives with TOPSIS.

#### Render Requested Reports

Render the requested markdown, JSON, or CSV outputs only after validation succeeds and ranking results are computed. Invalid models do not reach report rendering.

#### Render CSV Report

Render flat tabular CSV output for spreadsheet analysis and data exchange.

#### Render JSON Report

Render machine-readable JSON ranking output for automation, downstream processing, and reproducibility only when validation succeeds. If JSON output is requested and validation fails, the command may emit structured diagnostics as an error payload or via stderr, but that output is not a successful ranking report.

#### Render Markdown Report

Render narrative markdown output for human readers, including rankings, explanations, and scenario summaries.

#### Select Ranking Strategy

Select the ranking pipeline after computing scenario-local criterion weights with AHP. In v1, the design is built around an AHP + TOPSIS pipeline; v2 may add alternatives such as ELECTRE or TOPSIS followed by sensitivity analysis.

#### Select Requested Reports

Resolve which report definitions should run, applying any CLI filtering by report name or output target.

#### Reuse Shared Validation Flow

Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs. If validation fails, report generation stops immediately and no ranking report is produced.

#### Load CUE Config

Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.

#### Validate Config Model

Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data. For the `validate` command, this is the terminal result of the command.

## Validation flows

Graph view for validating an input config file.

### Input config validation graph

Text graph for the validate-config call chain.

- <a id="graph-node-call-validation-input-config"></a> Validate Input Config Call: Top-level CLI call flow for validating an input configuration file and returning validation results only, without scoring or report generation.
  - <a id="graph-node-call-validation-input-config-parse-args"></a> Parse Validation Arguments: Parse CLI arguments for the validate command, including the config path and output flags.
    - <a id="graph-node-call-validation-input-config-load-cue-config"></a> Load CUE Config: Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.
      - <a id="graph-node-call-validation-input-config-validate-model"></a> Validate Config Model: Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data. For the `validate` command, this is the terminal result of the command.
        - <a id="graph-node-call-validation-input-config-validate-model-check-structure"></a> Check Config Structure: Check that the loaded config matches the expected top-level shape, required sections, and field types after CUE evaluation.
          - <a id="graph-node-call-validation-input-config-validate-model-check-references"></a> Check Named References: Check that all named references resolve, including criteria names, scenario names, alternative names, and report focus selectors.
            - <a id="graph-node-call-validation-input-config-validate-model-check-pairwise-comparisons"></a> Check Pairwise Comparisons: Check that each scenario using AHP provides pairwise comparisons only between known active criteria, never compares a criterion with itself, and includes exactly one canonical comparison for every unordered pair of distinct active criteria. Reject duplicate comparisons, inverse duplicates, or any missing pair.
              - <a id="graph-node-call-validation-input-config-validate-model-check-evaluation-coverage"></a> Check Evaluation Coverage: Check that evaluations reference known scenarios and alternatives and provide supported v1 criterion values for each scenario's active criteria: measurable numbers, integer ordinals, or booleans with only `true` and `false` values.
                - <a id="graph-node-call-validation-input-config-validate-model-check-constraints"></a> Check Scenario Constraints: Check that each scenario constraint uses an operator and value compatible with the referenced criterion type: number criteria allow numeric values with `<=`, `>=`, `=`, or `!=`; ordinal criteria allow integer values with `<=`, `>=`, `=`, or `!=`; boolean criteria allow only `=` or `!=` with `true` or `false`. Invalid operator/type combinations must raise a validation error.
                  - <a id="graph-node-call-validation-input-config-validate-model-check-report-definitions"></a> Check Report Definitions: Check that report definitions use supported formats, valid focus selectors, and strictly validated report arguments. In v1 every report argument must use `key=value`, unknown arguments are validation errors, argument names must be allowed globally or for the selected format, format-specific arguments must match the report format, invalid values must be rejected, and duplicate keys are invalid unless explicitly defined otherwise.

### Input config validation notes

#### Validate Input Config Call

Top-level CLI call flow for validating an input configuration file and returning validation results only, without scoring or report generation.

#### Load CUE Config

Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.

#### Parse Validation Arguments

Parse CLI arguments for the validate command, including the config path and output flags.

#### Validate Config Model

Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data. For the `validate` command, this is the terminal result of the command.

#### Check Scenario Constraints

Check that each scenario constraint uses an operator and value compatible with the referenced criterion type: number criteria allow numeric values with `<=`, `>=`, `=`, or `!=`; ordinal criteria allow integer values with `<=`, `>=`, `=`, or `!=`; boolean criteria allow only `=` or `!=` with `true` or `false`. Invalid operator/type combinations must raise a validation error.

#### Check Evaluation Coverage

Check that evaluations reference known scenarios and alternatives and provide supported v1 criterion values for each scenario's active criteria: measurable numbers, integer ordinals, or booleans with only `true` and `false` values.

#### Check Pairwise Comparisons

Check that each scenario using AHP provides pairwise comparisons only between known active criteria, never compares a criterion with itself, and includes exactly one canonical comparison for every unordered pair of distinct active criteria. Reject duplicate comparisons, inverse duplicates, or any missing pair.

#### Check Named References

Check that all named references resolve, including criteria names, scenario names, alternative names, and report focus selectors.

#### Check Report Definitions

Check that report definitions use supported formats, valid focus selectors, and strictly validated report arguments. In v1 every report argument must use `key=value`, unknown arguments are validation errors, argument names must be allowed globally or for the selected format, format-specific arguments must match the report format, invalid values must be rejected, and duplicate keys are invalid unless explicitly defined otherwise.

#### Check Config Structure

Check that the loaded config matches the expected top-level shape, required sections, and field types after CUE evaluation.

