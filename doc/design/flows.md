# Execution Flows

Call-oriented CLI graphs. Normative semantics are defined in the specification.

## Report generation flows

Graph view for generating reports from a validated input config. Refer to the normative specification for scoring and failure behavior.

### Report generation graph

Text graph for the report-generation call chain, reusing the shared validation path.

- <a id="graph-node-call-reports-generate"></a> Generate Reports Call: Top-level report-generation call chain. See the normative specification for failure handling and output semantics.
  - <a id="graph-node-call-reports-generate-parse-args"></a> Parse Report Arguments: Parse CLI arguments for report generation, including the config path, requested report names, and output options.
    - <a id="graph-node-call-reports-generate-select-reports"></a> Select Requested Reports: Resolve which report definitions should run after CLI filtering.
      - <a id="graph-node-call-reports-generate-shared-validation"></a> Reuse Shared Validation Flow: Reuse the shared validation stage before scoring.
        - <a id="graph-node-call-reports-generate-build-ahp-inputs"></a> Build AHP Inputs: Prepare validated pairwise comparisons for AHP weight computation.
          - <a id="graph-node-call-reports-generate-compute-ahp-weights"></a> Compute Criteria Weights with AHP: Compute scenario-local criterion weights from validated pairwise comparisons.
            - <a id="graph-node-call-reports-generate-select-ranking-strategy"></a> Select Ranking Strategy: Select the ranking branch after AHP weighting.
              - <a id="graph-node-call-reports-generate-build-topsis-inputs"></a> Build TOPSIS Inputs: Assemble TOPSIS decision matrices from validated evaluations, polarity, and AHP-derived weights after filtering out alternatives excluded by scenario constraints.
                - <a id="graph-node-call-reports-generate-rank-alternatives-topsis"></a> Rank Alternatives with TOPSIS: Rank the alternatives that remain eligible after scenario-local constraint enforcement.
                  - <a id="graph-node-call-reports-generate-render-output"></a> Render Requested Reports: Render requested outputs after validation succeeds, constraints are enforced, and ranking results exist for eligible alternatives. Final aggregated rankings omit alternatives made ineligible by participating scenario constraints.
                    - <a id="graph-node-call-reports-generate-render-output-render-csv"></a> Render CSV Report: Render flat tabular CSV output for spreadsheet analysis and data exchange.
                    - <a id="graph-node-call-reports-generate-render-output-render-json"></a> Render JSON Report: Render JSON ranking output after successful validation, constraint enforcement, and scoring, while indicating scenario-level exclusions and omitting ineligible alternatives from the final aggregated ranking.
                    - <a id="graph-node-call-reports-generate-render-output-render-markdown"></a> Render Markdown Report: Render markdown rankings and scenario summaries, including explicit exclusion status when constraints remove an alternative and omitting ineligible alternatives from the final aggregated ranking.
              - <a id="graph-node-call-reports-generate-future-rank-electre"></a> Future Option: Rank with ELECTRE: Potential future branch for ELECTRE-based ranking.
              - <a id="graph-node-call-reports-generate-future-rank-topsis-sensitivity"></a> Future Option: TOPSIS with Sensitivity Analysis: Potential future branch that complements TOPSIS with sensitivity analysis.

## Validation flows

Graph view for validating an input config file. Refer to the normative specification for validation rules.

### Input config validation graph

Text graph for the validate-config call chain.

- <a id="graph-node-call-validation-input-config"></a> Validate Input Config Call: Top-level validate-command call chain. See the normative specification for authoritative command semantics.
  - <a id="graph-node-call-validation-input-config-parse-args"></a> Parse Validation Arguments: Parse CLI arguments for the validate command, including the config path and output flags.
    - <a id="graph-node-call-validation-input-config-load-cue-config"></a> Load CUE Config: Load and evaluate the CUE configuration package before validation.
      - <a id="graph-node-call-validation-input-config-validate-model"></a> Validate Config Model: Run the shared validation stage and emit diagnostics. See the normative specification for exact validation behavior.
        - <a id="graph-node-call-validation-input-config-validate-model-check-structure"></a> Check Config Structure: Check that the loaded config matches the expected top-level shape, required sections, and field types after CUE evaluation.
          - <a id="graph-node-call-validation-input-config-validate-model-check-references"></a> Check Named References: Check that all named references resolve, including criteria names, scenario names, alternative names, and report focus selectors.
            - <a id="graph-node-call-validation-input-config-validate-model-check-pairwise-comparisons"></a> Check Pairwise Comparisons: Validate AHP pairwise-comparison coverage and canonical comparison structure.
              - <a id="graph-node-call-validation-input-config-validate-model-check-evaluation-coverage"></a> Check Evaluation Coverage: Validate evaluation coverage and supported value forms for active criteria.
                - <a id="graph-node-call-validation-input-config-validate-model-check-constraints"></a> Check Scenario Constraints: Validate constraint operator and value compatibility, then apply constraints during scenario-local scoring to exclude violating alternatives before ranking.
                  - <a id="graph-node-call-validation-input-config-validate-model-check-report-definitions"></a> Check Report Definitions: Validate report definitions, focus selectors, and report arguments.

