package inspector

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func write(t *testing.T, root, path string) {
	t.Helper()
	name := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestInspect(t *testing.T) {
	root := t.TempDir()
	write(t, root, "composer.json")
	write(t, root, "AGENTS.md")
	write(t, root, ".claude/agents/reviewer.md")
	p, err := New().Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Frameworks) != 1 || p.Frameworks[0] != "PHP/Composer" {
		t.Fatalf("frameworks: %#v", p.Frameworks)
	}
	if len(p.Files) != 2 {
		t.Fatalf("files: %#v", p.Files)
	}
	if p.Files[0].Path != ".claude/agents/reviewer.md" {
		t.Fatalf("order: %#v", p.Files)
	}
}
func TestRejectEmptyAndFile(t *testing.T) {
	if _, err := New().Inspect(""); err == nil {
		t.Fatal("expected error")
	}
	root := t.TempDir()
	write(t, root, "plain.txt")
	if _, err := New().Inspect(filepath.Join(root, "plain.txt")); err == nil {
		t.Fatal("expected error")
	}
}
func TestIgnoreSymlink(t *testing.T) {
	root := t.TempDir()
	write(t, root, "outside.md")
	if err := os.MkdirAll(filepath.Join(root, ".claude", "agents"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "outside.md"), filepath.Join(root, ".claude", "agents", "escape.md")); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	p, err := New().Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Files) != 0 {
		t.Fatalf("unexpected symlink file: %#v", p.Files)
	}
}

func TestIgnoreSymlinkDirectories(t *testing.T) {
	for _, relative := range []string{".codex", ".claude", ".claude/agents"} {
		t.Run(relative, func(t *testing.T) {
			root, outside := t.TempDir(), t.TempDir()
			write(t, root, "AGENTS.md")
			// Every possible descendant is outside the selected project.
			write(t, outside, "config.toml")
			write(t, outside, "settings.json")
			write(t, outside, "settings.local.json")
			write(t, outside, "agents/reviewer.md")
			write(t, outside, "reviewer.md")
			link := filepath.Join(root, relative)
			if err := os.MkdirAll(filepath.Dir(link), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(outside, link); err != nil {
				t.Skipf("symlink unsupported: %v", err)
			}
			project, err := New().Inspect(root)
			if err != nil {
				t.Fatal(err)
			}
			want := []ConfigFile{{Path: "AGENTS.md", Provider: "codex", Kind: "instructions"}}
			if !reflect.DeepEqual(project.Files, want) {
				t.Fatalf("discovery escaped through %s: %#v", relative, project.Files)
			}
			if len(project.Diagnostics) != 1 || project.Diagnostics[0].Path != relative || project.Diagnostics[0].Code != "symlink" {
				t.Fatalf("missing symlink diagnostic: %#v", project.Diagnostics)
			}
			data, err := json.Marshal(project)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(data), outside) {
				t.Fatal("response disclosed the symlink target")
			}
		})
	}
}

func TestInspectAllSupportedPathsReadOnly(t *testing.T) {
	root := t.TempDir()
	wantFiles := []ConfigFile{
		{".claude/agents/Reviewer.MD", "claude", "agent"},
		{".claude/agents/writer.md", "claude", "agent"},
		{".claude/settings.json", "claude", "settings"},
		{".claude/settings.local.json", "claude", "settings"},
		{".codex/config.toml", "codex", "settings"},
		{"AGENTS.md", "codex", "instructions"},
		{"CLAUDE.md", "claude", "instructions"},
	}
	paths := []string{"composer.json", "go.mod", "package.json", ".claude/agents/notes.txt", "nested/AGENTS.md", ".claude/agents/nested/hidden.md"}
	for _, file := range wantFiles {
		paths = append(paths, file.Path)
	}
	for _, path := range paths {
		write(t, root, path)
	}
	project, err := New().Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	physical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if project.Root != physical || project.Name != filepath.Base(physical) {
		t.Fatalf("project metadata: %#v", project)
	}
	if !reflect.DeepEqual(project.Frameworks, []string{"PHP/Composer", "Go", "Node.js"}) {
		t.Fatalf("frameworks: %#v", project.Frameworks)
	}
	if !reflect.DeepEqual(project.Files, wantFiles) {
		t.Fatalf("files: %#v", project.Files)
	}
	if len(project.Diagnostics) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", project.Diagnostics)
	}
	// Discovery must not parse, execute or rewrite the arbitrary fixture text.
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "test" {
			t.Fatalf("discovery changed %s", path)
		}
	}
}

