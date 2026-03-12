export const binary = '.e2e-bin/seer';

export async function readGolden(name: string): Promise<string> {
  return Bun.file(`testdata/golden/${name}`).text();
}

export async function runSeer(args: string[]) {
  const proc = Bun.spawn([binary, ...args], {
    stdout: 'pipe',
    stderr: 'pipe',
  });

  const stdout = await new Response(proc.stdout).text();
  const stderr = await new Response(proc.stderr).text();
  const exitCode = await proc.exited;

  return { exitCode, stdout, stderr };
}
