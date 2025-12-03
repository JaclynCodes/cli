//go:build !windows

package findsh

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (cleanup func())
		wantErr bool
	}{
		{
			name: "sh found in PATH",
			setup: func() func() {
				// No setup needed for normal case
				return func() {}
			},
			wantErr: false,
		},
		{
			name: "sh not found in PATH",
			setup: func() func() {
				// Save original PATH
				originalPath := os.Getenv("PATH")
				// Set PATH to empty to simulate sh not being found
				os.Setenv("PATH", "")
				return func() {
					// Restore original PATH
					os.Setenv("PATH", originalPath)
				}
			},
			wantErr: true,
		},
		{
			name: "sh found in custom PATH",
			setup: func() func() {
				// Create a temporary directory with a mock sh executable
				tmpDir := t.TempDir()
				shPath := filepath.Join(tmpDir, "sh")
				
				// Create a mock sh file
				if err := os.WriteFile(shPath, []byte("#!/bin/bash\necho 'mock sh'\n"), 0755); err != nil {
					t.Fatalf("Failed to create mock sh: %v", err)
				}
				
				// Save original PATH and set new one
				originalPath := os.Getenv("PATH")
				os.Setenv("PATH", tmpDir)
				
				return func() {
					os.Setenv("PATH", originalPath)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup()
			defer cleanup()

			result, err := Find()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Find() expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Find() unexpected error: %v", err)
				return
			}
			
			if result == "" {
				t.Errorf("Find() returned empty path")
				return
			}
			
			// Verify the result is a valid path to sh
			if !strings.HasSuffix(result, "sh") {
				t.Errorf("Find() returned path that doesn't end with 'sh': %s", result)
			}
			
			// Verify the file exists and is executable
			info, err := os.Stat(result)
			if err != nil {
				t.Errorf("Find() returned non-existent path: %s", result)
				return
			}
			
			if info.Mode()&0111 == 0 {
				t.Errorf("Find() returned non-executable file: %s", result)
			}
		})
	}
}

func TestFindConsistency(t *testing.T) {
	// Test that Find() returns consistent results when called multiple times
	result1, err1 := Find()
	result2, err2 := Find()
	
	if err1 != err2 {
		t.Errorf("Find() returned different errors: %v vs %v", err1, err2)
	}
	
	if result1 != result2 {
		t.Errorf("Find() returned different paths: %s vs %s", result1, result2)
	}
}

func TestFindMatchesExecLookPath(t *testing.T) {
	// Test that our Find() function returns the same result as exec.LookPath("sh")
	ourResult, ourErr := Find()
	execResult, execErr := exec.LookPath("sh")
	
	if (ourErr == nil) != (execErr == nil) {
		t.Errorf("Find() and exec.LookPath() returned different error states: %v vs %v", ourErr, execErr)
		return
	}
	
	if ourErr == nil && ourResult != execResult {
		t.Errorf("Find() returned %s, exec.LookPath() returned %s", ourResult, execResult)
	}
}

// Benchmark to ensure Find() is reasonably fast
func BenchmarkFind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Find()
		if err != nil {
			b.Fatalf("Find() failed: %v", err)
		}
	}
}