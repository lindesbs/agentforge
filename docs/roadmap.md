# Roadmap

## Milestone 0 — Foundation
- [x] Product scope and architectural principles documented
- [x] Repository development instructions
- [ ] Verify toolchain versions and initialize Wails/Go/Vue scaffold
- [ ] Establish unit-test and CI workflows

## Milestone 1 — Read-only project inspector
- [ ] Open a local directory
- [ ] Detect project metadata and supported native agent files
- [ ] Show agent definitions and diagnostics
- [ ] Unit-test safe path discovery

## Milestone 2 — Safe native editor
- [ ] Provider capability model
- [ ] Codex configuration parsing, validation and diff preview
- [ ] Claude Code agent Markdown/frontmatter parsing, validation and diff preview
- [ ] Round-trip tests including comments and unknown settings
- [ ] Atomic writes, backups, conflict detection

## Milestone 3 — Global template library
- [ ] Create, tag and search local templates
- [ ] Explicit apply/copy into a project
- [ ] Preview conflicts and project-specific overrides

## Milestone 4 — Visual design
- [ ] Graphical relationships between definitions, without implying execution
- [ ] Provider compatibility indicators
- [ ] Import/export and portability diagnostics

## Acceptance criteria for first usable release
Open an existing project; discover native definitions; edit supported fields; validate and preview diffs; save without data loss; reuse a global template. No agent execution is present.
