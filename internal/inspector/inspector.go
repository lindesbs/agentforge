package inspector

import (
 "errors"
 "os"
 "path/filepath"
 "sort"
 "strings"
)

type ConfigFile struct {
 Path string `json:"path"`
 Provider string `json:"provider"`
 Kind string `json:"kind"`
}

type Project struct {
 Root string `json:"root"`
 Name string `json:"name"`
 Frameworks []string `json:"frameworks"`
 Files []ConfigFile `json:"files"`
}

type Service struct{}
func New() *Service { return &Service{} }

// Inspect accepts a local project root. No files are executed, and symlinks
// are not traversed. This first milestone only discovers names and metadata.
func (s *Service) Inspect(root string) (Project, error) {
 if strings.TrimSpace(root)=="" { return Project{}, errors.New("project path is required") }
 abs, err := filepath.Abs(root); if err!=nil {return Project{},err}
 info, err:=os.Lstat(abs); if err!=nil{return Project{},err}
 if info.Mode()&os.ModeSymlink!=0 {return Project{},errors.New("symlink project roots are not supported")}
 if !info.IsDir() {return Project{},errors.New("project path is not a directory")}
 p:=Project{Root:abs,Name:filepath.Base(abs),Frameworks:[]string{},Files:[]ConfigFile{}}
 isRegular:=func(rel string)bool{
  fi,e:=os.Lstat(filepath.Join(abs,rel))
  return e==nil && fi.Mode().IsRegular()
 }
 if isRegular("composer.json") {p.Frameworks=append(p.Frameworks,"PHP/Composer")}
 if isRegular("go.mod") {p.Frameworks=append(p.Frameworks,"Go")}
 if isRegular("package.json") {p.Frameworks=append(p.Frameworks,"Node.js")}
 if isRegular("AGENTS.md") {p.Files=append(p.Files,ConfigFile{"AGENTS.md","codex","instructions"})}
 if isRegular("CLAUDE.md") {p.Files=append(p.Files,ConfigFile{"CLAUDE.md","claude","instructions"})}
 if isRegular(filepath.Join(".codex","config.toml")) {p.Files=append(p.Files,ConfigFile{".codex/config.toml","codex","settings"})}
 if isRegular(filepath.Join(".claude","settings.json")) {p.Files=append(p.Files,ConfigFile{".claude/settings.json","claude","settings"})}
 if isRegular(filepath.Join(".claude","settings.local.json")) {p.Files=append(p.Files,ConfigFile{".claude/settings.local.json","claude","settings"})}
 agents:=filepath.Join(abs,".claude","agents")
 if di,e:=os.Lstat(agents);e==nil&&di.IsDir()&&di.Mode()&os.ModeSymlink==0 {
  entries,e:=os.ReadDir(agents);if e!=nil{return Project{},e}
  for _,entry:=range entries {
   if !strings.HasSuffix(strings.ToLower(entry.Name()),".md") {continue}
   if !isRegular(filepath.Join(".claude","agents",entry.Name())) {continue}
   p.Files=append(p.Files,ConfigFile{filepath.ToSlash(filepath.Join(".claude","agents",entry.Name())),"claude","agent"})
  }
 }
 sort.Slice(p.Files,func(i,j int)bool{return p.Files[i].Path<p.Files[j].Path})
 return p,nil
}
