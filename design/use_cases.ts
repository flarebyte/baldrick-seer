import type { UseCase } from './common.ts';

// Use cases for parsing a single source file (Go, Dart, TypeScript).
export const useCases: Record<string, UseCase> = {
  "technology.platform-selection": {
    name: "technology.platform-selection",
    title: "Technology Platform Selection",
    note:
      "Compare multiple technology platforms across different operational scenarios such as startup, scale-up, and enterprise maturity."
  },

  "technology.architecture-choice": {
    name: "technology.architecture-choice",
    title: "Software Architecture Decision",
    note:
      "Evaluate architectural approaches under different system growth conditions, performance requirements, and reliability expectations."
  },

  "technology.infrastructure-strategy": {
    name: "technology.infrastructure-strategy",
    title: "Infrastructure Strategy Planning",
    note:
      "Assess infrastructure alternatives such as cloud providers or deployment models under varying operational scenarios."
  },

  "vendor.supplier-selection": {
    name: "vendor.supplier-selection",
    title: "Supplier or Vendor Selection",
    note:
      "Evaluate competing suppliers using multiple criteria such as cost, reliability, and service quality under different operating conditions."
  },

  "vendor.service-provider-comparison": {
    name: "vendor.service-provider-comparison",
    title: "Service Provider Comparison",
    note:
      "Compare service providers where priorities may change depending on scale, regulatory environment, or organizational maturity."
  },

  "product.feature-prioritization": {
    name: "product.feature-prioritization",
    title: "Product Feature Prioritization",
    note:
      "Rank product features using multiple criteria such as user value, development effort, and strategic importance."
  },

  "product.roadmap-planning": {
    name: "product.roadmap-planning",
    title: "Product Roadmap Planning",
    note:
      "Evaluate product initiatives across different market or growth scenarios to support long-term planning."
  },

  "strategy.growth-scenario-evaluation": {
    name: "strategy.growth-scenario-evaluation",
    title: "Growth Scenario Evaluation",
    note:
      "Assess strategic options under different business growth trajectories such as startup, rapid expansion, or mature operations."
  },

  "strategy.investment-decision": {
    name: "strategy.investment-decision",
    title: "Strategic Investment Decision",
    note:
      "Compare investment alternatives considering financial return, risk exposure, and long-term strategic impact."
  },

  "planning.long-term-option-evaluation": {
    name: "planning.long-term-option-evaluation",
    title: "Long-Term Option Evaluation",
    note:
      "Evaluate alternatives that must perform well across multiple possible future environments."
  },

  "planning.lifecycle-decision": {
    name: "planning.lifecycle-decision",
    title: "Lifecycle Decision Support",
    note:
      "Compare options that must remain effective throughout different stages of organizational or system development."
  },

  "infrastructure.system-design-selection": {
    name: "infrastructure.system-design-selection",
    title: "System Design Selection",
    note:
      "Compare different system designs where trade-offs exist between cost, scalability, and reliability."
  },

  "policy.policy-option-analysis": {
    name: "policy.policy-option-analysis",
    title: "Policy Option Analysis",
    note:
      "Support evaluation of policy alternatives where multiple criteria such as impact, feasibility, and cost must be considered."
  },

  "decision.multi-criteria-ranking": {
    name: "decision.multi-criteria-ranking",
    title: "General Multi-Criteria Ranking",
    note:
      "Rank competing alternatives using multiple evaluation criteria and scenario-based priorities."
  },

  "decision.robust-choice-identification": {
    name: "decision.robust-choice-identification",
    title: "Robust Choice Identification",
    note:
      "Identify alternatives that perform consistently well across different scenarios or changing assumptions."
  }
};

export const getByName = (expectedName: string) =>
  Object.values(useCases).find(({ name }) => name === expectedName);

export const mustUseCases = new Set([
  ...Object.values(useCases).map(({ name }) => name),
]);

export const useCaseCatalogByName: Record<
  string,
  { name: string; title: string; note?: string }
> = Object.fromEntries(Object.values(useCases).map((u) => [u.name, u]));
