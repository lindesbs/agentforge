# Projektbezogene Fehler und Vermeidungsregeln — AgentForge

Hier bestätigte, vom Projekt abhängige Fehler dokumentieren. Für jeden Eintrag sind Ursache, verbindliche Vermeidungsregel und nachprüfbare Kontrolle erforderlich.

## Unlesbare Dropdown-Werte im dunklen Desktop-Theme

**Beobachtung:** Auf Linux/GTK/WebKit erscheinen die nativen Auswahlfelder der Erkenntnisbibliothek mit hellem Hintergrund und ebenfalls hellem Text. Ein Screenshot aus der laufenden Anwendung zeigt den schlechten Kontrast.

**Ursache (technische Einordnung):** Das CSS setzte zwar dunkle Farben für `select`, legte aber kein natives `color-scheme` fest. Browser-/GTK-eigene Darstellung von Auswahlfeldern und Dropdown-Optionen kann deshalb einen hellen Systemfarbmodus verwenden.

**Verbindliche Regel:** Bei dunklen Desktop-WebViews `color-scheme: dark` sowie kontrastreiche Farben für `select`, `option` und `optgroup` setzen. Sichtbare Fokuszustände beibehalten.

**Prüfung:** Vue-Build und Linux-Desktop-Build ausführen; anschließend auf KDE/Wayland und GTK/WebKit per visueller Kontrolle sowohl den geschlossenen Wert als auch die geöffnete Optionsliste testen. Ein erfolgreicher Build allein ersetzt die Sichtprüfung nicht.

