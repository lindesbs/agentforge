package learning

import (
 "errors"
 "fmt"
 "strings"
)

// ImportCandidate is only a preview; the source is read again before import.
// The source document remains unchanged.
type ImportCandidate struct {
 Category string `json:"category"`
 Index int `json:"index"`
 SourcePath string `json:"sourcePath"`
 SourceHash string `json:"sourceHash"`
 Entry Entry `json:"entry"`
 Duplicate bool `json:"duplicate"`
}
type ImportSkipped struct {
 SourcePath string `json:"sourcePath"`
 Heading string `json:"heading"`
 Reason string `json:"reason"`
}
type ImportScan struct {
 Candidates []ImportCandidate `json:"candidates"`
 Skipped []ImportSkipped `json:"skipped"`
}
type section struct {title string; lines []string}
func parseSections(text string) []section {
 var sections []section
 var current *section
 for _,line:=range strings.Split(strings.ReplaceAll(text,"\r\n","\n"),"\n") {
  trimmed:=strings.TrimSpace(line)
  // Accept headings used in existing learning documents, and generated templates.
  heading:=""
  if strings.HasPrefix(trimmed,"### ") {heading=strings.TrimSpace(strings.TrimPrefix(trimmed,"### "))}
  if strings.HasPrefix(trimmed,"## ") {heading=strings.TrimSpace(strings.TrimPrefix(trimmed,"## "))}
  if heading!="" {
   if current!=nil {sections=append(sections,*current)}
   current=&section{title:heading}
  } else if current!=nil {current.lines=append(current.lines,line)}
 }
 if current!=nil {sections=append(sections,*current)}
 return sections
}
func classifyLabel(text string)(string,string){
 value:=strings.TrimSpace(text)
 labels:=[]struct{prefix,key string}{
  {"**Sprache:**","language"},
  {"**Bestätigte Erkenntnis/Ursache:**","finding"},
  {"**Bestätigte Erkenntnis:**","finding"},
  {"**Erkenntnis:**","finding"},
  {"**Ursache:**","finding"},
  {"**Verbindliche Regel:**","rule"},
  {"**Regel:**","rule"},
  {"**Prüfung:**","check"},
  {"**Kontrolle:**","check"},
 }
 for _,item:=range labels {
  if strings.HasPrefix(value,item.prefix) {return item.key,strings.TrimSpace(strings.TrimPrefix(value,item.prefix))}
 }
 return "",""
}
func parseEntry(category string,s section)(Entry,error){
 e:=Entry{Category:category,Title:s.title,Language:"Allgemein"}
 fields:=map[string]string{}
 current:=""
 for _,line:=range s.lines {
  key,value:=classifyLabel(line)
  if key!="" {
   current=key
   if existing:=fields[key];existing!="" {return Entry{},fmt.Errorf("duplicate %s field",key)}
   fields[key]=value
  } else if current!="" && strings.TrimSpace(line)!="" {
   fields[current]+="\n"+strings.TrimSpace(line)
  }
 }
 if fields["language"]!="" {e.Language=strings.TrimSpace(fields["language"])}
 e.Finding=strings.TrimSpace(fields["finding"])
 e.Rule=strings.TrimSpace(fields["rule"])
 e.Check=strings.TrimSpace(fields["check"])
 if err:=validate(e);err!=nil{return Entry{},err}
 return e,nil
}
func duplicateEntry(existing []Entry,e Entry)bool {
 for _,item:=range existing {
  if item.Category==e.Category && item.Language==e.Language &&
   strings.EqualFold(strings.TrimSpace(item.Title),strings.TrimSpace(e.Title)){return true}
 }
 return false
}

// ScanProject imports nothing; it reads only the four fixed learning files.
func (l *Library) ScanProject(root string)(ImportScan,error){
 existing,err:=l.List();if err!=nil{return ImportScan{},err}
 result:=ImportScan{Candidates:[]ImportCandidate{},Skipped:[]ImportSkipped{}}
 for _,category:=range []string{"PROJEKT_OK","PROJEKT_NOK","ALLGEMEIN_OK","ALLGEMEIN_NOK"} {
  path,relative,err:=docPath(root,category);if err!=nil{return ImportScan{},err}
  data,hash,err:=loadDoc(path);if err!=nil{return ImportScan{},err}
  if hash==missingHash {continue}
  for index,sec:=range parseSections(string(data)) {
   entry,err:=parseEntry(category,sec)
   if err!=nil {
    result.Skipped=append(result.Skipped,ImportSkipped{SourcePath:relative,Heading:sec.title,Reason:err.Error()})
    continue
   }
   result.Candidates=append(result.Candidates,ImportCandidate{
    Category:category,Index:index,SourcePath:relative,SourceHash:hash,Entry:entry,
    Duplicate:duplicateEntry(existing,entry),
   })
  }
 }
 return result,nil
}

// ImportProjectEntry rejects stale previews and trusts the on-disk Markdown,
// never the content sent back by the browser.
func (l *Library) ImportProjectEntry(root,category string,index int,expectedHash,language string)(Entry,error){
 if _,err:=categoryPath(category);err!=nil{return Entry{},err}
 if index<0{return Entry{},errors.New("invalid section index")}
 if !languages[language]{return Entry{},errors.New("invalid selected language")}
 path,_,err:=docPath(root,category);if err!=nil{return Entry{},err}
 contents,hash,err:=loadDoc(path);if err!=nil{return Entry{},err}
 if hash==missingHash||hash!=expectedHash{return Entry{},errors.New("source changed; scan again")}
 sections:=parseSections(string(contents))
 if index>=len(sections){return Entry{},errors.New("source entry no longer exists")}
 entry,err:=parseEntry(category,sections[index]);if err!=nil{return Entry{},err}
 entry.Language=language
 if err:=validate(entry);err!=nil{return Entry{},err}
 return l.Save(entry)
}
