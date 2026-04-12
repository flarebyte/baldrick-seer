# Contributing

## Development setup

This repository uses:
- Go for the CLI and unit tests
- Bun + TypeScript for E2E tests
- Biome for TypeScript/tooling formatting and linting

Build the CLI locally with:

```sh
make build
```

The binary is written to:

```text
.e2e-bin/seer
```

## Main Make targets

Use these public commands for day-to-day development:

```sh
make build
make test
make test-unit
make test-race
make test-e2e
make lint
make format
make review
make coverage
```

The public targets delegate to more specific targets:

- `make build` -> `make build-go`
- `make test` -> `make test-go` and `make test-e2e`
- `make test-go` -> `make test-unit`
- `make coverage` -> `make coverage-go`
- `make lint` -> `make lint-go` and `make lint-e2e`
- `make format` -> `make format-go` and `make format-e2e`
- `make review` -> `make format`, `make test`, and `make lint`

Run `make help` to see the full target list.

## Test workflow

Recommended local verification before opening a change:

```sh
make review
make test-race
```

`make test-unit` runs verbose Go tests, writes a coverage profile to:

```text
tmp/test-unit.coverage.out
```

and prints a package/function summary.

`make coverage` additionally writes:

```text
tmp/test-unit.coverage.html
```

## Build metadata

The build supports deterministic metadata injection through linker flags:

```sh
make build VERSION=v1.0.0 COMMIT=$(git rev-parse --short HEAD) BUILD_DATE=2026-03-12T00:00:00Z
```

## Release preparation

Release preparation is intentionally local and explicit. The repository does not include CI workflows, goreleaser, containerization, or publishing automation.

Typical local release-prep flow:

```sh
make review
make test-race
make build-dist VERSION=... COMMIT=... BUILD_DATE=...
```

The release build currently targets:
- macOS ARM64
- Linux AMD64

## Design docs

The product and execution design lives under:

```text
doc/design/
```

If you need to regenerate derived design docs:

```sh
make doc-design
```

Decision records live under:

```text
doc/decision-meta/
doc/decision/
```

Each `*.seer.cue` file under `doc/decision-meta/` is treated as one decision config and generates a mirrored markdown report under `doc/decision/`.

To validate all decision configs and regenerate their markdown reports:

```sh
make doc-decision
```
