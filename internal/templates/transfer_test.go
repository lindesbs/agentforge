package templates

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/lindesbs/agentforge/internal/editor"
	"github.com/lindesbs/agentforge/internal/inspector"
)

func TestTaggedPortableRoundTrip(t *testing.T) {
	project := t.TempDir()
	content := "# Keep comments\r\nmodel = 'test'\r\n[unknown]\r\nvalue = true\r\n"
	write(t, project, ".codex/config.toml", content)
	library, _ := New(t.TempDir())
	item, err := library.CreateTagged(project, ".codex/config.toml", " Team defaults ", []string{" Go ", "review", "go"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(item.Tags, []string{"go", "review"}) {
		t.Fatalf("tags: %#v", item.Tags)
	}
	payload, err := library.Export(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(payload, item.ID) || strings.Contains(payload, project) {
		t.Fatal("local identity exported")
	}
	destination := filepath.Join(t.TempDir(), "not-created")
	review, err := ReviewImport(payload)
	if err != nil {
		t.Fatal(err)
	}
	if review.Template.Content != content {
		t.Fatal("round trip changed native content")
	}
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatal("review had side effects")
	}
	other, _ := New(destination)
	for i := 0; i < 2; i++ {
		imported, err := other.Import(payload)
		if err != nil || imported.ID == item.ID {
			t.Fatalf("%+v %v", imported, err)
		}
	}
	entries, err := other.List()
	if err != nil || len(entries) != 2 || entries[0].ID == entries[1].ID {
		t.Fatalf("%+v %v", entries, err)
	}
	if !reflect.DeepEqual(entries[0].Tags, item.Tags) {
		t.Fatal("tags lost")
	}
	info, err := os.Stat(filepath.Join(destination, entries[0].ID+".json"))
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatal("import permissions")
	}
}

func TestInvalidImportNeverCreatesLibrary(t *testing.T) {
	base := Envelope{"agentforge-template", 1, PortableTemplate{Name: "Example", Provider: "codex", Kind: "settings", SourcePath: ".codex/config.toml", Content: "model = 'ok'", Tags: []string{}}}
	cases := []func(*Envelope){
		func(e *Envelope) { e.Version = 2 },
		func(e *Envelope) { e.Template.SourcePath = "../AGENTS.md" },
		func(e *Envelope) { e.Template.SourcePath = "/AGENTS.md" },
		func(e *Envelope) { e.Template.Provider = "claude" },
		func(e *Envelope) { e.Template.Content = "model = [private-value" },
		func(e *Envelope) { e.Template.Content = strings.Repeat("x", maxTemplateSize+1) },
		func(e *Envelope) { e.Template.Tags = []string{strings.Repeat("x", 33)} },
	}
	for i, change := range cases {
		e := base
		change(&e)
		data, _ := json.Marshal(e)
		root := filepath.Join(t.TempDir(), "library")
		library, _ := New(root)
		if _, err := library.Import(string(data)); err == nil {
			t.Errorf("case %d accepted", i)
		} else if strings.Contains(err.Error(), "private-value") {
			t.Fatal("parser value leaked")
		}
		if _, err := os.Stat(root); !os.IsNotExist(err) {
			t.Fatalf("case %d wrote before validation", i)
		}
	}
	for _, payload := range []string{"{} {}", "{\"credential\":\"private-value\"}", strings.Repeat("x", maxEnvelopeSize+1)} {
		if _, err := ReviewImport(payload); err == nil || strings.Contains(err.Error(), "private-value") {
			t.Fatalf("unsafe review error: %v", err)
		}
	}
}

func TestLegacyTemplatesAndEscapedContent(t *testing.T) {
	project := t.TempDir()
	write(t, project, "AGENTS.md", strings.Repeat("\t", 600000))
	library, _ := New(t.TempDir())
	item, err := library.CreateFromProject(project, "AGENTS.md", "Large escaped template")
	if err != nil {
		t.Fatal(err)
	}
	item.Tags = nil
	data, _ := json.Marshal(item)
	if err := os.WriteFile(filepath.Join(library.root, item.ID+".json"), data, 0600); err != nil {
		t.Fatal(err)
	}
	got, err := library.Get(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Content != item.Content || got.Tags == nil {
		t.Fatal("legacy content or tags lost")
	}
	if _, err := library.Export(item.ID); err != nil {
		t.Fatal(err)
	}
}

func TestFirstUsableReleaseWorkflow(t *testing.T) {
	project := t.TempDir()
	relative := ".claude/agents/reviewer.md"
	original := "---\nname: reviewer # identity\ndescription: Reviews code\nfuture: true\n---\nKeep body.\n"
	write(t, project, relative, original)
	service := inspector.New()
	discovered, err := service.Inspect(project)
	if err != nil || len(discovered.Files) != 1 {
		t.Fatalf("%+v %v", discovered, err)
	}
	fields, err := editor.ReadAgentFields(project, relative)
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := editor.PrepareAgentFields(project, relative, fields.Hash, "auditor", "Checks changes")
	if err != nil {
		t.Fatal(err)
	}
	preview, err := editor.PreviewChange(project, relative, fields.Hash, prepared.Content)
	if err != nil || !preview.Changed {
		t.Fatalf("%+v %v", preview, err)
	}
	saved, err := editor.Save(project, relative, fields.Hash, prepared.Content)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(saved.Content, "# identity") || !strings.Contains(saved.Content, "future: true") {
		t.Fatal("source lost")
	}
	library, _ := New(t.TempDir())
	item, err := library.CreateTagged(project, relative, "Auditor", []string{"review"})
	if err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir()
	write(t, destination, relative, original)
	proposed, err := library.PrepareApply(item.ID, destination, relative)
	if err != nil {
		t.Fatal(err)
	}
	unchanged, _ := editor.Read(destination, relative)
	if unchanged.Content != original {
		t.Fatal("template applied without save")
	}
	if _, err := editor.PreviewChange(destination, relative, proposed.Hash, proposed.Content); err != nil {
		t.Fatal(err)
	}
	result, err := editor.Save(destination, relative, proposed.Hash, proposed.Content)
	if err != nil || result.Content != saved.Content {
		t.Fatalf("%+v %v", result, err)
	}
}
