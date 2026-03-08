# Validation Failed

The input config is invalid. No ranking was produced.

## Diagnostics

- `UNKNOWN_REFERENCE` at `evaluations[1].scenarioName`
  - Unknown scenario name: `regulated_growthh`
- `MISSING_EVALUATION_VALUE` at `evaluations[0].evaluations[1].values`
  - Missing value for active criterion: `speed`
- `INVALID_REPORT_ARGUMENT` at `reports[0].arguments[2]`
  - Unsupported argument: `include-score=yes`

## Next step

Fix the reported config errors and run validation again.
