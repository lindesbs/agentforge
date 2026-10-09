package editor

import "testing"

func TestValidateSettingsJSON(t *testing.T) {
 cases:=[]struct{input string; hasError bool}{
  {`{"permissions":{"allow":[]}}`,false},
  {`{"permissions":`,true},
  {`{"a":1} {"b":2}`,true},
  {`[1,2]`,true},
 }
 for _,tc:=range cases{
  issues:=Validate(".claude/settings.json",tc.input)
  hasError:=false
  for _,issue:=range issues {if issue.Severity=="error"{hasError=true}}
  if hasError!=tc.hasError {t.Fatalf("input %q issues %#v",tc.input,issues)}
 }
}

func TestValidateMarkdownAndTOML(t *testing.T) {
 if len(Validate("AGENTS.md","  "))==0 {t.Fatal("expected warning")}
 issues:=Validate(".codex/config.toml","model = 'test'")
 if len(issues)!=1||issues[0].Severity!="info" {t.Fatalf("unexpected TOML issues: %#v",issues)}
}
