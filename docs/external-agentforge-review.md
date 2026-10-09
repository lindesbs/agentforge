# Externes Repository: AgentForge/AgentForge — Ideenprüfung

**Quellprojekt:** https://github.com/AgentForge/AgentForge  
**Prüfstand:** 9. Oktober 2026, öffentliche README und Repository-Metadaten. Es besteht keine organisatorische oder technische Verbindung zu `lindesbs/agentforge`.

## Kerndifferenz

Das fremde `agentforge/agentforge` beschreibt ein Framework/API zum **Entwickeln, Orchestrieren und Ausführen** von KI-Agenten. Es verwendet einen Directed Acyclic Graph (DAG) aus Routinen und Subroutinen, gemeinsamem Kontext und ein Adapter-/Interface-Factory-Muster für APIs, Modelle und Dienste. Dokumentiert sind auch Speicher-/Reasoning-Bausteine, Docker-Services und ein Python/TypeScript-Stack.

**Unser AgentForge** ist eine lokale, nicht ausführende **Konfigurations- und Wissensverwaltung** mit optionaler künftiger Exchange-Website unter dem vorgesehenen Namen `agentforge.art`.

## Produktideen – Bewertung

| Idee aus dem fremden Projekt | Nutzen für unser Produkt | Entscheidung |
| --- | --- | --- |
| Modulare Subroutinen / komponierbare Bausteine | Vorlagen aus kleineren kompatiblen Regeln/Sammlungen zusammensetzen | Als **statische Komposition und Vorschau** vormerken |
| DAG/Graph-Visualisierung | Abhängigkeiten, Konflikte und Herkunft von Templates/Erkenntnissen sichtbar machen | **Visualisierung**, keine Workflow-Ausführung |
| Gemeinsamer Kontext | Metadaten zu Provider, Sprache, Stack, Gültigkeitsbereich, Tests, Versionen | In Bibliotheks-/Exchange-Datenmodell aufnehmen |
| Adapter/Interface-Factory | Entkopplung von Codex, Claude und künftig weiteren Providern | Weiterverfolgen als Provideradapter |
| Wissensspeicher | Projektübergreifende Suche, Kategorien und geprüfte Erkenntnisse | Zentraler Bestandteil von `agentforge.art` |
| Dockerisierte Services | Separater Web-Katalog, API, Datenbank und Worker für Moderation | Erst wenn Website-Backend implementiert wird |
| LLM-Inferenz / Reasoning / Agent-Runtime | Gehört nicht zum lokalen Konfigurationsmanager | **Nicht übernehmen** |
| Automatisches Ausführen von Plänen / Tools | Widerspricht dem Produkt-Sicherheitsziel | **Nicht übernehmen** |

## Lizenz- und Namensprüfung

Bei der Prüfung war im Repository-Stamm keine `LICENSE` sichtbar; ein Abruf von `LICENSE` über GitHub war nicht erfolgreich. Daraus darf keine verlässlich geprüfte Open-Source-Nutzungslizenz abgeleitet werden. **Keine Codeübernahme, kein Forken mit Übernahme in das Produkt und keine Wiederveröffentlichung fremder Assets**, bevor konkrete Lizenzrechte geklärt sind. Architektonische Ideen dürfen eigenständig entworfen und implementiert werden.

Die README verweist auf `agentforge.ai`. Der Produktname `AgentForge` wird außerdem von verschiedenen weiteren Projekten genutzt. `agentforge.art` ist daher **vor Kauf und Markenverwendung unabhängig auf kennzeichenrechtliche Konflikte zu prüfen**. GitHub-Repositorynamen und Domainregistrierung allein sind keine Markenfreigabe.

## Konkrete nächste Umsetzungsschritte

1. Austauschschema um `components`, `dependsOn` und `compatibility` erweitern, sobald konkrete Anwendungsfälle vorliegen.
2. Grafische Abhängigkeitsansicht für **statische** Daten entwerfen (keine Ausführungsschicht).
3. Komponenten-/Versionskonflikte in Import- und Diff-Vorschau markieren.
4. Rechte, Autor, Originalquelle und Lizenz pro Template/Erkenntnis dokumentieren.
5. Dokumentierte Nutzungsrechte für jeden importierten Community-Beitrag verlangen.
6. Community-Moderation und Nachweise/Reviewstatus vor automatischer Katalogübernahme etablieren.

**Quellen:** Öffentliches README https://github.com/AgentForge/AgentForge und Architektur-/Portalentwurf `docs/website-agentforge-art.md`. Die Aussagen über das fremde Repository sind eine zeitgebundene Momentaufnahme, keine kontinuierliche Überprüfung.
