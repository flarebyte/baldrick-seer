# Validation Failed

The input config is invalid. No ranking was produced.

## Diagnostics

Each diagnostic includes a machine-oriented `path` and a more readable `location`.

- `UNKNOWN_REFERENCE` at `evaluations[1].scenarioName`
  - Location: `evaluations/regulated_growthh/scenarioName`
  - Unknown scenario name: `regulated_growthh`
- `MISSING_EVALUATION_VALUE` at `evaluations[0].evaluations[1].values`
  - Location: `evaluations/lean_startup/provider_b/values/speed`
  - Missing value for active criterion: `speed`
- `INVALID_REPORT_ARGUMENT` at `reports[0].arguments[2]`
  - Location: `reports/hosting-choice-summary/arguments/include-score=yes`
  - Unsupported argument: `include-score=yes`

## Next step

Fix the reported config errors and run validation again.
