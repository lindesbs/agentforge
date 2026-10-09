# agentforge.art — Website- und Community-Konzept

> Status: geplantes Produkt und vorgesehene Domain; Registrierung, Markenrechte und tatsächliche Domainverfügbarkeit sind **nicht** bestätigt. Dieses Dokument ist ein Produktkonzept, keine Meldung über einen gestarteten Webdienst.

## 1. Vision und Positionierung

**AgentForge — The Art of Building Better AI Agents**

**Kernversprechen:** Vertrauenswürdige Konfigurationen, wiederverwendbare Agentenvorlagen und überprüfbares Entwicklungswissen entdecken, vergleichen, verwalten und mit ausdrücklicher Bestätigung in eigene Projekte übernehmen.

Zwei getrennte Produkte mit gemeinsamem Austauschformat:

- **AgentForge Desktop:** Lokal betriebener Konfigurationsmanager für Codex/Claude Code, offline nutzbar. Er verwaltet und validiert Dateien, zeigt Diffs und schreibt nur nach ausdrücklicher Bestätigung. **Keine Agentenausführung**.
- **AgentForge Exchange (agentforge.art):** Website, Katalog, Dokumentation, Community und später optionale API für Veröffentlichung und Import. Keine automatische Übertragung privater Projektinformationen.

**Zielgruppen:** Entwickler, Teams, Agent-/Prompt-Autoren, Open-Source-Maintainer, Auditoren sowie Organisationen, die verlässliche KI-Entwicklungsregeln versioniert dokumentieren wollen.

## 2. Informationsarchitektur der Website

| Pfad (Vorschlag) | Bereich | Zweck und Einstieg |
| --- | --- | --- |
| `/` | Startseite | Nutzenversprechen, Suche, aktuelle kuratierte Inhalte, Desktop-Download |
| `/desktop` | Desktop | Wails-Anwendung, Linux-Download, Changelog, lokale Sicherheit, Screenshots |
| `/templates` | Agent Templates | Suchbarer Katalog nach Anbieter, Sprache, Framework, Version, Lizenz |
| `/templates/{slug}` | Template-Detail | Beschreibung, Datei-Vorschau, Kompatibilität, Versionen, Herkunft, Lizenz |
| `/knowledge` | Erkenntnisbibliothek | Projekt-/Allgemein- und OK-/NOK-Filter, Sprachen, Kontrollen |
| `/knowledge/{slug}` | Erkenntnis-Detail | Ursache/Erkenntnis, Regel, Regressionstest, Revisionen, Belege |
| `/collections` | Stack-Sammlungen | Zusammengehörige Templates und Regeln, z. B. Go, Symfony, Contao |
| `/docs` | Dokumentation | Installation, Adapter, Dateiformate, Workflow, API und Sicherungsmodell |
| `/contribute` | Beitragen | Veröffentlichung, Review, Lizenzen, Qualitäts- und Moderationsrichtlinien |
| `/community` | Community | Diskussion, Fragen, Erfahrungsberichte, Vorschläge |
| `/about` | Über AgentForge | Abgrenzung zur Agent-Ausführung, Open-Source-Prinzipien, Roadmap |
| `/legal/*` | Rechtliches | Impressum, Datenschutz, Nutzungsbedingungen, Lizenzen |
| `/account` | Profil (später) | Eigene Veröffentlichungen, Entwürfe, Abos, gespeicherte Inhalte |

Navigation für MVP: **Templates · Wissen · Sammlungen · Desktop · Dokumentation**. Wenige klare Einstiege, mobile Bedienbarkeit und barrierefreie Tastaturnavigation.

## 3. Startseite: Inhalte und Seitenaufbau

1. **Hero:** „Bewährte KI-Agenten-Konfigurationen und Entwicklungswissen – transparent, versioniert und wiederverwendbar.“
2. **Suche:** Gemeinsame Suche über Vorlagen und Erkenntnisse mit Sprach-/Framework-Filtern und sofort verständlicher Ergebnisart.
3. **Zwei Hauptaktionen:** „Vorlagen entdecken“ und „AgentForge für Linux herunterladen“.
4. **Drei Nutzenbausteine:** Native Dateien behalten; Änderungen prüfen; aus bestätigten Lösungen lernen.
5. **Kuratierte Sammlungen:** Go, PHP/Symfony, TypeScript/Vue, Python; zunächst nur tatsächlich vorhandene Inhalte zeigen.
6. **Vertrauenssignale:** Lizenz, Quelle, Autor, Versionsdatum, Prüfnachweise, Kompatibilitätsstatus.
7. **Community-Einstieg:** „Eigene Vorlage einreichen“ mit ausdrücklicher Veröffentlichungsfreigabe.
8. **FAQ:** Führt AgentForge Agenten aus? Werden Dateien hochgeladen? Was bedeutet „bestätigt“?

**Tonalität:** Sachlich, entwicklerorientiert, international nutzbar; deutsche und englische Inhalte als priorisierte Sprachen.

