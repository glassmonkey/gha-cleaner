package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemovePath(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() (string, error)
		recursive bool
		wantErr   bool
		errMsg    string
	}{
		{
			name: "remove regular file",
			setup: func() (string, error) {
				tmpFile, err := os.CreateTemp("", "parm-test-*.txt")
				if err != nil {
					return "", err
				}
				tmpFile.Close()
				return tmpFile.Name(), nil
			},
			recursive: false,
			wantErr:   false,
		},
		{
			name: "remove nonexistent file",
			setup: func() (string, error) {
				return "/tmp/parm-nonexistent-file-123456", nil
			},
			recursive: false,
			wantErr:   true,
			errMsg:    "No such file or directory",
		},
		{
			name: "remove directory without recursive flag",
			setup: func() (string, error) {
				tmpDir, err := os.MkdirTemp("", "parm-test-dir-*")
				if err != nil {
					return "", err
				}
				return tmpDir, nil
			},
			recursive: false,
			wantErr:   true,
			errMsg:    "is a directory",
		},
		{
			name: "remove directory with recursive flag",
			setup: func() (string, error) {
				tmpDir, err := os.MkdirTemp("", "parm-test-dir-*")
				if err != nil {
					return "", err
				}
				// Create a file inside the directory
				testFile := filepath.Join(tmpDir, "test.txt")
				if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
					return "", err
				}
				return tmpDir, nil
			},
			recursive: true,
			wantErr:   false,
		},
		{
			name: "remove empty directory with recursive flag",
			setup: func() (string, error) {
				tmpDir, err := os.MkdirTemp("", "parm-test-empty-*")
				if err != nil {
					return "", err
				}
				return tmpDir, nil
			},
			recursive: true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := tt.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Clean up if the test doesn't remove it
			defer func() {
				if _, err := os.Stat(path); err == nil {
					os.RemoveAll(path)
				}
			}()

			err = removePath(path, tt.recursive)

			if tt.wantErr {
				if err == nil {
					t.Errorf("removePath() expected error but got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("removePath() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("removePath() unexpected error = %v", err)
				}
				// Verify the path was actually removed
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Errorf("removePath() did not remove the path")
				}
			}
		})
	}
}

func TestRemovePathsInParallel(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() ([]string, error)
		wantErr   bool
		setForce  bool
		setVerbose bool
	}{
		{
			name: "remove multiple files in parallel",
			setup: func() ([]string, error) {
				var paths []string
				for i := 0; i < 5; i++ {
					tmpFile, err := os.CreateTemp("", "parm-test-*.txt")
					if err != nil {
						return nil, err
					}
					tmpFile.Close()
					paths = append(paths, tmpFile.Name())
				}
				return paths, nil
			},
			wantErr: false,
		},
		{
			name: "remove with some nonexistent files and force flag",
			setup: func() ([]string, error) {
				tmpFile, err := os.CreateTemp("", "parm-test-*.txt")
				if err != nil {
					return nil, err
				}
				tmpFile.Close()
				return []string{
					tmpFile.Name(),
					"/tmp/nonexistent-123",
					"/tmp/nonexistent-456",
				}, nil
			},
			wantErr:  false,
			setForce: true,
		},
		{
			name: "remove with some nonexistent files without force flag",
			setup: func() ([]string, error) {
				tmpFile, err := os.CreateTemp("", "parm-test-*.txt")
				if err != nil {
					return nil, err
				}
				tmpFile.Close()
				return []string{
					tmpFile.Name(),
					"/tmp/nonexistent-789",
				}, nil
			},
			wantErr:  true,
			setForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global flags
			oldForce := force
			oldRecursive := recursive
			oldVerbose := verbose
			defer func() {
				force = oldForce
				recursive = oldRecursive
				verbose = oldVerbose
			}()

			force = tt.setForce
			recursive = true
			verbose = tt.setVerbose

			paths, err := tt.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Clean up any remaining paths
			defer func() {
				for _, path := range paths {
					os.RemoveAll(path)
				}
			}()

			err = removePathsInParallel(paths)

			if tt.wantErr {
				if err == nil {
					t.Errorf("removePathsInParallel() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("removePathsInParallel() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestExpandPath(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "parm-expand-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{"test1.txt", "test2.txt", "data.csv"}
	for _, file := range testFiles {
		path := filepath.Join(tmpDir, file)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	tests := []struct {
		name       string
		pattern    string
		wantCount  int
		wantErr    bool
	}{
		{
			name:      "expand wildcard pattern matching multiple files",
			pattern:   filepath.Join(tmpDir, "test*.txt"),
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "expand pattern with no matches",
			pattern:   filepath.Join(tmpDir, "nonexistent*.xyz"),
			wantCount: 1, // Should return the original pattern
			wantErr:   false,
		},
		{
			name:      "expand pattern matching all files",
			pattern:   filepath.Join(tmpDir, "*"),
			wantCount: 3,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, err := expandPath(tt.pattern)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expandPath() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expandPath() unexpected error = %v", err)
				}
				if len(matches) != tt.wantCount {
					t.Errorf("expandPath() got %d matches, want %d", len(matches), tt.wantCount)
				}
			}
		})
	}
}
