package editor

import (
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestReadPreviewSave(t *testing.T){
 root:=t.TempDir()
 name:=filepath.Join(root,"AGENTS.md")
 if err:=os.WriteFile(name,[]byte("first\n"),0600);err!=nil{t.Fatal(err)}
 doc,err:=Read(root,"AGENTS.md");if err!=nil{t.Fatal(err)}
 p,err:=PreviewChange(root,"AGENTS.md",doc.Hash,"second\n")
 if err!=nil||!p.Changed||!strings.Contains(p.Diff,"+second"){t.Fatalf("preview: %+v %v",p,err)}
 saved,err:=Save(root,"AGENTS.md",doc.Hash,"second\n");if err!=nil{t.Fatal(err)}
 if saved.Content!="second\n"||saved.Hash==doc.Hash{t.Fatal("save failed")}
 if _,err:=Save(root,"AGENTS.md",doc.Hash,"third\n");err==nil{t.Fatal("expected stale hash error")}
 entries,err:=os.ReadDir(root);if err!=nil{t.Fatal(err)}
 backupFound:=false
 for _,e:=range entries{if strings.HasPrefix(e.Name(),".agentforge-backup-"){backupFound=true}}
 if !backupFound{t.Fatal("expected recovery backup")}
}
func TestPathProtection(t *testing.T){
 root:=t.TempDir()
 if err:=os.WriteFile(filepath.Join(root,"AGENTS.md"),[]byte("ok"),0600);err!=nil{t.Fatal(err)}
 for _,p:=range []string{"../AGENTS.md","other.md","/tmp/AGENTS.md"}{
  if _,err:=Read(root,p);err==nil{t.Fatalf("accepted %s",p)}
 }
 if err:=os.MkdirAll(filepath.Join(root,".claude","agents"),0700);err!=nil{t.Fatal(err)}
 if err:=os.Symlink(filepath.Join(root,"AGENTS.md"),filepath.Join(root,".claude","agents","link.md"));err!=nil{t.Skip(err)}
 if _,err:=Read(root,".claude/agents/link.md");err==nil{t.Fatal("accepted symlink")}
}
func TestNoChange(t *testing.T){
 root:=t.TempDir()
 if err:=os.WriteFile(filepath.Join(root,"CLAUDE.md"),[]byte("hello"),0600);err!=nil{t.Fatal(err)}
 d,err:=Read(root,"CLAUDE.md");if err!=nil{t.Fatal(err)}
 p,err:=PreviewChange(root,d.Path,d.Hash,d.Content)
 if err!=nil||p.Changed{t.Fatalf("unexpected preview: %+v %v",p,err)}
}