## 4. Template-Katalog

**Mindestmetadaten:** Titel, Kurzbeschreibung, Autor/Organisation, Urheber-/Lizenzangaben, Provider (Codex/Claude/andere), Dateityp, unterstützte Provider-Versionen, Programmiersprache, Framework, Tags, Versionsnummer, Änderungsprotokoll und SHA-256 des Contents.

**Detailseite:** Lesbarer nativer Konfigurationstext mit Syntaxdarstellung, Änderungsvergleich zwischen Versionen, Einordnung (z. B. „getestet“ vs. „nicht verifiziert“), bekannte Einschränkungen, Hinweise zu möglichen Seiteneffekten und sicherer Download.

**Importfluss Desktop:** Suche → Metadaten prüfen → Inhalt lokal anzeigen → Provider-/Schema-Validierung → Diff gegen bestehende Datei → Nutzer bestätigt → gesicherte lokale Übernahme. Upload/Sync bleibt optional und standardmäßig ausgeschaltet.

## 5. Erkenntnisbibliothek

Vier verbindliche Kategorien:

- `PROJEKT_OK`: bewährte, projektspezifische Entscheidung.
- `PROJEKT_NOK`: projektspezifischer Fehler samt Ursache, Regel und Kontrolle.
- `ALLGEMEIN_OK`: projektunabhängige Best Practice.
- `ALLGEMEIN_NOK`: allgemeines Fehlerbild mit Vermeidungsregel.

**Zusatzattribute:** Sprache(n), Framework-/Tool-Version, Problemklasse, Status `unverified`/`community-reviewed`/`verified`, Quellen/Tests, Änderungsdatum, Geltungsbereich und ggf. Verfalls-/Reviewdatum.

**Darstellung:** Problem → bestätigte Ursache/Erkenntnis → verbindliche Regel → überprüfbare Kontrolle → Nachweise/Revisionen. Eintrag kann als Markdown-Vorschau in die passende der vier Projektdateien übernommen werden. Abweichende Projektarchitekturen müssen ausdrücklich bewertet werden.

**Wichtig:** Aus Community-Selbstauskunft darf nicht automatisch das Label „verifiziert“ folgen. Der Status braucht definierte Kriterien, prüfbare Nachweise und einen Reviewprozess.

## 6. Sammlungen und Lernpfade

Thematische Sammlungen kombinieren mehrere Dokumente statt ihre Inhalte blind zusammenzukopieren: z. B. „Go + Wails Desktop“, „Symfony/Contao“, „Vue 3 + TypeScript“ oder „Sicheres Konfigurationsmanagement“. Eine Sammlung enthält Voraussetzungen, Konflikte, kompatible Versionen, empfohlene Lesereihenfolge, Metadaten und Änderungshistorie.

Später optional: Abhängigkeitsgraph, der zeigt, welche Regeln und Templates zusammenpassen oder sich widersprechen. **Keine Ausführung von Workflows**.

## 7. Redaktions- und Community-Funktionen

**MVP:** Öffentlicher Lesekatalog, einreichbare Entwürfe, redaktionelle Freigabe, Versionshistorie, Inhaltsmeldung und dokumentierte Beiträge. Registrierung erst dort verlangen, wo sie nötig ist.

**Nach MVP:** Bewertungen mit konkreten Begründungen, Diskussionen, Follow-/Benachrichtigungsfunktionen, öffentliche Profile, kuratierte Sammlungen und API-Zugriffstoken.

**Moderation:** Rechteklärung, Secrets-/PII-Screening vor Veröffentlichung, Meldesystem, Quarantäne und Widerruf/De-Publikation. Änderungen an öffentlich geteilten Inhalten sind als neue Revision nachverfolgbar.

**Keine Gamification, die Quantität oder ungeprüfte „Best Practices“ über Qualität belohnt.**

## 8. Technische Architektur (Vorschlag)

- **Website:** responsive, serverseitig oder statisch vorgerenderte Katalogseiten für gute Indexierbarkeit. Frontend kann Vue 3 bleiben; Rendering-Konzept gesondert entscheiden.
- **Backend:** getrennte Exchange-API, wahlweise Go oder Symfony. Entscheidung erst nach API-/Datenmodell-Review; keine Kopplung des Desktop-Kerns an HTTP.
- **Datenhaltung:** PostgreSQL für Nutzer, Metadaten, Reviewstände und revisionsfeste Beiträge. Inhalte unveränderlich speichern; optional Objektablage für große Anhänge.
- **Suche:** zunächst PostgreSQL-Volltextsuche; später dedizierter Index bei Bedarf.
- **Betrieb:** Linux, Caddy für HTTPS, Container oder systemd, migrationsbasierte Deployments, Backup/Restore-Tests und Monitoring.
- **API:** `/api/v1` mit öffentlichem Lesen und authentifiziertem, limitiertem Schreiben. OpenAPI-Spezifikation, Pagination und ETags.
- **Austauschformat:** `schemas/exchange-item-v1.schema.json` als Ausgangspunkt; Lizenz, Herkunft, Version, Sprache, Geltungsbereich und Integrität vor einem echten Import/Export implementieren.
- **Sicherheit:** nie Remote-Agenten-Code ausführen; CSP, sichere Markdown-Ausgabe, CSRF-/Autorisierungschecks, Uploadgrenzen, Rate-Limits, sichere Session-Verwaltung, nachvollziehbare Moderation.
- **Datenschutz:** keine Telemetrie oder automatischen Projekt-Uploads; Einwilligung für öffentliche Beiträge, minimale Datenerhebung, Lösch- und Exportfunktion.

