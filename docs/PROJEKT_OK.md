# Projektbezogene bestätigte Lösungen — AgentForge

Hier nur nachgewiesene Lösungen dokumentieren, die von AgentForges Architektur, Konfiguration oder Konventionen abhängen. Jeder Eintrag braucht eine verbindliche Regel und eine überprüfbare Kontrolle.

## Vorlagenprüfung, Bibliotheksimport und Projektübernahme trennen

**Erkenntnis:** Ein versionierter JSON-Umschlag kann vor jedem Dateieingriff auf Format, Grenzen, Provider/Dateityp und native Syntax geprüft werden. Ein Import mit neuer ID überschreibt keine Bibliothekseinträge; Projektübernahme bleibt ein gesonderter, hashgeschützter Entwurf mit Diff.

**Regel:** Importprüfung bleibt seiteneffektfrei, Import revalidiert serverseitig. Export enthält keine Bibliotheks-ID oder absoluten Projektwurzelpfade. Native Inhalte bleiben unverändert und können sensible Daten enthalten; vor Kopieren und Importieren darauf hinweisen. Definitionskarten verwenden ausschließlich entdeckte Pfade und ausdrücklich geladene Inhalte, niemals implizite Referenzzugriffe.

**Prüfung:** `TestTaggedPortableRoundTrip`, `TestInvalidImportNeverCreatesLibrary`, `TestFirstUsableReleaseWorkflow` und die Browserfälle in `tests/roadmap.spec.ts`. Die Test-Bridge kopiert Listen wie die echte serialisierte Wails-Grenze, statt gemeinsame veränderliche Objekte zurückzugeben.

## Linux-Desktop-Build

**Erkenntnis:** Die AgentForge-Binärdatei benötigte bei der bestehenden Wails-v2-Konfiguration sowohl das Tag `production` als auch `webkit2_41`. Ohne `production` entstand eine Binärdatei, die beim Start einen Wails-Tag-Fehler ausgab.

**Regel:** Für den Linux-Produktionsbuild `CGO_ENABLED=1 go build -tags 'production,webkit2_41'` verwenden.

**Prüfung:** `make build`, anschließend `./build/bin/agentforge` starten. Die CI kontrolliert zusätzlich den bekannten Wails-Stub-Text.

## Reproduzierbare Build-Abhängigkeiten

**Erkenntnis:** Mit vollständigem `go.mod`, `go.sum` und `frontend/package-lock.json` bestehen die 16 Go-Tests einschließlich Race Detector, die Vue-Typprüfung und der Linux-Desktop-Build. Wiederholtes `make setup-go` und `npm ci` verändern diese Dateien nicht.

**Regel:** Reguläres Setup verwendet `go mod download` mit `go mod verify` sowie `npm ci`; Tests und Desktop-Kompilierung verwenden `-mod=readonly`. `go mod tidy` und `npm install` bleiben bewussten Abhängigkeitsänderungen vorbehalten, deren Manifest- und Lockdateiänderungen gemeinsam geprüft werden.

**Prüfung:** `make setup-go setup-frontend test`, `go test -mod=readonly -race ./internal/...` und `make build`; Prüfsummen von `go.mod`, `go.sum` und `frontend/package-lock.json` vor und nach dem Setup vergleichen.

## Begrenzte Projekt-Erkennung mit Diagnosen

**Erkenntnis:** Die explizite Prüfung von `.codex`, `.claude` und `.claude/agents` vor dem Zugriff auf deren Kinder verhindert die nachgewiesene Erkennung durch verlinkte Elternverzeichnisse. Ein nicht lesbares Agentenverzeichnis unterbricht nicht mehr die gesamte Inspektion; zugängliche Konfigurationsdateien bleiben im Ergebnis.

**Regel:** Nur die dokumentierten Projektpfade durchsuchen. Übersprungene Links, falsche Dateitypen und Zugriffsprobleme mit relativen Pfaden und festen Meldungen anzeigen; keine Inhalte oder Linkziele in Diagnosen übernehmen. Leere Listen im Bridge-Ergebnis als `[]` serialisieren.

**Prüfung:** `TestIgnoreSymlinkDirectories`, `TestBrokenSymlinksAreReported`, `TestUnexpectedPathTypes`, `TestUnreadableAgentsPreserveOtherResults` und `TestEmptyProjectUsesEmptyArrays` ausführen; `TestInspectAllSupportedPathsReadOnly` prüft die erkannten Pfade und unveränderten Dateiinhalte.

## Native Definitionen schreibgeschützt beschreiben

**Erkenntnis:** Ein reiner Provider-Parser kann aus dem bereits geöffneten Dokument YAML-Frontmatter, ausgewählte TOML-Einstellungen und Rollen lesen, ohne Kommentare, Zeilenenden oder unbekannte Einstellungen neu zu schreiben. Die Detailansicht unterstützt auch mehrzeilige und korrekt zitierte Beschreibungen. Feste Diagnosemeldungen verhindern, dass Quelltext aus Parserfehlern in Zusammenfassungen gelangt.

**Regel:** Strukturierte Inspektion und natives Schreiben getrennt halten. Nur ausdrücklich unterstützte Felder anzeigen, referenzierte Konfigurationen nicht öffnen und Grenzen der Prüfung im UI nennen. Die Übersicht beschreibt den geladenen Dateistand und wird nach Laden beziehungsweise Speichern aktualisiert.

**Prüfung:** `go test -mod=readonly -race ./internal/providers` ausführen; insbesondere `TestClaudeFrontmatterVariants`, `TestMalformedClaudeDefinitions`, `TestCodexSettingsAndRoles`, `TestDiagnosticsDoNotExposeSource` und `TestInspectLoadedConfigPreservesSource`.

## UI-Verträge mit gerenderten Zuständen prüfen

**Erkenntnis:** Die Skills `ux-ui-audit` und `design-system-review` aus `aditya-ariosity/ux-ui-skills` führten zu getrennten Arbeitsbereichen, einem Dirty-State-Dialog und gemeinsamen semantischen Tokens. Browserprüfungen deckten zusätzlich eine unpassende Landmark-Rolle und eine Lücke im Fokusumlauf des Dialogs auf; beide wurden vor Abschluss korrigiert.

**Regel:** Quelltextprüfung durch gerenderte Zustände, Tastaturabläufe und konkrete Fehlerfälle ergänzen. Wails-Fixtures bleiben ausschließlich im Testverzeichnis. Erfolgreiche Browser-Automation nicht als Nachweis für native Dateizugriffe oder vollständige Barrierefreiheit ausgeben.

**Prüfung:** `npm run test:ui`, `npm run typecheck` und `npm run build`; Screenshots bei 320, 768, 1100 und 1280 Pixeln prüfen. Skill-Revision, Zustandsverträge und verbleibende Prüflücken stehen in `docs/ui-ux-review.md`.
