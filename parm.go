// Package parm implements parm (parallel rm) - a rm-compatible library with parallel execution.
// parm provides core functionality for removing files and directories in parallel using goroutines.
package parm

import (
	"fmt"
	"os"
	"sync"
)

// Config holds configuration for removal operations
type Config struct {
	Recursive bool
	Force     bool
	Verbose   bool
}

// RemovePathsInParallel removes multiple paths in parallel.
// It uses goroutines to remove paths concurrently while maintaining synchronization.
func RemovePathsInParallel(paths []string, cfg Config) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	hasError := false

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if err := removePath(p, cfg.Recursive); err != nil {
				if !cfg.Force {
					mu.Lock()
					fmt.Fprintf(os.Stderr, "parm: cannot remove '%s': %v\n", p, err)
					hasError = true
					mu.Unlock()
				}
			} else if cfg.Verbose {
				mu.Lock()
				fmt.Printf("removed '%s'\n", p)
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()

	if hasError && !cfg.Force {
		return fmt.Errorf("some files could not be removed")
	}
	return nil
}

// removePath removes a single file or directory.
func removePath(path string, recursive bool) error {
	// Get file info to check if it exists and if it's a directory
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("No such file or directory")
		}
		return err
	}

	// If it's a directory and -r is not set, return error
	if info.IsDir() && !recursive {
		return fmt.Errorf("is a directory")
	}

	// If it's a directory with -r flag, use RemoveAll
	// Otherwise use Remove
	if info.IsDir() && recursive {
		return os.RemoveAll(path)
	}

	return os.Remove(path)
}
