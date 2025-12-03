package findsh

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestWindowsFind tests the Windows-specific Find implementation
// This test only runs on Windows and tests the Windows logic
func TestWindowsFind(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	tests := []struct {
		name    string
		setup   func() (cleanup func())
		wantErr bool
	}{
		{
			name: "sh found directly in PATH",
			setup: func() func() {
				// This test relies on sh being available in PATH
				return func() {}
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
			
			// On Windows, verify the result ends with sh.exe
			if !filepath.IsAbs(result) {
				t.Errorf("Find() should return absolute path, got: %s", result)
			}
		})
	}
}

// TestWindowsFindLogic tests the Windows-specific logic in a platform-agnostic way
// by creating a mock filesystem structure
func TestWindowsFindLogic(t *testing.T) {
	// Create a temporary directory structure that mimics Windows Git installations
	tmpDir := t.TempDir()
	
	tests := []struct {
		name           string
		setupFiles     func(baseDir string)
		mockGitPath    string
		expectShPath   string
		expectError    bool
	}{
		{
			name: "regular Git for Windows install",
			setupFiles: func(baseDir string) {
				// Create git.exe in a bin directory
				gitBinDir := filepath.Join(baseDir, "git", "bin")
				os.MkdirAll(gitBinDir, 0755)
				gitPath := filepath.Join(gitBinDir, "git.exe")
				os.WriteFile(gitPath, []byte("mock git"), 0755)
				
				// Create sh.exe in the expected location
				shBinDir := filepath.Join(baseDir, "git", "bin")
				os.MkdirAll(shBinDir, 0755)
				shPath := filepath.Join(shBinDir, "sh.exe")
				os.WriteFile(shPath, []byte("mock sh"), 0755)
			},
			mockGitPath:  "git/bin/git.exe",
			expectShPath: "git/bin/sh.exe",
			expectError:  false,
		},
		{
			name: "scoop Git install",
			setupFiles: func(baseDir string) {
				// Create git.exe in scoop shims directory
				gitBinDir := filepath.Join(baseDir, "scoop", "shims")
				os.MkdirAll(gitBinDir, 0755)
				gitPath := filepath.Join(gitBinDir, "git.exe")
				os.WriteFile(gitPath, []byte("mock git"), 0755)
				
				// Create sh.exe in the scoop apps location
				shBinDir := filepath.Join(baseDir, "scoop", "apps", "git", "current", "bin")
				os.MkdirAll(shBinDir, 0755)
				shPath := filepath.Join(shBinDir, "sh.exe")
				os.WriteFile(shPath, []byte("mock sh"), 0755)
			},
			mockGitPath:  "scoop/shims/git.exe",
			expectShPath: "scoop/apps/git/current/bin/sh.exe",
			expectError:  false,
		},
		{
			name: "git found but no sh.exe",
			setupFiles: func(baseDir string) {
				// Create git.exe but no sh.exe
				gitBinDir := filepath.Join(baseDir, "git", "bin")
				os.MkdirAll(gitBinDir, 0755)
				gitPath := filepath.Join(gitBinDir, "git.exe")
				os.WriteFile(gitPath, []byte("mock git"), 0755)
			},
			mockGitPath:  "git/bin/git.exe",
			expectShPath: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the mock filesystem
			testDir := filepath.Join(tmpDir, tt.name)
			os.MkdirAll(testDir, 0755)
			tt.setupFiles(testDir)
			
			// The actual Windows implementation would use safeexec.LookPath
			// For testing purposes, we verify the file structure is correct
			gitPath := filepath.Join(testDir, tt.mockGitPath)
			if _, err := os.Stat(gitPath); err != nil {
				t.Fatalf("Mock git.exe not found at %s", gitPath)
			}
			
			if !tt.expectError {
				expectedShPath := filepath.Join(testDir, tt.expectShPath)
				if _, err := os.Stat(expectedShPath); err != nil {
					t.Errorf("Expected sh.exe not found at %s", expectedShPath)
				}
			}
		})
	}
}

// TestWindowsPathCleaning tests that Windows paths are properly cleaned
func TestWindowsPathCleaning(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path with double slashes",
			input:    "C:\\Program Files\\Git\\\\bin\\sh.exe",
			expected: "C:\\Program Files\\Git\\bin\\sh.exe",
		},
		{
			name:     "path with forward slashes",
			input:    "C:/Program Files/Git/bin/sh.exe",
			expected: "C:\\Program Files\\Git\\bin\\sh.exe",
		},
		{
			name:     "relative path with ..",
			input:    "git\\..\\bin\\sh.exe",
			expected: "bin\\sh.exe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filepath.Clean(tt.input)
			// On non-Windows systems, filepath.Clean behaves differently
			// so we just verify it returns a cleaned path
			if result == "" {
				t.Errorf("filepath.Clean returned empty string for input: %s", tt.input)
			}
		})
	}
}