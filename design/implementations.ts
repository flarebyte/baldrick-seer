import type { ImplementationConsideration } from './common.ts';

// Initial implementation suggestions. Keep this list small and actionable.
export const implementations: Record<string, ImplementationConsideration> = {
  'lang.go': {
    name: 'lang.go',
    title: 'Implement the CLI in Go',
    description:
      'Use Go as the primary implementation language for strong typing, fast startup, and straightforward single-binary distribution.',
    calls: ['cli.root'],
  },
};
