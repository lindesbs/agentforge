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

## UI/UX design skills

For interface and interaction work, use the relevant stable skills from
[`aditya-ariosity/ux-ui-skills`](https://github.com/aditya-ariosity/ux-ui-skills):
`ux-ui-audit` for the actual user journey and `design-system-review` for shared
tokens and control behavior. Read the selected skill and its required references
before applying it. The reviewed upstream revision, local decisions, evidence
and component contracts are recorded in `docs/ui-ux-review.md`.

- Keep project editing and the learning library distinct; preserve state across views.
- Protect unsaved native and structured drafts before replacing them.
- Reuse semantic tokens from `frontend/src/style.css` and native controls.
- Inspect rendered desktop and narrow layouts, keyboard behavior and error recovery.
- Run `cd frontend && npm run test:ui` for consequential UI changes. Browser tests
  use a test-only Wails fixture, so they do not establish native OS or filesystem behavior.
- Do not claim accessibility conformance from automated checks alone.

## Verbindliches Entwicklungs-Lernprotokoll

Für jede Entwicklungsaufgabe sind `docs/PROJEKT_OK.md`, `docs/PROJEKT_NOK.md`, `docs/ALLGEMEIN_OK.md` und `docs/ALLGEMEIN_NOK.md` zu berücksichtigen. Lies vor Änderungen die relevanten Einträge aller vier Dateien.

- Dokumentiere nach Diagnose oder Review **bestätigte** Fehlerursachen mit Vermeidungsregel und Prüfung in der zutreffenden `*_NOK.md`.
- Dokumentiere nach erfolgreicher Lösung bewährte Entscheidungen mit überprüfbarer Kontrolle in der zutreffenden `*_OK.md`.
- Vor jedem Commit klären, ob neue Erkenntnisse gemeinsam mit dem Code dokumentiert werden müssen.
- Der Geltungsbereich richtet sich nach der tatsächlichen Übertragbarkeit, nicht nach dem Entstehungsprojekt.
- Vorhandene Einträge ergänzen statt Duplikate anzulegen; veraltete Regeln aktualisieren.
- Nie Passwörter, Tokens, vollständige DSNs oder personenbezogene Daten dokumentieren.
- Unbestätigte Vermutungen nicht als Regeln übernehmen.

Siehe `docs/learning-protocol.md` für Kategorien, Format und Verwendung der sprachbezogenen Vorlagenbibliothek.
