package learning

import (
 "crypto/rand"
 "encoding/hex"
 "encoding/json"
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "sort"
 "strings"
 "time"
 "unicode/utf8"
)

const maxField = 8192
var categories = map[string]bool{"PROJEKT_OK":true,"PROJEKT_NOK":true,"ALLGEMEIN_OK":true,"ALLGEMEIN_NOK":true}
var languages = map[string]bool{
 "Allgemein":true,"Go":true,"PHP":true,"TypeScript":true,"JavaScript":true,
 "Python":true,"Rust":true,"Kotlin":true,"Java":true,"C#":true,"C++":true,
 "Dart":true,"SQL":true,"Shell":true,"Andere":true,
}

// Entry stores a confirmed development insight, not a hypothesis.
type Entry struct {
 ID string `json:"id"`
 Category string `json:"category"`
 Language string `json:"language"`
 Title string `json:"title"`
 Finding string `json:"finding"`
 Rule string `json:"rule"`
 Check string `json:"check"`
 CreatedAt string `json:"createdAt"`
}

type Library struct{ root string }

func New(root string)(*Library,error){
 if root=="" {
  dir,err:=os.UserConfigDir();if err!=nil{return nil,err}
  root=filepath.Join(dir,"agentforge","learning")
 }
 absolute,err:=filepath.Abs(root);if err!=nil{return nil,err}
 return &Library{root:absolute},nil
}
func (l *Library) ensure()error{
 if err:=os.MkdirAll(l.root,0700);err!=nil{return err}
 info,err:=os.Lstat(l.root);if err!=nil{return err}
 if !info.IsDir()||info.Mode()&os.ModeSymlink!=0{return errors.New("unsafe library directory")}
 return nil
}
func validate(e Entry)error {
 if !categories[e.Category] {return errors.New("unknown insight category")}
 if !languages[e.Language] {return errors.New("unknown language")}
 for name,val:=range map[string]string{"title":e.Title,"finding":e.Finding,"rule":e.Rule,"check":e.Check} {
  if strings.TrimSpace(val)==""||len(val)>maxField||!utf8.ValidString(val) {return fmt.Errorf("%s is required and must be valid UTF-8 up to %d bytes",name,maxField)}
 }
 if len(e.Title)>120 {return errors.New("title exceeds 120 bytes")}
 return nil
}
func validID(id string)bool{if len(id)!=32{return false};_,err:=hex.DecodeString(id);return err==nil}
func (l *Library) load(id string)(Entry,error){
 if !validID(id){return Entry{},errors.New("invalid ID")}
 path:=filepath.Join(l.root,id+".json")
 fi,err:=os.Lstat(path);if err!=nil{return Entry{},err}
 if !fi.Mode().IsRegular()||fi.Size()>maxField*5 {return Entry{},errors.New("invalid insight file")}
 b,err:=os.ReadFile(path);if err!=nil{return Entry{},err}
 var e Entry
 if err=json.Unmarshal(b,&e);err!=nil{return Entry{},err}
 if e.ID!=id {return Entry{},errors.New("insight ID mismatch")}
 if err=validate(e);err!=nil{return Entry{},err}
 return e,nil
}
func (l *Library) List()([]Entry,error){
 if err:=l.ensure();err!=nil{return nil,err}
 entries,err:=os.ReadDir(l.root);if err!=nil{return nil,err}
 result:=make([]Entry,0,len(entries))
 for _,ent:=range entries {
  name:=ent.Name()
  if !strings.HasSuffix(name,".json")||!ent.Type().IsRegular(){continue}
  id:=strings.TrimSuffix(name,".json");if !validID(id){continue}
  item,err:=l.load(id);if err!=nil{return nil,err}
  result=append(result,item)
 }
 sort.Slice(result,func(i,j int)bool{
  if result[i].Language!=result[j].Language{return result[i].Language<result[j].Language}
  if result[i].Category!=result[j].Category{return result[i].Category<result[j].Category}
  return strings.ToLower(result[i].Title)<strings.ToLower(result[j].Title)
 })
 return result,nil
}
func (l *Library) Save(input Entry)(Entry,error){
 input.ID=""
 input.CreatedAt=""
 if err:=validate(input);err!=nil{return Entry{},err}
 if err:=l.ensure();err!=nil{return Entry{},err}
 existing,err:=l.List();if err!=nil{return Entry{},err}
 for _,item:=range existing {
  if item.Language==input.Language && item.Category==input.Category && strings.EqualFold(strings.TrimSpace(item.Title),strings.TrimSpace(input.Title)) {
   return Entry{},errors.New("an insight with this category, language and title already exists; extend the existing entry instead")
  }
 }
 buf:=make([]byte,16);if _,err:=rand.Read(buf);err!=nil{return Entry{},err}
 input.ID=hex.EncodeToString(buf)
 input.CreatedAt=time.Now().UTC().Format(time.RFC3339)
 data,err:=json.MarshalIndent(input,"","  ");if err!=nil{return Entry{},err}
 file,err:=os.OpenFile(filepath.Join(l.root,input.ID+".json"),os.O_CREATE|os.O_EXCL|os.O_WRONLY,0600)
 if err!=nil{return Entry{},err}
 path:=file.Name()
 if _,err=file.Write(data);err!=nil{file.Close();os.Remove(path);return Entry{},err}
 if err=file.Sync();err!=nil{file.Close();os.Remove(path);return Entry{},err}
 if err=file.Close();err!=nil{os.Remove(path);return Entry{},err}
 return input,nil
}

// Markdown formats an insight for a project learning document.
// It does not write to a project, run code or change AGENTS.md.
func Markdown(e Entry)string{
 return "### "+e.Title+"\n\n"+
 "**Sprache:** "+e.Language+"\n\n"+
 "**Bestätigte Erkenntnis/Ursache:** "+e.Finding+"\n\n"+
 "**Verbindliche Regel:** "+e.Rule+"\n\n"+
 "**Prüfung:** "+e.Check+"\n"
}
