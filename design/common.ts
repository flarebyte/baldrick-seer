import { appendFile, mkdir, writeFile } from 'node:fs/promises';
import { stringify as yamlStringify } from 'bun';

/**
 * Canonical description of a capability or behavior we want to support.
 * - name: short machine-friendly identifier
 * - title: concise human-friendly description
 * - note: optional extra context or constraint
 */
export type UseCase = {
  name: string;
  title: string;
  note?: string;
};

/**
 * Canonical description of a risk in the design.
 * - name: short machine-friendly identifier
 * - title: concise human-friendly description
 * - description: what could go wrong and why
 * - mitigation: actions or controls to reduce impact/likelihood
 * - calls: names of flow calls this risk applies to
 */
export type Risk = {
  name: string;
  title: string;
  description: string;
  mitigation: string;
  calls?: string[];
};

/**
 * Canonical implementation suggestion for the current design.
 * - name: short machine-friendly identifier
 * - title: concise recommendation statement
 * - description: why this choice is suggested
 * - calls: names of flow calls this suggestion applies to
 */
export type ImplementationConsideration = {
  name: string;
  title: string;
  description: string;
  calls?: string[];
};

/**
 * A single step in the flow-tree that the CLI will execute.
 * - name/title: brief identifiers for the step
 * - note: optional clarifying detail
 * - directory: logical Go package or path where code will live
 * - level: depth for tree rendering
 * - useCases: names of use-cases this step satisfies
 */
export type ComponentCall = {
  name: string;
  title: string;
  note?: string;
  // When true, render full details only on first appearance in FLOW_DESIGN,
  // later appearances are rendered as a lightweight reference line.
  displayOnce?: boolean;
  directory?: string;
  level: number;
  useCases?: string[];
  signature?: {
    input?: string;
    success?: string;
    failure?: string;
  };
  suggest?: {
    pkg?: string;
    func?: string;
    file?: string;
  };
};

/**
 * One stickie per idea by default
 */
export type Stickie = {
  //format for name: blackboard-export-format
  name: string;
  // Keep each note concise and focused on what is not obvious, with emphasis
  // on the “why” (intent, trade-offs, constraints, and rationale). Prefer neutral,
  // implementation-agnostic phrasing.
  note?: string;
  //code in programming language
  code?: string;
  // Use labels from this controlled set when
  // relevant: usecase, example, flow, design, implementation, decision, security, operations,
  // compliance, glossary, principle, validation, howto, faq, library, wip, stable.
  labels: string[];
};

export type ArgumentScope = 'h3-section' | 'note' | 'renderer';

export type ArgumentValueType =
  | 'string'
  | 'string[]'
  | 'boolean'
  | 'int'
  | 'float'
  | 'enum';

export type ArgumentDefinition = {
  // Machine-friendly key used in free-form arguments, e.g. "format-csv"
  name: string;
  title: string;
  description: string;
  valueType: ArgumentValueType;
  required?: boolean;
  defaultValue?: string | string[] | boolean | number;
  allowedValues?: string[];
  // Where this argument is valid
  scopes: ArgumentScope[];
  // Optional renderer compatibility hints
  renderers?: string[];
};

export type ArgumentRegistry = {
  version: string;
  arguments: ArgumentDefinition[];
};

export type OrderingStrategy =
  | 'by-name'
  | 'by-title'
  | 'by-label-then-name'
  | 'source-order';

export type OrderingPolicy = {
  // Comparator contract for `by-label-then-name`:
  // primaryLabel = lexicographically smallest label in `labels` (or empty string if no labels),
  // then compare by note `name`.
  notes: OrderingStrategy;
  // Comparator contract for `by-name` on relationships:
  // compare by (`from`, `to`, labelsSortedJoined), where labelsSortedJoined is
  // relationship labels sorted lexicographically and joined with a stable separator (`|`).
  relationships: OrderingStrategy;
  // Comparator contract for `by-title`:
  // compare lowercase `title`, then original index as stable tie-breaker.
  sections: OrderingStrategy;
  // Comparator contract for `by-name`:
  // compare by argument key/name.
  arguments: OrderingStrategy;
};

export type DiagnosticSeverity = 'error' | 'warning' | 'info';

export type Diagnostic = {
  code: string;
  message: string;
  severity: DiagnosticSeverity;
  // Logical component/call that produced this diagnostic
  source: string;
  // Canonical machine-readable location inside user config, e.g.
  // "reports[0].sections[0].sections[0].notes[2].arguments[1]"
  location?: string;
  // Optional human-readable context (does not replace canonical `location`)
  reportTitle?: string;
  sectionTitle?: string;
  // Optional entity linkage
  noteName?: string;
  relationship?: {
    from: string;
    to: string;
  };
  argumentName?: string;
  labelValue?: string;
  subjectLabel?: string;
  edgeLabel?: string;
  counterpartLabel?: string;
};