func TestEmptyProjectUsesEmptyArrays(t *testing.T) {
	project, err := New().Inspect(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(project)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"frameworks", "files", "diagnostics"} {
		if !strings.Contains(string(data), `"`+field+`":[]`) {
			t.Fatalf("%s must serialize as an empty array: %s", field, data)
		}
	}
}

func TestInvalidProjectRoots(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{"", " \t", filepath.Join(root, "missing")} {
		if _, err := New().Inspect(path); err == nil {
			t.Fatalf("accepted invalid root %q", path)
		}
	}
	link := filepath.Join(t.TempDir(), "project")
	if err := os.Symlink(root, link); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	if _, err := New().Inspect(link); err == nil {
		t.Fatal("accepted a symlink project root")
	}
}

func TestRootAncestorIsCanonicalized(t *testing.T) {
	parent := t.TempDir()
	write(t, parent, "project/AGENTS.md")
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(parent, alias); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	project, err := New().Inspect(filepath.Join(alias, "project"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(filepath.Join(parent, "project"))
	if err != nil {
		t.Fatal(err)
	}
	if project.Root != want || len(project.Files) != 1 {
		t.Fatalf("canonical project: %#v", project)
	}
}

func TestUnexpectedPathTypes(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".codex")
	write(t, root, ".claude/agents")
	write(t, root, "AGENTS.md/child")
	write(t, root, "CLAUDE.md")
	project, err := New().Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []Diagnostic{
		{".claude/agents", "not-directory", "warning", "Expected a directory; skipped."},
		{".codex", "not-directory", "warning", "Expected a directory; skipped."},
		{"AGENTS.md", "not-file", "warning", "Expected a regular file; skipped."},
	}
	if !reflect.DeepEqual(project.Diagnostics, want) {
		t.Fatalf("diagnostics: %#v", project.Diagnostics)
	}
	if len(project.Files) != 1 || project.Files[0].Path != "CLAUDE.md" {
		t.Fatalf("valid files lost: %#v", project.Files)
	}
}

func TestBrokenSymlinksAreReported(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude", "agents"), 0700); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"go.mod", "AGENTS.md", ".claude/agents/broken.md"} {
		if err := os.Symlink(filepath.Join(root, "absent"), filepath.Join(root, path)); err != nil {
			t.Skipf("symlink unsupported: %v", err)
		}
	}
	project, err := New().Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(project.Files) != 0 || len(project.Frameworks) != 0 || len(project.Diagnostics) != 3 {
		t.Fatalf("broken links: %#v", project)
	}
	for _, d := range project.Diagnostics {
		if d.Code != "symlink" {
			t.Fatalf("diagnostic: %#v", d)
		}
	}
}

func TestUnreadableAgentsPreserveOtherResults(t *testing.T) {
	root := t.TempDir()
	write(t, root, "AGENTS.md")
	write(t, root, ".claude/agents/hidden.md")
	agents := filepath.Join(root, ".claude", "agents")
	if err := os.Chmod(agents, 0000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(agents, 0700)
	if _, err := os.ReadDir(agents); err == nil {
		t.Skip("filesystem or current user does not enforce directory permissions")
	}
	project, err := New().Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(project.Files) != 1 || project.Files[0].Path != "AGENTS.md" {
		t.Fatalf("valid files lost: %#v", project.Files)
	}
	if len(project.Diagnostics) != 1 || project.Diagnostics[0].Path != ".claude/agents" || project.Diagnostics[0].Code != "unreadable" {
		t.Fatalf("diagnostics: %#v", project.Diagnostics)
	}
}
