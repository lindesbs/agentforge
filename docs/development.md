# Local development

## Linux build with Make

From the repository root:

```sh
make install-system-deps  # Optional: installs Arch/CachyOS or Debian/Ubuntu system packages via sudo
make setup                # Checks prerequisites, installs Go/npm dependencies and Wails CLI
make test                 # Go tests and TypeScript checks
make build                # Generates build/bin/agentforge (Wails production tag)
./build/bin/agentforge
```

For iterative desktop development run `make dev`. `make doctor` only checks prerequisites, and `make help` lists targets.

`make build` compiles for the **current Linux CPU architecture** (usually amd64 on x86-64); it does not cross-compile. The installed GTK/WebKit system packages remain runtime requirements. The Wails production build requires both `production` and `webkit2_41` tags. A build with only `webkit2_41` may compile but exit at runtime with `Wails applications will not build without the correct build tags.` Rebuild using `make build` if you encounter this message. `make install-system-deps` is the only target that invokes sudo. Setup downloads and verifies the checked-in Go module graph and uses `npm ci` for the checked-in frontend lockfile. Tests and desktop builds use `-mod=readonly` so incomplete module metadata fails instead of silently changing the dependency graph.

When intentionally changing dependencies, run `go mod tidy` or `npm install` in `frontend` as appropriate and review the resulting manifest and checksum/lockfile changes together. Ordinary setup, CI and Wails startup should not regenerate these files.


## Prerequisites
Go 1.23+, Node.js and npm, Wails v2 CLI and platform-specific Wails prerequisites (GTK/WebKit on Linux). Follow the official Wails v2 installation instructions.

## Run
```sh
cd frontend
npm ci
npm run build
cd ..
go test ./internal/...
wails dev
```

The checked-in `frontend/dist/index.html` is a placeholder solely to satisfy Go embed before the first frontend build. Vite replaces it with the compiled interface.

## UI interaction checks

From `frontend`, run `npx playwright install --with-deps chromium` once, then
`npm run test:ui`. Tests start Vite automatically, exercise the real Vue UI with
a test-only Wails bridge, measure token contrast and run automated accessibility
checks. Use `PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH` only when deliberately testing
an existing Chromium installation. Screenshots and traces go to `build/ui-tests`.
The suite also covers template search/tags, compatibility, import review and
failure recovery, explicit template application, and definition-map navigation.
These tests supplement Go tests; they do not substitute for native OS or
filesystem validation. See `docs/ui-ux-review.md` for the adopted UI/UX skills,
the interaction contracts and validation limits.

The workspace and learning library are separate views; changing views retains
their form state. Replacing an edited file, reloading, applying a template or
opening another project requires a discard decision when drafts are dirty.
The review/save flow still requires a fresh diff and an explicit save.

## Current scope
Discovery takes an explicit directory path, reads file names and basic framework
markers, and lists recognized project-local configuration files. Discovery does
not parse prompts or modify files. Explicitly opening a file enables inspection,
native editing, supported structured Claude fields and reviewed saving. Nothing
executes project commands or automatically scans the home directory.

## Known limitations
- A native OS directory chooser is available via the Browse button; manual paths remain supported.
- The inspector rejects a symlink selected as the project root and canonicalizes ancestor aliases of the selected directory. It checks provider directories before discovering their children, skips symlinks inside the project and does not recurse through arbitrary project files. Discovery assumes directories are not concurrently replaced; these checks are not a defense against a malicious process changing the tree during a scan.
- Discovery diagnostics show project-relative paths for skipped links, unexpected file types and inaccessible paths. Other accessible results remain visible. Missing optional configuration files are normal and produce no warning. Diagnostics do not expose file contents, link targets or raw filesystem errors.
- Codex/Claude provider paths are an initial subset, not a complete provider compatibility guarantee.
- GitHub Actions now runs Go inspector tests, Vue typechecking/build, and a Linux desktop compile. Successful CI status must be verified on GitHub; macOS and Windows packaging are not yet covered.

## Configuration editor (early MVP)

Click a discovered configuration file to inspect its supported structured fields and edit its **raw UTF-8 text**. AgentForge checks the file's SHA-256 digest before preview and before save, displays a line-oriented change preview, creates a recovery backup alongside the file, then replaces the original using a same-directory temporary file. Only known project-local configuration paths can be edited. The read-only detail view parses a documented subset; structured writing and save validation retain the limitations below.

**Security limitations:** backups contain the original file contents and should be protected like the originals. This implementation does not yet prevent all concurrent write races or hard-link attacks. Do not use the editor on untrusted projects or security-sensitive configuration files until those issues are addressed. Saves within one app process are serialized; external writers can still race the final hash check and rename. Structured validation covers only the documented subset, not complete provider semantics.

## Structured Claude agent properties (experimental)

For existing `.claude/agents/*.md` files with simple single-line YAML frontmatter,
the UI edits existing `name` and `description` using YAML source positions,
preserving comments, spacing around values, CRLF and unknown settings. Quoted
values are decoded correctly. Multiline styles, anchors, aliases, explicit tags,
missing fields and malformed/duplicate mappings are refused by this editor;
use the native editor instead. Changes still require diff review and saving.

