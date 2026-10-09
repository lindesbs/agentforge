# Projektbezogene Fehler und Vermeidungsregeln — AgentForge

Hier bestätigte, vom Projekt abhängige Fehler dokumentieren. Für jeden Eintrag sind Ursache, verbindliche Vermeidungsregel und nachprüfbare Kontrolle erforderlich.

## Unterschiedliche Prüfungen zwischen Detailansicht und Speichern

**Ursache:** Der bisherige Save-Validator behandelte TOML nur informativ und prüfte Claude-Frontmatter nicht mit dem bereits verfügbaren Provider-Parser. Außerdem ersetzte der einfache Feldeditor komplette YAML-Zeilen; dadurch gingen Inline-Kommentare und CRLF an diesen Stellen verloren.

**Regel:** Detailansicht und Entwurfsprüfung verwenden denselben dokumentierten Parser-Umfang. Speichern blockiert Fehler serverseitig. Strukturierte Änderungen ersetzen ausschließlich ein nachweislich einfaches skalares Quelltextsegment; komplexe Formen bleiben im nativen Editor.

**Prüfung:** `TestProviderDraftValidationAndRoundTrip`, `TestScalarPreservesCommentsCRLFAndUnknownFields` und `TestUnsafeStructuredFormsAreRejected`; ungültige Entwürfe verändern weder Original noch Backupbestand, gültige Änderungen erhalten unbekannte Einstellungen und Kommentare.

## Ungespeicherte Editorentwürfe gingen beim Wechsel verloren

**Ursache:** `clearEditor()` lief vor dem Laden eines anderen Dokuments beziehungsweise Projekts. Eine Dirty-State-Prüfung fehlte; der Browser-Audittest reproduzierte den Verlust eines bearbeiteten `AGENTS.md`-Entwurfs beim Wechsel zu `.codex/config.toml`.

**Regel:** Vor dem Ersetzen nativer oder noch nicht angewendeter strukturierter Änderungen explizit nachfragen. Abbrechen und Escape erhalten den Entwurf. Erst nach erfolgreichem Laden den bisherigen Dokumentzustand ersetzen; Fehler dürfen ihn nicht löschen. Reiner Ansichtswechsel zur Bibliothek erhält beide Formulare ohne Rückfrage.

**Prüfung:** `npm run test:ui` in `frontend`; insbesondere Dateiwechsel, Reload, Template-Anwendung, fehlgeschlagene Projektinspektion und Bibliotheksnavigation. Nach Änderungen an strukturierten Feldern darf eine vorherige Speichervorschau nicht mehr freigegeben sein.
