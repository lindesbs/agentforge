package main

import (
 "embed"
 "log"

 "github.com/lindesbs/agentforge/internal/inspector"
 "github.com/wailsapp/wails/v2"
 "github.com/wailsapp/wails/v2/pkg/options"
 "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
 app := &App{inspector: inspector.New()}
 if err := wails.Run(&options.App{
  Title: "AgentForge",
  Width: 1100,
  Height: 750,
  AssetServer: &assetserver.Options{Assets: assets},
  Bind: []interface{}{app},
 }); err != nil { log.Fatal(err) }
}
