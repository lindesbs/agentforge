# Local development

## Prerequisites
Go 1.23+, Node.js and npm, Wails v2 CLI and platform-specific Wails prerequisites (GTK/WebKit on Linux). Follow the official Wails v2 installation instructions.

## Run
```sh
cd frontend
npm install
npm run build
cd ..
go test ./internal/...
wails dev
```

The checked-in `frontend/dist/index.html` is a placeholder solely to satisfy Go embed before the first frontend build. Vite replaces it with the compiled interface.

## Current scope
The initial inspector takes an explicit directory path, reads file names and basic framework markers, and lists recognized project-local configuration files. It does not parse agent prompts, modify files, execute commands, or automatically scan the home directory.

## Known limitations
- A native OS directory chooser is available via the Browse button; manual paths remain supported.
- The inspector intentionally does not follow symlinks and does not recurse through arbitrary project files.
- Codex/Claude provider paths are an initial subset, not a complete provider compatibility guarantee.
- GitHub Actions now runs Go inspector tests, Vue typechecking/build, and a Linux desktop compile. Successful CI status must be verified on GitHub; macOS and Windows packaging are not yet covered.
