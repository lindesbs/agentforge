package main

import (
 "context"
 "errors"

 "github.com/lindesbs/agentforge/internal/inspector"
 "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
 ctx context.Context
 inspector *inspector.Service
}

func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// SelectProjectDirectory opens the native OS directory chooser.
// Cancellation returns an empty path without changing any project files.
func (a *App) SelectProjectDirectory() (string, error) {
 if a.ctx == nil { return "", errors.New("desktop application is not ready") }
 return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
  Title: "Select an AgentForge project",
 })
}

// InspectProject only reads metadata and configuration file names.
func (a *App) InspectProject(path string) (inspector.Project, error) {
 return a.inspector.Inspect(path)
}
