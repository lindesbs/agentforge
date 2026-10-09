package editor

import (
 "errors"
 "strings"
)

// AgentFields represents supported simple Claude Code agent frontmatter fields.
// Complex YAML values remain untouched and cannot be edited in this view.
type AgentFields struct {
 Path string `json:"path"`
 Hash string `json:"hash"`
 Name string `json:"name"`
 Description string `json:"description"`
}

func readScalar(content, key string) (string,error) {
 lines:=strings.Split(content,"\n")
 if len(lines)<3||strings.TrimSpace(lines[0])!="---" {return "",errors.New("YAML frontmatter not found")}
 for i:=1;i<len(lines);i++ {
  line:=lines[i]
  if strings.TrimSpace(line)=="---" {break}
  if strings.HasPrefix(line,key+":") {
   raw:=strings.TrimSpace(strings.TrimPrefix(line,key+":"))
   if raw==""||strings.HasPrefix(raw,"|")||strings.HasPrefix(raw,">") {return "",errors.New("complex frontmatter value unsupported")}
   return strings.Trim(raw,"\"'"),nil
  }
 }
 return "",nil
}
func ReadAgentFields(root,relative string)(AgentFields,error) {
 if !strings.HasPrefix(filepathSlash(relative),".claude/agents/") {return AgentFields{},errors.New("only Claude agent definitions supported")}
 d,err:=Read(root,relative);if err!=nil{return AgentFields{},err}
 name,err:=readScalar(d.Content,"name");if err!=nil{return AgentFields{},err}
 description,err:=readScalar(d.Content,"description");if err!=nil{return AgentFields{},err}
 return AgentFields{Path:relative,Hash:d.Hash,Name:name,Description:description},nil
}
func filepathSlash(s string)string{return strings.ReplaceAll(s,"\\","/")}

func updateScalar(content,key,value string)(string,error) {
 if strings.ContainsAny(value,"\r\n") {return "",errors.New("multiline values not supported")}
 if strings.TrimSpace(value)=="" {return "",errors.New("field must not be empty")}
 lines:=strings.Split(content,"\n")
 if len(lines)<3||strings.TrimSpace(lines[0])!="---" {return "",errors.New("YAML frontmatter not found")}
 quoted:="'"+strings.ReplaceAll(value,"'","''")+"'"
 for i:=1;i<len(lines);i++ {
  if strings.TrimSpace(lines[i])=="---" {return "",errors.New("field not found; add it in raw editor")}
  if strings.HasPrefix(lines[i],key+":") {
   old:=strings.TrimSpace(strings.TrimPrefix(lines[i],key+":"))
   if old==""||strings.HasPrefix(old,"|")||strings.HasPrefix(old,">"){return "",errors.New("complex YAML field cannot be edited here")}
   lines[i]=key+": "+quoted
   return strings.Join(lines,"\n"),nil
  }
 }
 return "",errors.New("frontmatter not terminated")
}
func PrepareAgentFields(root,relative,hash,name,description string)(Document,error) {
 fields,err:=ReadAgentFields(root,relative);if err!=nil{return Document{},err}
 if fields.Hash!=hash{return Document{},errors.New("file changed on disk")}
 doc,err:=Read(root,relative);if err!=nil{return Document{},err}
 content:=doc.Content
 if name!=fields.Name {content,err=updateScalar(content,"name",name);if err!=nil{return Document{},err}}
 if description!=fields.Description {content,err=updateScalar(content,"description",description);if err!=nil{return Document{},err}}
 return Document{Path:relative,Hash:hash,Content:content},nil
}
