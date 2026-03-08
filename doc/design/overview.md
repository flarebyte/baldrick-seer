# Design Overview

High-level overview of the baldrick-seer decision-model design.

## Scope

Primary use cases and example design anchors.

### Primary use cases

#### General Multi-Criteria Ranking (v1)

Run a scenario-based CLI evaluation that ranks named alternatives from a structured config and emits decision reports for humans and tools.

#### Robust Choice Identification (v2)

Identify alternatives that remain strong across scenarios and under changing assumptions, likely using robustness or sensitivity-oriented post-analysis.

#### System Design Selection (v1)

Compare system designs where trade-offs exist between cost, scalability, and reliability.

#### Infrastructure Strategy Planning (v1)

Assess infrastructure alternatives such as cloud providers or deployment models under varying operational scenarios.

#### Technology Platform Selection (v1)

Compare multiple technology platforms across operational scenarios such as startup, scale-up, and enterprise maturity.

### Reference examples

#### Hosting Choice Example

Minimal scenario-based MCDA example that compares two hosting providers across lean-startup and regulated-growth scenarios using equal-average aggregation.

#### Platform Selection Example

Scenario-based MCDA example that compares candidate platforms across startup, unicorn, and established-enterprise contexts using weighted-average aggregation.