## 9. SEO, Barrierefreiheit und Auffindbarkeit

- Verständliche, eindeutige Texte und beschreibende URLs; sprechende Filterseiten nur bei eigenständigem Nutzen indexieren.
- Semantisches HTML, Tastaturbedienung, Fokuszustände, verständliche Form-Fehler und WCAG-2.2-AA als Gestaltungsziel.
- XML-Sitemap, Canonicals, hreflang für DE/EN, strukturierte Daten nur entsprechend sichtbarem Inhalt und passenden Schema.org-Typen.
- Klare Anbieter-/Versionsinformationen und maschinenlesbare Lizenzangaben; keine versteckten KI-Inhalte oder Cloaking.
- Performancebudget, minimierte Drittanbieter-Skripte, Datenschutz-freundliche Analytik, redaktionelle Qualitätskontrolle.

## 10. Designrichtung

**Werkstatt statt Bot-Marktplatz:** ruhiges, zugängliches Interface mit präzisen Code-/Dokumentenansichten, klaren Kategorien und visueller Nachvollziehbarkeit. Bildsprache von „Forge/Werkstatt“ sparsam einsetzen, keine austauschbaren Roboter-Stockbilder. Hell-/Dunkelmodus mit ausreichend Kontrast.

Beispielclaim: **“The Art of Building Better AI Agents.”**  
Alternative: **“Templates, Knowledge, Confidence.”**

## 11. Veröffentlichung und Domain

`agentforge.art` ist **Arbeitstitel und gewünschte Domain**, nicht als registriert bestätigt.

Vor öffentlicher Nutzung:
1. Verfügbarkeit und Renewal-Kosten der Domain überprüfen.
2. Unabhängige Namens-/Markenrecherche (DPMA, EUIPO und internationale Nutzung), insbesondere wegen anderer „AgentForge“-Projekte und deren Nutzung von `agentforge.ai`.
3. Logo-/Wortmarke und Lizenzrechte klären; keinerlei Zugehörigkeit zu fremden AgentForge-Projekten suggerieren.
4. Betreiberangaben und Rechtstexte erstellen; DNS/TLS, Backups und Monitoring vorbereiten.

## 12. Umsetzungsphasen und Abnahmekriterien

| Phase | Lieferumfang | Abnahmekriterium |
| --- | --- | --- |
| 0 — Konzept | Informationsarchitektur, Rollenmodell, Daten- und API-Entwurf, Markenprüfung | Freigegebener Umfang, Risiken und Nutzungsrechte dokumentiert |
| 1 — Öffentliche Website | Startseite, Docs, Linux-Download, statischer/lesender Katalog | Mobile/Screenreader-Checks, PageSpeed, saubere Versionierung und Downloads |
| 2 — Suche und Austausch | Kategorien, Filter, Template-Detail, Lernprotokoll, versioniertes Exportformat | Suche und Importvorschau funktionieren ohne schreibende Seiteneffekte |
| 3 — Beiträge und Review | Registrierung, Drafts, Moderation, Lizenzen, Missbrauchsmeldung | Keine Veröffentlichung ohne aktive Zustimmung und Freigabe |
| 4 — Desktop-Anbindung | Optionaler Katalogzugriff, gezielter Download, lokale Diff-Vorschau | Offline-Modus bleibt voll funktionsfähig; keine Auto-Uploads |

## 13. Nicht-Ziele für den Erststart

- Kein Hosting oder Ausführen fremder KI-Agenten.
- Keine automatische Übernahme unüberprüfter Regeln.
- Kein heimlicher Upload lokaler `AGENTS.md` oder Projektverzeichnisse.
- Keine komplexe soziale Plattform vor funktionierendem Katalog, Rechteprüfung und Moderation.
- Keine Gleichsetzung des Community-Feedbacks mit geprüfter technischer Qualität.

## Referenzen

- Eigenes Architekturkonzept: `docs/exchange-platform.md`
- Lernprotokoll: `docs/learning-protocol.md`
- Austauschschema: `schemas/exchange-item-v1.schema.json`
- Inspiration/Abgrenzung: https://github.com/AgentForge/AgentForge — separates externes Framework, kein Teil dieses Projekts.
