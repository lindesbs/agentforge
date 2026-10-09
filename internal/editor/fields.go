package editor

import (
	"errors"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

type AgentFields struct {
	Path        string `json:"path"`
	Hash        string `json:"hash"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func filepathSlash(s string) string { return strings.ReplaceAll(s, "\\", "/") }

// scalarLocation uses source locations but never serializes the YAML mapping.
// Unsupported styles are refused rather than silently reformatting the file.
func scalarLocation(content, key string) (*yaml.Node, []string, error) {
	lines := strings.Split(content, "\n")
	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return nil, nil, errors.New("YAML frontmatter not found")
	}
	end := 1
	for end < len(lines) && strings.TrimSpace(lines[end]) != "---" {
		end++
	}
	if end == len(lines) {
		return nil, nil, errors.New("frontmatter not terminated")
	}
	source := strings.Join(lines[1:end], "\n")
	var mapping map[string]any
	var node yaml.Node
	if yaml.Unmarshal([]byte(source), &mapping) != nil || yaml.Unmarshal([]byte(source), &node) != nil || len(node.Content) != 1 || node.Content[0].Kind != yaml.MappingNode {
		return nil, nil, errors.New("frontmatter must be a valid mapping without duplicate keys")
	}
	children := node.Content[0].Content
	for i := 0; i < len(children); i += 2 {
		if children[i].Value != key {
			continue
		}
		value := children[i+1]
		if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || value.Anchor != "" || value.Style&(yaml.LiteralStyle|yaml.FoldedStyle|yaml.TaggedStyle) != 0 || strings.ContainsAny(value.Value, "\r\n") || value.Line != children[i].Line {
			return nil, nil, errors.New("complex YAML field must be edited in the native editor")
		}
		// Folded quoted/plain scalars may span lines without a newline in Value.
		line := strings.TrimSuffix(lines[value.Line], "\r")
		start := len(string([]rune(line)[:value.Column-1]))
		var single yaml.Node
		if yaml.Unmarshal([]byte(line[start:]), &single) != nil || len(single.Content) != 1 || single.Content[0].Value != value.Value {
			return nil, nil, errors.New("multiline fields must be edited in the native editor")
		}
		return value, lines, nil
	}
	return nil, nil, errors.New("field not found; add it in the native editor")
}

func readScalar(content, key string) (string, error) {
	node, _, err := scalarLocation(content, key)
	if err != nil {
		return "", err
	}
	return node.Value, nil
}

func ReadAgentFields(root, relative string) (AgentFields, error) {
	if !strings.HasPrefix(filepathSlash(relative), ".claude/agents/") {
		return AgentFields{}, errors.New("only Claude agent definitions supported")
	}
	doc, err := Read(root, relative)
	if err != nil {
		return AgentFields{}, err
	}
	name, err := readScalar(doc.Content, "name")
	if err != nil {
		return AgentFields{}, err
	}
	description, err := readScalar(doc.Content, "description")
	if err != nil {
		return AgentFields{}, err
	}
	return AgentFields{doc.Path, doc.Hash, name, description}, nil
}

func updateScalar(content, key, value string) (string, error) {
	if !utf8.ValidString(value) || strings.ContainsAny(value, "\r\n\x00") || strings.TrimSpace(value) == "" {
		return "", errors.New("field must be a non-empty single-line string")
	}
	node, lines, err := scalarLocation(content, key)
	if err != nil {
		return "", err
	}
	line := lines[node.Line]
	start := len(string([]rune(line)[:node.Column-1]))
	end := len(strings.TrimRight(line, " \t\r"))
	// Find a comment outside the existing quoted scalar; preserve spacing too.
	quote := byte(0)
	for i := start; i < len(line); i++ {
		c := line[i]
		if quote == '"' && c == '\\' {
			i++
			continue
		}
		if quote != 0 {
			if c == quote {
				if quote == '\'' && i+1 < len(line) && line[i+1] == '\'' {
					i++
					continue
				}
				quote = 0
			}
			continue
		}
		if i == start && (c == '\'' || c == '"') {
			quote = c
			continue
		}
		if c == '#' && i > start && (line[i-1] == ' ' || line[i-1] == '\t') {
			end = len(strings.TrimRight(line[:i], " \t"))
			break
		}
	}
	lines[node.Line] = line[:start] + "'" + strings.ReplaceAll(value, "'", "''") + "'" + line[end:]
	result := strings.Join(lines, "\n")
	actual, err := readScalar(result, key)
	if err != nil || actual != value {
		return "", errors.New("field cannot be safely represented; use the native editor")
	}
	return result, nil
}

func PrepareAgentFields(root, relative, expectedHash, name, description string) (Document, error) {
	doc, err := Read(root, relative)
	if err != nil {
		return Document{}, err
	}
	if !strings.HasPrefix(filepathSlash(relative), ".claude/agents/") {
		return Document{}, errors.New("only Claude agent definitions supported")
	}
	if doc.Hash != expectedHash {
		return Document{}, errors.New("file changed on disk")
	}
	content := doc.Content
	for _, field := range [][2]string{{"name", name}, {"description", description}} {
		old, err := readScalar(content, field[0])
		if err != nil {
			return Document{}, err
		}
		if old != field[1] {
			content, err = updateScalar(content, field[0], field[1])
			if err != nil {
				return Document{}, err
			}
		}
	}
	return Document{Path: relative, Hash: expectedHash, Content: content}, nil
}
