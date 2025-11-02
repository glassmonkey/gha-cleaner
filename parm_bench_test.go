package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// setupBenchFiles creates a temporary directory with the specified number of files
func setupBenchFiles(b *testing.B, numFiles int) string {
	b.Helper()
	tmpDir, err := os.MkdirTemp("", "parm-bench-files-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	for i := 0; i < numFiles; i++ {
		file := filepath.Join(tmpDir, fmt.Sprintf("file-%d.txt", i))
		if err := os.WriteFile(file, []byte("benchmark test data"), 0644); err != nil {
			os.RemoveAll(tmpDir)
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	return tmpDir
}

// setupBenchDirs creates a temporary directory with the specified number of subdirectories
func setupBenchDirs(b *testing.B, numDirs int) string {
	b.Helper()
	tmpDir, err := os.MkdirTemp("", "parm-bench-dirs-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	for i := 0; i < numDirs; i++ {
		dir := filepath.Join(tmpDir, fmt.Sprintf("dir-%d", i))
		if err := os.MkdirAll(dir, 0755); err != nil {
			os.RemoveAll(tmpDir)
			b.Fatalf("Failed to create test dir: %v", err)
		}
		// Add a file in each directory
		file := filepath.Join(dir, "file.txt")
		if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
			os.RemoveAll(tmpDir)
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	return tmpDir
}

// setupBenchNestedDirs creates nested directory structure
func setupBenchNestedDirs(b *testing.B, depth, filesPerDir int) string {
	b.Helper()
	tmpDir, err := os.MkdirTemp("", "parm-bench-nested-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	var createNested func(string, int) error
	createNested = func(parent string, level int) error {
		if level <= 0 {
			return nil
		}

		// Create files in current directory
		for i := 0; i < filesPerDir; i++ {
			file := filepath.Join(parent, fmt.Sprintf("file-%d.txt", i))
			if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
				return err
			}
		}

		// Create subdirectories
		for i := 0; i < 3; i++ {
			subDir := filepath.Join(parent, fmt.Sprintf("subdir-%d", i))
			if err := os.MkdirAll(subDir, 0755); err != nil {
				return err
			}
			if err := createNested(subDir, level-1); err != nil {
				return err
			}
		}
		return nil
	}

	if err := createNested(tmpDir, depth); err != nil {
		os.RemoveAll(tmpDir)
		b.Fatalf("Failed to create nested structure: %v", err)
	}

	return tmpDir
}

// Benchmark parm with different numbers of files
func BenchmarkParmRemoveFiles100(b *testing.B) {
	benchmarkParmRemoveFiles(b, 100)
}

func BenchmarkParmRemoveFiles1000(b *testing.B) {
	benchmarkParmRemoveFiles(b, 1000)
}

func BenchmarkParmRemoveFiles10000(b *testing.B) {
	benchmarkParmRemoveFiles(b, 10000)
}

func benchmarkParmRemoveFiles(b *testing.B, numFiles int) {
	// Set flags for the benchmark
	oldRecursive := recursive
	oldForce := force
	oldVerbose := verbose
	defer func() {
		recursive = oldRecursive
		force = oldForce
		verbose = oldVerbose
	}()

	recursive = true
	force = true
	verbose = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tmpDir := setupBenchFiles(b, numFiles)
		b.StartTimer()

		// Get all files to remove
		files, err := filepath.Glob(filepath.Join(tmpDir, "*"))
		if err != nil {
			b.Fatalf("Failed to glob files: %v", err)
		}

		if err := removePathsInParallel(files); err != nil {
			b.Fatalf("Failed to remove files: %v", err)
		}

		b.StopTimer()
		os.RemoveAll(tmpDir)
	}
}

// Benchmark OS rm command with different numbers of files
func BenchmarkOSRmRemoveFiles100(b *testing.B) {
	benchmarkOSRmRemoveFiles(b, 100)
}

func BenchmarkOSRmRemoveFiles1000(b *testing.B) {
	benchmarkOSRmRemoveFiles(b, 1000)
}

func BenchmarkOSRmRemoveFiles10000(b *testing.B) {
	benchmarkOSRmRemoveFiles(b, 10000)
}

func benchmarkOSRmRemoveFiles(b *testing.B, numFiles int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tmpDir := setupBenchFiles(b, numFiles)
		b.StartTimer()

		// Use OS rm command
		cmd := exec.Command("rm", "-rf", tmpDir)
		if err := cmd.Run(); err != nil {
			b.Fatalf("Failed to run rm: %v", err)
		}
	}
}

// Benchmark parm with directories
func BenchmarkParmRemoveDirs100(b *testing.B) {
	benchmarkParmRemoveDirs(b, 100)
}

func BenchmarkParmRemoveDirs1000(b *testing.B) {
	benchmarkParmRemoveDirs(b, 1000)
}

func benchmarkParmRemoveDirs(b *testing.B, numDirs int) {
	oldRecursive := recursive
	oldForce := force
	oldVerbose := verbose
	defer func() {
		recursive = oldRecursive
		force = oldForce
		verbose = oldVerbose
	}()

	recursive = true
	force = true
	verbose = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tmpDir := setupBenchDirs(b, numDirs)
		b.StartTimer()

		// Get all directories to remove
		dirs, err := filepath.Glob(filepath.Join(tmpDir, "*"))
		if err != nil {
			b.Fatalf("Failed to glob dirs: %v", err)
		}

		if err := removePathsInParallel(dirs); err != nil {
			b.Fatalf("Failed to remove dirs: %v", err)
		}

		b.StopTimer()
		os.RemoveAll(tmpDir)
	}
}

