# Configuration model

## Principles
1. Native provider files remain authoritative.
2. Each adapter declares supported paths, fields, validations and write capabilities.
3. The UI must show unsupported or unknown settings without deleting them.
4. A global library is independent of project files; applying templates requires an explicit preview.
5. Agent relationships in the UI are descriptive configuration metadata, not an execution graph.

## Initial discovery targets
- Codex: `AGENTS.md`, `.codex/config.toml` where supported, and user-level configuration when explicitly selected.
- Claude Code: `.claude/agents/*.md`, `CLAUDE.md`, and applicable settings files.
Actual semantics and precedence are provider-version-specific and must be verified against official documentation before implementing writes.

## Proposed adapter contract (conceptual)
- Detect(projectRoot) -> candidate files
- Read(path) -> raw document + structured supported fields + diagnostics
- Validate(change) -> issues and compatibility warnings
- Preview(change) -> exact diff and potential information loss
- Save(change, expectedHash) -> atomic update or conflict

## Templates
Global template records have a stable identifier, name, description, tags, provider compatibility and source content. Templates are never implicitly synchronized with projects.

## Security and integrity
No automatic command execution. Guard approved roots and symlinks. Use optimistic locking. Prefer lossless editing; refuse unsafe serialization unless the user explicitly approves a complete rewrite after seeing the diff.
