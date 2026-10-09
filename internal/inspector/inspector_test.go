package inspector

import (
 "os"
 "path/filepath"
 "testing"
)

func write(t *testing.T, root,path string) {
 t.Helper()
 name:=filepath.Join(root,path)
 if err:=os.MkdirAll(filepath.Dir(name),0700);err!=nil {t.Fatal(err)}
 if err:=os.WriteFile(name,[]byte("test"),0600);err!=nil {t.Fatal(err)}
}

func TestInspect(t *testing.T) {
 root:=t.TempDir()
 write(t,root,"composer.json")
 write(t,root,"AGENTS.md")
 write(t,root,".claude/agents/reviewer.md")
 p,err:=New().Inspect(root)
 if err!=nil {t.Fatal(err)}
 if len(p.Frameworks)!=1||p.Frameworks[0]!="PHP/Composer" {t.Fatalf("frameworks: %#v",p.Frameworks)}
 if len(p.Files)!=2 {t.Fatalf("files: %#v",p.Files)}
 if p.Files[0].Path!=".claude/agents/reviewer.md" {t.Fatalf("order: %#v",p.Files)}
}
func TestRejectEmptyAndFile(t *testing.T) {
 if _,err:=New().Inspect("");err==nil {t.Fatal("expected error")}
 root:=t.TempDir();write(t,root,"plain.txt")
 if _,err:=New().Inspect(filepath.Join(root,"plain.txt"));err==nil {t.Fatal("expected error")}
}
func TestIgnoreSymlink(t *testing.T) {
 root:=t.TempDir()
 write(t,root,"outside.md")
 if err:=os.MkdirAll(filepath.Join(root,".claude","agents"),0700);err!=nil {t.Fatal(err)}
 if err:=os.Symlink(filepath.Join(root,"outside.md"),filepath.Join(root,".claude","agents","escape.md"));err!=nil {t.Skipf("symlink unsupported: %v",err)}
 p,err:=New().Inspect(root);if err!=nil {t.Fatal(err)}
 if len(p.Files)!=0 {t.Fatalf("unexpected symlink file: %#v",p.Files)}
}
