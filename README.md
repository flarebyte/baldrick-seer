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

## Install

Recommended install method:

```sh
brew install flarebyte/tap/baldrick-seer
```

Check that the CLI is available:

```sh
seer -v
seer
```

## Build from source

```sh
make build
```

The local development CLI binary is written to:

```text
.e2e-bin/seer
```

If you build from source and want to run `seer` directly, add `.e2e-bin/` to your shell `PATH`.

Release binaries for supported operating systems are written to:

```text
build/
```

`make build-dist` currently produces:
- `seer-darwin-arm64`
- `seer-linux-amd64`
- `checksums.txt`

You can embed build metadata when needed:

```sh
make build VERSION=v1.0.0 COMMIT=$(git rev-parse --short HEAD) BUILD_DATE=2026-03-12T00:00:00Z
```

## CLI usage

```sh
seer validate --config testdata/config/minimal.cue
seer report generate --config testdata/config/minimal.cue
seer validate --config testdata/config_split
```

`--config` accepts either a single `.cue` file or a directory containing a CUE package.

## Quick start examples

If you want a working model to copy and adapt, start with:
- [examples/hello-world.cue](./examples/hello-world.cue) for the smallest single-file markdown example
- [examples/hello-world-json.cue](./examples/hello-world-json.cue) for a minimal JSON report example
- [examples/hello-world-package](./examples/hello-world-package) for a split CUE package loaded from a directory

Try them with:

```sh
seer validate --config examples/hello-world.cue
seer report generate --config examples/hello-world.cue
seer report generate --config examples/hello-world-json.cue
seer validate --config examples/hello-world-package
seer report generate --config examples/hello-world-package
```

These examples are intentionally small:
- one criterion
- one alternative
- one scenario
- one report

They are meant as a quick-start baseline, not a full demonstration of the v1 feature set.

## Output formats

`seer report generate` currently renders:
- Markdown as a standalone, source-backed decision report
- JSON as a standalone, source-backed machine-readable report
- CSV for spreadsheet and analytics workflows

The output is deterministic across repeated runs for the same input.

## Validation behavior

`seer validate` checks the loaded model before any ranking or report generation. Current validation covers:
- top-level structure
- duplicate names
- reference resolution
- AHP pairwise comparison completeness and shape
- evaluation coverage and value types
- scenario constraint definitions
- report definition and argument validation

Validation failures return deterministic error output with concise remediation guidance.

## Learn more

See:
- [Overview](./doc/design/overview.md)
- [Examples](./doc/design/examples.md)
- [Glossary](./doc/design/glossary.md)
- [Implementation Notes](./doc/design/implementation.md)
- [Execution Flows](./doc/design/flows.md)
- [Use Cases](./doc/design/use-cases.md)
- [Contributing](./CONTRIBUTING.md)

## Contributing

Contributor workflow, development commands, coverage, and release-preparation notes are documented in [CONTRIBUTING.md](./CONTRIBUTING.md).

## Why this shape

Many decision tools are either too simple for real trade-offs or too heavy for day-to-day use.

The goal here is a middle ground:
- explicit enough to be reliable
- readable enough to edit
- structured enough to automate
- small enough to fit in a human or AI head
