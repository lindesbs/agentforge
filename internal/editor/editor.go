package editor

import (
 "crypto/sha256"
 "encoding/hex"
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "strings"
 "unicode/utf8"
)

const maxSize int64 = 1024 * 1024

type Document struct {
 Path string `json:"path"`
 Content string `json:"content"`
 Hash string `json:"hash"`
}
type Preview struct {
 Path string `json:"path"`
 Diff string `json:"diff"`
 Changed bool `json:"changed"`
}

func hash(s string) string { sum:=sha256.Sum256([]byte(s));return hex.EncodeToString(sum[:]) }

// resolve permits only recognized project-local configuration files.
// Every path component is checked with Lstat to reject symlink traversal.
func resolve(root, relative string) (string,error) {
 if root==""||relative==""||filepath.IsAbs(relative) {return "",errors.New("invalid project or file path")}
 cleaned:=filepath.Clean(filepath.FromSlash(relative))
 if cleaned=="."||cleaned==".."||strings.HasPrefix(cleaned,".."+string(os.PathSeparator)){return "",errors.New("path escapes project")}
 allowed:= cleaned=="AGENTS.md"||cleaned=="CLAUDE.md"||
 cleaned==filepath.Join(".codex","config.toml")||
 cleaned==filepath.Join(".claude","settings.json")||
 cleaned==filepath.Join(".claude","settings.local.json")||
 (filepath.Dir(cleaned)==filepath.Join(".claude","agents")&&strings.EqualFold(filepath.Ext(cleaned),".md"))
 if !allowed {return "",errors.New("unsupported configuration path")}
 abs,err:=filepath.Abs(root);if err!=nil{return "",err}
 current:=abs
 info,err:=os.Lstat(current);if err!=nil{return "",err}
 if !info.IsDir()||info.Mode()&os.ModeSymlink!=0{return "",errors.New("invalid project root")}
 for _,part:=range strings.Split(cleaned,string(os.PathSeparator)){
  current=filepath.Join(current,part)
  info,err=os.Lstat(current);if err!=nil{return "",err}
  if info.Mode()&os.ModeSymlink!=0{return "",errors.New("symlink paths are not allowed")}
 }
 if !info.Mode().IsRegular(){return "",errors.New("not a regular file")}
 return current,nil
}
func read(root,relative string)(Document,error){
 path,err:=resolve(root,relative);if err!=nil{return Document{},err}
 info,err:=os.Stat(path);if err!=nil{return Document{},err}
 if info.Size()>maxSize{return Document{},errors.New("file exceeds 1 MiB limit")}
 data,err:=os.ReadFile(path);if err!=nil{return Document{},err}
 if int64(len(data))>maxSize||!utf8.Valid(data){return Document{},errors.New("file too large or not UTF-8")}
 content:=string(data)
 return Document{Path:relative,Content:content,Hash:hash(content)},nil
}
func Read(root,relative string)(Document,error){return read(root,relative)}

// Preview performs a simple line-based change preview. It never writes files.
func PreviewChange(root,relative,expectedHash,newContent string)(Preview,error){
 old,err:=read(root,relative);if err!=nil{return Preview{},err}
 if old.Hash!=expectedHash{return Preview{},errors.New("file changed on disk; reload before editing")}
 if int64(len(newContent))>maxSize||!utf8.ValidString(newContent){return Preview{},errors.New("new content too large or not UTF-8")}
 if old.Content==newContent{return Preview{Path:relative,Changed:false},nil}
 before:=strings.Split(old.Content,"\n")
 after:=strings.Split(newContent,"\n")
 prefix:=0
 for prefix<len(before)&&prefix<len(after)&&before[prefix]==after[prefix]{prefix++}
 suffix:=0
 for suffix<len(before)-prefix&&suffix<len(after)-prefix&&before[len(before)-1-suffix]==after[len(after)-1-suffix]{suffix++}
 var b strings.Builder
 fmt.Fprintf(&b,"--- a/%s\n+++ b/%s\n",relative,relative)
 for _,line:=range before[prefix:len(before)-suffix]{b.WriteString("-"+line+"\n")}
 for _,line:=range after[prefix:len(after)-suffix]{b.WriteString("+"+line+"\n")}
 return Preview{Path:relative,Diff:b.String(),Changed:true},nil
}

// Save creates a same-directory backup before replacing a file atomically.
// File mode is preserved. Changes must be previewed by the UI first.
func Save(root,relative,expectedHash,newContent string)(Document,error){
 preview,err:=PreviewChange(root,relative,expectedHash,newContent);if err!=nil{return Document{},err}
 if !preview.Changed{return read(root,relative)}
 path,err:=resolve(root,relative);if err!=nil{return Document{},err}
 original,err:=os.ReadFile(path);if err!=nil{return Document{},err}
 if hash(string(original))!=expectedHash{return Document{},errors.New("file changed on disk; reload")}
 info,err:=os.Stat(path);if err!=nil{return Document{},err}
 dir:=filepath.Dir(path)
 backup,err:=os.CreateTemp(dir,".agentforge-backup-*");if err!=nil{return Document{},err}
 backupPath:=backup.Name()
 if _,err=backup.Write(original);err!=nil{backup.Close();os.Remove(backupPath);return Document{},err}
 if err=backup.Chmod(info.Mode().Perm());err!=nil{backup.Close();os.Remove(backupPath);return Document{},err}
 if err=backup.Sync();err!=nil{backup.Close();os.Remove(backupPath);return Document{},err}
 if err=backup.Close();err!=nil{os.Remove(backupPath);return Document{},err}
 tmp,err:=os.CreateTemp(dir,".agentforge-write-*");if err!=nil{return Document{},err}
 tmpPath:=tmp.Name()
 defer os.Remove(tmpPath)
 if _,err=tmp.WriteString(newContent);err!=nil{tmp.Close();return Document{},err}
 if err=tmp.Chmod(info.Mode().Perm());err!=nil{tmp.Close();return Document{},err}
 if err=tmp.Sync();err!=nil{tmp.Close();return Document{},err}
 if err=tmp.Close();err!=nil{return Document{},err}
 // Re-check before replacement. External writes between this check and rename
 // remain a limitation; a future platform-specific lock can close that window.
 current,err:=read(root,relative);if err!=nil{return Document{},err}
 if current.Hash!=expectedHash{return Document{},errors.New("file changed before save; reload")}
 if err=os.Rename(tmpPath,path);err!=nil{return Document{},err}
 return read(root,relative)
}
