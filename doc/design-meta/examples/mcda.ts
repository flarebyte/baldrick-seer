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
      importanceWeight: 0.4,
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
      importanceWeight: 0.35,
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
      importanceWeight: 0.25,
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
