import { promises as fs } from 'node:fs';

export async function readFileSafe(path: string): Promise<string> {
  try {
    return await fs.readFile(path, 'utf8');
  } catch {
    return '';
  }
}

export async function readVersionFromProjectYAML(
  path = 'main.project.yaml',
): Promise<string> {
  const raw = await readFileSafe(path);
  if (!raw) return '';
  const lines = raw.split(/\r?\n/);
  let inTags = false;
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const trimmed = line.trim();
    if (!inTags) {
      if (/^tags\s*:\s*$/.test(trimmed)) inTags = true;
      continue;
    }
    if (/^\S/.test(line)) break;
    const match = line.match(/^\s*version\s*:\s*(.+)\s*$/);
    if (!match) continue;

    let version = match[1].trim();
    if (
      (version.startsWith('"') && version.endsWith('"')) ||
      (version.startsWith("'") && version.endsWith("'"))
    ) {
      version = version.slice(1, -1);
    }
    return version;
  }
  return '';
}

export async function runChecked(
  cmd: string[],
  opts: { cwd?: string; env?: Record<string, string | undefined> } = {},
): Promise<void> {
  const proc = Bun.spawn(cmd, {
    cwd: opts.cwd,
    env: opts.env,
    stdout: 'inherit',
    stderr: 'inherit',
  });
  const exitCode = await proc.exited;
  if (exitCode !== 0) {
    throw new Error(`Command failed (${exitCode}): ${cmd.join(' ')}`);
  }
}
