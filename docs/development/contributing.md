# Contributing

## Development Setup

```bash
# Clone
git clone https://github.com/grokify/vibium-go
cd vibium-go

# Install dependencies
go mod download

# Install clicker
npm install -g vibium

# Run tests
go test -v ./...
```

## Code Style

- Use `gofmt` for formatting
- Follow standard Go conventions
- Run `golangci-lint run` before committing

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: resolve bug
docs: update documentation
refactor: restructure code
test: add tests
chore: maintenance
```

## Pull Requests

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run linter and tests
5. Submit PR with description

## Project Structure

```
vibium-go/
├── *.go              # Core client library
├── cmd/vibium/       # CLI
├── mcp/              # MCP server
├── script/           # Script format
├── integration/      # Integration tests
├── docs/             # Documentation
└── examples/         # Example code
```

## Adding Features

### New Element Method

1. Add method to `element.go`
2. Add recording to `mcp/recorder.go`
3. Add MCP tool to `mcp/tools_*.go`
4. Register in `mcp/server.go`
5. Add to script actions in `script/types.go`
6. Regenerate schema
7. Update documentation

### New MCP Tool

1. Add handler to `mcp/tools_*.go`
2. Register in `mcp/server.go`
3. Add recording call if applicable
4. Update `docs/reference/mcp-tools.md`

### New Script Action

1. Add action to `script/types.go`
2. Add to `AllActions()` function
3. Regenerate schema: `go run ./cmd/genscriptschema > script/vibium-script.schema.json`
4. Add handler in `cmd/vibium/cmd/run.go`
5. Update `docs/guide/scripts.md`

## Regenerating Schema

```bash
go run ./cmd/genscriptschema > script/vibium-script.schema.json
schemago lint script/vibium-script.schema.json
```

## Documentation

Documentation uses MkDocs with Material theme.

### Local Preview

```bash
pip install mkdocs-material
mkdocs serve
```

### Build

```bash
mkdocs build
```

## Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create PR for review
4. Merge to main
5. Tag release: `git tag vX.Y.Z`
6. Push tag: `git push origin vX.Y.Z`
