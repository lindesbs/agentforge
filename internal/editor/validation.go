package editor

import (
	"errors"
	"github.com/lindesbs/agentforge/internal/providers"
)

// Issue contains fixed messages, never raw parser errors or source values.
type Issue struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// Validate checks the same supported syntax/field subset as loaded-file inspection.
func Validate(relative, content string) []Issue {
	issues := make([]Issue, 0)
	for _, issue := range providers.Inspect(relative, content).Issues {
		issues = append(issues, Issue{issue.Severity, issue.Message})
	}
	return issues
}

func validateForSave(relative, content string) error {
	for _, issue := range Validate(relative, content) {
		if issue.Severity == "error" {
			return errors.New(issue.Message)
		}
	}
	return nil
}
