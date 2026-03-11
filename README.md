# baldrick-seer

Scenario-based decision support for comparing alternatives across multiple criteria and multiple future contexts.

![Baldrick Seer Hero](./baldrick-seer-hero.jpg)

## What it is

`baldrick-seer` is designed for decisions where one answer must survive more than one situation.

Instead of asking "what is the best option overall?", it lets you ask:
- which option is best for a startup?
- which option is best for a regulated environment?
- which option still looks good when the context changes?

The v1 design is centered on a CLI that reads a structured CUE configuration, validates it, computes weights with AHP, ranks alternatives with TOPSIS, and emits reports for both humans and tools.

## What you should expect from v1

### Core decision workflow

- Define a decision problem with named alternatives, criteria, scenarios, and evaluations.
- Reuse the same alternatives across several scenarios instead of rebuilding the model each time.
- Keep scenario setup separate from scenario evaluations, so re-evaluating a model is easy.
- Validate the input config before any ranking happens.
- Derive criteria weights from pairwise comparisons with AHP.
- Rank alternatives with TOPSIS.
- Aggregate scenario results with practical v1 strategies such as equal or weighted averaging.

### Input model

- CUE is the configuration source of truth.
- The model is designed to stay readable for both humans and AI.
- Criteria, scenarios, alternatives, reports, and references use meaningful names instead of opaque ids.
- Scenarios can define:
  - active criteria
  - pairwise preferences
  - hard constraints
  - their own alternative evaluations
- Report definitions can carry CLI-style `key=value` arguments for later customization.

### Output

- Human-readable markdown reports.
- Machine-readable JSON output.
- CSV output for spreadsheet or analytics workflows.
- Deterministic output ordering, so repeated runs are stable and reviewable.
- Explainable and traceable results rather than bare rankings only.

### Validation

The v1 design expects the CLI to catch problems early, including:
- invalid config structure
- unknown references
- broken pairwise comparisons
- missing evaluation coverage
- invalid constraints
- invalid report definitions

### Typical v1 use cases

- technology platform selection
- infrastructure strategy planning
- system design selection
- supplier or service-provider comparison
- product feature prioritization
- scenario-based ranking of strategic options

## What is not the focus of v1

These are intentionally treated as later extensions in the current design:
- robustness-focused post-analysis
- sensitivity analysis as part of the ranking workflow
- ELECTRE as an alternative ranking method
- more general method extensibility beyond the initial AHP + TOPSIS pipeline

## Example shape

The current design examples use a scenario-based MCDA model with:
- one problem definition
- a criteria catalog
- a list of alternatives
- a list of scenarios
- a separate evaluation section keyed by scenario name
- a list of requested reports

See:
- [Overview](./doc/design/overview.md)
- [Examples](./doc/design/examples.md)
- [Glossary](./doc/design/glossary.md)
- [Implementation Notes](./doc/design/implementation.md)
- [Execution Flows](./doc/design/flows.md)
- [Use Cases](./doc/design/use-cases.md)

## Status

The repository now includes the full v1 CLI pipeline with:
- Go CLI implementation
- Cobra command wiring
- CUE loading and validation
- AHP weighting
- TOPSIS scenario ranking
- cross-scenario aggregation
- markdown, JSON, and CSV rendering
- Bun + TypeScript end-to-end tests

## Local development

```sh
make build
make test
make test-unit
make coverage
make test-race
make test-e2e
make lint
make format
```

`make test-unit` runs verbose Go tests, writes a coverage profile to `tmp/test-unit.coverage.out`, and prints a package/function coverage summary.

`make coverage` additionally writes an HTML report to `tmp/test-unit.coverage.html`.

The build target embeds deterministic build metadata inputs through linker flags:

```sh
make build VERSION=v1.0.0 COMMIT=$(git rev-parse --short HEAD) BUILD_DATE=2026-03-11T00:00:00Z
```

## CLI usage

```sh
./.e2e-bin/seer validate --config testdata/config/minimal.cue
./.e2e-bin/seer report generate --config testdata/config/minimal.cue
./.e2e-bin/seer validate --config testdata/config_split
```

`--config` accepts either a single `.cue` file or a directory containing a CUE package.

## Release preparation

Keep release preparation local and explicit:
- run `make format`
- run `make lint`
- run `make test-unit`
- run `make test-race`
- run `make test-e2e`
- build with explicit metadata using `make build VERSION=... COMMIT=... BUILD_DATE=...`

This repository intentionally does not include CI workflows, goreleaser, containers, or publishing automation at this stage.

## Design references

See:
- [Overview](./doc/design/overview.md)
- [Examples](./doc/design/examples.md)
- [Glossary](./doc/design/glossary.md)
- [Implementation Notes](./doc/design/implementation.md)
- [Execution Flows](./doc/design/flows.md)
- [Use Cases](./doc/design/use-cases.md)

## Why this shape

Many decision tools are either too simple for real trade-offs or too heavy for day-to-day use.

The goal here is a middle ground:
- explicit enough to be reliable
- readable enough to edit
- structured enough to automate
- small enough to fit in a human or AI head
