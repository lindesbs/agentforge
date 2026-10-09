# AgentForge

**Visual AI Agent Configuration Manager**

AgentForge is a planned local, cross-platform application to inspect projects and manage native AI coding-agent configurations. It **does not execute agents**.

## Goals

- Open a local development project and detect supported agent configuration files.
- Inspect, edit, validate and safely save native configurations for Codex and Claude Code.
- Maintain a reusable global library of agent templates with project-specific overrides.
- Show configuration differences before saving and preserve unknown settings.
- Keep configuration files versionable in Git.

## Proposed stack

Go + Wails + Vue 3 + TypeScript. The visual editor may use Vue Flow.

## Status

Early planning and architecture. No executable application is shipped yet.

## Development

See `docs/architecture.md`, `docs/roadmap.md`, `docs/configuration.md` and `AGENTS.md` in the foundation pull request.

## Scope boundary

No model calls, agent execution, task scheduling or process orchestration. AgentForge edits configuration; the actual agent tools remain responsible for running agents.
