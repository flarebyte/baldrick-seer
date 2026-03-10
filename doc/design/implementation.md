# Implementation Considerations

Engineering guidance, conventions, and implementation choices for the CLI.

## Engineering conventions

Implementation rules intended to keep the codebase readable, testable, and deterministic.

### Code structure

#### Deterministic Output Ordering (v1)

Guarantee deterministic ordering in generated outputs so repeated runs produce stable markdown, JSON, and CSV artifacts.

#### Guard Clauses and Early Returns (v1)

Prefer early returns and guard clauses for error handling so failure paths stay short, obvious, and easy to test.

#### I/O and Core Logic Separation (v1)

Separate filesystem, terminal, and config-loading I/O from core decision logic so the computation pipeline stays testable and deterministic.

#### Named Predicates over Boolean Soup (v1)

Replace tangled boolean expressions with named predicates so validation and ranking rules read like domain logic instead of control noise.

#### Small Single-Purpose Functions (v1)

Keep functions small and single-purpose so validation, weighting, ranking, and rendering logic remain easy to understand and reuse.

#### Tiny Structs over Long Parameter Lists (v1)

Use small focused structs to carry grouped inputs instead of long parameter lists that are brittle and hard to read.

## Implementation stack

Primary languages, libraries, and tools chosen for the first release.

### Runtime and tooling

#### Cobra Command and Argument Parsing (v1)

Use Cobra for CLI command structure and argument parsing so command behavior and report argument handling follow one consistent model, while keeping report arguments extensible in representation but strictly validated against documented v1 argument definitions.

#### Go CLI Implementation (v1)

Implement the production CLI in Go so the tool remains fast, portable, and straightforward to distribute.

#### CUE as Configuration Source of Truth (v1)

Use CUE as the configuration source of truth so schema, defaults, validation, and concrete config evaluation live in one place.

#### Extensible Decision Methods (v2)

Extend the v1 AHP + TOPSIS pipeline so additional MCDA methods can be added later without replacing the overall CLI shape. Future methods such as ELECTRE, PROMETHEE, or VIKOR may require additional optional metadata or method-specific configuration beyond the v1 model.

#### Bun and TypeScript for E2E Tests (v1)

Implement end-to-end tests in TypeScript with Bun so CLI scenarios can be expressed tersely while staying fast to run in CI.

## User experience and output

Guidance for readable, reproducible, and automatable execution.

### CLI and explainability

#### Sensitivity and Robustness Analysis (v2)

Add post-ranking analysis that tests how changes in criteria importance or scenario assumptions affect the final result so users can judge stability.

#### Structured Output for Automation (v1)

Provide machine-readable output such as JSON in addition to human-readable summaries, while keeping validation diagnostics distinct from successful ranking reports.

#### Readable CLI Output (v1)

Present results in a clear terminal-friendly format with summaries, tables, and scenario breakdowns.

#### Explainable Results (v1)

Explain ranking outputs in terms of criteria influence and scenario differences.

#### Traceable Decision Process (v1)

Show the reasoning path from inputs to outputs, including scenario-local criterion weights, scenario aggregation weights, and contribution of each factor.

#### Reproducible Decision Runs (v1)

Running the same model with the same inputs should always produce identical results for auditing and comparison.

#### Guidance for Model Creation (v2)

Provide richer prompts and guidance that help users define criteria, comparisons, and scenario descriptions with fewer modeling errors.

