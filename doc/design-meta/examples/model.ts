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
  valueType?: "number" | "ordinal" | "boolean" | "text";
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
  | OrdinalCriterionValue
  | TextCriterionValue;

export interface NumericCriterionValue {
  kind: "number";
  value: number;
  estimated?: boolean;
  source?: "human" | "ai" | "hybrid" | "measured" | "imported";
  justification?: string;
}

export interface BooleanCriterionValue {
  kind: "boolean";
  value: boolean;
  source?: "human" | "ai" | "hybrid" | "measured" | "imported";
  justification?: string;
}

export interface OrdinalCriterionValue {
  kind: "ordinal";
  value: number;
  label?: string;
  source?: "human" | "ai" | "hybrid" | "measured" | "imported";
  justification?: string;
}

export interface TextCriterionValue {
  kind: "text";
  value: string;
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
  value: number | boolean | string;
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
