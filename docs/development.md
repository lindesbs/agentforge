# Local development

## Projekt direkt beim Start öffnen

AgentForge akzeptiert beim Start ein optionales Projektverzeichnis als **erstes Positionsargument**:

```sh
./build/bin/agentforge ~/GolandProjects/mein-projekt
./build/bin/agentforge .
./build/bin/agentforge /absoluter/pfad
./build/bin/agentforge --help
```

Ohne Argument öffnet AgentForge die normale Projektauswahl. Relative Pfade werden gegenüber dem aktuellen Arbeitsverzeichnis aufgelöst. Bei ungültigem Pfad, Datei statt Verzeichnis, symbolischem Link als Projektwurzel oder mehreren Positionsargumenten wird der Start mit einer Fehlermeldung abgebrochen. Das angegebene Projekt wird im UI automatisch eingelesen, ohne es zu verändern.

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

`make build` compiles for the **current Linux CPU architecture** (usually amd64 on x86-64); it does not cross-compile. The installed GTK/WebKit system packages remain runtime requirements. The Wails production build requires both \`production\` and \`webkit2_41\` tags. A build with only \`webkit2_41\` may compile but exit at runtime with \`Wails applications will not build without the correct build tags.\` Rebuild using \`make build\` if you encounter this message. `make install-system-deps` is the only target that invokes sudo. The `setup-go` target runs `go mod tidy`, which can update `go.mod` and `go.sum`; commit the resulting lock/checksum changes as appropriate.


## Prerequisites
Go 1.23+, Node.js and npm, Wails v2 CLI and platform-specific Wails prerequisites (GTK/WebKit on Linux). Follow the official Wails v2 installation instructions.

## Run
```sh
cd frontend
npm install
npm run build
cd ..
go test ./internal/...
wails dev
```

The checked-in `frontend/dist/index.html` is a placeholder solely to satisfy Go embed before the first frontend build. Vite replaces it with the compiled interface.

## Current scope
The initial inspector takes an explicit directory path, reads file names and basic framework markers, and lists recognized project-local configuration files. It does not parse agent prompts, modify files, execute commands, or automatically scan the home directory.

## Known limitations
- A native OS directory chooser is available via the Browse button; manual paths remain supported.
- The inspector intentionally does not follow symlinks and does not recurse through arbitrary project files.
- Codex/Claude provider paths are an initial subset, not a complete provider compatibility guarantee.
- GitHub Actions now runs Go inspector tests, Vue typechecking/build, and a Linux desktop compile. Successful CI status must be verified on GitHub; macOS and Windows packaging are not yet covered.

## Configuration editor (early MVP)

Click a discovered configuration file to edit its **raw UTF-8 text**. AgentForge checks the file's SHA-256 digest before preview and before save, displays a line-oriented change preview, creates a recovery backup alongside the file, then replaces the original using a same-directory temporary file. Only known project-local configuration paths can be edited. It does **not yet parse structured agent properties** or validate provider-specific semantics.

**Security limitations:** backups contain the original file contents and should be protected like the originals. This implementation does not yet prevent all concurrent write races or hard-link attacks. Do not use the editor on untrusted projects or security-sensitive configuration files until those issues are addressed. A future release will add structured validation and stronger filesystem safeguards.

## Structured Claude agent properties (experimental)

For existing `.claude/agents/*.md` files with simple single-line YAML frontmatter, the UI can edit existing `name` and `description` keys without rewriting other lines. The generated native text is shown in the same mandatory preview before saving. Complex or absent fields must be edited in the raw editor. This is not a general-purpose YAML parser or provider semantic validator.

## CI artifacts

GitHub Actions runs frontend build/typecheck, Go unit tests, and a Linux desktop compilation. The desktop job executes `go mod tidy` before building and uploads `agentforge-linux-amd64` when compilation succeeds. Dependency lockfiles should be committed in a follow-up change for reproducibility.

## Configuration diagnostics

The editor reports JSON syntax errors for Claude settings files and blocks saving invalid JSON. Markdown files receive non-blocking diagnostics for empty content and missing Claude agent frontmatter; TOML currently has informational diagnostics only, **not syntax validation**. These checks are not complete provider-schema validation. Recovery backups are created with owner-only permissions (0600), even if the original file is more widely readable. Known limitations of concurrent write races and hard links remain; avoid untrusted or security-critical directories.

## Global template library (MVP)

Global templates are local copies stored under the operating system's user configuration directory (typically `~/.config/agentforge/templates` on Linux). Each template contains a name, provider, kind, original relative path and the native file text. Template files have owner-only permissions (0600); the library folder has mode 0700 when created.

In an opened project, choose a discovered file and save its **current on-disk contents** as a global template. To reuse a template, select a compatible provider/kind entry and preview it against an **existing** target file. Applying changes the editor draft, not the target file; the existing diff preview and explicit save step are still required.

Limitations: the first library supports existing configuration targets only, not creation of new agent files. It does not yet support rename, deletion, tags, import/export, template search or project synchronization. Sensitive settings can be copied into templates: review source contents before saving a template. No agent execution occurs.

## Entwicklungs-Lernprotokoll und Sprachvorlagen

Die vier verpflichtenden Dokumente `docs/PROJEKT_OK.md`, `docs/PROJEKT_NOK.md`, `docs/ALLGEMEIN_OK.md` und `docs/ALLGEMEIN_NOK.md` sowie der Ablauf sind in `docs/learning-protocol.md` beschrieben. Vor jeder Codeänderung relevante Einträge lesen; bestätigte Erkenntnisse vor dem Commit in der passenden Datei ergänzen.

Die lokale Erkenntnisbibliothek im Desktop-UI ist von der bisherigen Konfigurationsvorlagenbibliothek getrennt. Neue Erkenntnisse lassen sich mit Kategorie, Sprache, Titel, nachgewiesener Ursache/Erkenntnis, verbindlicher Regel und überprüfbarer Kontrolle speichern. Eine Ansicht erzeugt Markdown zum gezielten Einfügen in die zuständige Dokumentationsdatei. Der Export ist **manuell**; es gibt keine automatische Veränderung von Projektdateien und keine Ausführung von Agenten.

Der Benutzerkonfigurationsordner enthält die Einträge unter `agentforge/learning` als einzelne JSON-Dateien mit Dateimodus 0600. Sprachkategorien sind etwa Allgemein, Go, PHP, TypeScript, Python, Rust, Kotlin, Java, Shell und SQL. Projektspezifische Regeln werden nicht allein durch Auswahl einer Programmiersprache automatisch zu allgemeinen Regeln.

## Bestehende Erkenntnisse importieren

In „Entwicklungs-Erkenntnisse“ nach Auswahl eines Projekts auf **„Projektdokumente durchsuchen“** klicken. Importiert werden nur ausdrücklich bestätigte und einzeln ausgewählte Einträge der vier Markdown-Protokolle, deren Überschriften und Pflichtfelder erkannt werden. Ungültige oder unvollständige Einträge bleiben unverändert und werden als übersprungene Abschnitte angezeigt. Backend und UI übernehmen nichts automatisch; ursprüngliche Dokumente werden nie geändert. Weitere Einzelheiten: `docs/learning-protocol.md`.
