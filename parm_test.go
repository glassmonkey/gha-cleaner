package parm

import (
	"os"
	"testing"
)

func TestRemovePathsInParallel(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() ([]string, error)
		wantErr    bool
		setForce   bool
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

			cfg := Config{
				Recursive: true,
				Force:     tt.setForce,
				Verbose:   tt.setVerbose,
			}

			err = RemovePathsInParallel(paths, cfg)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RemovePathsInParallel() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("RemovePathsInParallel() unexpected error = %v", err)
				}
			}
		})
	}
}
