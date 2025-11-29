// Package main implements the parm CLI application.
package main

import (
	"fmt"
	"os"

	"github.com/glassmonkey/gha-cleaner"
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

	cfg := parm.Config{
		Recursive: recursive,
		Force:     force,
		Verbose:   verbose,
	}

	return parm.RemovePathsInParallel(args, cfg)
}
