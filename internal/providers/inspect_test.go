package providers_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/lindesbs/agentforge/internal/editor"
	"github.com/lindesbs/agentforge/internal/providers"
)

func hasIssue(d providers.Definition, code string) bool {
	for _, issue := range d.Issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func TestClaudeAgentFields(t *testing.T) {
	content := "---\nname: reviewer\ndescription: >-\n  Reviews code\n  with care.\nmodel: sonnet\ntools: [Read, Grep]\nunknown: preserve-me\n---\nPrivate prompt text.\n"
	d := providers.Inspect(".claude/agents/reviewer.md", content)
	want := []providers.Field{{"name", "reviewer"}, {"description", "Reviews code with care."}, {"model", "sonnet"}, {"tools", "Read, Grep"}}
	if d.Provider != "claude" || d.Kind != "agent" || !reflect.DeepEqual(d.Fields, want) {
		t.Fatalf("unexpected definition: %#v", d)
	}
	if len(d.Agents) != 1 || d.Agents[0].Name != "reviewer" || len(d.Issues) != 1 || !hasIssue(d, "limited-validation") {
		t.Fatalf("unexpected agents or issues: %#v", d)
	}
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	for _, hidden := range []string{"preserve-me", "Private prompt text"} {
		if strings.Contains(string(data), hidden) {
			t.Fatalf("unexpected content in summary: %s", hidden)
		}
	}
}

func TestClaudeFrontmatterVariants(t *testing.T) {
	for _, tc := range []struct{ name, content, description string }{
		{"quotes", "---\nname: reviewer\ndescription: 'Checks developer''s code' # comment\ntools: Read, Grep\n---\nPrompt", "Checks developer's code"},
		{"crlf-bom", "\ufeff---\r\nname: reviewer\r\ndescription: \"Checks \\\"quoted\\\" code\"\r\n---\r\nPrompt", "Checks \"quoted\" code"},
		{"literal", "---\nname: reviewer\ndescription: |-\n  First line\n  Second line\n---\nPrompt", "First line\nSecond line"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := providers.Inspect(".claude/agents/reviewer.MD", tc.content)
			if len(d.Agents) != 1 || d.Agents[0].Description != tc.description {
				t.Fatalf("description not decoded: %#v", d)
			}
		})
	}
}

func TestMalformedClaudeDefinitions(t *testing.T) {
	for _, tc := range []struct{ name, content, code string }{
		{"no-frontmatter", "# Prompt", "missing-frontmatter"},
		{"no-end", "---\nname: reviewer\ndescription: review", "unterminated-frontmatter"},
		{"syntax", "---\nname: [\n---\nPrompt", "invalid-yaml"},
		{"duplicate", "---\nname: first\nname: second\n---\nPrompt", "invalid-yaml"},
		{"sequence", "---\n- item\n---\nPrompt", "invalid-yaml"},
		{"trailing-yaml", "---\nname: reviewer\ndescription: Reviews\n...\ntrailing-value\n---\nPrompt", "invalid-yaml"},
		{"recursive-alias", "---\nname: reviewer\ndescription: Reviews\nunknown: &loop {self: *loop}\n---\nPrompt", "invalid-yaml"},
		{"missing", "---\nname: reviewer\n---\nPrompt", "missing-field"},
		{"empty-name", "---\nname: ''\ndescription: Reviews\n---\nPrompt", "invalid-field"},
		{"numeric-name", "---\nname: 123\ndescription: Reviews\n---\nPrompt", "invalid-field"},
		{"mapped-description", "---\nname: reviewer\ndescription: {nested: value}\n---\nPrompt", "invalid-field"},
		{"tools", "---\nname: reviewer\ndescription: Reviews\ntools: [Read, 12]\n---\nPrompt", "invalid-tools"},
		{"empty-prompt", "---\nname: reviewer\ndescription: Reviews\n---\n", "empty-prompt"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := providers.Inspect(".claude/agents/reviewer.md", tc.content)
			if !hasIssue(d, tc.code) {
				t.Fatalf("missing %s: %#v", tc.code, d)
			}
		})
	}
}

func TestCodexSettingsAndRoles(t *testing.T) {
	content := `# Preserve formatting and unknown settings.
model = "example-model"
model_provider = "example-provider"
model_reasoning_effort = "high"
sandbox_mode = "workspace-write"
api_key = "hidden-setting"
[agents]
max_threads = 4
[agents.writer]
description = "Writes code"
config_file = "agents/writer.toml"
[agents.reviewer]
description = "Reviews code"
`
	d := providers.Inspect(".codex/config.toml", content)
	want := []providers.Agent{{"reviewer", "Reviews code", ""}, {"writer", "Writes code", "agents/writer.toml"}}
	if d.Provider != "codex" || d.Kind != "settings" || len(d.Fields) != 4 || !reflect.DeepEqual(d.Agents, want) {
		t.Fatalf("unexpected Codex definition: %#v", d)
	}
	if len(d.Issues) != 1 || !hasIssue(d, "limited-validation") {
		t.Fatalf("issues: %#v", d.Issues)
	}
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "hidden-setting") {
		t.Fatal("unknown setting leaked")
	}
}

