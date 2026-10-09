# Architecture (proposal)

## Purpose
AgentForge is a cross-platform desktop configuration manager for AI coding agents. It opens an existing project and edits the supported **native configuration files**. It never starts or controls agents.

## Layers
1. **UI:** Vue 3 + TypeScript; project browser, agent editor, validation diagnostics, diff preview, global template library.
2. **Desktop bridge:** Wails; narrow, typed APIs for opening projects, inspecting configurations, validating changes, previewing diffs and saving.
3. **Application services (Go):** project discovery, template management, change planning and safe persistence.
4. **Domain (Go):** projects, configuration documents, agent definitions, provider capabilities, validation issues and changesets.
5. **Provider adapters:** Codex and Claude Code, each responsible for identifying supported paths and semantics. Unsupported features must be surfaced, not silently discarded.
6. **Filesystem:** bounded roots, symlink-safe path validation, optimistic concurrency, atomic replacement and recovery.

## Data ownership
- Native files are the source of truth for project configuration.
- Global templates are stored in the user's app configuration directory; copying or applying one to a project is an explicit operation.
- App state stores recent project paths and UI preferences, not duplicated canonical agent definitions.
- A future optional portable interchange format must never overwrite provider-specific settings implicitly.

## Editing strategy
1. Discover project paths without executing project content.
2. Read files and capture hashes.
3. Parse known structures, retaining original source and unknown fields.
4. Show supported fields, diagnostics and proposed modifications.
5. Preview an exact diff.
6. Check hashes and permitted paths again; write atomically with recovery strategy.

## Initial project detection
Look for composer.json (Symfony/Contao hints), go.mod, package.json, .git, AGENTS.md, .codex/config.toml and .claude/agents/*.md. Detection informs the UI but must not run any project commands.

## Explicit non-goals
Agent execution, terminal embedding, LLM inference, workflow scheduling, process monitoring, token billing, autonomous Git changes, and cloud accounts.

## Design decisions to validate
- Lossless TOML and Markdown/frontmatter editing: research available libraries and test round trips before choosing one.
- Cross-platform symlink/path behavior and atomic replacement guarantees.
- Provider version compatibility matrix and capability-based editor fields.
