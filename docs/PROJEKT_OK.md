# Projektbezogene bestätigte Lösungen — AgentForge

Hier nur nachgewiesene Lösungen dokumentieren, die von AgentForges Architektur, Konfiguration oder Konventionen abhängen. Jeder Eintrag braucht eine verbindliche Regel und eine überprüfbare Kontrolle.

## Linux-Desktop-Build

**Erkenntnis:** Die AgentForge-Binärdatei benötigte bei der bestehenden Wails-v2-Konfiguration sowohl das Tag `production` als auch `webkit2_41`. Ohne `production` entstand eine Binärdatei, die beim Start einen Wails-Tag-Fehler ausgab.

**Regel:** Für den Linux-Produktionsbuild `CGO_ENABLED=1 go build -tags 'production,webkit2_41'` verwenden.

**Prüfung:** `make build`, anschließend `./build/bin/agentforge` starten. Die CI kontrolliert zusätzlich den bekannten Wails-Stub-Text.
