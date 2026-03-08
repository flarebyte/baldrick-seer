# Execution Flows

Call-oriented flows for CLI validation and report generation.

## Report generation flows

Graph view for generating reports from a validated input config.

### Report generation graph

Text graph for the report-generation call chain, reusing the shared validation path.

- <a id="graph-node-call-reports-generate"></a> Generate Reports Call: Top-level CLI call flow for generating reports from an input decision model.
  - <a id="graph-node-call-reports-generate-parse-args"></a> Parse Report Arguments: Parse CLI arguments for report generation, including the config path, requested report names, and output options.
    - <a id="graph-node-call-reports-generate-shared-validation"></a> Reuse Shared Validation Flow: Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs.
      - <a id="graph-node-call-reports-generate-compute-ahp-weights"></a> Compute Criteria Weights with AHP: Transform pairwise scenario preferences into normalized criteria weights using Analytic Hierarchy Process.
        - <a id="graph-node-call-reports-generate-rank-alternatives-topsis"></a> Rank Alternatives with TOPSIS: Use the validated evaluations and AHP-derived weights to rank alternatives with TOPSIS.
          - <a id="graph-node-call-reports-generate-render-output"></a> Render Requested Reports: Render the requested markdown, JSON, or CSV reports from the computed ranking results.

### Report generation notes

#### Generate Reports Call

Top-level CLI call flow for generating reports from an input decision model.

#### Compute Criteria Weights with AHP

Transform pairwise scenario preferences into normalized criteria weights using Analytic Hierarchy Process.

#### Parse Report Arguments

Parse CLI arguments for report generation, including the config path, requested report names, and output options.

#### Rank Alternatives with TOPSIS

Use the validated evaluations and AHP-derived weights to rank alternatives with TOPSIS.

#### Render Requested Reports

Render the requested markdown, JSON, or CSV reports from the computed ranking results.

#### Reuse Shared Validation Flow

Reuse the same CUE loading and model validation path as the dedicated validate command before any scoring runs.

#### Load CUE Config

Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.

#### Validate Config Model

Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data.

## Validation flows

Graph view for validating an input config file.

### Input config validation graph

Text graph for the validate-config call chain.

- <a id="graph-node-call-validation-input-config"></a> Validate Input Config Call: Top-level CLI call flow for validating an input configuration file before any decision analysis runs.
  - <a id="graph-node-call-validation-input-config-parse-args"></a> Parse Validation Arguments: Parse CLI arguments for the validate command, including the config path and output flags.
    - <a id="graph-node-call-validation-input-config-load-cue-config"></a> Load CUE Config: Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.
      - <a id="graph-node-call-validation-input-config-validate-model"></a> Validate Config Model: Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data.

### Input config validation notes

#### Validate Input Config Call

Top-level CLI call flow for validating an input configuration file before any decision analysis runs.

#### Load CUE Config

Load and evaluate the CUE configuration package so the CLI works with a concrete validated config value.

#### Parse Validation Arguments

Parse CLI arguments for the validate command, including the config path and output flags.

#### Validate Config Model

Run structural and graph validation on the loaded config and emit diagnostics for any invalid references or incomplete model data.

