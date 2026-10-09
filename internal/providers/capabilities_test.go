package providers_test

import (
	"github.com/lindesbs/agentforge/internal/providers"
	"strings"
	"testing"
)

func TestCapabilitiesMatchSupportedDefinitions(t *testing.T) {
	files := []string{"AGENTS.md", ".codex/config.toml", "CLAUDE.md", ".claude/settings.json", ".claude/agents/reviewer.md"}
	capabilities := providers.Capabilities()
	if len(capabilities) != len(files) {
		t.Fatal("capability model incomplete")
	}
	for _, file := range files {
		definition := providers.Inspect(file, "")
		found := false
		for _, capability := range capabilities {
			if capability.Provider == definition.Provider && capability.Kind == definition.Kind {
				found = true
				if capability.Format == "" || capability.Validation == "" || capability.StructuredFields == nil {
					t.Fatal("missing capability contract")
				}
			}
		}
		if !found {
			t.Fatalf("missing %s", file)
		}
	}
}

func TestPortabilityReferencesAreAdvisoryAndRedacted(t *testing.T) {
	for _, reference := range []string{"../private.toml", "/private/config.toml", "C:\\private\\config.toml", "~/.codex/config.toml", "${PRIVATE}/config.toml", "agents/reviewer.toml"} {
		issues := providers.Portability(".codex/config.toml", "[agents.reviewer]\nconfig_file = '"+reference+"'\n")
		found := false
		for _, issue := range issues {
			if strings.Contains(issue.Message, reference) {
				t.Fatal("reference leaked in diagnostic")
			}
			if issue.Code == "unbundled-reference" || issue.Code == "nonportable-reference" {
				found = true
				if issue.Severity != "warning" {
					t.Fatal("portability should be advisory")
				}
			}
		}
		if !found {
			t.Fatalf("missing reference diagnostic for %s", reference)
		}
	}
}
