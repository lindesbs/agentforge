package templates

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

 "github.com/lindesbs/agentforge/internal/editor"
)

const maxTemplateSize = 1024 * 1024

// Template is a local reusable snapshot. No automatic project synchronization.
type Template struct {
 ID string `json:"id"`
 Name string `json:"name"`
 Provider string `json:"provider"`
 Kind string `json:"kind"`
 SourcePath string `json:"sourcePath"`
 Content string `json:"content"`
 CreatedAt string `json:"createdAt"`
}
type Summary struct {
 ID string `json:"id"`
 Name string `json:"name"`
 Provider string `json:"provider"`
 Kind string `json:"kind"`
 SourcePath string `json:"sourcePath"`
 CreatedAt string `json:"createdAt"`
}
type Library struct{ root string }

// New opens an app-local library. Pass an empty path to use the OS user config directory.
func New(root string) (*Library,error) {
 if root=="" {
  config,err:=os.UserConfigDir();if err!=nil{return nil,err}
  root=filepath.Join(config,"agentforge","templates")
 }
 abs,err:=filepath.Abs(root);if err!=nil{return nil,err}
 return &Library{root:abs},nil
}
func (l *Library) ensureRoot()error {
 if err:=os.MkdirAll(l.root,0700);err!=nil{return err}
 info,err:=os.Lstat(l.root);if err!=nil{return err}
 if !info.IsDir()||info.Mode()&os.ModeSymlink!=0 {return errors.New("invalid template directory")}
 return nil
}
func validID(id string)bool {
 if len(id)!=32{return false}
 _,err:=hex.DecodeString(id);return err==nil
}
func (l *Library) read(id string)(Template,error) {
 if !validID(id){return Template{},errors.New("invalid template ID")}
 if err:=l.ensureRoot();err!=nil{return Template{},err}
 path:=filepath.Join(l.root,id+".json")
 info,err:=os.Lstat(path);if err!=nil{return Template{},err}
 if !info.Mode().IsRegular(){return Template{},errors.New("template is not a regular file")}
 if info.Size()>maxTemplateSize+4096{return Template{},errors.New("template too large")}
 data,err:=os.ReadFile(path);if err!=nil{return Template{},err}
 var result Template
 if err:=json.Unmarshal(data,&result);err!=nil{return Template{},err}
 if result.ID!=id{return Template{},errors.New("template ID mismatch")}
 return result,nil
}
func (l *Library) Get(id string)(Template,error){return l.read(id)}
func (l *Library) List()([]Summary,error) {
 if err:=l.ensureRoot();err!=nil{return nil,err}
 entries,err:=os.ReadDir(l.root);if err!=nil{return nil,err}
 result:=make([]Summary,0)
 for _,entry:=range entries {
  id:=strings.TrimSuffix(entry.Name(),".json")
  if !strings.HasSuffix(entry.Name(),".json")||!validID(id)||!entry.Type().IsRegular(){continue}
  item,err:=l.read(id);if err!=nil{return nil,fmt.Errorf("template %s: %w",id,err)}
  result=append(result,Summary{item.ID,item.Name,item.Provider,item.Kind,item.SourcePath,item.CreatedAt})
 }
 sort.Slice(result,func(i,j int)bool{
  if result[i].Name==result[j].Name{return result[i].ID<result[j].ID}
  return strings.ToLower(result[i].Name)<strings.ToLower(result[j].Name)
 })
 return result,nil
}
func describe(path string)(provider,kind string,err error){
 normalized:=filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
 switch {
 case normalized=="AGENTS.md":return "codex","instructions",nil
 case normalized==".codex/config.toml":return "codex","settings",nil
 case normalized=="CLAUDE.md":return "claude","instructions",nil
 case normalized==".claude/settings.json"||normalized==".claude/settings.local.json":return "claude","settings",nil
 case strings.HasPrefix(normalized,".claude/agents/") && strings.HasSuffix(strings.ToLower(normalized),".md") && !strings.Contains(strings.TrimPrefix(normalized,".claude/agents/"),"/"):return "claude","agent",nil
 }
 return "","",errors.New("unsupported configuration type")
}
// CreateFromProject snapshots an existing project-local config into the global library.
func (l *Library) CreateFromProject(projectRoot,relative,name string)(Template,error){
 name=strings.TrimSpace(name)
 if name==""||len([]rune(name))>120{return Template{},errors.New("template name must be 1–120 characters")}
 provider,kind,err:=describe(relative);if err!=nil{return Template{},err}
 doc,err:=editor.Read(projectRoot,relative);if err!=nil{return Template{},err}
 if len(doc.Content)>maxTemplateSize||!utf8.ValidString(doc.Content){return Template{},errors.New("invalid template content")}
 if err:=l.ensureRoot();err!=nil{return Template{},err}
 idBytes:=make([]byte,16);if _,err:=rand.Read(idBytes);err!=nil{return Template{},err}
 item:=Template{
  ID:hex.EncodeToString(idBytes),Name:name,Provider:provider,Kind:kind,
  SourcePath:relative,Content:doc.Content,CreatedAt:time.Now().UTC().Format(time.RFC3339),
 }
 bytes,err:=json.MarshalIndent(item,"","  ");if err!=nil{return Template{},err}
 path:=filepath.Join(l.root,item.ID+".json")
 f,err:=os.OpenFile(path,os.O_WRONLY|os.O_CREATE|os.O_EXCL,0600);if err!=nil{return Template{},err}
 if _,err=f.Write(bytes);err!=nil{f.Close();os.Remove(path);return Template{},err}
 if err=f.Sync();err!=nil{f.Close();os.Remove(path);return Template{},err}
 if err=f.Close();err!=nil{os.Remove(path);return Template{},err}
 return item,nil
}
// PrepareApply only returns proposed native content. The existing editor is
// responsible for diff preview, optimistic locking, and explicit save.
func (l *Library) PrepareApply(id,projectRoot,targetPath string)(editor.Document,error){
 item,err:=l.read(id);if err!=nil{return editor.Document{},err}
 provider,kind,err:=describe(targetPath);if err!=nil{return editor.Document{},err}
 if provider!=item.Provider||kind!=item.Kind {return editor.Document{},errors.New("template and target configuration types differ")}
 target,err:=editor.Read(projectRoot,targetPath);if err!=nil{return editor.Document{},err}
 if _,err:=editor.PreviewChange(projectRoot,targetPath,target.Hash,item.Content);err!=nil{return editor.Document{},err}
 return editor.Document{Path:targetPath,Hash:target.Hash,Content:item.Content},nil
}
