# Roadmap

## Milestone 0 — Foundation
- [x] Product scope and architectural principles documented
- [x] Repository development instructions
- [x] Verify toolchain versions and initialize Wails/Go/Vue scaffold
- [x] Establish unit-test and CI workflows

Foundation validation: Go 1.23.12, Node.js 24.19.0 and npm 11.9.0;
16 Go tests pass with the race detector, Vue typechecking and production
asset compilation pass, and the Linux desktop binary compiles against
GTK 3.24.49 / WebKit2GTK 2.54.0 with `production,webkit2_41`.
A local Docker/Xvfb startup check confirms that the AgentForge window opens;
this is not an end-to-end test of editing or saving through the UI.
Dependency manifests/checksums are checked in and ordinary setup uses
`go mod download`, `go mod verify` and `npm ci` without regenerating them.
CI workflows are defined; their remote execution was not verified here.

## Milestone 1 — Read-only project inspector
- [x] Open a local directory
- [x] Detect project metadata and supported native agent files
- [x] Show agent definitions and diagnostics
- [x] Unit-test safe path discovery

The inspector lists supported configuration paths and PHP/Composer, Go and
Node.js markers. Discovery diagnostics report skipped links, unexpected path
types and inaccessible directories while retaining other results. Tests cover
linked provider directories and files, broken links, invalid roots, canonical
root paths, unreadable agent directories and read-only discovery. Selecting a
file shows a read-only definition: Claude agent frontmatter, selected Codex
settings and declared Codex roles, with syntax and selected-field diagnostics.
Discovery itself does not read configuration contents; detail inspection uses
only the explicitly opened document and never follows configuration references.
This is a supported subset, not complete provider-schema validation. Concurrent
replacement of directories during discovery remains outside the path checks.

## Milestone 2 — Safe native editor
- [x] Provider capability model
- [x] Codex configuration parsing, validation and diff preview
- [x] Claude Code agent Markdown/frontmatter parsing, validation and diff preview
- [x] Round-trip tests including comments and unknown settings
- [x] Atomic writes, backups, conflict detection

Implemented for the documented provider subset: the loaded view and draft/save
path share syntax and selected-field checks. Invalid TOML, YAML frontmatter and
settings JSON block writes. Simple Claude field edits preserve inline comments,
CRLF, unknown keys and body text; complex styles require native editing.
Same-directory replacement preserves file mode and creates a private backup.
Hash checks reject stale versions; app-local saves are serialized. This is not
a cross-process filesystem compare-and-swap or protection against hostile tree
replacement. Use trusted project directories; see the security limitations below.

## Milestone 3 — Global template library
- [x] Create, tag and search local templates
- [x] Explicit apply/copy into a project
- [x] Preview conflicts and project-specific overrides

Templates snapshot the on-disk file, with normalized tags and name/tag/path
search plus provider filtering. Applying to an existing same-provider/kind target
requires a dirty-draft decision, matching file hash, full replacement diff and
separate save. The UI explicitly warns that project overrides are replaced, not
merged. New target-file creation and automatic merging are not implemented.

## Milestone 4 — Visual design
- [x] Graphical relationships between definitions, without implying execution
- [x] Provider compatibility indicators
- [x] Import/export and portability diagnostics

The definition map groups paths by provider and shows declared roles/references
only for explicitly opened files. Semantic nested lists and visible connectors
express containment, never execution. References remain unopened and unverified.
Provider support and template compatibility are explicit. Versioned template
JSON can be reviewed, copied out and pasted in; import creates a new library ID
only after validation and confirmation. Portability warnings identify supported
reference risks and require manual review of unknown fields and secrets.

## Release verification and remaining limits

All original feature checklist items above are implemented within their stated
MVP boundaries, not a claim of complete vendor support or release hardening.
`TestFirstUsableReleaseWorkflow` exercises discovery, structured editing, preview,
safe save and template reuse against real temporary files. Browser interaction
tests use a synthetic Wails bridge; native build/start checks are separate.

Current local validation: 44 Go tests pass with the race detector and a fresh
run; 18 browser tests pass; Vue/TypeScript checking and production assets build
successfully. The updated Linux desktop binary compiles and opens its AgentForge
window under Xvfb. New library/map layouts were inspected at 320, 768 and 1280px.

Remaining verification/hardening: full native end-to-end saves and OS chooser,
screen-reader testing, Windows/macOS packaging, persistent draft recovery, and
protection against malicious concurrent filesystem changes or hard-link attacks.
Remote CI execution and user task testing are not established by local checks.

## Acceptance criteria for first usable release
Open an existing project; discover native definitions; edit supported fields; validate and preview diffs; save without data loss; reuse a global template. No agent execution is present.
