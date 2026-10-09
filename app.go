package main

import (
	"context"
	"errors"

	"github.com/lindesbs/agentforge/internal/editor"
	"github.com/lindesbs/agentforge/internal/inspector"
	"github.com/lindesbs/agentforge/internal/learning"
	"github.com/lindesbs/agentforge/internal/providers"
	"github.com/lindesbs/agentforge/internal/templates"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	inspector *inspector.Service
}

func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// SelectProjectDirectory opens the native OS directory chooser.
// Cancellation returns an empty path without changing any project files.
func (a *App) SelectProjectDirectory() (string, error) {
	if a.ctx == nil {
		return "", errors.New("desktop application is not ready")
	}
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select an AgentForge project",
	})
}

// InspectProject only reads metadata and configuration file names.
func (a *App) InspectProject(path string) (inspector.Project, error) {
	return a.inspector.Inspect(path)
}

func (a *App) ReadConfig(root, path string) (editor.Document, error) {
	return editor.Read(root, path)
}
func (a *App) PreviewConfig(root, path, hash, content string) (editor.Preview, error) {
	return editor.PreviewChange(root, path, hash, content)
}
func (a *App) SaveConfig(root, path, hash, content string) (editor.Document, error) {
	return editor.Save(root, path, hash, content)
}

func (a *App) ReadAgentFields(root, path string) (editor.AgentFields, error) {
	return editor.ReadAgentFields(root, path)
}
func (a *App) PrepareAgentFields(root, path, hash, name, description string) (editor.Document, error) {
	return editor.PrepareAgentFields(root, path, hash, name, description)
}

func (a *App) ValidateConfig(path, content string) []editor.Issue {
	return editor.Validate(path, content)
}

// DescribeConfig inspects text already opened by the user; it performs no I/O.
func (a *App) DescribeConfig(path, content string) providers.Definition {
	return providers.Inspect(path, content)
}

func (a *App) ProviderCapabilities() []providers.Capability { return providers.Capabilities() }
func (a *App) CheckPortability(path, content string) []providers.Issue {
	return providers.Portability(path, content)
}

func (a *App) ListTemplates() ([]templates.Summary, error) {
	library, err := templates.New("")
	if err != nil {
		return nil, err
	}
	return library.List()
}
func (a *App) SaveProjectAsTemplate(root, path, name string) (templates.Summary, error) {
	library, err := templates.New("")
	if err != nil {
		return templates.Summary{}, err
	}
	item, err := library.CreateFromProject(root, path, name)
	if err != nil {
		return templates.Summary{}, err
	}
	return templates.Summarize(item), nil
}
func (a *App) SaveTaggedTemplate(root, path, name string, tags []string) (templates.Summary, error) {
	library, err := templates.New("")
	if err != nil {
		return templates.Summary{}, err
	}
	item, err := library.CreateTagged(root, path, name, tags)
	if err != nil {
		return templates.Summary{}, err
	}
	return templates.Summarize(item), nil
}
func (a *App) ReviewTemplateImport(payload string) (templates.ImportReview, error) {
	return templates.ReviewImport(payload)
}
func (a *App) ImportTemplate(payload string) (templates.Summary, error) {
	library, err := templates.New("")
	if err != nil {
		return templates.Summary{}, err
	}
	item, err := library.Import(payload)
	if err != nil {
		return templates.Summary{}, err
	}
	return templates.Summarize(item), nil
}
func (a *App) ExportTemplate(id string) (string, error) {
	library, err := templates.New("")
	if err != nil {
		return "", err
	}
	return library.Export(id)
}
func (a *App) PrepareTemplateApply(id, root, targetPath string) (editor.Document, error) {
	library, err := templates.New("")
	if err != nil {
		return editor.Document{}, err
	}
	return library.PrepareApply(id, root, targetPath)
}

func (a *App) ListLearnings() ([]learning.Entry, error) {
	library, err := learning.New("")
	if err != nil {
		return nil, err
	}
	return library.List()
}
func (a *App) SaveLearning(entry learning.Entry) (learning.Entry, error) {
	library, err := learning.New("")
	if err != nil {
		return learning.Entry{}, err
	}
	return library.Save(entry)
}
func (a *App) FormatLearning(entry learning.Entry) string {
	return learning.Markdown(entry)
}
