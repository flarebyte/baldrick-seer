import type { ImplementationConsideration } from "./common.ts";

// Initial implementation suggestions. Keep this list small and actionable.
export const implementations: Record<string, ImplementationConsideration> = {
  "model.structure": {
    name: "model.structure",
    title: "Decision Model Structure",
    description:
      "Define a clear structure for representing the decision problem, including criteria, alternatives, and scenarios. The model should remain understandable to humans and easily generated or modified by AI agents.",
  },

  "scenario.isolation": {
    name: "scenario.isolation",
    title: "Scenario Isolation",
    description:
      "Each scenario should be evaluated independently with its own criteria priorities and candidate evaluations. Avoid mixing scenario data in a single decision matrix to preserve conceptual clarity.",
  },

  "input.format": {
    name: "input.format",
    title: "Human and AI Friendly Input Format",
    description:
      "Use a semantic data format such as JSON or YAML that is easy for humans and AI systems to read and generate. Avoid positional or matrix-heavy formats that are difficult to validate or explain.",
  },

  "criteria.pairwise.clarity": {
    name: "criteria.pairwise.clarity",
    title: "Clear Representation of Pairwise Judgments",
    description:
      "Represent pairwise comparisons explicitly with named criteria rather than positional matrices. This makes the model easier to interpret, validate, and generate using AI assistance.",
  },

  "model.validation": {
    name: "model.validation",
    title: "Model Validation",
    description:
      "Validate the input model before computation. This includes checking that all criteria referenced in scenarios exist, that comparisons are complete, and that alternative evaluations include the required values.",
  },

  "model.incomplete.data": {
    name: "model.incomplete.data",
    title: "Handling Incomplete Information",
    description:
      "The system should detect missing comparisons or evaluation values and provide helpful feedback. In some cases, it may suggest or infer missing values, but such assumptions should be clearly marked.",
  },

  "scenario.aggregation.policy": {
    name: "scenario.aggregation.policy",
    title: "Scenario Aggregation Strategy",
    description:
      "Define how results from multiple scenarios are combined into a final decision. Possible strategies include equal averaging, weighted scenarios, or robustness-focused approaches.",
  },

  "decision.explainability": {
    name: "decision.explainability",
    title: "Explainable Results",
    description:
      "Provide explanations for the ranking results, including the influence of criteria and scenario differences. This helps users understand why a particular alternative is recommended.",
  },

  "decision.traceability": {
    name: "decision.traceability",
    title: "Traceable Decision Process",
    description:
      "Ensure the system can show the reasoning path from inputs to outputs. This includes displaying scenario weights, criteria importance, and the contribution of each factor.",
  },

  "analysis.robustness": {
    name: "analysis.robustness",
    title: "Sensitivity and Robustness Analysis",
    description:
      "Support the ability to test how changes in criteria importance or scenario assumptions affect the final ranking. This helps determine whether the decision is stable or fragile.",
  },

  "scenario.constraints": {
    name: "scenario.constraints",
    title: "Constraint Enforcement",
    description:
      "Allow scenarios to define minimum requirements or constraints for alternatives, such as mandatory compliance levels. Alternatives that violate these constraints may be excluded before ranking.",
  },

  "ux.model.guidance": {
    name: "ux.model.guidance",
    title: "Guidance for Model Creation",
    description:
      "Provide help or prompts that guide users in defining criteria, comparisons, and scenario descriptions. This reduces modeling errors and improves the quality of the decision model.",
  },

  "cli.output.readability": {
    name: "cli.output.readability",
    title: "Readable CLI Output",
    description:
      "Results should be displayed in a clear and structured format suitable for terminal environments. Tables, summaries, and scenario breakdowns help users quickly understand the outcome.",
  },

  "cli.output.machine": {
    name: "cli.output.machine",
    title: "Structured Output for Automation",
    description:
      "In addition to human-readable output, the CLI should support machine-readable formats such as JSON so results can be consumed by other systems or AI agents.",
  },

  "criteria.semantic.consistency": {
    name: "criteria.semantic.consistency",
    title: "Consistent Criteria Interpretation",
    description:
      "Ensure that criteria maintain consistent meaning across scenarios. Differences in importance are acceptable, but the interpretation of a criterion should remain stable.",
  },

  "system.extensibility.methods": {
    name: "system.extensibility.methods",
    title: "Extensible Decision Methods",
    description:
      "Design the system so additional MCDA methods can be added later without redesigning the data model or CLI interface.",
  },

  "model.documentation": {
    name: "model.documentation",
    title: "Model Documentation",
    description:
      "Allow decision models to include descriptions, notes, and justifications for comparisons and values. This improves transparency and makes the model easier to review.",
  },

  "execution.reproducibility": {
    name: "execution.reproducibility",
    title: "Reproducible Decision Runs",
    description:
      "Ensure that running the same model with the same inputs always produces identical results. This is important for auditing and comparing decisions over time.",
  },
};


export type AlgorithmReference = {
  name: string;
  title: string;
  description: string;
  typicalUse?: string;
};

export const algorithmReferences: Record<string, AlgorithmReference> = {
  "mcda.ahp": {
    name: "mcda.ahp",
    title: "Analytic Hierarchy Process (AHP)",
    description:
      "A structured decision-making method that derives criteria weights from pairwise comparisons. It converts subjective judgments about relative importance into a consistent numerical weighting system.",
    typicalUse:
      "Used to compute criteria importance before ranking alternatives with another MCDA method."
  },

  "mcda.topsis": {
    name: "mcda.topsis",
    title: "Technique for Order Preference by Similarity to Ideal Solution (TOPSIS)",
    description:
      "Ranks alternatives based on their distance from an ideal best solution and a worst solution. The preferred alternative is the one closest to the ideal and farthest from the negative ideal.",
    typicalUse:
      "Commonly used to rank alternatives once criteria weights are known."
  },

  "mcda.electre": {
    name: "mcda.electre",
    title: "ELECTRE Outranking Method",
    description:
      "An outranking-based MCDA approach that determines whether one alternative sufficiently dominates another using concordance and discordance measures.",
    typicalUse:
      "Useful when decisions involve strong conflicts between criteria or when certain criteria should act as veto conditions."
  },

  "mcda.promethee": {
    name: "mcda.promethee",
    title: "PROMETHEE (Preference Ranking Organization Method for Enrichment Evaluation)",
    description:
      "An outranking MCDA method that compares alternatives pairwise using preference functions and produces partial or complete rankings.",
    typicalUse:
      "Useful when a transparent ranking of alternatives is needed with clear preference modeling."
  },

  "mcda.vikor": {
    name: "mcda.vikor",
    title: "VIKOR Compromise Ranking Method",
    description:
      "A multi-criteria ranking method focused on identifying a compromise solution that balances group utility and individual regret among alternatives.",
    typicalUse:
      "Often used when decision-makers seek a compromise solution among conflicting criteria."
  },

  "analysis.sensitivity": {
    name: "analysis.sensitivity",
    title: "Sensitivity Analysis",
    description:
      "A technique that evaluates how changes in criteria weights or input values affect the ranking of alternatives.",
    typicalUse:
      "Used to test the robustness of MCDA results and identify parameters that significantly influence the decision."
  },

  "analysis.robustness": {
    name: "analysis.robustness",
    title: "Robustness Analysis",
    description:
      "Evaluates how stable a decision remains across variations in assumptions, scenarios, or parameter ranges.",
    typicalUse:
      "Useful when decisions must remain valid under uncertain future conditions."
  }
};