## CI artifacts

GitHub Actions runs frontend build/typecheck, Go unit tests with the race detector, and a Linux desktop compilation. Jobs use the checked-in `go.mod`, `go.sum` and `frontend/package-lock.json`; the desktop job verifies Go modules before building and uploads `agentforge-linux-amd64` when compilation succeeds. Local validation does not establish that a GitHub Actions run passed.

## Configuration diagnostics

The **loaded-file detail view** uses YAML, TOML and JSON parsers to report syntax problems and selected provider field errors. Claude agent details support quoted and multiline names/descriptions, model and string/list tool declarations. Codex details display selected model/sandbox settings and declared role names, descriptions and configuration references. Referenced files are never opened. Empty instruction documents are identified without including prompt text in the structured summary. The field list and limits are documented in `docs/configuration.md`.

These details describe the loaded file, not unsaved edits. They refresh when a
file is opened, reloaded or saved. The draft is checked again on preview and by
the backend before saving. Unknown settings remain untouched; raw parser errors
(which may contain source values) are never returned by validation/inspection.

The **draft/save validator** shares the provider parser: invalid TOML, malformed
or missing required Claude agent frontmatter, and invalid settings JSON block
saving. Empty instruction Markdown is a warning. Neither path implements complete
provider-schema validation. Recovery backups use owner-only permissions (0600);
external concurrent write races and hard links remain known limitations.

## Global template library (MVP)

Global templates are local copies stored under the operating system's user configuration directory (typically `~/.config/agentforge/templates` on Linux). Each template contains a name, provider, kind, original relative path and the native file text. Template files have owner-only permissions (0600); the library folder has mode 0700 when created.

In an opened project, choose a discovered file and save its **current on-disk contents** as a global template. To reuse a template, select a compatible provider/kind entry and preview it against an **existing** target file. Applying changes the editor draft, not the target file; the existing diff preview and explicit save step are still required.

The separate **Template library** view supports case-insensitive name/tag/path
search and provider filtering. Use at most 12 tags, each up to 32 characters;
tags are trimmed, lowercased and deduplicated. Existing untagged records remain
readable. Compatibility requires the same provider and configuration kind.
Applying is full-file replacement, not a merge; project overrides appear as
removed/added lines, and stale target hashes block replacement.

**Import/export:** Review export exposes version-1 `agentforge-template` JSON in
a read-only text field for manual copying to a local file. Paste that JSON into
Import a template, review native text and portability diagnostics, then choose
Import reviewed template. Editing the payload invalidates approval. The backend
revalidates at import and creates a new private record, never overwriting another
ID or writing a project file. IDs, timestamps and absolute source-project roots
are not exported. Native contents may include credentials; exports are neither
redacted nor uploaded. Referenced files are not bundled.

Limits: existing configuration targets only; no new agent-file creation, rename,
deletion, automatic merge/conversion, synchronization or native import/export
file picker. Invalid native content is blocked on import and project save.
Snapshotting on-disk content into a template remains non-destructive even when
the original needs repair. No agent execution occurs.

## Definition map and compatibility

The **Definition map** is a responsive, read-only containment tree. It groups
discovered files by provider and shows declared roles and references from loaded
snapshots only. Opening a node uses the existing draft-protected file editor.
It never follows a reference, infers execution order or resolves vendor-specific
configuration precedence. Reopen files to refresh snapshots after external edits.

The provider capability model lists format, validation scope, available simple
fields and template support. Portability diagnostics flag supported role paths
that are absolute, home/environment-dependent or outside the destination, and
warn that relative references are not bundled. Unknown settings always need
manual review; these diagnostics are not a secret scanner or portability guarantee.

## Entwicklungs-Lernprotokoll und Sprachvorlagen

Die vier verpflichtenden Dokumente `docs/PROJEKT_OK.md`, `docs/PROJEKT_NOK.md`, `docs/ALLGEMEIN_OK.md` und `docs/ALLGEMEIN_NOK.md` sowie der Ablauf sind in `docs/learning-protocol.md` beschrieben. Vor jeder Codeänderung relevante Einträge lesen; bestätigte Erkenntnisse vor dem Commit in der passenden Datei ergänzen.

Die lokale Erkenntnisbibliothek im Desktop-UI ist von der bisherigen Konfigurationsvorlagenbibliothek getrennt. Neue Erkenntnisse lassen sich mit Kategorie, Sprache, Titel, nachgewiesener Ursache/Erkenntnis, verbindlicher Regel und überprüfbarer Kontrolle speichern. Eine Ansicht erzeugt Markdown zum gezielten Einfügen in die zuständige Dokumentationsdatei. Der Export ist **manuell**; es gibt keine automatische Veränderung von Projektdateien und keine Ausführung von Agenten.

Der Benutzerkonfigurationsordner enthält die Einträge unter `agentforge/learning` als einzelne JSON-Dateien mit Dateimodus 0600. Sprachkategorien sind etwa Allgemein, Go, PHP, TypeScript, Python, Rust, Kotlin, Java, Shell und SQL. Projektspezifische Regeln werden nicht allein durch Auswahl einer Programmiersprache automatisch zu allgemeinen Regeln.
