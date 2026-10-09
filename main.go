package main

import (
 "embed"
 "io/fs"
 "log"
 "os"
 "errors"
 "flag"
 "fmt"

 "github.com/lindesbs/agentforge/internal/launch"

 "github.com/lindesbs/agentforge/internal/inspector"
 "github.com/wailsapp/wails/v2"
 "github.com/wailsapp/wails/v2/pkg/options"
 "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var embeddedAssets embed.FS

func main() {
 selected, parseErr := launch.ProjectPath(os.Args[1:])
 if errors.Is(parseErr,flag.ErrHelp) {
  fmt.Fprintln(os.Stdout,"Usage: agentforge [PROJECT_DIRECTORY]")
  fmt.Fprintln(os.Stdout,"Open an optional existing local project; relative paths are resolved from the current working directory.")
  return
 }
 if parseErr!=nil {
  log.Fatalf("Invalid project argument: %v (usage: agentforge [PROJECT_DIRECTORY])",parseErr)
 }
 assets, err := fs.Sub(embeddedAssets, "frontend/dist")
 if err != nil { log.Fatal(err) }
 app := &App{inspector: inspector.New(), startupProject: selected}
 if err := wails.Run(&options.App{
  Title: "AgentForge",
  Width: 1100,
  Height: 750,
  AssetServer: &assetserver.Options{Assets: assets},
  Bind: []interface{}{app},
  OnStartup: app.startup,
 }); err != nil { log.Fatal(err) }
}
