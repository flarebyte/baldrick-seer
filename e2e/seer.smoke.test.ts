import { describe, expect, test } from 'bun:test';

const binary = '.e2e-bin/seer';
const config = 'testdata/config/minimal.cue';

describe('seer CLI smoke test', () => {
  test('runs validate', async () => {
    const proc = Bun.spawn([binary, 'validate', '--config', config], {
      stdout: 'pipe',
      stderr: 'pipe',
    });

    const stdout = await new Response(proc.stdout).text();
    const stderr = await new Response(proc.stderr).text();
    const exitCode = await proc.exited;

    expect(exitCode).toBe(0);
    expect(stdout).toBe('validate: ok\n');
    expect(stderr).toBe('');
  });

  test('runs report generate', async () => {
    const proc = Bun.spawn([binary, 'report', 'generate', '--config', config], {
      stdout: 'pipe',
      stderr: 'pipe',
    });

    const stdout = await new Response(proc.stdout).text();
    const stderr = await new Response(proc.stderr).text();
    const exitCode = await proc.exited;

    expect(exitCode).toBe(0);
    expect(stdout).toBe('report generate: ok\n');
    expect(stderr).toBe('');
  });
});
