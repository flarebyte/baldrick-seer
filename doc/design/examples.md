# Example Models

Scenario-based MCDA examples translated from the TypeScript design fixtures.

## Example outputs

Illustrative output shapes for the hosting-choice example across the main v1 report formats.

### Illustrative note

#### Illustrative Output Examples

The example outputs in this section are illustrative only. They demonstrate output structure, field layout, and report formats for markdown, JSON, and CSV; the shown scores and rankings are not normative fixtures and should not be read as exact or reproducible algorithmic results from the example models.

### Illustrative reports

#### Hosting Choice CSV Output

| alternative | rank | scenario | score |
| --- | --- | --- | --- |
| provider_a | 1 | lean_startup | 0.61 |
| provider_b | 2 | lean_startup | 0.39 |
| provider_b | 1 | regulated_growth | 0.72 |
| provider_a | 2 | regulated_growth | 0.28 |
| provider_b | 1 | overall | 0.555 |
| provider_a | 2 | overall | 0.445 |

#### Hosting Choice JSON Output

```json
{
  "problemName": "hosting-choice",
  "reportName": "hosting-choice-results",
  "format": "json",
  "aggregation": {
    "method": "equal_average"
  },
  "scenarioResults": [
    {
      "scenarioName": "lean_startup",
      "ranking": [
        { "alternativeName": "provider_a", "rank": 1, "score": 0.61 },
        { "alternativeName": "provider_b", "rank": 2, "score": 0.39 }
      ]
    },
    {
      "scenarioName": "regulated_growth",
      "ranking": [
        { "alternativeName": "provider_b", "rank": 1, "score": 0.72 },
        { "alternativeName": "provider_a", "rank": 2, "score": 0.28 }
      ]
    }
  ],
  "finalRanking": [
    { "alternativeName": "provider_b", "rank": 1, "score": 0.555 },
    { "alternativeName": "provider_a", "rank": 2, "score": 0.445 }
  ]
}
```

#### Hosting Choice Markdown Output

```markdown
# Hosting Choice Summary

Illustrative example only. The scores and rankings below are shown to demonstrate report structure and narrative layout, not as exact normative results.

## Result

`Provider B` appears as the overall choice in this illustrative example, with an example aggregated score of `0.555`.

## Why

- `Provider A` is stronger on cost in the lean-startup scenario.
- `Provider B` is materially stronger on speed and compliance.
- With equal scenario aggregation in this example layout, `Provider B` leads on combined score.

## Scenario ranking

### Lean Startup

1. `Provider A` (`0.61`)
2. `Provider B` (`0.39`)

### Regulated Growth

1. `Provider B` (`0.72`)
2. `Provider A` (`0.28`)

## Final ranking

1. `Provider B` (`0.555`)
2. `Provider A` (`0.445`)
```

### Validation failures

#### Validation Failure JSON Output

```json
{
  "status": "invalid",
  "problemName": "hosting-choice",
  "diagnostics": [
    {
      "code": "UNKNOWN_REFERENCE",
      "severity": "error",
      "path": "evaluations[1].scenarioName",
      "location": "evaluations/regulated_growthh/scenarioName",
      "message": "Unknown scenario name: regulated_growthh"
    },
    {
      "code": "MISSING_EVALUATION_VALUE",
      "severity": "error",
      "path": "evaluations[0].evaluations[1].values",
      "location": "evaluations/lean_startup/provider_b/values/speed",
      "message": "Missing value for active criterion: speed"
    },
    {
      "code": "INVALID_REPORT_ARGUMENT",
      "severity": "error",
      "path": "reports[0].arguments[2]",
      "location": "reports/hosting-choice-summary/arguments/include-score=yes",
      "message": "Unsupported argument: include-score=yes"
    }
  ]
}
```

#### Validation Failure Markdown Output

```markdown
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
```

## Hosting choice

Minimal example showing the same shape with fewer criteria and alternatives.

### Model

#### Hosting Choice Example

Minimal scenario-based MCDA example that compares two hosting providers across lean-startup and regulated-growth scenarios using equal-average aggregation.

### Scenarios