// Benchmark OS rm with directories
func BenchmarkOSRmRemoveDirs100(b *testing.B) {
	benchmarkOSRmRemoveDirs(b, 100)
}

func BenchmarkOSRmRemoveDirs1000(b *testing.B) {
	benchmarkOSRmRemoveDirs(b, 1000)
}

func benchmarkOSRmRemoveDirs(b *testing.B, numDirs int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tmpDir := setupBenchDirs(b, numDirs)
		b.StartTimer()

		// Use OS rm command
		cmd := exec.Command("rm", "-rf", tmpDir)
		if err := cmd.Run(); err != nil {
			b.Fatalf("Failed to run rm: %v", err)
		}
	}
}

// Benchmark nested directory structures
func BenchmarkParmRemoveNested(b *testing.B) {
	oldRecursive := recursive
	oldForce := force
	oldVerbose := verbose
	defer func() {
		recursive = oldRecursive
		force = oldForce
		verbose = oldVerbose
	}()

	recursive = true
	force = true
	verbose = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tmpDir := setupBenchNestedDirs(b, 4, 10) // depth=4, 10 files per dir
		b.StartTimer()

		if err := removePathsInParallel([]string{tmpDir}); err != nil {
			b.Fatalf("Failed to remove nested dirs: %v", err)
		}
	}
}

func BenchmarkOSRmRemoveNested(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tmpDir := setupBenchNestedDirs(b, 4, 10) // depth=4, 10 files per dir
		b.StartTimer()

		// Use OS rm command
		cmd := exec.Command("rm", "-rf", tmpDir)
		if err := cmd.Run(); err != nil {
			b.Fatalf("Failed to run rm: %v", err)
		}
	}
}

// Benchmark comparison: parallel vs sequential removal
func BenchmarkParmParallelVsSequential(b *testing.B) {
	numFiles := 1000

	b.Run("Parallel", func(b *testing.B) {
		oldRecursive := recursive
		oldForce := force
		oldVerbose := verbose
		defer func() {
			recursive = oldRecursive
			force = oldForce
			verbose = oldVerbose
		}()

		recursive = true
		force = true
		verbose = false

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			tmpDir := setupBenchFiles(b, numFiles)
			files, _ := filepath.Glob(filepath.Join(tmpDir, "*"))
			b.StartTimer()

			removePathsInParallel(files)

			b.StopTimer()
			os.RemoveAll(tmpDir)
		}
	})

	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			tmpDir := setupBenchFiles(b, numFiles)
			files, _ := filepath.Glob(filepath.Join(tmpDir, "*"))
			b.StartTimer()

			// Sequential removal
			for _, file := range files {
				os.Remove(file)
			}

			b.StopTimer()
			os.RemoveAll(tmpDir)
		}
	})
}

