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

## Read-only definition inspection

Selecting a discovered file passes its already-loaded UTF-8 text (maximum 1 MiB)
to `internal/providers`. Project discovery remains metadata-only. Inspection does
not write, serialize configuration back to disk, open references or execute any
agent instructions. The original text remains the editor's source of truth.

| Native file | Displayed fields and checks |
| --- | --- |
| `.claude/agents/*.md` | YAML frontmatter `name`, `description`, optional `model` and string/list `tools`; quoted and multiline strings; missing/invalid required fields, malformed or duplicate YAML keys, missing delimiter and empty prompt diagnostics |
| `.codex/config.toml` | TOML syntax; optional string `model`, `model_provider`, `model_reasoning_effort`, `sandbox_mode`; map-valued `[agents.<role>]` entries with optional string `description` and `config_file` |
| `.claude/settings.json`, `.claude/settings.local.json` | Exactly one JSON object; optional string `model`; no environment, hook or credential values in the summary |
| `AGENTS.md`, `CLAUDE.md` | Provider/document kind and an empty-document warning; instruction text is not duplicated in the summary |

Unknown fields are not interpreted or displayed by the summary. Diagnostics use
fixed messages and known field names, never raw parser errors or arbitrary source
values. The read-only view does not validate model identifiers, permission policy
semantics, all provider options, JSON duplicate keys, role reference targets or
configuration precedence. The same supported checks now run for drafts and
before writes: error-level findings block saving; warnings remain advisory.

Field shapes were checked against the upstream [Codex configuration schema](https://github.com/openai/codex/blob/main/codex-rs/core/config.schema.json)
and an official [Claude Code agent definition](https://github.com/anthropics/claude-code/blob/main/plugins/plugin-dev/agents/agent-creator.md).
These are compatibility references for the supported subset, not pinned vendor
schemas or a guarantee that every future provider feature is understood.

## Implemented capability model and adapter boundaries

`providers.Capabilities()` lists provider/kind, native format, validation scope,
simple structured fields and template support. All known types use native-text
editing; only simple Claude `name`/`description` fields have structured writing.
The UI exposes these differences rather than suggesting provider equivalence.

The architecture follows these boundaries:
- Detect(projectRoot) -> candidate files
- Read(path) -> raw document + structured supported fields + diagnostics
- Validate(change) -> issues and compatibility warnings
- Preview(change) -> exact diff and potential information loss
- Save(change, expectedHash) -> atomic update or conflict

## Templates
Global template records contain an ID, name, normalized tags, provider, kind,
canonical relative source path, creation timestamp and native content. The
version-1 portable envelope excludes local IDs/timestamps. Import rejects unknown
envelope fields, unsupported versions, traversal paths, mismatched provider/kind,
oversized content and malformed native configuration. Validation only previews;
separate import creates a new library entry. Project application requires an
existing compatible target, matching content hash, full diff and explicit save.
Templates are never implicitly synchronized or merged with projects.

## Security and integrity
No automatic command execution. Guard approved roots and symlinks. Use optimistic locking. Prefer lossless editing; refuse unsafe serialization unless the user explicitly approves a complete rewrite after seeing the diff.
