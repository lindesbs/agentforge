package editor

import (
 "testing"
 "os"
 "path/filepath"
)

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

func TestInvalidJSONCannotBeSaved(t *testing.T) {
 root:=t.TempDir()
 dir:=filepath.Join(root,".claude")
 if err:=os.MkdirAll(dir,0700);err!=nil {t.Fatal(err)}
 name:=filepath.Join(dir,"settings.json")
 if err:=os.WriteFile(name,[]byte("{}"),0600);err!=nil {t.Fatal(err)}
 doc,err:=Read(root,".claude/settings.json");if err!=nil {t.Fatal(err)}
 if _,err:=Save(root,doc.Path,doc.Hash,`{"invalid":`);err==nil {t.Fatal("invalid JSON was accepted")}
 unchanged,err:=os.ReadFile(name);if err!=nil {t.Fatal(err)}
 if string(unchanged)!="{}" {t.Fatalf("file was modified: %s",unchanged)}
}

func TestBackupHasPrivatePermissions(t *testing.T) {
 root:=t.TempDir()
 name:=filepath.Join(root,"AGENTS.md")
 if err:=os.WriteFile(name,[]byte("original"),0644);err!=nil {t.Fatal(err)}
 doc,err:=Read(root,"AGENTS.md");if err!=nil {t.Fatal(err)}
 if _,err:=Save(root,doc.Path,doc.Hash,"updated");err!=nil {t.Fatal(err)}
 entries,err:=os.ReadDir(root);if err!=nil {t.Fatal(err)}
 found:=false
 for _,entry:=range entries {
  if len(entry.Name())<18||entry.Name()[:18]!=".agentforge-backup" {continue}
  info,err:=entry.Info();if err!=nil {t.Fatal(err)}
  if info.Mode().Perm()!=0600 {t.Fatalf("backup mode %o; expected 600",info.Mode().Perm())}
  found=true
 }
 if !found {t.Fatal("backup missing")}
}
