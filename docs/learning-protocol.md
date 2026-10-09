# Entwicklungs-Lernprotokoll

## Klassifikation

| Geltungsbereich | Bewährte Lösung | Fehler / Vermeidungsregel |
| --- | --- | --- |
| Projektspezifisch | `docs/PROJEKT_OK.md` | `docs/PROJEKT_NOK.md` |
| Projektunabhängig | `docs/ALLGEMEIN_OK.md` | `docs/ALLGEMEIN_NOK.md` |

Ausschlaggebend ist die tatsächliche Übertragbarkeit einer Erkenntnis, nicht der Ort ihrer Entstehung.

## Verbindlicher Entwicklungsablauf

1. **Vor der Implementierung:** Relevante Einträge aller vier Dateien lesen und bei Planung und Umsetzung berücksichtigen.
2. **Nach Diagnose oder Review:** Bestätigte Fehlerursachen im passenden `*_NOK.md` ergänzen, mit Ursache, Regel und Prüfung.
3. **Nach erfolgreicher Lösung:** Bestätigte Vorgehensweisen im passenden `*_OK.md` ergänzen.
4. **Vor jedem Commit:** Auf neue Erkenntnisse prüfen und die Dokumentation gegebenenfalls zusammen mit dem Code aktualisieren.
5. **Datenschutz:** Keine Passwörter, Zugangsdaten, vollständigen DSNs, Tokens oder personenbezogenen Daten dokumentieren.

## Qualitätsregeln

- Nur bestätigte Erkenntnisse übernehmen; Vermutungen zunächst untersuchen.
- Projektbezogene Details gehören nur in `PROJEKT_*`, allgemeine nur in `ALLGEMEIN_*`.
- Bestehende Einträge ergänzen, nicht duplizieren.
- Für jede Regel eine überprüfbare Kontrolle nennen.
- Veraltete Einträge bei technischen Änderungen korrigieren.
- Programmiersprachen sind ein zusätzliches Ordnungskriterium, **kein Ersatz** für den Geltungsbereich.

## Bibliotheksfunktion in AgentForge

Die lokale Erkenntnisbibliothek speichert bestätigte Einträge nach Kategorie und Programmiersprache (z. B. Go, PHP, TypeScript, Python, Rust oder Allgemein). Zu jedem Eintrag gehören Titel, Erkenntnis/Ursache, verbindliche Regel und Kontrolle. Das ist eine Vorlagenbibliothek: Speichern eines Eintrags dort verändert **keine** Projektdateien. Eine erzeugte Markdown-Vorlage kann anschließend in die zuständige der vier Dateien übernommen werden.

Der Abschnitt `Standard-Stack für neue Projekte` aus einer AGENTS.md muss projektspezifisch ergänzt werden, solange keine konkreten Standard-Stack-Vorgaben vorliegen.

## Vorhandene Projekt-Erkenntnisse importieren

Nach Auswahl oder CLI-Start eines Projekts lässt sich in der Desktop-App unter „Entwicklungs-Erkenntnisse“ die Aktion **„Projektdokumente durchsuchen“** aufrufen. AgentForge liest ausschließlich die vorhandenen vier `docs/PROJEKT_*.md`- und `docs/ALLGEMEIN_*.md`-Dateien. Es folgen eine Vorschau und eine individuelle Übernahme je erkanntem Eintrag in die lokale Bibliothek. Die Projektdateien werden beim Import **nicht verändert**.

Der Import unterstützt vorhandene Markdown-Abschnitte mit `##` oder `###` und den beschrifteten Feldern `**Erkenntnis:**` beziehungsweise `**Bestätigte Erkenntnis/Ursache:**`, `**Regel:**` beziehungsweise `**Verbindliche Regel:**` sowie `**Prüfung:**`. Die Sprache wird aus `**Sprache:**` gelesen oder kann vor der Übernahme ausgewählt werden (Standard: `Allgemein`).

Unvollständige Abschnitte werden **nicht automatisch als bestätigte Erkenntnisse importiert**; die Oberfläche zeigt die Gründe an. Doppelte Kombinationen aus Kategorie, Sprache und Titel werden abgewiesen. Für jeden Import wird die Quelldatei erneut gelesen und mit dem Vorschau-Hash verglichen. Importierte Daten können vertrauliche Informationen enthalten – Inhalte vor einer späteren öffentlichen Veröffentlichung stets manuell prüfen.

Derzeit gibt es keinen generischen Markdown-/AGENTS.md-Freitextimport, keine automatische KI-Extraktion und keinen automatischen Upload zur Exchange-Plattform.