#### Lean Startup Scenario

Scenario emphasizing cost and speed when budget pressure matters more than peak performance.

#### Regulated Growth Scenario

Scenario adding compliance as a dominant concern alongside cost and speed.

## Platform selection

Three-scenario example for selecting a platform across company maturity stages.

### Model

#### Platform Selection Example

Scenario-based MCDA example that compares candidate platforms across startup, unicorn, and established-enterprise contexts using weighted-average aggregation.

### Scenarios

#### Established Enterprise Scenario

Mature-organization scenario where reliability and compliance dominate cost and future flexibility.

#### Startup Scenario

Early-stage scenario where cost and time to market matter more than enterprise controls.

#### Unicorn Scenario

Hyper-growth scenario where scalability becomes dominant while cost still matters.

## TypeScript source examples

Readable draft examples for the future CUE input model.

### Model types

#### TypeScript Input Model

```ts
export type Name = string;

export type CriterionPolarity = "benefit" | "cost";
export type ReportFormat = "markdown" | "json" | "csv";

export type ScenarioAggregationMethod =
  | "weighted_average"
  | "equal_average"
  | "maximin"
  | "minimax_regret";

export type PairwiseStrength =
  | 1
  | 2
  | 3
  | 4
  | 5
  | 6
  | 7
  | 8
  | 9;

export interface McdaModel {
  modelType: "scenario_based_mcda";
  version: "1.0";
  problem: ProblemDefinition;
  reports: ReportDefinition[];
  criteriaCatalog: CriterionDefinition[];
  alternatives: AlternativeDefinition[];
  scenarios: ScenarioDefinition[];
  evaluations: ScenarioEvaluationDefinition[];
  aggregation: ScenarioAggregationDefinition;
}

export interface ReportDefinition {
  name: Name;
  title: string;
  description?: string;
  format: ReportFormat;
  /**
   * Optional report parameters using the same key=value convention as CLI args.
   * These are intended to be parsed with the same Cobra-based argument handling
   * used by the CLI so report-level customization stays consistent.
   * In v1 the representation remains extensible, but validation is strict:
   * every entry must use key=value form, only documented arguments are allowed,
   * format-specific arguments must match the report format, invalid values are
   * rejected, and duplicate keys are invalid unless explicitly defined
   * otherwise by the spec.
   */
  arguments?: string[];
  focus?: ReportFocusDefinition;
}

export interface ReportFocusDefinition {
  scenarioNames?: Name[];
  alternativeNames?: Name[];
  criterionNames?: Name[];
}

export interface ProblemDefinition {
  name: Name;
  title: string;
  goal: string;
  description?: string;
  owner?: string;
  notes?: string[];
}

export interface CriterionDefinition {
  name: Name;
  title: string;
  description?: string;
  polarity: CriterionPolarity;
  unit?: string;
  valueType?: "number" | "ordinal" | "boolean";
  /**
   * For ordinal criteria in v1, document the meaning of each integer level,
   * for example "1=poor, 2=fair, 3=good, 4=excellent".
   */
  scaleGuidance?: string;
}

export interface AlternativeDefinition {
  name: Name;
  title: string;
  description?: string;
  labels?: string[];
}

export interface ScenarioDefinition {
  name: Name;
  title: string;
  description?: string;

  /**
   * Human meaning of the scenario, for example:
   * "Early-stage company with budget pressure and fast experimentation."
   */
  narrative?: string;

  /**
   * Criteria used in this scenario.
   * Allows each scenario to activate only relevant criteria.
   */
  activeCriteria: ScenarioCriterionRef[];

  /**
   * AHP-style semantic comparisons, human/AI friendly.
   * These are inputs only. The engine can transform them into a matrix.
   */
  preferences?: ScenarioPreferences;

  /**
   * Optional hard rules for filtering before scoring.
   * Helpful when some scenarios have non-negotiable requirements.
   */
  constraints?: ScenarioConstraint[];
}

export interface ScenarioEvaluationDefinition {
  scenarioName: Name;
  description?: string;

  /**
   * Raw measurements or judgments for alternatives under the referenced scenario.
   * Keeping this outside the scenario definition makes re-evaluation easier.
   */
  evaluations: AlternativeScenarioEvaluation[];
}

export interface ScenarioCriterionRef {
  criterionName: Name;
  description?: string;
}

export interface ScenarioPreferences {
  method: "ahp_pairwise";
  scale: "saaty_1_9";
  /**
   * For v1, scenarios using AHP must provide exactly one comparison for every
   * unordered pair of distinct active criteria. Duplicate comparisons, inverse
   * duplicates, and self-comparisons are invalid.
   */
  comparisons: PairwiseComparison[];
}

export interface PairwiseComparison {
  /**
   * Canonical v1 direction: name the criterion judged more important for this
   * unordered pair.
   */
  moreImportantCriterionName: Name;
  /**
   * Canonical v1 direction: name the criterion judged less important for this
   * unordered pair.
   */
  lessImportantCriterionName: Name;
  strength: PairwiseStrength;
  justification?: string;
  source?: "human" | "ai" | "hybrid";
  confidence?: "low" | "medium" | "high";
}

export interface AlternativeScenarioEvaluation {
  alternativeName: Name;
  values: Record<Name, CriterionValue>;
  description?: string;
  evidence?: EvidenceRef[];
}

export type CriterionValue =
  | NumericCriterionValue
  | BooleanCriterionValue
  | OrdinalCriterionValue;

export interface NumericCriterionValue {
  kind: "number";
  /**
   * Numeric values represent measurable quantities and are used directly in
   * the decision matrix in v1.
   */
  value: number;
  estimated?: boolean;
  source?: "human" | "ai" | "hybrid" | "measured" | "imported";
  justification?: string;
}

export interface BooleanCriterionValue {
  kind: "boolean";
  /**
   * Boolean values are normalized before scoring in v1: true = 1, false = 0.
   */
  value: boolean;
  source?: "human" | "ai" | "hybrid" | "measured" | "imported";
  justification?: string;
}

export interface OrdinalCriterionValue {
  kind: "ordinal";
  /**
   * Ordinal values are integer levels in v1 and are treated numerically after
   * validation.
   */
  value: number;
  label?: string;
  source?: "human" | "ai" | "hybrid" | "measured" | "imported";
  justification?: string;
}

export interface EvidenceRef {
  label: string;
  detail?: string;
}

export interface ScenarioConstraint {
  criterionName: Name;
  operator: "<=" | ">=" | "=" | "!=";
  /**
   * Constraint values must match the referenced criterion type in v1.
   * - number criteria: numeric values with <=, >=, =, or !=
   * - ordinal criteria: integer values with <=, >=, =, or !=
   * - boolean criteria: true/false values with = or != only
   */
  value: number | boolean;
  justification?: string;
}

export interface ScenarioAggregationDefinition {
  method: ScenarioAggregationMethod;

  /**
   * Optional explicit scenario aggregation weights.
   * Recommended when method = weighted_average.
   * Keys are scenario names and this is the single source of truth for
   * cross-scenario weighting in v1.
   */
  scenarioWeights?: Record<Name, number>;

  /**
   * Optional policy note for humans / AI agents.
   */
  rationale?: string;
}
```

