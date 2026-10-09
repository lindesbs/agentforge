package inspector

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ConfigFile struct {
	Path     string `json:"path"`
	Provider string `json:"provider"`
	Kind     string `json:"kind"`
}

// Diagnostic describes discovery problems without including file contents or
// symlink targets. Path is always relative to the selected project.
type Diagnostic struct {
	Path     string `json:"path"`
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Project struct {
	Root        string       `json:"root"`
	Name        string       `json:"name"`
	Frameworks  []string     `json:"frameworks"`
	Files       []ConfigFile `json:"files"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Service struct{}

func New() *Service { return &Service{} }

// Inspect accepts a local project root and only discovers names and metadata.
// The selected root is canonicalized; symlinks below it are skipped. Checks
// assume the directory tree is not concurrently replaced during discovery.
func (s *Service) Inspect(root string) (Project, error) {
	if strings.TrimSpace(root) == "" {
		return Project{}, errors.New("project path is required")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return Project{}, err
	}
	info, err := os.Lstat(abs)
	if err != nil {
		return Project{}, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return Project{}, errors.New("symlink project roots are not supported")
	}
	if !info.IsDir() {
		return Project{}, errors.New("project path is not a directory")
	}
	// Resolve aliases in the explicitly selected root (e.g. /var on macOS)
	// once, so downstream reads use the same physical project directory.
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return Project{}, err
	}
	p := Project{Root: abs, Name: filepath.Base(abs), Frameworks: []string{}, Files: []ConfigFile{}, Diagnostics: []Diagnostic{}}
	warn := func(path, code, message string) {
		p.Diagnostics = append(p.Diagnostics, Diagnostic{Path: path, Code: code, Severity: "warning", Message: message})
	}
	// Call only after verifying the parent directory. Lstat on a whole path
	// alone would still follow symlinks in its intermediate components.
	matches := func(rel string, directory bool) bool {
		fi, e := os.Lstat(filepath.Join(abs, filepath.FromSlash(rel)))
		if os.IsNotExist(e) {
			return false
		}
		if e != nil {
			warn(rel, "unreadable", "Path could not be inspected; skipped.")
			return false
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			warn(rel, "symlink", "Symbolic link skipped; links are not followed inside the project.")
			return false
		}
		if directory {
			if fi.IsDir() {
				return true
			}
			warn(rel, "not-directory", "Expected a directory; skipped.")
		} else {
			if fi.Mode().IsRegular() {
				return true
			}
			warn(rel, "not-file", "Expected a regular file; skipped.")
		}
		return false
	}
	addFile := func(rel, provider, kind string) {
		if matches(rel, false) {
			p.Files = append(p.Files, ConfigFile{rel, provider, kind})
		}
	}
	for _, marker := range []struct{ path, label string }{
		{"composer.json", "PHP/Composer"}, {"go.mod", "Go"}, {"package.json", "Node.js"},
	} {
		if matches(marker.path, false) {
			p.Frameworks = append(p.Frameworks, marker.label)
		}
	}
	addFile("AGENTS.md", "codex", "instructions")
	addFile("CLAUDE.md", "claude", "instructions")
	if matches(".codex", true) {
		addFile(".codex/config.toml", "codex", "settings")
	}
	if matches(".claude", true) {
		addFile(".claude/settings.json", "claude", "settings")
		addFile(".claude/settings.local.json", "claude", "settings")
		if matches(".claude/agents", true) {
			entries, e := os.ReadDir(filepath.Join(abs, ".claude", "agents"))
			if e != nil {
				warn(".claude/agents", "unreadable", "Agent directory could not be listed; skipped.")
			} else {
				for _, entry := range entries {
					if strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
						addFile(".claude/agents/"+entry.Name(), "claude", "agent")
					}
				}
			}
		}
	}
	sort.Slice(p.Files, func(i, j int) bool { return p.Files[i].Path < p.Files[j].Path })
	sort.Slice(p.Diagnostics, func(i, j int) bool { return p.Diagnostics[i].Path < p.Diagnostics[j].Path })
	return p, nil
}
