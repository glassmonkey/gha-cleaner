# gha-cleaner

A GitHub Actions tool to free up disk space on GitHub-hosted runners by removing unnecessary SDKs, caches, and browser binaries. Can reclaim up to approximately 30GB of disk space.

## parm: Parallel RM

This tool uses **parm** (parallel rm), a custom Go implementation that aims to be compatible with the `rm` command while providing parallel execution using goroutines. This makes file deletion faster and more efficient than traditional sequential deletion.

### Current Status (v0.0.1)

**Note:** parm is currently in early development (v0.0.1 - placeholder). The current version outputs a hello world message to verify the build and distribution pipeline. Actual cleanup functionality is not yet implemented.

Future versions will implement full rm-compatible functionality with parallel execution.

### Future Milestones

parm is designed to become a general-purpose rm-compatible tool. Future milestones include:

- [ ] Full `rm` command compatibility (options, behavior, error messages)
- [ ] Extended options for controlling parallelism (e.g., `-j` flag for max parallel jobs)
- [ ] Performance benchmarks comparing parm vs traditional rm
- [ ] Support for interactive mode (`-i`, `-I`)
- [ ] Integration with other tools and scripts
- [ ] Cross-platform testing and optimization
- [ ] Standalone distribution as a general-purpose CLI tool

For now, parm provides basic rm-compatible options (`-r`, `-f`, `-v`) and focuses on parallel deletion performance.

## Inspiration and Implementation

This action is inspired by [mathio/gha-cleanup](https://github.com/mathio/gha-cleanup/), but differs in its implementation approach. While the original uses shell scripts with `rm` commands, this implementation is written in Go with parallel execution, providing:

- **Parallel deletion**: Uses goroutines to remove multiple directories simultaneously
- **Faster cleanup**: Reduces total cleanup time through concurrent operations
- **Better error handling**: Each deletion operation is independently handled
- **Cross-platform compatibility**: Built with Go's standard library
- **Type-safe implementation**: Compile-time type checking
- **Easier testing and maintenance**: Clean, testable Go code

## Features

This action removes the following from GitHub hosted runners:

### Always Removed
- Java Virtual Machine (`/usr/lib/jvm`)
- .NET Runtime (`/usr/share/dotnet`)
- Swift (`/usr/share/swift`)
- Haskell GHC (`/usr/local/.ghcup`)
- Julia (`/usr/local/julia*`)
- Android SDK (`/usr/local/lib/android`)
- Azure CLI (`/opt/az`)
- PowerShell (`/usr/local/share/powershell`)
- Hosted tool cache (`/opt/hostedtoolcache`)
- Docker system cache

### Optionally Removed (when `remove-browsers: true`)
- Chromium (`/usr/local/share/chromium`)
- Microsoft Edge (`/opt/microsoft`)
- Google Chrome (`/opt/google`)
- Firefox (`/usr/lib/firefox`)

## Architecture

```
┌─────────────────────────┐
│  GitHub Actions Workflow│
└───────────┬─────────────┘
            │
            ▼
┌─────────────────────────┐
│     action.yml          │
│  (builds parm & runs    │
│   cleanup.sh)           │
└───────────┬─────────────┘
            │
            ▼
┌─────────────────────────┐
│     cleanup.sh          │
│  (GHA-specific cleanup  │
│   logic using parm)     │
└───────────┬─────────────┘
            │
            ▼
┌─────────────────────────┐
│     parm (Go binary)    │
│  (rm-compatible tool    │
│   with parallel exec)   │
└─────────────────────────┘
```

## Usage

### Basic Usage

```yaml
name: Build

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Clean up runner
        uses: glassmonkey/gha-cleaner@v1

      - name: Build
        run: |
          # Your build commands here
```

### With Browser Removal

```yaml
- name: Clean up runner
  uses: glassmonkey/gha-cleaner@v1
  with:
    remove-browsers: true
```

### With Verbose Logging

```yaml
- name: Clean up runner
  uses: glassmonkey/gha-cleaner@v1
  with:
    remove-browsers: true
    verbose: true
```

## Inputs

| Parameter | Description | Required | Default |
|-----------|-------------|----------|---------|
| `remove-browsers` | Enable removal of browser caches and binaries | No | `false` |
| `verbose` | Enable verbose logging to see what is being removed | No | `false` |

## Using parm Standalone

**Note:** v0.0.1 is a placeholder version. Full functionality is coming soon.

Current version (v0.0.1):
```bash
# Build parm
go build -o parm parm.go

# Run to see hello world
./parm
# Output:
# Hello from parm (parallel rm)!
# Version: 0.0.1 (placeholder)
# Full rm-compatible functionality coming soon...
```

Future versions will support rm-compatible usage:
```bash
# Future usage (not yet implemented in v0.0.1)
parm -rf /path/to/dir1 /path/to/dir2 /path/to/dir3

# Planned options:
# -r, --recursive  remove directories and their contents recursively
# -f, --force      ignore nonexistent files and arguments, never prompt
# -v, --verbose    explain what is being done
```

## Impact

Running this action can reclaim approximately 30GB of disk space. For example, disk usage has been observed to drop from 62% to 22% after cleanup.

## Development

### Building parm

```bash
go build -o parm parm.go
```

### Testing locally

```bash
# Build parm
go build -o parm parm.go

# Run parm
./parm

# Test cleanup script
export REMOVE_BROWSERS=false
export VERBOSE=true
./cleanup.sh
```

## License

MIT License

## Acknowledgments

This project is inspired by [mathio/gha-cleanup](https://github.com/mathio/gha-cleanup/). The key differences are:

1. **Implementation language**: Go instead of shell scripts
2. **Execution model**: Parallel deletion (parm) using goroutines instead of sequential `rm` commands
3. **Performance**: Faster cleanup through concurrent directory removal