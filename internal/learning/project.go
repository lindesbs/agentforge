package learning

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

const maxDocumentBytes = 1024 * 1024
const missingHash = "missing"

type ProjectPreview struct {
 Path string `json:"path"`
 BeforeHash string `json:"beforeHash"`
 Proposed string `json:"proposed"`
 Diff string `json:"diff"`
 Changed bool `json:"changed"`
}

func digest(data []byte) string {
 value:=sha256.Sum256(data)
 return hex.EncodeToString(value[:])
}

func permittedProject(root string)(string,error) {
 if strings.TrimSpace(root)=="" {return "",errors.New("project root is required")}
 abs,err:=filepath.Abs(root);if err!=nil{return "",err}
 fi,err:=os.Lstat(abs);if err!=nil{return "",err}
 if !fi.IsDir()||fi.Mode()&os.ModeSymlink!=0{return "",errors.New("invalid project root")}
 return abs,nil
}
func categoryPath(category string)(string,error) {
 if !categories[category] {return "",errors.New("unknown learning category")}
 return filepath.ToSlash(filepath.Join("docs",category+".md")),nil
}
func docPath(root,category string)(string,string,error) {
 abs,err:=permittedProject(root);if err!=nil{return "","",err}
 rel,err:=categoryPath(category);if err!=nil{return "","",err}
 docs:=filepath.Join(abs,"docs")
 info,err:=os.Lstat(docs)
 if err==nil {
  if !info.IsDir()||info.Mode()&os.ModeSymlink!=0 {return "","",errors.New("docs must be a real directory")}
 } else if !os.IsNotExist(err){return "","",err}
 path:=filepath.Join(abs,filepath.FromSlash(rel))
 info,err=os.Lstat(path)
 if err==nil {
  if !info.Mode().IsRegular(){return "","",errors.New("learning document must be a regular file")}
 } else if !os.IsNotExist(err){return "","",err}
 return path,rel,nil
}
func loadDoc(path string)([]byte,string,error){
 fi,err:=os.Lstat(path)
 if os.IsNotExist(err){return nil,missingHash,nil}
 if err!=nil{return nil,"",err}
 if !fi.Mode().IsRegular()||fi.Size()>maxDocumentBytes {return nil,"",errors.New("invalid or oversized document")}
 data,err:=os.ReadFile(path);if err!=nil{return nil,"",err}
 if len(data)>maxDocumentBytes||!utf8.Valid(data){return nil,"",errors.New("learning document must be UTF-8 and <= 1 MiB")}
 return data,digest(data),nil
}
func appendLearning(old []byte,entry Entry)(string,error){
 if err:=validate(entry);err!=nil{return "",err}
 marker:="### "+entry.Title
 for _,line:=range strings.Split(string(old),"\n") {
  if strings.EqualFold(strings.TrimSpace(line),marker){
   return "",errors.New("document already contains an entry with this heading; update it manually instead")
  }
 }
 prefix:=string(old)
 if strings.TrimSpace(prefix)=="" {prefix="# "+entry.Category+"\n"}
 if !strings.HasSuffix(prefix,"\n"){prefix+="\n"}
 prefix+="\n"+Markdown(entry)+"\n"
 if len(prefix)>maxDocumentBytes{return "",errors.New("resulting document exceeds size limit")}
 return prefix,nil
}

// PreviewProject computes the complete proposed Markdown file without modifying the project.
// The project source remains authoritative; no implicit overwrite or synchronization.
func PreviewProject(root string,entry Entry)(ProjectPreview,error){
 path,rel,err:=docPath(root,entry.Category);if err!=nil{return ProjectPreview{},err}
 old,hash,err:=loadDoc(path);if err!=nil{return ProjectPreview{},err}
 proposed,err:=appendLearning(old,entry);if err!=nil{return ProjectPreview{},err}
 return ProjectPreview{
  Path:rel,BeforeHash:hash,Proposed:proposed,
  Diff:fmt.Sprintf("--- a/%s\n+++ b/%s\n%s",rel,rel,addedLines(string(old),proposed)),
  Changed:proposed!=string(old),
 },nil
}
func addedLines(old,next string)string{
 suffix:=strings.TrimPrefix(next,old)
 if old!=""&&!strings.HasPrefix(next,old) {suffix=next}
 var b strings.Builder
 for _,line:=range strings.SplitAfter(suffix,"\n") {
  if line!="" {b.WriteString("+"+strings.TrimSuffix(line,"\n")+"\n")}
 }
 return b.String()
}

// ApplyProject rechecks the original content and an exact proposed preview,
// then writes through a temporary file. Existing documents are backed up privately.
func ApplyProject(root string,entry Entry,expectedHash,expectedProposed string)(ProjectPreview,error){
 plan,err:=PreviewProject(root,entry);if err!=nil{return ProjectPreview{},err}
 if plan.BeforeHash!=expectedHash||plan.Proposed!=expectedProposed {return ProjectPreview{},errors.New("project document or preview changed; review again")}
 path,_,err:=docPath(root,entry.Category);if err!=nil{return ProjectPreview{},err}
 old,hash,err:=loadDoc(path);if err!=nil{return ProjectPreview{},err}
 if hash!=expectedHash{return ProjectPreview{},errors.New("concurrent change detected")}
 dir:=filepath.Dir(path)
 if err=os.MkdirAll(dir,0700);err!=nil{return ProjectPreview{},err}
 // Re-check the directory after creation, to reject a symlink substitution.
 if _,_,err=docPath(root,entry.Category);err!=nil{return ProjectPreview{},err}
 // A dedicated backup is retained if an existing document is replaced.
 if expectedHash!=missingHash {
  backup,err:=os.CreateTemp(dir,".agentforge-learning-backup-*");if err!=nil{return ProjectPreview{},err}
  backupName:=backup.Name()
  if _,err=backup.Write(old);err!=nil {backup.Close();os.Remove(backupName);return ProjectPreview{},err}
  if err=backup.Sync();err!=nil{backup.Close();os.Remove(backupName);return ProjectPreview{},err}
  if err=backup.Close();err!=nil {os.Remove(backupName);return ProjectPreview{},err}
 }
 tmp,err:=os.CreateTemp(dir,".agentforge-learning-write-*");if err!=nil{return ProjectPreview{},err}
 tmpName:=tmp.Name()
 defer os.Remove(tmpName)
 if _,err=tmp.WriteString(expectedProposed);err!=nil{tmp.Close();return ProjectPreview{},err}
 if err=tmp.Sync();err!=nil{tmp.Close();return ProjectPreview{},err}
 if err=tmp.Close();err!=nil{return ProjectPreview{},err}
 _,latest,err:=loadDoc(path);if err!=nil{return ProjectPreview{},err}
 if latest!=expectedHash{return ProjectPreview{},errors.New("document changed during save")}
 if expectedHash==missingHash {
  // Hard-link creation fails when another process created the file first.
  if err=os.Link(tmpName,path);err!=nil{return ProjectPreview{},fmt.Errorf("unable to safely create document: %w",err)}
 } else {
  if err=os.Rename(tmpName,path);err!=nil{return ProjectPreview{},err}
 }
 return plan,nil
}
