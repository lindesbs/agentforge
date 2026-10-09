package templates

import (
 "os"
 "path/filepath"
 "testing"

 "github.com/lindesbs/agentforge/internal/editor"
)

func write(t *testing.T,root,path,content string) {
 t.Helper()
 full:=filepath.Join(root,filepath.FromSlash(path))
 if err:=os.MkdirAll(filepath.Dir(full),0700);err!=nil{t.Fatal(err)}
 if err:=os.WriteFile(full,[]byte(content),0600);err!=nil{t.Fatal(err)}
}
func TestLibraryRoundTrip(t *testing.T) {
 project:=t.TempDir()
 write(t,project,".claude/agents/reviewer.md","---\nname: reviewer\n---\nbody\n")
 library,err:=New(t.TempDir());if err!=nil{t.Fatal(err)}
 saved,err:=library.CreateFromProject(project,".claude/agents/reviewer.md","Review agent")
 if err!=nil{t.Fatal(err)}
 items,err:=library.List();if err!=nil{t.Fatal(err)}
 if len(items)!=1||items[0].ID!=saved.ID {t.Fatalf("unexpected summaries %#v",items)}
 target:=t.TempDir();write(t,target,".claude/agents/other.md","old")
 prepared,err:=library.PrepareApply(saved.ID,target,".claude/agents/other.md")
 if err!=nil{t.Fatal(err)}
 if prepared.Content!=saved.Content{t.Fatal("content differs")}
 current,err:=editor.Read(target,prepared.Path);if err!=nil{t.Fatal(err)}
 if current.Content!="old"{t.Fatal("PrepareApply unexpectedly modified the project")}
}
func TestRejectIncompatibleOrMissing(t *testing.T) {
 project:=t.TempDir();write(t,project,"AGENTS.md","hello");write(t,project,"CLAUDE.md","world")
 library,err:=New(t.TempDir());if err!=nil{t.Fatal(err)}
 saved,err:=library.CreateFromProject(project,"AGENTS.md","Instructions");if err!=nil{t.Fatal(err)}
 if _,err:=library.PrepareApply(saved.ID,project,"CLAUDE.md");err==nil{t.Fatal("should reject provider mismatch")}
 if _,err:=library.Get("../escape");err==nil{t.Fatal("should reject invalid id")}
 if _,err:=library.CreateFromProject(project,"missing.md","Nope");err==nil{t.Fatal("should reject unsupported type")}
}