export type GraphIntegrityPolicy = {
  // Behavior when a relationship references a note that does not exist
  missingNode: 'error' | 'warning' | 'ignore';
  // Behavior when notes are disconnected from all report roots/sections
  orphanNode: 'error' | 'warning' | 'ignore';
  // Behavior when duplicate note names are detected
  duplicateNoteName: 'error' | 'warning' | 'ignore';
  // Behavior when referenced labels (for example label-filter arguments) are unknown
  unknownRelationshipLabel: 'error' | 'warning' | 'ignore';
  // Whether cross-report references are allowed
  crossReportReference: 'allow' | 'disallow';
};

export type ValidatedNote = {
  name: string;
  title: string;
  description?: string;
  labels: string[];
  arguments?: string[];
  location?: string;
};

export type ValidatedRelationship = {
  from: string;
  to: string;
  labels: string[];
  location?: string;
};

export type ValidatedH3Section = {
  title: string;
  description: string;
  arguments?: string[];
};

export type ValidatedH2Section = {
  title: string;
  description: string;
  sections: ValidatedH3Section[];
};

export type ValidatedReport = {
  title: string;
  filepath: string;
  sections: ValidatedH2Section[];
};

export type ValidatedApp = {
  notes: ValidatedNote[];
  relationships: ValidatedRelationship[];
  reports: ValidatedReport[];
  // Ordering policy can remain generation-time if desired.
  orderingPolicy?: OrderingPolicy;
  graphIntegrityPolicy: GraphIntegrityPolicy;
  argumentRegistry: ArgumentRegistry;
  diagnostics: Diagnostic[];
};

/**
 * Carry indentation depth while walking the design tree.
 */
export type FlowContext = {
  level: number;
};

export type Field = {
  name: string;
  title: string;
  description: string;
  //Typescript kind: ex: string
  tsKind: string;
  //Go Language kind
  goKind: string;
  //Dart Language kind
  dartKind: string;
  category?: string;
  required: boolean;
};

export type Entity = {
  name: string;
  title: string;
  description: string;
  fields: Field[];
};

/**
 * Increase nesting level for child calls.
 */
export const incrContext = (flowContext: FlowContext) => ({
  level: flowContext.level + 1,
});

/**
 * Extract all referenced use-case names from a call list.
 */
export const toUseCaseSet = (calls: ComponentCall[]) => {
  const allUseCases = calls
    .flatMap(({ useCases }) => useCases)
    .filter((useCase) => typeof useCase === 'string');
  return new Set(allUseCases);
};

/**
 * Reset the generated design report.
 */
export const resetReport = async () => {
  await writeFile('doc/FLOW_DESIGN.md', '');
};

/**
 * Append a single line to the generated design report.
 */
export const appendToReport = async (line: string) => {
  await appendFile('doc/FLOW_DESIGN.md', `${line}\n`, 'utf8');
};

/**
 * Render the flow-tree as indented lines of titles.
 */
export const displayCallsAsText = async (calls: ComponentCall[]) => {
  const seen = new Set<string>();
  for (let i = 0; i < calls.length; i += 1) {
    const call = calls[i];
    const spaces = ' '.repeat(call.level * 2);

    if (call.displayOnce && seen.has(call.name)) {
      await appendToReport(`${spaces}${call.title} [${call.name}] (ref)`);
      const startLevel = call.level;
      while (i + 1 < calls.length && calls[i + 1].level > startLevel) {
        i += 1;
      }
      continue;
    }

    if (call.displayOnce) {
      seen.add(call.name);
    }

    await appendToReport(`${spaces}${call.title}`);
  }
};

/**
 * Render a detailed view of the call tree with notes and suggestions.
 */
