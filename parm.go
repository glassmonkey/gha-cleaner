// Package parm implements parm (parallel rm) - a rm-compatible library with parallel execution.
// parm provides core functionality for removing files and directories in parallel using goroutines.
package parm

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Config holds configuration for removal operations
type Config struct {
	Recursive bool
	Force     bool
	Verbose   bool
}

// RemovalTarget represents a single target to be removed
type RemovalTarget struct {
	Path  string
	IsDir bool
}

// RemovePathsInParallel removes multiple paths in parallel.
// It processes paths through three stages: parsing, validation, and removal.
func RemovePathsInParallel(paths []string, cfg Config) error {
	// Stage 1: Collect and validate removal targets
	targets, hadCollectionErrors := collectRemovalTargets(paths, cfg)

	// Stage 2: Remove targets in parallel
	removalErr := removeInParallel(targets, cfg)

	// Return collection errors if force is not set and there were errors
	if hadCollectionErrors && !cfg.Force {
		return fmt.Errorf("some files could not be removed")
	}

	return removalErr
}

// collectRemovalTargets validates and collects removal targets.
// It expands glob patterns and validates each path.
// Returns the list of targets and a flag indicating if validation errors occurred.
func collectRemovalTargets(patterns []string, cfg Config) ([]RemovalTarget, bool) {
	var targets []RemovalTarget
	hadErrors := false

	for _, pattern := range patterns {
		// Try to expand glob patterns
		matches, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parm: cannot remove '%s': %v\n", pattern, err)
			hadErrors = true
			continue
		}

		// If no matches, treat the pattern as a literal path
		if len(matches) == 0 {
			matches = []string{pattern}
		}

		// Process each matched path
		for _, path := range matches {
			info, err := os.Lstat(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parm: cannot remove '%s': %v\n", path, err)
				hadErrors = true
				continue
			}

			// Skip directories if recursive flag is not set
			if info.IsDir() && !cfg.Recursive {
				fmt.Fprintf(os.Stderr, "parm: cannot remove '%s': is a directory\n", path)
				hadErrors = true
				continue
			}

			targets = append(targets, RemovalTarget{
				Path:  path,
				IsDir: info.IsDir(),
			})
		}
	}

	return targets, hadErrors
}

// removeInParallel removes collected targets in parallel.
func removeInParallel(targets []RemovalTarget, cfg Config) error {
	if len(targets) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	hasError := false

	for _, target := range targets {
		wg.Add(1)
		go func(t RemovalTarget) {
			defer wg.Done()
			if err := removePath(t.Path); err != nil {
				if !cfg.Force {
					mu.Lock()
					fmt.Fprintf(os.Stderr, "parm: cannot remove '%s': %v\n", t.Path, err)
					hasError = true
					mu.Unlock()
				}
			} else if cfg.Verbose {
				mu.Lock()
				fmt.Printf("removed '%s'\n", t.Path)
				mu.Unlock()
			}
		}(target)
	}

	wg.Wait()

	if hasError && !cfg.Force {
		return fmt.Errorf("some files could not be removed")
	}
	return nil
}

// removePath removes a single file or directory.
// The path must have already been validated and approved for removal.
func removePath(path string) error {
	// Check if the path is a directory
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return os.RemoveAll(path)
	}

	return os.Remove(path)
}