### Scenario-based examples

#### Hosting Choice TypeScript Example

```ts
import { McdaModel } from "./model";

export const minimalScenarioMcda: McdaModel = {
  modelType: "scenario_based_mcda",
  version: "1.0",

  problem: {
    name: "hosting-choice",
    title: "Hosting Choice",
    goal: "Choose the best hosting provider across business scenarios"
  },

  reports: [
    {
      name: "hosting-choice-summary",
      title: "Hosting Choice Summary",
      description: "Human-readable summary of rankings and scenario trade-offs.",
      format: "markdown",
      arguments: ["include-scenarios=all", "top-alternatives=2", "include-scores=true"]
    },
    {
      name: "hosting-choice-results",
      title: "Hosting Choice Results",
      description: "Structured ranking output for downstream tooling.",
      format: "json",
      arguments: ["include-evidence=false", "pretty=true"]
    },
    {
      name: "hosting-choice-scenario-scores",
      title: "Hosting Choice Scenario Scores",
      description: "Flat scenario and alternative scores for spreadsheet-style analysis.",
      format: "csv",
      arguments: ["columns=scenario,alternative,score,rank", "header=true"]
    }
  ],

  criteriaCatalog: [
    {
      name: "cost",
      title: "Cost",
      description: "Estimated monthly hosting spend.",
      polarity: "cost",
      unit: "USD/month",
      valueType: "number"
    },
    {
      name: "speed",
      title: "Speed",
      description: "Overall delivery and runtime responsiveness.",
      polarity: "benefit",
      unit: "score",
      valueType: "number"
    },
    {
      name: "compliance",
      title: "Compliance",
      description: "Ability to satisfy governance and regulatory expectations.",
      polarity: "benefit",
      unit: "score",
      valueType: "number"
    }
  ],

  alternatives: [
    { name: "provider_a", title: "Provider A" },
    { name: "provider_b", title: "Provider B" }
  ],

  scenarios: [
    {
      name: "lean_startup",
      title: "Lean Startup",
      description: "Early-stage context where budget discipline matters more than peak capability.",
      activeCriteria: [
        { criterionName: "cost" },
        { criterionName: "speed" }
      ],
      preferences: {
        method: "ahp_pairwise",
        scale: "saaty_1_9",
        comparisons: [
          {
            moreImportantCriterionName: "cost",
            lessImportantCriterionName: "speed",
            strength: 3,
            justification: "Budget matters more than peak performance."
          }
        ]
      }
    },
    {
      name: "regulated_growth",
      title: "Regulated Growth",
      description: "Scaling context where compliance becomes a first-order requirement.",
      activeCriteria: [
        { criterionName: "cost" },
        { criterionName: "speed" },
        { criterionName: "compliance" }
      ],
      preferences: {
        method: "ahp_pairwise",
        scale: "saaty_1_9",
        comparisons: [
          {
            moreImportantCriterionName: "compliance",
            lessImportantCriterionName: "cost",
            strength: 5
          },
          {
            moreImportantCriterionName: "compliance",
            lessImportantCriterionName: "speed",
            strength: 4
          },
          {
            moreImportantCriterionName: "speed",
            lessImportantCriterionName: "cost",
            strength: 2
          }
        ]
      }
    }
  ],

  evaluations: [
    {
      scenarioName: "lean_startup",
      description: "Assessment of hosting options for a lean startup context.",
      evaluations: [
        {
          alternativeName: "provider_a",
          values: {
            cost: { kind: "number", value: 100 },
            speed: { kind: "number", value: 70 }
          }
        },
        {
          alternativeName: "provider_b",
          values: {
            cost: { kind: "number", value: 180 },
            speed: { kind: "number", value: 90 }
          }
        }
      ]
    },
    {
      scenarioName: "regulated_growth",
      description: "Assessment of hosting options when compliance becomes critical.",
      evaluations: [
        {
          alternativeName: "provider_a",
          values: {
            cost: { kind: "number", value: 100 },
            speed: { kind: "number", value: 70 },
            compliance: { kind: "number", value: 60 }
          }
        },
        {
          alternativeName: "provider_b",
          values: {
            cost: { kind: "number", value: 180 },
            speed: { kind: "number", value: 90 },
            compliance: { kind: "number", value: 92 }
          }
        }
      ]
    }
  ],

  aggregation: {
    method: "equal_average"
  }
};
```

