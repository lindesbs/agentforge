# AgentForge Exchange — Architektur für das zukünftige Webportal

**Ergänzende Dokumente:** [Website-Konzept agentforge.art](website-agentforge-art.md). Die Domain ist geplant, aber nicht als registriert bestätigt.

## Produktaufteilung

**AgentForge Desktop** bleibt eine lokal laufende Wails-Anwendung. Eine Domain stellt kein Desktopprogramm im Browser bereit. Die Domain wird später das Webportal **AgentForge Exchange** hosten: öffentlicher Katalog, Dokumentation, Community, Downloadbereich und eine optionale Publishing-API. Offline-Nutzung der Desktop-App bleibt möglich.

## Inhaltstypen

1. **Agent-Konfigurationsvorlagen**: Anbieter (`codex`, `claude`), unterstützte Versionen/Dateitypen, neutrale Beschreibung, Tags, Programmiersprache/Framework, native Konfigurationsdatei, Changelog.
2. **Bestätigte Entwicklungs-Erkenntnisse**: `PROJEKT_OK`, `PROJEKT_NOK`, `ALLGEMEIN_OK`, `ALLGEMEIN_NOK`, Sprache/Framework, Erkenntnis beziehungsweise Ursache, verbindliche Regel, überprüfbare Kontrolle und Quellen/Prüfbelege.
3. **Sammlungen**: kuratierte Gruppen mehrerer kompatibler Vorlagen oder Erkenntnisse für einen Technologie-Stack.

Projektbezogene Einträge sind standardmäßig **privat**. Eine Veröffentlichung ist erst nach einer bewussten Entscheidung und gegebenenfalls nach Anonymisierung möglich. Die Klassifikation darf nicht stillschweigend von projektspezifisch auf allgemein geändert werden.

## Austauschformat v1

`schemas/exchange-item-v1.schema.json` beschreibt ein versioniertes, portables Manifest für Export/Import. Es enthält keine Anmeldedaten und erzwingt Metadaten zu Lizenz, Sprache, Inhaltstyp, Revision und Veröffentlichungsabsicht. Der eigentliche native Inhalt bleibt Text und darf nicht ausgeführt werden. Nur UTF-8 und Größenlimits akzeptieren.

Beim Desktop-Import sind Vertrauensquelle und Änderungen zu prüfen: `download -> anzeigen -> validieren -> Diff -> bestätigen -> lokale Speicherung`. Ein importiertes Objekt darf niemals automatisch überschreiben, Agenten ausführen oder ohne ausdrückliche Auswahl neue Aufgaben starten.

## Backend-Konzept

Eine mögliche Anfangsarchitektur:
- **Frontend**: öffentliche responsive Website mit Katalog, Volltextsuche, Sprach-/Provider-Filtern, Versionsanzeige und Download.
- **REST-API**: JSON über HTTPS, versioniert unter `/api/v1`; öffentliche lesende Endpunkte, Schreibzugriff ausschließlich nach Anmeldung.
- **Persistenz**: relationale Datenbank (z. B. PostgreSQL); revisionssichere Inhalte, Bewertungen, Moderation, Provenienz, Meldungen und optionale Prüfbelege.
- **Dateiablage**: unveränderliche Revisionen über Content-Hash; gerenderte Vorschau und Rohtext separat, kein ausführbarer Upload.
- **Betrieb**: Caddy als TLS-Reverse-Proxy, isolierter Anwendungsdienst, Backups, Rate Limits, Audit-Logging und Monitoring.

Die Implementierung des Backend-Stacks bleibt offen. Go oder Symfony sind beide möglich; das öffentliche Portal ist **noch nicht deployed** und besitzt **noch keine konfigurierte Domain**.

## API-Entwurf (noch nicht implementiert)

| Methode | Pfad | Verhalten |
| --- | --- | --- |
| GET | `/api/v1/items` | Katalogfilter: `type`, `language`, `provider`, `q`, `license`, `page` |
| GET | `/api/v1/items/{id}` | Metadaten, Status, aktuelle Revision |
| GET | `/api/v1/items/{id}/versions/{version}` | Unveränderliche Revision |
| POST | `/api/v1/items` | Entwurf mit expliziter Bestätigung hochladen |
| POST | `/api/v1/items/{id}/versions` | Neue Revision statt stiller Überschreibung |
| POST | `/api/v1/items/{id}/reports` | Missbrauch/Geheimnis melden |
| POST | `/api/v1/items/{id}/reviews` | Moderierte Bewertung und Prüfnachweise |

Die Routen sind ein Entwicklungsvertrag, kein aktiver Webdienst.

## Veröffentlichung und Sicherheit

- **Opt-in und Vorschau** vor jeder Veröffentlichung. Niemals automatisch Dateien aus einem Projekt uploaden.
- Projektpfade, Hostnamen, personenbezogene Daten, Tokens und Secrets vorab mit Regeln und manueller Prüfung aufspüren. Automatische Prüfung ersetzt keinen Nutzerentscheid.
- Sichere Textdarstellung ohne ungeprüftes HTML/Markdown; keine Ausführung von Konfigurationen oder Skripten.
- Transparente Lizenz pro Artefakt und verifizierbare Herkunft/Version; kein automatisches Kopieren von fremden Inhalten unter neuer Lizenz.
- Melde- und Moderationsprozess, Quarantäne problematischer Inhalte, Revisionshistorie und Löschkonzept.
- Begrenzte Upload-Größe, Rate Limits, CSRF-/Authentifizierungs-/Autorisierungsprüfungen, kein ungefilterter HTML-Render.
- Trennung von **verifiziert** (Prüfnachweise vorhanden) und **nicht überprüft** (Community-Inhalt).
- Datenschutzkonzept, Impressum, Nutzungsbedingungen und Community-Regeln vor öffentlichem Start erstellen.

## Meilensteine

**A. Lokal (jetzt):** Bibliothek, Klassifikation, Diff-Vorschau und ausdrückliche Übernahme in ein Projekt.  
**B. Portables Austauschformat:** Schema, deterministische Serialisierung, JSON-Validierung, Import-/Export-Tests.  
**C. Lesendes Portal:** Katalog, Suche, Detailseiten, Versionen, anonyme Downloads.  
**D. Veröffentlichungsfunktion:** Accounts, Upload-Quarantäne, Review, Redaktionsworkflow.  
**E. Desktop-Synchronisation (optional):** manuelle Suche und gezielter Download ohne automatische Schreibzugriffe.

## Noch zu entscheiden

Domainname, Hosting und Betreiberangaben; Nutzungsmodell und Lizenzen; Moderation und Veröffentlichungsrechte; Stack für das öffentliche Backend; Umfang an Community-Funktionen. Keine dieser Entscheidungen blockiert die lokale AgentForge-Entwicklung.
