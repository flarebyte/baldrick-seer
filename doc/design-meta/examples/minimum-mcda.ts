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
