import type { UseCase } from './common.ts';

// Use cases for parsing a single source file (Go, Dart, TypeScript).
export const useCases: Record<string, UseCase> = {
  'cli.report.generate': {
    name: 'cli.report.generate',
    title: 'Generate design reports from configured notes and relationships',
    note: 'This is the primary end-to-end report generation use case.',
  },
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
