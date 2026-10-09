package editor

import (
 "os"
 "path/filepath"
 "strings"
 "testing"
)
func TestStructuredClaudeAgent(t *testing.T){
 root:=t.TempDir();dir:=filepath.Join(root,".claude","agents")
 if err:=os.MkdirAll(dir,0700);err!=nil{t.Fatal(err)}
 original:="---\nname: reviewer\ndescription: Reviews code\nmodel: sonnet\n---\n\nDo not modify body.\n"
 if err:=os.WriteFile(filepath.Join(dir,"reviewer.md"),[]byte(original),0600);err!=nil{t.Fatal(err)}
 fields,err:=ReadAgentFields(root,".claude/agents/reviewer.md");if err!=nil{t.Fatal(err)}
 if fields.Name!="reviewer"||fields.Description!="Reviews code"{t.Fatalf("%+v",fields)}
 next,err:=PrepareAgentFields(root,fields.Path,fields.Hash,"auditor","Checks code")
 if err!=nil{t.Fatal(err)}
 if !strings.Contains(next.Content,"model: sonnet")||!strings.Contains(next.Content,"Do not modify body."){t.Fatal("lost unedited content")}
 if !strings.Contains(next.Content,"name: 'auditor'"){t.Fatal("name not changed")}
 if _,err:=PreviewChange(root,fields.Path,fields.Hash,next.Content);err!=nil{t.Fatal(err)}
}
func TestRejectComplexField(t *testing.T){
 _,err:=updateScalar("---\nname: |\n  multiline\n---","name","new")
 if err==nil{t.Fatal("expected error")}
}
