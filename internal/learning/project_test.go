package learning

import (
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestProjectPreviewThenApply(t *testing.T) {
 root:=t.TempDir()
 e:=entry()
 preview,err:=PreviewProject(root,e);if err!=nil{t.Fatal(err)}
 if preview.Path!="docs/ALLGEMEIN_NOK.md"||preview.BeforeHash!=missingHash||!preview.Changed{t.Fatalf("unexpected preview %#v",preview)}
 if _,err:=os.Stat(filepath.Join(root,filepath.FromSlash(preview.Path)));!os.IsNotExist(err){t.Fatalf("preview wrote file: %v",err)}
 if _,err:=ApplyProject(root,e,preview.BeforeHash,preview.Proposed);err!=nil{t.Fatal(err)}
 b,err:=os.ReadFile(filepath.Join(root,filepath.FromSlash(preview.Path)));if err!=nil{t.Fatal(err)}
 if string(b)!=preview.Proposed||!strings.Contains(string(b),"**Prüfung:**"){t.Fatal("saved wrong content")}
 if _,err:=PreviewProject(root,e);err==nil{t.Fatal("duplicate learning accepted")}
}

func TestExistingDocConflictAndBackup(t *testing.T) {
 root:=t.TempDir();dir:=filepath.Join(root,"docs")
 if err:=os.MkdirAll(dir,0700);err!=nil{t.Fatal(err)}
 path:=filepath.Join(dir,"ALLGEMEIN_NOK.md")
 if err:=os.WriteFile(path,[]byte("# Existing\n"),0600);err!=nil{t.Fatal(err)}
 e:=entry();p,err:=PreviewProject(root,e);if err!=nil{t.Fatal(err)}
 if err:=os.WriteFile(path,[]byte("# Concurrent\n"),0600);err!=nil{t.Fatal(err)}
 if _,err:=ApplyProject(root,e,p.BeforeHash,p.Proposed);err==nil{t.Fatal("stale preview accepted")}
 if err:=os.WriteFile(path,[]byte("# Existing\n"),0600);err!=nil{t.Fatal(err)}
 if _,err:=ApplyProject(root,e,p.BeforeHash,p.Proposed);err!=nil{t.Fatal(err)}
 entries,err:=os.ReadDir(dir);if err!=nil{t.Fatal(err)}
 found:=false
 for _,file:=range entries {
  if strings.HasPrefix(file.Name(),".agentforge-learning-backup-") {
   found=true;info,err:=file.Info();if err!=nil{t.Fatal(err)}
   if info.Mode().Perm()!=0600 {t.Fatal("backup permissions are not private")}
  }
 }
 if !found {t.Fatal("recovery backup not created")}
}
func TestRejectSymlinkDocs(t *testing.T) {
 root:=t.TempDir();outside:=t.TempDir()
 if err:=os.Symlink(outside,filepath.Join(root,"docs"));err!=nil{t.Skip(err)}
 if _,err:=PreviewProject(root,entry());err==nil{t.Fatal("accepted symlink")}
}
