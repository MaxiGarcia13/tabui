# Contributing

## Setup

Requires Go 1.24+ and [golangci-lint](https://golangci-lint.run/) for linting.

Once per clone, point Git at the tracked hooks (runs `golangci-lint` before each commit):

```bash
git config core.hooksPath .githooks
```

## Commands

```bash
# run unit tests
go test ./...

# lint (see .golangci.yml)
golangci-lint run ./...

# build the binary
go build -o tabui ./cmd/tabui

# print version (const in cmd/tabui/main.go)
./tabui -version

# install onto PATH
go install ./cmd/tabui
```

Try the package-scripts picker (interactive TTY required):

```bash
go run ./cmd/tabui --line="pnpm run " --cwd=/path/to/node/project
```

On success, the chosen value is printed on stdout (for shell insert). Esc / Ctrl+C cancels with no output.

Provider tests live next to the code (e.g. `internal/package_scripts/`, fixtures in `testdata/`).

## Layout

```text
cmd/tabui/           # binary entry — wire only
internal/<feature>/  # providers, matching, domain types
internal/ui/         # Charm TUI (Bubble Tea / Lip Gloss)
internal/<name>/     # shared helpers only when reused (named packages)
shell/               # zsh / bash keybindings
```

## Providers

Under `internal/<feature>/`:

1. **Match** — does this provider apply to the cmdline context?
2. **Suggest** — return suggestions; never render

| Field | Role |
|-------|------|
| `Label` | Primary (name, filename, subcommand) |
| `Detail` | Secondary (path, script body, description) |
| `Value` | Token inserted into the line |

The TUI only renders `label`, `detail`, and returns the chosen `Value`. Keep filesystem / manifest logic out of `internal/ui`.
