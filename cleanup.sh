#!/bin/bash
# GitHub Actions runner cleanup script
# Note: v0.0.1 is a placeholder version for testing the build and distribution pipeline

set -e

# Parse options
REMOVE_BROWSERS="${REMOVE_BROWSERS:-false}"
VERBOSE="${VERBOSE:-false}"

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARM="$SCRIPT_DIR/parm"

echo "=========================================="
echo "gha-cleaner v0.0.1 - Testing parm binary"
echo "=========================================="
echo ""

# Run parm (currently a hello world placeholder)
if [ -f "$PARM" ]; then
    echo "Running parm:"
    "$PARM" || true
    echo ""
else
    echo "Error: parm binary not found at $PARM"
    exit 1
fi

echo "Configuration:"
echo "  REMOVE_BROWSERS: $REMOVE_BROWSERS"
echo "  VERBOSE: $VERBOSE"
echo ""

echo "=========================================="
echo "Note: Actual cleanup functionality will be"
echo "implemented in future versions."
echo ""
echo "Future cleanup targets will include:"
echo "  - Java, .NET, Swift, Haskell, Julia"
echo "  - Android SDK, Azure CLI, PowerShell"
echo "  - Hosted tool cache"
if [ "$REMOVE_BROWSERS" == "true" ]; then
    echo "  - Browsers (Chromium, Edge, Chrome, Firefox)"
fi
echo "  - Docker system cache"
echo "=========================================="
