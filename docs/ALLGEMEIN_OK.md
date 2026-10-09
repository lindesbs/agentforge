# Allgemeine bewährte Vorgehensweisen

Nur bestätigte Erkenntnisse aufnehmen, die unabhängig von der konkreten Projektarchitektur übertragbar sind.

## Konfigurationsänderungen kontrolliert übernehmen

**Erkenntnis:** Eine Vorschau, ein Änderungsstand-Check und eine ausdrückliche Speicherung trennen die Prüfung vom eigentlichen Dateieingriff.

**Regel:** Vor dem Speichern ein Diff zeigen und die zugrunde liegende Dateiversion überprüfen; Vorschau darf die Originaldatei nicht verändern.

**Prüfung:** Regressionstest: Vorschau belässt Originaldatei unverändert; abweichender Datei-Hash verhindert Übernahme.
