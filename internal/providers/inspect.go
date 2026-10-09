// Package providers describes selected native configuration fields without
// modifying source text, executing prompts or following referenced files.
package providers

import (
	"encoding/json"
	"io"
	"path"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type Issue struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Agent struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ConfigFile  string `json:"configFile"`
}

type Definition struct {
	Provider string  `json:"provider"`
	Kind     string  `json:"kind"`
	Fields   []Field `json:"fields"`
	Agents   []Agent `json:"agents"`
	Issues   []Issue `json:"issues"`
}

func (d *Definition) issue(code, severity, message string) {
	d.Issues = append(d.Issues, Issue{code, severity, message})
}

// stringField is only called with fixed, supported field names. Parser errors
// and arbitrary unknown keys/values must never be interpolated into issues.
func (d *Definition) stringField(values map[string]any, key string, required bool) string {
	value, exists := values[key]
	if !exists {
		if required {
			d.issue("missing-field", "error", "Required field '"+key+"' is missing.")
		}
		return ""
	}
	text, ok := value.(string)
	if !ok || (required && strings.TrimSpace(text) == "") {
		message := "Field '" + key + "' must be a string."
		if required {
			message = "Field '" + key + "' must be a non-empty string."
		}
		d.issue("invalid-field", "error", message)
		return ""
	}
	return text
}

func (d *Definition) addString(values map[string]any, key string, required bool) string {
	value := d.stringField(values, key, required)
	if value != "" {
		d.Fields = append(d.Fields, Field{key, value})
	}
	return value
}

// Inspect parses only the explicitly supplied document. Its output is an
// allowlist of descriptive fields, never a dump of settings or prompt text.
// It is not a complete provider schema validator and is not used to rewrite
// source or to authorize saving a configuration.
func Inspect(relative, content string) Definition {
	d := Definition{Fields: []Field{}, Agents: []Agent{}, Issues: []Issue{}}
	relative = strings.ReplaceAll(relative, "\\", "/")
	switch {
	case relative == "AGENTS.md":
		d.Provider, d.Kind = "codex", "instructions"
	case relative == "CLAUDE.md":
		d.Provider, d.Kind = "claude", "instructions"
	case relative == ".codex/config.toml":
		d.Provider, d.Kind = "codex", "settings"
	case relative == ".claude/settings.json" || relative == ".claude/settings.local.json":
		d.Provider, d.Kind = "claude", "settings"
	case path.Dir(relative) == ".claude/agents" && strings.EqualFold(path.Ext(relative), ".md"):
		d.Provider, d.Kind = "claude", "agent"
	default:
		d.issue("unsupported-path", "info", "Structured inspection is unavailable for this path.")
		return d
	}
	if len(content) > 1024*1024 || !utf8.ValidString(content) {
		d.issue("invalid-document", "error", "Inspection requires UTF-8 text no larger than 1 MiB.")
		return d
	}
	if d.Kind == "instructions" {
		if strings.TrimSpace(content) == "" {
			d.issue("empty-instructions", "warning", "The instruction document is empty.")
		}
		return d
	}
	if d.Kind == "agent" {
		d.inspectClaudeAgent(content)
	} else if d.Provider == "codex" {
		d.inspectCodex(content)
	} else {
		var values map[string]any
		if err := json.Unmarshal([]byte(content), &values); err != nil || values == nil {
			d.issue("invalid-json", "error", "Claude settings must contain exactly one valid JSON object.")
			return d
		}
		d.addString(values, "model", false)
	}
	d.issue("limited-validation", "info", "Only syntax and selected field types are checked. Other provider options and referenced files are not validated.")
	return d
}

func (d *Definition) inspectClaudeAgent(content string) {
	lines := strings.Split(strings.ReplaceAll(strings.TrimPrefix(content, "\ufeff"), "\r\n", "\n"), "\n")
	if strings.TrimSpace(lines[0]) != "---" {
		d.issue("missing-frontmatter", "error", "Claude agent definitions need YAML frontmatter.")
		return
	}
	end := 1
	for end < len(lines) && strings.TrimSpace(lines[end]) != "---" {
		end++
	}
	if end == len(lines) {
		d.issue("unterminated-frontmatter", "error", "YAML frontmatter has no closing delimiter.")
		return
	}
	var values map[string]any
	decoder := yaml.NewDecoder(strings.NewReader(strings.Join(lines[1:end], "\n")))
	if err := decoder.Decode(&values); err != nil && err != io.EOF {
		d.issue("invalid-yaml", "error", "Agent frontmatter must be a valid YAML mapping without duplicate keys.")
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		d.issue("invalid-yaml", "error", "Agent frontmatter must contain exactly one YAML mapping.")
		return
	}
	name := d.addString(values, "name", true)
	description := d.addString(values, "description", true)
	d.addString(values, "model", false)
	if value, exists := values["tools"]; exists {
		var tools []string
		valid := true
		switch value := value.(type) {
		case string:
			tools = append(tools, value)
		case []any:
			for _, item := range value {
				text, ok := item.(string)
				if !ok || strings.TrimSpace(text) == "" {
					valid = false
					break
				}
				tools = append(tools, text)
			}
		default:
			valid = false
		}
		if valid {
			d.Fields = append(d.Fields, Field{"tools", strings.Join(tools, ", ")})
		} else {
			d.issue("invalid-tools", "error", "Field 'tools' must be a string or a list of strings.")
		}
	}
	if name != "" && description != "" {
		d.Agents = append(d.Agents, Agent{Name: name, Description: description})
	}
	if strings.TrimSpace(strings.Join(lines[end+1:], "\n")) == "" {
		d.issue("empty-prompt", "warning", "The agent has no Markdown instructions after its frontmatter.")
	}
}

func (d *Definition) inspectCodex(content string) {
	var values map[string]any
	if err := toml.Unmarshal([]byte(content), &values); err != nil {
		d.issue("invalid-toml", "error", "Codex settings contain invalid TOML or duplicate keys.")
		return
	}
	for _, key := range []string{"model", "model_provider", "model_reasoning_effort", "sandbox_mode"} {
		d.addString(values, key, false)
	}
	value, exists := values["agents"]
	if !exists {
		return
	}
	roles, ok := value.(map[string]any)
	if !ok {
		d.issue("invalid-agents", "error", "Codex 'agents' must be a TOML table.")
		return
	}
	names := make([]string, 0, len(roles))
	for name := range roles {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		// Scalar controls (e.g. max_threads) are not role definitions.
		role, ok := roles[name].(map[string]any)
		if !ok {
			continue
		}
		d.Agents = append(d.Agents, Agent{
			Name: name, Description: d.stringField(role, "description", false),
			ConfigFile: d.stringField(role, "config_file", false),
		})
	}
}
