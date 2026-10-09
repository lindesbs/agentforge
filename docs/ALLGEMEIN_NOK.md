# Allgemeine Fehlerbilder und Vermeidungsregeln

Nur bestätigte, projektunabhängig relevante Fehlerursachen aufnehmen. Jeder Eintrag muss Ursache, verbindliche Vermeidungsregel und überprüfbare Kontrolle enthalten.

## Lstat prüft nur die letzte Pfadkomponente

**Ursache:** `os.Lstat` auf einem vollständigen Dateipfad folgt weiterhin symbolischen Links in dessen Elternverzeichnissen. Die ursprüngliche Projekt-Erkennung konnte dadurch Dateien außerhalb des gewählten Projekts auflisten, wenn `.claude` oder `.codex` verlinkt war. Regressionstests haben beide Fälle vor der Korrektur reproduziert.

**Regel:** Innerhalb einer festgelegten Verzeichnisgrenze jede Verzeichniskomponente prüfen, bevor auf deren Kinder zugegriffen wird. Nur reguläre Dateien als Konfigurationsdateien akzeptieren. Pfadbasierte Einzelprüfungen schützen nicht gegen einen gleichzeitigen Austausch von Verzeichnissen; diese Einschränkung ausdrücklich dokumentieren.

**Prüfung:** Regressionstests für verlinkte Elternverzeichnisse, direkte Dateilinks und defekte Links ausführen. Außerhalb liegende Dateien dürfen weder in der Ergebnisliste noch als Linkziel in Diagnosen erscheinen. In AgentForge: `go test -mod=readonly -race ./internal/inspector`.
