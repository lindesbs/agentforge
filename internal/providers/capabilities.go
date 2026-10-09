package providers

// Capability describes AgentForge's supported subset, not every vendor feature.
type Capability struct {
	Provider         string   `json:"provider"`
	Kind             string   `json:"kind"`
	Format           string   `json:"format"`
	StructuredFields []string `json:"structuredFields"`
	Validation       string   `json:"validation"`
	TemplateSupport  bool     `json:"templateSupport"`
}

func Capabilities() []Capability {
	return []Capability{
		{"codex", "instructions", "Markdown", []string{}, "Non-empty content advisory", true},
		{"codex", "settings", "TOML", []string{}, "Syntax and selected field types", true},
		{"claude", "instructions", "Markdown", []string{}, "Non-empty content advisory", true},
		{"claude", "settings", "JSON", []string{}, "Object syntax and selected field types", true},
		{"claude", "agent", "Markdown + YAML", []string{"name", "description"}, "Frontmatter syntax and selected field types", true},
	}
}