#### Platform Selection TypeScript Example

```ts
import { McdaModel } from "./model";

export const exampleScenarioBasedMcda: McdaModel = {
  modelType: "scenario_based_mcda",
  version: "1.0",

  problem: {
    name: "platform-selection",
    title: "Platform Selection",
    goal: "Select the best platform across different company growth scenarios",
    description:
      "Evaluate the same candidate platforms for startup, unicorn, and established-enterprise contexts."
  },

  reports: [
    {
      name: "platform-selection-decision-brief",
      title: "Platform Selection Decision Brief",
      description: "Narrative report for humans comparing the leading platforms across scenarios.",
      format: "markdown",
      arguments: ["include-scenarios=all", "top-alternatives=3", "explain=true"],
      focus: {
        scenarioNames: ["startup", "unicorn", "established"]
      }
    },
    {
      name: "platform-selection-machine-results",
      title: "Platform Selection Machine Results",
      description: "Structured data for automation, reproducibility, and downstream processing.",
      format: "json",
      arguments: ["include-evidence=true", "include-weights=true", "pretty=true"]
    },
    {
      name: "platform-selection-scenario-matrix",
      title: "Platform Selection Scenario Matrix",
      description: "Tabular scenario, criterion, and alternative output for analytics workflows.",
      format: "csv",
      arguments: ["columns=scenario,alternative,criterion,value,score,rank", "header=true"],
      focus: {
        criterionNames: ["cost", "time_to_market", "scalability", "reliability", "compliance"],
        alternativeNames: ["platform_a", "platform_b", "platform_c"]
      }
    }
  ],

  criteriaCatalog: [
    {
      name: "cost",
      title: "Cost",
      description: "Total operating and implementation cost",
      polarity: "cost",
      unit: "USD/month",
      valueType: "number"
    },
    {
      name: "time_to_market",
      title: "Time to Market",
      description: "How quickly the platform can be adopted",
      polarity: "cost",
      unit: "weeks",
      valueType: "number"
    },
    {
      name: "scalability",
      title: "Scalability",
      description: "Ability to support rapid growth",
      polarity: "benefit",
      unit: "score",
      valueType: "number",
      scaleGuidance: "1 to 100, higher is better"
    },
    {
      name: "reliability",
      title: "Reliability",
      description: "Expected operational reliability",
      polarity: "benefit",
      unit: "score",
      valueType: "number",
      scaleGuidance: "1 to 100, higher is better"
    },
    {
      name: "compliance",
      title: "Compliance",
      description: "Ability to satisfy governance and regulatory requirements",
      polarity: "benefit",
      unit: "score",
      valueType: "number",
      scaleGuidance: "1 to 100, higher is better"
    }
  ],

  alternatives: [
    {
      name: "platform_a",
      title: "Platform A",
      description: "Fast to adopt and relatively inexpensive"
    },
    {
      name: "platform_b",
      title: "Platform B",
      description: "Balanced option with strong scalability"
    },
    {
      name: "platform_c",
      title: "Platform C",
      description: "Enterprise-oriented option with strong reliability and compliance"
    }
  ],

  scenarios: [
    {
      name: "startup",
      title: "Startup",
      description: "Small company with budget pressure and need for rapid experimentation",
      narrative:
        "In this scenario, low cost and fast deployment matter more than enterprise controls.",
      activeCriteria: [
        { criterionName: "cost", description: "Startup budgets remain tightly constrained." },
        {
          criterionName: "time_to_market",
          description: "Delivery speed is critical for fast product iteration."
        },
        {
          criterionName: "scalability",
          description: "Growth headroom matters, but less than near-term execution."
        },
        {
          criterionName: "reliability",
          description: "Operational resilience matters, though less than cost and speed."
        }
      ],
      preferences: {
        method: "ahp_pairwise",
        scale: "saaty_1_9",
        comparisons: [
          {
            moreImportantCriterionName: "cost",
            lessImportantCriterionName: "reliability",
            strength: 4,
            justification: "Budget pressure is significant at startup stage.",
            source: "hybrid",
            confidence: "high"
          },
          {
            moreImportantCriterionName: "time_to_market",
            lessImportantCriterionName: "reliability",
            strength: 3,
            justification: "Speed is important for product iteration.",
            source: "hybrid",
            confidence: "high"
          },
          {
            moreImportantCriterionName: "cost",
            lessImportantCriterionName: "scalability",
            strength: 2,
            justification: "Scalability matters, but near-term survival matters more.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "time_to_market",
            lessImportantCriterionName: "cost",
            strength: 2,
            justification: "A slightly faster launch is preferred over marginal savings.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "scalability",
            lessImportantCriterionName: "reliability",
            strength: 2,
            justification: "Growth readiness slightly outweighs mature reliability.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "time_to_market",
            lessImportantCriterionName: "scalability",
            strength: 2,
            justification: "Immediate execution is slightly more important than future growth capacity.",
            source: "hybrid",
            confidence: "medium"
          }
        ]
      }
    },

    {
      name: "unicorn",
      title: "Unicorn",
      description: "Rapidly scaling company with strong growth pressure",
      narrative:
        "In this scenario, scalability becomes dominant, while cost still matters but less than growth readiness.",
      activeCriteria: [
        { criterionName: "cost", description: "Cost still matters, but it no longer dominates." },
        {
          criterionName: "scalability",
          description: "This scenario optimizes for rapid expansion under heavy load."
        },
        {
          criterionName: "reliability",
          description: "Service continuity remains important while scaling aggressively."
        },
        {
          criterionName: "compliance",
          description: "Governance matters, but it is not yet the dominant driver."
        }
      ],
      preferences: {
        method: "ahp_pairwise",
        scale: "saaty_1_9",
        comparisons: [
          {
            moreImportantCriterionName: "scalability",
            lessImportantCriterionName: "cost",
            strength: 5,
            justification: "Growth capacity dominates cost concerns.",
            source: "hybrid",
            confidence: "high"
          },
          {
            moreImportantCriterionName: "scalability",
            lessImportantCriterionName: "reliability",
            strength: 3,
            justification: "Scalability is moderately more important than reliability during hyper-growth.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "reliability",
            lessImportantCriterionName: "cost",
            strength: 2,
            justification: "Service continuity is slightly more important than savings.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "scalability",
            lessImportantCriterionName: "compliance",
            strength: 3,
            justification: "Compliance matters, but scale pressure is stronger in this phase.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "reliability",
            lessImportantCriterionName: "compliance",
            strength: 2,
            justification: "Reliability slightly outweighs compliance during rapid expansion.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "compliance",
            lessImportantCriterionName: "cost",
            strength: 2,
            justification: "As the company grows, governance matters more than pure cost.",
            source: "hybrid",
            confidence: "medium"
          }
        ]
      }
    },

    {
      name: "established",
      title: "Established Enterprise",
      description: "Mature organization with governance, reliability, and compliance needs",
      narrative:
        "In this scenario, operational stability and compliance dominate speed and startup efficiency.",
      activeCriteria: [
        {
          criterionName: "cost",
          description: "Cost is secondary to enterprise readiness and risk reduction."
        },
        {
          criterionName: "reliability",
          description: "Stable operations are a primary concern for established organizations."
        },
        {
          criterionName: "compliance",
          description: "Regulatory and governance fit is critical in this scenario."
        },
        {
          criterionName: "scalability",
          description: "The platform should continue to support large-scale operations."
        }
      ],
      preferences: {
        method: "ahp_pairwise",
        scale: "saaty_1_9",
        comparisons: [
          {
            moreImportantCriterionName: "reliability",
            lessImportantCriterionName: "cost",
            strength: 5,
            justification: "Operational stability strongly outweighs cost in mature environments.",
            source: "hybrid",
            confidence: "high"
          },
          {
            moreImportantCriterionName: "compliance",
            lessImportantCriterionName: "cost",
            strength: 4,
            justification: "Compliance obligations are critical.",
            source: "hybrid",
            confidence: "high"
          },
          {
            moreImportantCriterionName: "reliability",
            lessImportantCriterionName: "scalability",
            strength: 3,
            justification: "Reliability is moderately more important than future scaling potential.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "compliance",
            lessImportantCriterionName: "scalability",
            strength: 3,
            justification: "Governance is moderately more important than scaling capacity.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "reliability",
            lessImportantCriterionName: "compliance",
            strength: 2,
            justification: "Reliability is slightly more important than compliance.",
            source: "hybrid",
            confidence: "medium"
          },
          {
            moreImportantCriterionName: "scalability",
            lessImportantCriterionName: "cost",
            strength: 2,
            justification: "Future flexibility is still slightly more important than raw cost savings.",
            source: "hybrid",
            confidence: "medium"
          }
        ]
      },
      constraints: [
        {
          criterionName: "compliance",
          operator: ">=",
          value: 70,
          justification: "Enterprise scenario requires minimum compliance readiness."
        }
      ]
    }
  ],

  evaluations: [
    {
      scenarioName: "startup",
      description: "Assessment of candidate platforms for an early-stage startup context.",
      evaluations: [
        {
          alternativeName: "platform_a",
          values: {
            cost: { kind: "number", value: 8000, source: "imported" },
            time_to_market: { kind: "number", value: 4, source: "human" },
            scalability: { kind: "number", value: 70, source: "hybrid" },
            reliability: { kind: "number", value: 68, source: "hybrid" }
          }
        },
        {
          alternativeName: "platform_b",
          values: {
            cost: { kind: "number", value: 12000, source: "imported" },
            time_to_market: { kind: "number", value: 6, source: "human" },
            scalability: { kind: "number", value: 88, source: "hybrid" },
            reliability: { kind: "number", value: 80, source: "hybrid" }
          }
        },
        {
          alternativeName: "platform_c",
          values: {
            cost: { kind: "number", value: 18000, source: "imported" },
            time_to_market: { kind: "number", value: 10, source: "human" },
            scalability: { kind: "number", value: 82, source: "hybrid" },
            reliability: { kind: "number", value: 92, source: "hybrid" }
          }
        }
      ]
    },
    {
      scenarioName: "unicorn",
      description: "Assessment of candidate platforms during a hyper-growth phase.",
      evaluations: [
        {
          alternativeName: "platform_a",
          values: {
            cost: { kind: "number", value: 10000, source: "imported" },
            scalability: { kind: "number", value: 72, source: "hybrid" },
            reliability: { kind: "number", value: 70, source: "hybrid" },
            compliance: { kind: "number", value: 58, source: "hybrid" }
          }
        },
        {
          alternativeName: "platform_b",
          values: {
            cost: { kind: "number", value: 14000, source: "imported" },
            scalability: { kind: "number", value: 93, source: "hybrid" },
            reliability: { kind: "number", value: 84, source: "hybrid" },
            compliance: { kind: "number", value: 76, source: "hybrid" }
          }
        },
        {
          alternativeName: "platform_c",
          values: {
            cost: { kind: "number", value: 21000, source: "imported" },
            scalability: { kind: "number", value: 87, source: "hybrid" },
            reliability: { kind: "number", value: 94, source: "hybrid" },
            compliance: { kind: "number", value: 91, source: "hybrid" }
          }
        }
      ]
    },
    {
      scenarioName: "established",
      description: "Assessment of candidate platforms for a mature enterprise context.",
      evaluations: [
        {
          alternativeName: "platform_a",
          values: {
            cost: { kind: "number", value: 11000, source: "imported" },
            reliability: { kind: "number", value: 69, source: "hybrid" },
            compliance: { kind: "number", value: 60, source: "hybrid" },
            scalability: { kind: "number", value: 72, source: "hybrid" }
          }
        },
        {
          alternativeName: "platform_b",
          values: {
            cost: { kind: "number", value: 15000, source: "imported" },
            reliability: { kind: "number", value: 86, source: "hybrid" },
            compliance: { kind: "number", value: 80, source: "hybrid" },
            scalability: { kind: "number", value: 90, source: "hybrid" }
          }
        },
        {
          alternativeName: "platform_c",
          values: {
            cost: { kind: "number", value: 22000, source: "imported" },
            reliability: { kind: "number", value: 96, source: "hybrid" },
            compliance: { kind: "number", value: 95, source: "hybrid" },
            scalability: { kind: "number", value: 85, source: "hybrid" }
          }
        }
      ]
    }
  ],

  aggregation: {
    method: "weighted_average",
    scenarioWeights: {
      startup: 0.4,
      unicorn: 0.35,
      established: 0.25
    },
    rationale:
      "Startup and unicorn phases are weighted more heavily because near-term growth is strategically more important than long-term enterprise maturity."
  }
};
```