func TestCodexDiagnostics(t *testing.T) {
	if hasIssue(providers.Inspect(".codex/config.toml", "model = ''"), "invalid-field") {
		t.Fatal("optional string fields must not invent a non-empty schema constraint")
	}
	for _, tc := range []struct{ content, code string }{
		{"model = [", "invalid-toml"},
		{"model = 'first'\nmodel = 'second'", "invalid-toml"},
		{"model = 12", "invalid-field"},
		{"agents = 'invalid'", "invalid-agents"},
		{"[agents.reviewer]\ndescription = 12", "invalid-field"},
	} {
		d := providers.Inspect(".codex/config.toml", tc.content)
		if !hasIssue(d, tc.code) {
			t.Fatalf("missing %s: %#v", tc.code, d)
		}
	}
}

func TestClaudeSettingsDiagnostics(t *testing.T) {
	for _, content := range []string{`{`, `[]`, `null`, `{} {}`, `{"model":12}`} {
		d := providers.Inspect(".claude/settings.local.json", content)
		if !hasIssue(d, "invalid-json") && !hasIssue(d, "invalid-field") {
			t.Fatalf("invalid settings accepted: %#v", d)
		}
	}
	d := providers.Inspect(".claude/settings.json", `{"model":"sonnet","env":{"KEY":"private-setting"}}`)
	if !reflect.DeepEqual(d.Fields, []providers.Field{{"model", "sonnet"}}) {
		t.Fatalf("fields: %#v", d.Fields)
	}
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "private-setting") {
		t.Fatal("environment setting leaked")
	}
}

func TestDiagnosticsDoNotExposeSource(t *testing.T) {
	for _, tc := range []struct{ path, content string }{
		{".claude/agents/agent.md", "---\nname: reviewer\ndescription: [private-sentinel\n---\n"},
		{".claude/agents/agent.md", "---\nprivate-sentinel: a\nprivate-sentinel: b\n---\n"},
		{".codex/config.toml", "private-sentinel = ["},
		{".claude/settings.json", `{"private-sentinel":`},
	} {
		d := providers.Inspect(tc.path, tc.content)
		data, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "private-sentinel") {
			t.Fatal("parser diagnostic disclosed source")
		}
		if len(d.Issues) == 0 {
			t.Fatal("missing diagnostic")
		}
	}
}

func TestInstructionDocumentsAndUnsupportedPaths(t *testing.T) {
	for _, file := range []string{"AGENTS.md", "CLAUDE.md"} {
		d := providers.Inspect(file, "# Private instructions")
		if d.Kind != "instructions" || len(d.Fields) != 0 || len(d.Issues) != 0 {
			t.Fatalf("instructions: %#v", d)
		}
		if !hasIssue(providers.Inspect(file, " \n"), "empty-instructions") {
			t.Fatal("missing empty warning")
		}
	}
	for _, file := range []string{"other.md", ".claude/agents/nested/agent.md", "../.codex/config.toml"} {
		if !hasIssue(providers.Inspect(file, ""), "unsupported-path") {
			t.Fatalf("accepted %s", file)
		}
	}
	if providers.Inspect(`.claude\agents\reviewer.md`, "").Provider != "claude" {
		t.Fatal("Windows separators not recognized")
	}
}

func TestInspectionLimitsAndEmptyArrays(t *testing.T) {
	for _, content := range []string{strings.Repeat("x", 1024*1024+1), string([]byte{0xff})} {
		if !hasIssue(providers.Inspect("AGENTS.md", content), "invalid-document") {
			t.Fatal("invalid document accepted")
		}
	}
	data, err := json.Marshal(providers.Inspect("AGENTS.md", "instructions"))
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"fields", "agents", "issues"} {
		if !strings.Contains(string(data), `"`+key+`":[]`) {
			t.Fatalf("missing array %s: %s", key, data)
		}
	}
}

func TestInspectLoadedConfigPreservesSource(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".claude", "agents")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	content := "---\r\n# Preserve this comment\r\nname: reviewer\r\ndescription: Reviews code\r\nunknown: [value]\r\n---\r\nPrivate prompt\r\n"
	file := filepath.Join(dir, "reviewer.md")
	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	doc, err := editor.Read(root, ".claude/agents/reviewer.md")
	if err != nil {
		t.Fatal(err)
	}
	d := providers.Inspect(doc.Path, doc.Content)
	if len(d.Agents) != 1 {
		t.Fatalf("agent not inspected: %#v", d)
	}
	after, err := editor.Read(root, doc.Path)
	if err != nil {
		t.Fatal(err)
	}
	if after.Hash != doc.Hash || after.Content != content {
		t.Fatal("inspection changed the native file")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatal("read-only inspection created files")
	}
}