// setupGHALikeSDK creates a directory structure similar to GHA SDK directories
// Mimics structure of /usr/lib/jvm, /usr/share/dotnet, /usr/local/lib/android, etc.
func setupGHALikeSDK(b *testing.B, sdkName string) string {
	b.Helper()
	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("parm-bench-gha-%s-*", sdkName))
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create SDK-like structure with:
	// - bin/ directory with executable files
	// - lib/ directory with library files
	// - include/ directory with header files
	// - share/ directory with shared resources
	// - Multiple version subdirectories

	versions := []string{"11.0.1", "17.0.2", "21.0.1"}
	for _, ver := range versions {
		verDir := filepath.Join(tmpDir, ver)

		// Create bin directory with "executables"
		binDir := filepath.Join(verDir, "bin")
		if err := os.MkdirAll(binDir, 0755); err != nil {
			os.RemoveAll(tmpDir)
			b.Fatalf("Failed to create bin dir: %v", err)
		}
		for i := 0; i < 50; i++ {
			file := filepath.Join(binDir, fmt.Sprintf("tool-%d", i))
			if err := os.WriteFile(file, []byte("executable content"), 0755); err != nil {
				os.RemoveAll(tmpDir)
				b.Fatalf("Failed to create file: %v", err)
			}
		}

		// Create lib directory with many "library" files
		libDir := filepath.Join(verDir, "lib")
		if err := os.MkdirAll(libDir, 0755); err != nil {
			os.RemoveAll(tmpDir)
			b.Fatalf("Failed to create lib dir: %v", err)
		}
		for i := 0; i < 200; i++ {
			file := filepath.Join(libDir, fmt.Sprintf("lib-%d.so", i))
			if err := os.WriteFile(file, make([]byte, 1024*10), 0644); err != nil {
				os.RemoveAll(tmpDir)
				b.Fatalf("Failed to create lib file: %v", err)
			}
		}

		// Create include directory with header files
		includeDir := filepath.Join(verDir, "include")
		if err := os.MkdirAll(includeDir, 0755); err != nil {
			os.RemoveAll(tmpDir)
			b.Fatalf("Failed to create include dir: %v", err)
		}
		for i := 0; i < 100; i++ {
			file := filepath.Join(includeDir, fmt.Sprintf("header-%d.h", i))
			if err := os.WriteFile(file, []byte("header content"), 0644); err != nil {
				os.RemoveAll(tmpDir)
				b.Fatalf("Failed to create header file: %v", err)
			}
		}

		// Create share directory with nested structure
		shareDir := filepath.Join(verDir, "share", "resources")
		if err := os.MkdirAll(shareDir, 0755); err != nil {
			os.RemoveAll(tmpDir)
			b.Fatalf("Failed to create share dir: %v", err)
		}
		for i := 0; i < 150; i++ {
			file := filepath.Join(shareDir, fmt.Sprintf("resource-%d.xml", i))
			if err := os.WriteFile(file, []byte("resource data"), 0644); err != nil {
				os.RemoveAll(tmpDir)
				b.Fatalf("Failed to create resource file: %v", err)
			}
		}
	}

	return tmpDir
}

// setupGHALikeEnvironment creates multiple SDK directories like a real GHA runner
func setupGHALikeEnvironment(b *testing.B) []string {
	b.Helper()

	sdks := []string{"jvm", "dotnet", "android", "swift", "ghcup"}
	var paths []string

	for _, sdk := range sdks {
		path := setupGHALikeSDK(b, sdk)
		paths = append(paths, path)
	}

	return paths
}

// BenchmarkGHAScenario simulates real GHA cleanup scenario
func BenchmarkGHAScenarioParm(b *testing.B) {
	oldRecursive := recursive
	oldForce := force
	oldVerbose := verbose
	defer func() {
		recursive = oldRecursive
		force = oldForce
		verbose = oldVerbose
	}()

	recursive = true
	force = true
	verbose = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		sdkPaths := setupGHALikeEnvironment(b)
		b.StartTimer()

		// Remove all SDK directories in parallel (real scenario)
		if err := removePathsInParallel(sdkPaths); err != nil {
			b.Fatalf("Failed to remove SDK paths: %v", err)
		}

		b.StopTimer()
		// Cleanup any remaining files
		for _, path := range sdkPaths {
			os.RemoveAll(path)
		}
	}
}

func BenchmarkGHAScenarioOSRm(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		sdkPaths := setupGHALikeEnvironment(b)
		b.StartTimer()

		// Remove all SDK directories sequentially (like original rm -rf in shell)
		for _, path := range sdkPaths {
			cmd := exec.Command("rm", "-rf", path)
			if err := cmd.Run(); err != nil {
				b.Fatalf("Failed to run rm: %v", err)
			}
		}
	}
}

// BenchmarkSingleLargeSDK tests a single large SDK directory
func BenchmarkSingleLargeSDKParm(b *testing.B) {
	oldRecursive := recursive
	oldForce := force
	oldVerbose := verbose
	defer func() {
		recursive = oldRecursive
		force = oldForce
		verbose = oldVerbose
	}()

	recursive = true
	force = true
	verbose = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		sdkPath := setupGHALikeSDK(b, "android")
		b.StartTimer()

		if err := removePathsInParallel([]string{sdkPath}); err != nil {
			b.Fatalf("Failed to remove SDK: %v", err)
		}

		b.StopTimer()
		os.RemoveAll(sdkPath)
	}
}

func BenchmarkSingleLargeSDKOSRm(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		sdkPath := setupGHALikeSDK(b, "android")
		b.StartTimer()

		cmd := exec.Command("rm", "-rf", sdkPath)
		if err := cmd.Run(); err != nil {
			b.Fatalf("Failed to run rm: %v", err)
		}
	}
}
