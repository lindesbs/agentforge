package learning

import (
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func entry()Entry {return Entry{Category:"ALLGEMEIN_NOK",Language:"Go",Title:"Fehler beim Build",Finding:"Fehlendes production-Tag erzeugt Stub",Rule:"Wails mit production kompilieren",Check:"Binärdatei starten und CI prüfen"}}
func TestSaveAndList(t *testing.T){
 dir:=t.TempDir()
 lib,err:=New(dir);if err!=nil{t.Fatal(err)}
 saved,err:=lib.Save(entry());if err!=nil{t.Fatal(err)}
 if saved.ID==""||saved.CreatedAt==""{t.Fatal("missing persisted metadata")}
 all,err:=lib.List();if err!=nil{t.Fatal(err)}
 if len(all)!=1||all[0].Title!=saved.Title{t.Fatalf("unexpected entries %#v",all)}
 info,err:=os.Stat(filepath.Join(dir,saved.ID+".json"));if err!=nil{t.Fatal(err)}
 if info.Mode().Perm()!=0600{t.Fatalf("insecure mode: %o",info.Mode().Perm())}
 if _,err=lib.Save(entry());err==nil{t.Fatal("duplicate not detected")}
 text:=Markdown(saved)
 if !strings.Contains(text,"**Prüfung:**")||!strings.Contains(text,"Wails"){t.Fatal("missing markdown fields")}
}
func TestRejectUnconfirmedOrInvalid(t *testing.T){
 lib,_:=New(t.TempDir())
 bad:=entry();bad.Category="HYPOTHESE"
 if _,err:=lib.Save(bad);err==nil{t.Fatal("invalid classification accepted")}
 bad=entry();bad.Check=""
 if _,err:=lib.Save(bad);err==nil{t.Fatal("missing test accepted")}
 bad=entry();bad.Language="../../escape"
 if _,err:=lib.Save(bad);err==nil{t.Fatal("unknown language accepted")}
}
