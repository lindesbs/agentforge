package main

import (
 "context"
 "github.com/lindesbs/agentforge/internal/inspector"
)

type App struct {
 ctx context.Context
 inspector *inspector.Service
}
func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// InspectProject reads metadata and configuration file names only.
// It never executes project code or returns configuration contents.
func (a *App) InspectProject(path string) (inspector.Project, error) {
 return a.inspector.Inspect(path)
}
