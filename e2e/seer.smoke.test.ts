import { describe, expect, test } from 'bun:test';

const binary = '.e2e-bin/seer';
const config = 'testdata/config/minimal.cue';

async function readGolden(name: string): Promise<string> {
  return Bun.file(`testdata/golden/${name}`).text();
}

async function runSeer(args: string[]) {
  const proc = Bun.spawn([binary, ...args], {
    stdout: 'pipe',
    stderr: 'pipe',
  });

  const stdout = await new Response(proc.stdout).text();
  const stderr = await new Response(proc.stderr).text();
  const exitCode = await proc.exited;

  return { exitCode, stdout, stderr };
}

describe('seer CLI smoke test', () => {
  test('runs validate', async () => {
    const { exitCode, stdout, stderr } = await runSeer([
      'validate',
      '--config',
      config,
    ]);

    expect(exitCode).toBe(0);
    expect(stdout).toBe(await readGolden('validate_success.stdout.golden'));
    expect(stderr).toBe('');
  });

  test('runs report generate', async () => {
    const { exitCode, stdout, stderr } = await runSeer([
      'report',
      'generate',
      '--config',
      config,
    ]);

    expect(exitCode).toBe(0);
    expect(stdout).toBe(
      await readGolden('report_generate_success.stdout.golden'),
    );
    expect(stderr).toBe('');
  });

  test('fails without config', async () => {
    const { exitCode, stdout, stderr } = await runSeer(['validate']);

    expect(exitCode).toBe(1);
    expect(stdout).toBe('');
    expect(stderr).toBe(await readGolden('missing_config.stderr.golden'));
  });

  test('is deterministic across repeated runs', async () => {
    const first = await runSeer(['validate', '--config', config]);
    const second = await runSeer(['validate', '--config', config]);

    expect(first.exitCode).toBe(0);
    expect(second.exitCode).toBe(0);
    expect(first.stdout).toBe(second.stdout);
    expect(first.stderr).toBe(second.stderr);
  });
});
