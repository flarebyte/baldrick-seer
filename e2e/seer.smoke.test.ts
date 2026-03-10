import { describe, expect, test } from 'bun:test';

const binary = '.e2e-bin/seer';
const config = 'testdata/config/minimal.cue';

async function readGolden(name: string): Promise<string> {
  return Bun.file(`testdata/golden/${name}`).text();
}

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
    expect(stdout).toBe(await readGolden('validate_success.stdout.golden'));
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
      await readGolden('report_generate_success.stdout.golden'),
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
    expect(stderr).toBe(await readGolden('missing_config.stderr.golden'));
  });

  test('is deterministic across repeated runs', async () => {
    const proc1 = Bun.spawn([binary, 'validate', '--config', config], {
      stdout: 'pipe',
      stderr: 'pipe',
    });
    const proc2 = Bun.spawn([binary, 'validate', '--config', config], {
      stdout: 'pipe',
      stderr: 'pipe',
    });

    const stdout1 = await new Response(proc1.stdout).text();
    const stderr1 = await new Response(proc1.stderr).text();
    const exitCode1 = await proc1.exited;

    const stdout2 = await new Response(proc2.stdout).text();
    const stderr2 = await new Response(proc2.stderr).text();
    const exitCode2 = await proc2.exited;

    expect(exitCode1).toBe(0);
    expect(exitCode2).toBe(0);
    expect(stdout1).toBe(stdout2);
    expect(stderr1).toBe(stderr2);
  });
});
