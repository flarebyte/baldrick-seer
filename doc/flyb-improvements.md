# flyb Improvements

Short list of backward-compatible improvements that would make `flyb` easier to use every day for both humans and AI agents.

## 1. Load full CUE packages, not just one file

If `--config` points to `app.cue`, `flyb` should resolve the full CUE package in that directory, including sibling files and standard CUE imports.

Why it matters:
- makes large configs maintainable
- lets users split concerns without introducing generators
- matches normal CUE expectations

Backward-compatible path:
- keep `--config path/to/app.cue`
- when the file is CUE, load the containing package instead of the single file only

## 2. Add explicit support for config directories

Allow `--config doc/design-meta` in addition to `--config doc/design-meta/app.cue`.

Why it matters:
- simpler mental model
- clearer for humans and AI
- makes package-based layouts obvious

Backward-compatible path:
- keep current file behavior
- add directory support as an extra accepted form

## 3. Make markdown generation fully deterministic

Generated markdown should always replace the target file cleanly and deterministically.

Why it matters:
- avoids confusion when files appear partially stale
- improves trust in generated docs
- makes AI edit loops and CI comparisons much safer

Backward-compatible path:
- keep current command
- guarantee full overwrite, stable ordering, and stable section ordering

## 4. Add report targeting

Support generating or validating a subset of reports.

Examples:
- `flyb generate markdown --config ... --report glossary`
- `flyb validate --config ... --report flows`

Why it matters:
- faster feedback during editing
- easier for AI agents to iterate on one document
- less noise during everyday work

Backward-compatible path:
- default remains “all reports”
- add optional report filters

## 5. Improve graph controls for documentation flows

The graph renderer should expose a few more controls without changing the core model.

Useful additions:
- stable child ordering
- optional depth limit
- optional include/exclude labels
- explicit branch priority
- optional hidden helper nodes

Why it matters:
- large flow docs stay readable
- future-only branches can be shown without drowning the main path

## 6. Treat report-linked notes differently from true graph orphans

Today a note can appear in reports and still be warned as an orphan if it has no relationships.

Why it matters:
- filepath-backed source notes are useful even without graph edges
- current warning is noisy for documentation-oriented notes

Backward-compatible path:
- keep orphan detection
- add a softer class such as `DOC_ONLY_NOTE`
- or suppress orphan warnings for notes that are referenced by reports

## 7. Add clearer machine-friendly diagnostics

Diagnostics are already structured, but they could be easier to act on automatically.

Useful additions:
- suggested fixes
- note/report names in dedicated fields
- absolute/normalized paths
- optional related nodes

Why it matters:
- AI agents can repair configs more reliably
- humans get faster feedback without reading the full model

## 8. Show note labels in generated markdown when requested

Allow an optional rendering mode that prints note labels such as `v1`, `v2`, `method`, `implementation`, or `future`.

Why it matters:
- helps generated docs carry release scope and intent directly
- avoids encoding too much into titles

Backward-compatible path:
- default stays unchanged
- add an argument such as `show-labels=true`

## 9. Add a built-in flatten/export command

Support exporting a resolved CUE config as a single normalized CUE file.

Why it matters:
- helps teams keep modular sources while shipping one flattened file
- useful for debugging and review

Backward-compatible path:
- add a new command, for example `flyb export cue`

## 10. Preserve the current design philosophy

The main goal should be to improve packaging, determinism, and iteration speed without changing the core `flyb` idea:
- simple declarative notes
- relationships for graphs
- report-oriented markdown generation
- lightweight configuration model
