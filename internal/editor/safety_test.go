package editor

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestScalarPreservesCommentsCRLFAndUnknownFields(t *testing.T) {
	for _, source := range []string{"'Sam''s # reviewer'", "\"Sam's # reviewer\"", "Sam's reviewer"} {
		original := "---\r\n# heading\r\nname: " + source + "  # keep name comment\r\ndescription: Reviews code # purpose\r\ncustom: {future: true}\r\n---\r\nUnchanged body\r\n"
		next, err := updateScalar(original, "name", "new 'reviewer' #1")
		if err != nil {
			t.Fatal(err)
		}
		expected := strings.Replace(original, "name: "+source, "name: 'new ''reviewer'' #1'", 1)
		if next != expected {
			t.Fatalf("unexpected rewrite\ngot: %q\nwant: %q", next, expected)
		}
		got, err := readScalar(next, "name")
		if err != nil || got != "new 'reviewer' #1" {
			t.Fatalf("%q %v", got, err)
		}
	}
}

func TestUnsafeStructuredFormsAreRejected(t *testing.T) {
	for _, source := range []string{"name: >\n  folded", "name: 'multi\n  line'", "name: &shared reviewer", "name: [reviewer]", "name: first\nname: second", "name: !!str reviewer"} {
		if _, err := updateScalar("---\n"+source+"\n---\nbody", "name", "replacement"); err == nil {
			t.Errorf("accepted %q", source)
		}
	}
}

func TestProviderDraftValidationAndRoundTrip(t *testing.T) {
	for _, tc := range []struct{ path, original, updated, invalid string }{
		{".codex/config.toml", "# keep\r\nmodel = 'old' # inline\r\n[future]\r\nunknown = [1, 2]\r\n", "# keep\r\nmodel = 'new' # inline\r\n[future]\r\nunknown = [1, 2]\r\n", "model = [secret-value"},
		{".claude/agents/reviewer.md", "---\nname: reviewer\ndescription: old # keep\nfuture: true\n---\nBody\n", "---\nname: reviewer\ndescription: new # keep\nfuture: true\n---\nBody\n", "---\nname: [secret-value\n---\nBody"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			root := t.TempDir()
			name := filepath.Join(root, tc.path)
			if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(name, []byte(tc.original), 0640); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(name, 0640); err != nil {
				t.Fatal(err)
			}
			doc, err := Read(root, tc.path)
			if err != nil {
				t.Fatal(err)
			}
			for _, issue := range Validate(tc.path, tc.invalid) {
				if strings.Contains(issue.Message, "secret-value") {
					t.Fatal("source leaked")
				}
			}
			if _, err := Save(root, tc.path, doc.Hash, tc.invalid); err == nil {
				t.Fatal("invalid draft saved")
			}
			if _, err := PreviewChange(root, tc.path, doc.Hash, tc.updated); err != nil {
				t.Fatal(err)
			}
			unchanged, _ := os.ReadFile(name)
			if string(unchanged) != tc.original {
				t.Fatal("preview wrote data")
			}
			saved, err := Save(root, tc.path, doc.Hash, tc.updated)
			if err != nil {
				t.Fatal(err)
			}
			if saved.Content != tc.updated {
				t.Fatal("native text changed")
			}
			info, _ := os.Stat(name)
			if info.Mode().Perm() != 0640 {
				t.Fatal("mode lost")
			}
			backups, _ := filepath.Glob(filepath.Join(filepath.Dir(name), ".agentforge-backup-*"))
			if len(backups) != 1 {
				t.Fatalf("backups: %v", backups)
			}
			before, _ := os.ReadFile(backups[0])
			if string(before) != tc.original {
				t.Fatal("backup changed")
			}
		})
	}
}

func TestConcurrentSavesRejectStaleWriter(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	doc, err := Read(root, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, content := range []string{"first", "second"} {
		wg.Add(1)
		go func(value string) { defer wg.Done(); _, err := Save(root, doc.Path, doc.Hash, value); results <- err }(content)
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("%d concurrent writes succeeded", success)
	}
}
