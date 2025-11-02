# gha-cleaner

A GitHub Actions tool to free up disk space on GitHub-hosted runners by removing unnecessary SDKs, caches, and browser binaries. Can reclaim up to approximately 30GB of disk space.

## Acknowledgments

This project is inspired by [mathio/gha-cleanup](https://github.com/mathio/gha-cleanup/). The key differences are:

1. **Implementation language**: Go instead of shell scripts
2. **Execution model**: Parallel deletion (parm) using goroutines instead of sequential `rm` commands
3. **Performance**: Faster cleanup through concurrent directory removal

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

### Testing the Action

This repository includes automated tests for the action using a reusable workflow pattern:

#### Workflow Structure

- **_verify.yml** (Reusable workflow): Core verification logic that can be called by other workflows (prefixed with `_` to indicate it's a shared workflow)
- **ci.yml**: Runs verification automatically on PRs
- **manual.yml**: Allows manual execution via workflow_dispatch

#### 1. Continuous Integration (Automatic)

The `.github/workflows/ci.yml` workflow runs automatically on pull requests when relevant files are changed. It verifies three scenarios in parallel:
- Default settings
- Verbose mode
- Browser removal with verbose mode

#### 2. Manual Execution

The `.github/workflows/manual.yml` workflow can be triggered manually from the Actions tab:

1. Go to the "Actions" tab in GitHub
2. Select "Manual Run" workflow
3. Click "Run workflow"
4. Choose your configuration:
   - **Runner OS**: Select which runner to use (ubuntu-latest, macos-latest, etc.)
   - **Remove Browsers**: Enable/disable browser removal
   - **Verbose**: Enable/disable verbose logging

#### Benefits of Reusable Workflows

- **DRY principle**: Verification logic is defined once in _verify.yml
- **Consistency**: All workflows use the same verification procedure
- **Maintainability**: Updates to verification logic only need to be made in one place
- **Flexibility**: Easy to add new scenarios by calling the reusable workflow with different parameters

## Security

### Security Considerations for Contributors

This project includes test workflows that execute code from pull requests, including external forks. Please be aware of the following security considerations:

#### CI Workflow

The `ci.yml` workflow runs code from PRs using `pull_request` trigger with `uses: ./`. This means:

- ✅ **Limited permissions**: Workflows have read-only access to repository contents
- ⚠️ **Code execution risk**: Malicious code in PRs will be executed during tests
- ⚠️ **sudo usage**: The cleanup.sh script uses `sudo`, which could be exploited

#### Recommended Security Practices

For repository maintainers:

1. **Enable manual approval for external PRs**:
   - Go to Settings → Actions → General
   - Under "Fork pull request workflows from outside collaborators"
   - Select "Require approval for first-time contributors"

2. **Review PR code before approving**:
   - Always review `action.yml`, `parm.go`, `cleanup.sh` changes
   - Check for suspicious commands or network calls
   - Look for attempts to exfiltrate data or escalate privileges

3. **Monitor workflow runs**:
   - Watch for unexpected behavior during test runs
   - Check for unusual network activity or long-running jobs

4. **Keep dependencies pinned**:
   - Use specific commit SHAs for GitHub Actions (e.g., `actions/checkout@<sha>`)
   - Currently using tags for convenience, but SHAs are more secure

#### For Users of This Action

When using this action in your workflows:

- Review the source code before using it
- Pin to specific versions/tags instead of `@main`
- Be aware that the action removes system files (by design)
- Test in a safe environment first

### Reporting Security Issues

If you discover a security vulnerability, please email the maintainer directly instead of opening a public issue. See SECURITY.md for details.

## License

MIT License