# Execution Flows

Call-oriented flows for CLI execution and validation.

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

