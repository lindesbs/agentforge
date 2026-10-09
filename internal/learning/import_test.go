package learning

import (
 "os"
 "path/filepath"
 "strings"
 "testing"
)
func TestImportExistingProjectDocuments(t *testing.T){
 project:=t.TempDir()
 if err:=os.MkdirAll(filepath.Join(project,"docs"),0700);err!=nil{t.Fatal(err)}
 original:="# Allgemein OK\n\n## Bekannter Ansatz\n\n**Erkenntnis:** Das Verhalten wurde bestätigt.\n\n**Regel:** Die Prüfung vor der Änderung ausführen.\n\n**Prüfung:** Regressionstest für den betroffenen Pfad.\n"
 source:=filepath.Join(project,"docs","ALLGEMEIN_OK.md")
 if err:=os.WriteFile(source,[]byte(original),0600);err!=nil{t.Fatal(err)}
 lib,err:=New(t.TempDir());if err!=nil{t.Fatal(err)}
 scan,err:=lib.ScanProject(project);if err!=nil{t.Fatal(err)}
 if len(scan.Candidates)!=1||scan.Candidates[0].Entry.Title!="Bekannter Ansatz" {t.Fatalf("unexpected scan: %+v",scan)}
 candidate:=scan.Candidates[0]
 if candidate.Entry.Language!="Allgemein"||candidate.Duplicate {t.Fatalf("unexpected import candidate: %+v",candidate)}
 saved,err:=lib.ImportProjectEntry(project,candidate.Category,candidate.Index,candidate.SourceHash,"Go")
 if err!=nil{t.Fatal(err)}
 if saved.Language!="Go"||saved.Finding!="Das Verhalten wurde bestätigt."{t.Fatalf("unexpected saved entry: %+v",saved)}
 after,err:=os.ReadFile(source);if err!=nil{t.Fatal(err)}
 if string(after)!=original{t.Fatal("project document modified by import")}
 if _,err:=lib.ImportProjectEntry(project,candidate.Category,candidate.Index,candidate.SourceHash,"Go");err==nil{t.Fatal("duplicate accepted")}
 if err:=os.WriteFile(source,[]byte(original+"\nchanged\n"),0600);err!=nil{t.Fatal(err)}
 if _,err:=lib.ImportProjectEntry(project,candidate.Category,candidate.Index,candidate.SourceHash,"PHP");err==nil{t.Fatal("stale source accepted")}
}
func TestImportSkipsUnstructuredSections(t *testing.T){
 project:=t.TempDir();os.Mkdir(filepath.Join(project,"docs"),0700)
 source:=filepath.Join(project,"docs","PROJEKT_NOK.md")
 data:="## Notiz ohne Regel\n\nNur ein Symptom.\n\n### Formatierter Eintrag\n\n**Sprache:** TypeScript\n\n**Bestätigte Erkenntnis/Ursache:** Etwas bestätigt.\n\n**Verbindliche Regel:** Immer prüfen.\n\n**Prüfung:** npm run build\n"
 if err:=os.WriteFile(source,[]byte(data),0600);err!=nil{t.Fatal(err)}
 lib,_:=New(t.TempDir())
 scan,err:=lib.ScanProject(project);if err!=nil{t.Fatal(err)}
 if len(scan.Candidates)!=1||len(scan.Skipped)!=1{t.Fatalf("unexpected scan: %+v",scan)}
 if scan.Candidates[0].Entry.Language!="TypeScript"{t.Fatal("language not recognized")}
}
func TestImportRejectsSymlinkAndTraversal(t *testing.T){
 project:=t.TempDir();external:=t.TempDir()
 if err:=os.Symlink(external,filepath.Join(project,"docs"));err!=nil{t.Skip(err)}
 lib,_:=New(t.TempDir())
 if _,err:=lib.ScanProject(project);err==nil{t.Fatal("followed symlink docs")}
 if _,err:=lib.ImportProjectEntry(project,"../../evil",0,strings.Repeat("0",64),"Go");err==nil{t.Fatal("invalid category accepted")}
}
