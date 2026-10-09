package templates

import (
	"encoding/json"
	"errors"
	"io"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/lindesbs/agentforge/internal/providers"
)

// JSON escaping can expand a byte sixfold; metadata is separately bounded.
const maxEnvelopeSize = 6*maxTemplateSize + 8192

type PortableTemplate struct {
	Name       string   `json:"name"`
	Provider   string   `json:"provider"`
	Kind       string   `json:"kind"`
	SourcePath string   `json:"sourcePath"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
}
type Envelope struct {
	Format   string           `json:"format"`
	Version  int              `json:"version"`
	Template PortableTemplate `json:"template"`
}
type ImportReview struct {
	Template PortableTemplate  `json:"template"`
	Issues   []providers.Issue `json:"issues"`
}

func Summarize(item Template) Summary {
	return Summary{ID: item.ID, Name: item.Name, Provider: item.Provider, Kind: item.Kind, SourcePath: item.SourcePath, CreatedAt: item.CreatedAt, Tags: item.Tags}
}
func portable(item Template) PortableTemplate {
	return PortableTemplate{item.Name, item.Provider, item.Kind, item.SourcePath, item.Content, item.Tags}
}
func fromPortable(item PortableTemplate) Template {
	return Template{Name: item.Name, Provider: item.Provider, Kind: item.Kind, SourcePath: item.SourcePath, Content: item.Content, Tags: item.Tags}
}

func normalize(item *Template) error {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" || len([]rune(item.Name)) > 120 || !utf8.ValidString(item.Name) || strings.ContainsFunc(item.Name, unicode.IsControl) {
		return errors.New("template name must be 1–120 printable characters")
	}
	// Export paths are canonical project-relative hints, never destination paths.
	if item.SourcePath != path.Clean(item.SourcePath) || strings.Contains(item.SourcePath, "\\") || len(item.SourcePath) > 512 {
		return errors.New("template needs a canonical supported relative source path")
	}
	provider, kind, err := describe(item.SourcePath)
	if err != nil || item.Provider != provider || item.Kind != kind {
		return errors.New("template provider, kind and source path must match")
	}
	if len(item.Content) > maxTemplateSize || !utf8.ValidString(item.Content) {
		return errors.New("template content must be UTF-8, up to 1 MiB")
	}
	if len(item.Tags) > 12 {
		return errors.New("use at most 12 tags")
	}
	tags := make([]string, 0, len(item.Tags))
	seen := map[string]bool{}
	for _, tag := range item.Tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		if !utf8.ValidString(tag) || len([]rune(tag)) > 32 || strings.ContainsFunc(tag, unicode.IsControl) {
			return errors.New("tags must be printable and at most 32 characters")
		}
		if !seen[tag] {
			tags = append(tags, tag)
			seen[tag] = true
		}
	}
	item.Tags = tags
	return nil
}

// ReviewImport is pure: it never creates a directory, template, or project file.
func ReviewImport(payload string) (ImportReview, error) {
	if len(payload) > maxEnvelopeSize || !utf8.ValidString(payload) {
		return ImportReview{}, errors.New("import is too large or not UTF-8")
	}
	decoder := json.NewDecoder(strings.NewReader(payload))
	decoder.DisallowUnknownFields()
	var envelope Envelope
	if err := decoder.Decode(&envelope); err != nil {
		return ImportReview{}, errors.New("invalid template envelope; use AgentForge JSON export format")
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF || envelope.Format != "agentforge-template" || envelope.Version != 1 {
		return ImportReview{}, errors.New("unsupported template format or version")
	}
	item := fromPortable(envelope.Template)
	if err := normalize(&item); err != nil {
		return ImportReview{}, err
	}
	return ImportReview{portable(item), providers.Portability(item.SourcePath, item.Content)}, nil
}

// Import always creates a new library item, never overwrites an existing ID.
func (l *Library) Import(payload string) (Template, error) {
	review, err := ReviewImport(payload)
	if err != nil {
		return Template{}, err
	}
	for _, issue := range review.Issues {
		if issue.Severity == "error" {
			return Template{}, errors.New(issue.Message)
		}
	}
	return l.store(fromPortable(review.Template))
}

func (l *Library) Export(id string) (string, error) {
	item, err := l.Get(id)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(Envelope{"agentforge-template", 1, portable(item)}, "", "  ")
	return string(data), err
}
