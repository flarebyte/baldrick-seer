import { McdaModel } from "./model";

export const minimalScenarioMcda: McdaModel = {
  modelType: "scenario_based_mcda",
  version: "1.0",

  problem: {
    id: "hosting-choice",
    name: "Hosting Choice",
    goal: "Choose the best hosting provider across business scenarios"
  },

  criteriaCatalog: [
    { id: "cost", name: "Cost", polarity: "cost", unit: "USD/month", valueType: "number" },
    { id: "speed", name: "Speed", polarity: "benefit", unit: "score", valueType: "number" },
    { id: "compliance", name: "Compliance", polarity: "benefit", unit: "score", valueType: "number" }
  ],

  alternatives: [
    { id: "a", name: "Provider A" },
    { id: "b", name: "Provider B" }
  ],

  scenarios: [
    {
      id: "lean_startup",
      name: "Lean Startup",
      activeCriteria: [
        { criterionId: "cost" },
        { criterionId: "speed" }
      ],
      preferences: {
        method: "ahp_pairwise",
        scale: "saaty_1_9",
        comparisons: [
          {
            moreImportantCriterionId: "cost",
            lessImportantCriterionId: "speed",
            strength: 3,
            justification: "Budget matters more than peak performance."
          }
        ]
      },
      evaluations: [
        {
          alternativeId: "a",
          values: {
            cost: { kind: "number", value: 100 },
            speed: { kind: "number", value: 70 }
          }
        },
        {
          alternativeId: "b",
          values: {
            cost: { kind: "number", value: 180 },
            speed: { kind: "number", value: 90 }
          }
        }
      ]
    },
    {
      id: "regulated_growth",
      name: "Regulated Growth",
      activeCriteria: [
        { criterionId: "cost" },
        { criterionId: "speed" },
        { criterionId: "compliance" }
      ],
      preferences: {
        method: "ahp_pairwise",
        scale: "saaty_1_9",
        comparisons: [
          {
            moreImportantCriterionId: "compliance",
            lessImportantCriterionId: "cost",
            strength: 5
          },
          {
            moreImportantCriterionId: "compliance",
            lessImportantCriterionId: "speed",
            strength: 4
          },
          {
            moreImportantCriterionId: "speed",
            lessImportantCriterionId: "cost",
            strength: 2
          }
        ]
      },
      evaluations: [
        {
          alternativeId: "a",
          values: {
            cost: { kind: "number", value: 100 },
            speed: { kind: "number", value: 70 },
            compliance: { kind: "number", value: 60 }
          }
        },
        {
          alternativeId: "b",
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