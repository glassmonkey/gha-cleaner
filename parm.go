// Package main implements parm (parallel rm) - a rm-compatible tool with parallel execution.
// parm aims to be compatible with rm command but executes deletions in parallel using goroutines.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
)

var (
	recursive bool
	force     bool
	verbose   bool
)

var rootCmd = &cobra.Command{
	Use:   "parm [flags] FILE...",
	Short: "Parallel rm - remove files and directories in parallel",
	Long: `parm is a rm-compatible tool that removes files and directories in parallel using goroutines.
This tool aims to be compatible with rm command while providing better performance through parallel execution.`,
	Example: `  parm -rf /path/to/dir1 /path/to/dir2
  parm -v file1.txt file2.txt
  parm -r -f /usr/local/lib/android`,
	RunE: runRemove,
	SilenceUsage: true,
}

func init() {
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "remove directories and their contents recursively")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "ignore nonexistent files and arguments, never prompt")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "explain what is being done")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRemove(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		if !force {
			return fmt.Errorf("missing operand")
		}
		return nil
	}

	return removePathsInParallel(args)
}

func removePathsInParallel(paths []string) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	hasError := false

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if err := removePath(p, recursive); err != nil {
				if !force {
					mu.Lock()
					fmt.Fprintf(os.Stderr, "parm: cannot remove '%s': %v\n", p, err)
					hasError = true
					mu.Unlock()
				}
			} else if verbose {
				mu.Lock()
				fmt.Printf("removed '%s'\n", p)
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()

	if hasError && !force {
		return fmt.Errorf("some files could not be removed")
	}
	return nil
}

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

// expandPath expands wildcards and returns list of matching paths
func expandPath(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		// If no matches, return the original pattern
		return []string{pattern}, nil
	}
	return matches, nil
}
