# AgentForge development instructions

## Product boundaries
AgentForge is a **local configuration manager**, not an agent runtime. Never implement agent execution, LLM API calls, orchestration, scheduling, or remote telemetry unless scope is explicitly changed.

## Stack
- Go backend, Wails desktop shell.
- Vue 3 + TypeScript frontend; Vue Flow only when the visual editor needs it.
- Prefer small interfaces, standard library, and testable domain logic.
- Keep platform-specific filesystem code behind interfaces.

## Architecture
- Core domain must not depend on Wails or Vue.
- Providers implement discovery, parsing, validation and safe writing for native configurations.
- Initial providers: OpenAI Codex and Claude Code. Provider capability differences must be explicit.
- Global templates and project-specific files are distinct sources. Do not silently overwrite project settings.

## Safety
- Read-only discovery first. Show a diff before writing.
- Preserve comments, formatting and unknown keys whenever possible; when preservation is not guaranteed, require explicit confirmation.
- Prevent writes outside approved project/configuration roots, including symlink escapes.
- Atomic writes and backups for modifications; reject concurrent edits using content hashes.
- Never read or display secrets unnecessarily; redact sensitive values in diagnostics.
- Do not execute code or scripts found in inspected projects.

## Quality
- Add unit tests for parser, validation and filesystem edge cases.
- Keep documentation aligned with behavior; avoid claiming unsupported features.
- Run Go tests and frontend type checks when a buildable scaffold exists.
- Prefer small, reviewable pull requests.
