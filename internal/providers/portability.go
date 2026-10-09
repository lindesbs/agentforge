package providers

import "strings"

// Portability only inspects supported fields; it never opens referenced files.
func Portability(relative, content string) []Issue {
	definition := Inspect(relative, content)
	issues := append([]Issue{}, definition.Issues...)
	for _, agent := range definition.Agents {
		if agent.ConfigFile == "" {
			continue
		}
		ref := strings.ReplaceAll(agent.ConfigFile, "\\", "/")
		if strings.HasPrefix(ref, "/") || strings.HasPrefix(ref, "~") || strings.Contains(ref, ":") || strings.Contains(ref, "$") || strings.Contains("/"+ref+"/", "/../") {
			issues = append(issues, Issue{"nonportable-reference", "warning", "A role references a machine-specific or outside-project path. Review it in the destination project."})
		} else {
			issues = append(issues, Issue{"unbundled-reference", "warning", "A role references another configuration file. Referenced files are not included or verified."})
		}
	}
	issues = append(issues, Issue{"manual-portability-review", "warning", "Review native content for credentials, local paths and project-specific settings. Unknown fields are preserved, not checked for portability."})
	return issues
}
