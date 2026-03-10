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
    expect(stdout).toBe(
      'status: ok\ncommand: validate\nmessage: validate stub ok\n',
    );
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
    expect(stdout).toBe(
      'status: ok\ncommand: report generate\nmessage: report generate stub ok\n',
    );
    expect(stderr).toBe('');
  });

  test('fails without config', async () => {
    const proc = Bun.spawn([binary, 'validate'], {
      stdout: 'pipe',
      stderr: 'pipe',
    });

    const stdout = await new Response(proc.stdout).text();
    const stderr = await new Response(proc.stderr).text();
    const exitCode = await proc.exited;

    expect(exitCode).toBe(1);
    expect(stdout).toBe('');
    expect(stderr).toBe('status: error\nmessage: config flag is required\n');
  });
});
