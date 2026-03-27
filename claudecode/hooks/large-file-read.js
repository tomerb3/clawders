#!/usr/bin/env node
/**
 * PreToolUse hook for Read operations.
 * Reminds Claude to use chunked reading (offset/limit) for files > 2000 lines.
 *
 * From /home/baum/.claude/CLAUDE.md:
 * - When reading large files, run `wc -l` first to check the line count
 * - If over 2000 lines, use 'offset' and 'limit' parameters on the Read tool
 *   to read in chunks rather than attempting to read the entire file at once
 */

const fs = require('fs');
const path = require('path');

const input = JSON.parse(fs.readFileSync(0, 'utf-8'));

const LINE_THRESHOLD = 2000;

function countLines(filePath) {
  try {
    const stats = fs.statSync(filePath);
    // Quick estimate: if file is huge, read only a sample
    if (stats.size > 10 * 1024 * 1024) {
      // For large files, sample first 1MB to estimate line count
      const fd = fs.openSync(filePath, 'r');
      const buffer = Buffer.alloc(1024 * 1024);
      const bytesRead = fs.readSync(fd, buffer, 0, buffer.length, 0);
      fs.closeSync(fd);
      const sample = buffer.slice(0, bytesRead).toString('utf-8');
      const sampleLines = sample.split('\n').length;
      const estLines = Math.ceil((sampleLines / bytesRead) * stats.size);
      return estLines;
    }
    const content = fs.readFileSync(filePath, 'utf-8');
    return content.split('\n').length;
  } catch {
    return 0;
  }
}

function main() {
  // Only process Read tool
  if (input.tool_name !== 'Read') {
    console.log(JSON.stringify({ continue: true }));
    return;
  }

  const filePath = input.tool_input?.file_path;
  if (!filePath) {
    console.log(JSON.stringify({ continue: true }));
    return;
  }

  // Resolve to absolute path
  const absolutePath = path.isAbsolute(filePath)
    ? filePath
    : path.join(input.cwd || process.cwd(), filePath);

  // Check if file exists and get line count
  let lineCount = 0;
  try {
    lineCount = countLines(absolutePath);
  } catch {
    // File might not exist yet, skip
    console.log(JSON.stringify({ continue: true }));
    return;
  }

  if (lineCount > LINE_THRESHOLD) {
    const reminder = `📖 LARGE FILE READ: ${filePath} (${lineCount} lines)

CLAUDE.md instructs:
• Run \`wc -l\` first to verify line count
• For files > 2000 lines: use Read tool's \`offset\` and \`limit\` parameters
• Read in chunks (e.g., offset: 0, limit: 2000) rather than entire file

Example: Read with offset=0 limit=2000, then offset=2000 limit=2000, etc.`;

    console.log(JSON.stringify({
      continue: true,
      suppressOutput: false,
      systemMessage: reminder
    }));
    return;
  }

  // Small file - continue without message
  console.log(JSON.stringify({ continue: true }));
}

main();
