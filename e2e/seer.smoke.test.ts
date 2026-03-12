import { describe, expect, test } from 'bun:test';

import { readGolden, runSeer } from './helpers';

const minimalConfig = 'testdata/config/minimal.cue';
const splitConfig = 'testdata/config_split';
const validReportConfig = 'testdata/config/valid_report.cue';
const invalidReferenceConfig = 'testdata/config/invalid_reference.cue';
const malformedConfig = 'testdata/config/malformed.cue';
const emptyConfigDir = 'testdata/config_empty';

describe('seer CLI E2E', () => {
  describe('validate', () => {
    test('runs with a single cue file', async () => {
      const { exitCode, stdout, stderr } = await runSeer([
        'validate',
        '--config',
        minimalConfig,
      ]);

      expect(exitCode).toBe(0);
      expect(stdout).toBe(await readGolden('validate_success.stdout.golden'));
      expect(stderr).toBe('');
    });

    test('runs with a cue directory package', async () => {
      const { exitCode, stdout, stderr } = await runSeer([
        'validate',
        '--config',
        splitConfig,
      ]);

      expect(exitCode).toBe(0);
      expect(stdout).toBe(await readGolden('validate_success.stdout.golden'));
      expect(stderr).toBe('');
    });

    test('fails without config', async () => {
      const { exitCode, stdout, stderr } = await runSeer(['validate']);

      expect(exitCode).toBe(1);
      expect(stdout).toBe('');
      expect(stderr).toBe(await readGolden('missing_config.stderr.golden'));
    });

    test('fails with semantic validation errors', async () => {
      const { exitCode, stdout, stderr } = await runSeer([
        'validate',
        '--config',
        invalidReferenceConfig,
      ]);

      expect(exitCode).toBe(1);
      expect(stdout).toBe('');
      expect(stderr).toBe(await readGolden('invalid_validation.stderr.golden'));
    });
  });

  describe('report generate', () => {
    test('renders the default minimal report', async () => {
      const { exitCode, stdout, stderr } = await runSeer([
        'report',
        'generate',
        '--config',
        minimalConfig,
      ]);

      expect(exitCode).toBe(0);
      expect(stdout).toBe(
        await readGolden('report_generate_success.stdout.golden'),
      );
      expect(stderr).toBe('');
    });

    test('renders the multi-report fixture', async () => {
      const { exitCode, stdout, stderr } = await runSeer([
        'report',
        'generate',
        '--config',
        validReportConfig,
      ]);

      expect(exitCode).toBe(0);
      expect(stdout).toBe(
        await readGolden('report_generate_valid_report.stdout.golden'),
      );
      expect(stderr).toBe('');
    });

    test('fails for malformed cue input', async () => {
      const { exitCode, stdout, stderr } = await runSeer([
        'report',
        'generate',
        '--config',
        malformedConfig,
      ]);

      expect(exitCode).toBe(1);
      expect(stdout).toBe('');
      expect(stderr).toBe(await readGolden('malformed_cue.stderr.golden'));
    });

    test('fails for empty cue directories', async () => {
      const { exitCode, stdout, stderr } = await runSeer([
        'report',
        'generate',
        '--config',
        emptyConfigDir,
      ]);

      expect(exitCode).toBe(1);
      expect(stdout).toBe('');
      expect(stderr).toBe(await readGolden('directory_path.stderr.golden'));
    });
  });

  describe('determinism', () => {
    test('keeps validate success stable across repeated runs', async () => {
      const first = await runSeer(['validate', '--config', minimalConfig]);
      const second = await runSeer(['validate', '--config', minimalConfig]);

      expect(first).toEqual(second);
    });

    test('keeps report generation stable across repeated runs', async () => {
      const first = await runSeer([
        'report',
        'generate',
        '--config',
        validReportConfig,
      ]);
      const second = await runSeer([
        'report',
        'generate',
        '--config',
        validReportConfig,
      ]);

      expect(first).toEqual(second);
    });

    test('keeps validation failures stable across repeated runs', async () => {
      const first = await runSeer([
        'validate',
        '--config',
        invalidReferenceConfig,
      ]);
      const second = await runSeer([
        'validate',
        '--config',
        invalidReferenceConfig,
      ]);

      expect(first).toEqual(second);
    });
  });
});
