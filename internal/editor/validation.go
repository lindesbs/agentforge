package editor

import (
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "strings"
)

// Issue is a non-destructive warning or error for a proposed configuration.
// It does not claim comprehensive vendor-specific schema validation.
type Issue struct {
 Severity string `json:"severity"`
 Message string `json:"message"`
}

func Validate(relative, content string) []Issue {
 issues:=make([]Issue,0)
 if strings.HasSuffix(relative,".json") {
  decoder:=json.NewDecoder(strings.NewReader(content))
  var value any
  if err:=decoder.Decode(&value);err!=nil {
   issues=append(issues,Issue{"error",fmt.Sprintf("Invalid JSON: %v",err)})
   return issues
  }
  var extra any
  if err:=decoder.Decode(&extra);err!=io.EOF {
   if err==nil {issues=append(issues,Issue{"error","JSON must contain exactly one top-level value"})} else {issues=append(issues,Issue{"error",fmt.Sprintf("Invalid trailing JSON: %v",err)})}
  }
  if _,ok:=value.(map[string]any);!ok {
   issues=append(issues,Issue{"error","Settings must be a JSON object"})
  }
 } else if strings.HasSuffix(relative,".md") {
  if strings.TrimSpace(content)=="" {issues=append(issues,Issue{"warning","Markdown configuration is empty"})}
  if strings.HasPrefix(filepathSlash(relative),".claude/agents/") {
   if !strings.HasPrefix(content,"---\n") && !strings.HasPrefix(content,"---\r\n") {
    issues=append(issues,Issue{"warning","Claude agent has no YAML frontmatter; structured editing is unavailable"})
   }
  }
 } else if strings.HasSuffix(relative,".toml") {
  issues=append(issues,Issue{"info","TOML syntax and provider-specific options are not yet validated"})
 }
 return issues
}

func validateForSave(relative,content string)error {
 for _,issue:=range Validate(relative,content) {
  if issue.Severity=="error" {return errors.New(issue.Message)}
 }
 return nil
}