export const displayCallsDetailed = async (calls: ComponentCall[]) => {
  const seen = new Set<string>();
  for (let i = 0; i < calls.length; i += 1) {
    const call = calls[i];
    const base = ' '.repeat(call.level * 2);

    if (call.displayOnce && seen.has(call.name)) {
      await appendToReport(`${base}${call.title} [${call.name}]`);
      await appendToReport(
        `${base}  - ref: see first occurrence above for full subtree`,
      );
      const startLevel = call.level;
      while (i + 1 < calls.length && calls[i + 1].level > startLevel) {
        i += 1;
      }
      continue;
    }

    if (call.displayOnce) {
      seen.add(call.name);
    }

    await appendToReport(`${base}${call.title} [${call.name}]`);
    if (call.note) {
      await appendToReport(`${base}  - note: ${call.note}`);
    }
    const pkg = call.directory || call.suggest?.pkg;
    if (pkg) {
      await appendToReport(`${base}  - pkg: ${pkg}`);
    }
    if (call.suggest?.func) {
      await appendToReport(`${base}  - func: ${call.suggest.func}`);
    }
    if (call.signature?.input) {
      await appendToReport(`${base}  - input: ${call.signature.input}`);
    }
    if (call.signature?.success) {
      await appendToReport(`${base}  - success: ${call.signature.success}`);
    }
    if (call.signature?.failure) {
      await appendToReport(`${base}  - failure: ${call.signature.failure}`);
    }
    if (call.suggest?.file) {
      await appendToReport(`${base}  - file: ${call.suggest.file}`);
    }
  }
};

/**
 * Pure set difference: items in A not in B.
 */
export const getSetDifference = (
  setA: Set<string>,
  setB: Set<string>,
): Set<string> => {
  return new Set([...setA].filter((item) => !setB.has(item)));
};

/**
 * Render bullet points suitable for markdown.
 */
export const toBulletPoints = (lines: string[]) =>
  lines.map((line) => `  - ${line}`).join('\n');

/**
 * Convert a free-form title/name into a dash-lower slug suitable for stickie names and filenames.
 * Rules: lowercase, replace non [a-z0-9] with '-', merge repeated separators, trim edges.
 */
export const toStickieName = (s: string): string =>
  s
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '');

/**
 * Create a stickie YAML file for a section under `notes/`.
 * Non-fatal: any error while writing the stickie is swallowed so docs generation continues.
 */
export const writeSectionStickie = async (
  title: string,
  lines: string[] | string,
  outDir = 'notes',
) => {
  try {
    await mkdir(outDir, { recursive: true });
    const name = toStickieName(title);
    const note = Array.isArray(lines) ? toBulletPoints(lines) : lines;
    const stickie: Stickie = {
      name,
      note,
      labels: ['design'],
    };
    const yaml = yamlStringify(stickie, { indent: 2 });
    await writeFile(`${outDir}/${name}.stickie.yaml`, yaml, 'utf8');
  } catch (_) {
    // Best-effort: ignore errors (e.g., sandboxed environments)
  }
};

/**
 * Convenience to add a titled section to the report.
 */
export const appendSection = async (
  title: string,
  lines: string[] | string,
) => {
  await appendToReport(`\n## ${title}`);
  if (Array.isArray(lines)) {
    await appendToReport(toBulletPoints(lines));
  } else {
    await appendToReport(lines);
  }
  // Also emit a stickie for this section (best-effort)
  await writeSectionStickie(title, lines);
};

/**
 * Build human-readable lines for a set of use-case names using a catalog.
 */
export const toUseCaseLines = (
  useCaseNames: Set<string>,
  catalogByName: Record<string, UseCase>,
): string[] => {
  return [...useCaseNames].map((name) => {
    const uc = catalogByName[name];
    if (!uc) return name;
    return uc.note ? `${uc.title} — ${uc.note}` : uc.title;
  });
};

/**
 * Append a small heading then bullet list for use-cases (title + note).
 */
export const appendUseCases = async (
  heading: string,
  useCaseNames: Set<string>,
  catalogByName: Record<string, UseCase>,
) => {
  await appendToReport(`${heading}\n`);
  await appendToReport(
    toBulletPoints(toUseCaseLines(useCaseNames, catalogByName)),
  );
  await appendToReport('\n');
};

/**
 * Append a key/value list as a section, preserving entry order.
 */
export const appendKeyValueList = async (
  title: string,
  entries: Array<[string, string | string[]]>,
) => {
  await appendToReport(`\n## ${title}`);
  const lines = entries.map(([k, v]) => {
    const value = Array.isArray(v) ? v.join(', ') : v;
    return `${k}: ${value}`;
  });
  await appendToReport(toBulletPoints(lines));
};

/**
 * Baseline risk catalogue for the project.
 * Add or refine entries as the design evolves.
 */
// Risk records live in docs/designs/risks.mts

export const toDashCase = (str: string): string => {
  return str
    .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
    .replace(/[\s_]+/g, '-')
    .replace(/[^a-zA-Z0-9-]/g, '')
    .replace(/^-+|-+$/g, '')
    .toLowerCase();
};

export const toCamelCase = (str: string): string => {
  return str
    .toLowerCase()
    .replace(/[^a-zA-Z0-9]+(.)/g, (_match, char) => char.toUpperCase());
};
