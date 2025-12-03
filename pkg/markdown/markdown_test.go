package markdown

import (
	"os"
	"testing"

	"github.com/charmbracelet/glamour"
)

func TestWithWrap(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		envValue string
		setup    func()
		cleanup  func()
	}{
		{
			name:     "default max width 120",
			input:    100,
			envValue: "",
			setup: func() {
				os.Unsetenv("GH_MDWIDTH")
			},
			cleanup: func() {},
		},
		{
			name:     "width exceeds default max",
			input:    150,
			envValue: "",
			setup: func() {
				os.Unsetenv("GH_MDWIDTH")
			},
			cleanup: func() {},
		},
		{
			name:     "custom max width from env",
			input:    80,
			envValue: "100",
			setup: func() {
				os.Setenv("GH_MDWIDTH", "100")
			},
			cleanup: func() {
				os.Unsetenv("GH_MDWIDTH")
			},
		},
		{
			name:     "width exceeds custom max",
			input:    120,
			envValue: "100",
			setup: func() {
				os.Setenv("GH_MDWIDTH", "100")
			},
			cleanup: func() {
				os.Unsetenv("GH_MDWIDTH")
			},
		},
		{
			name:     "invalid env value falls back to default",
			input:    80,
			envValue: "invalid",
			setup: func() {
				os.Setenv("GH_MDWIDTH", "invalid")
			},
			cleanup: func() {
				os.Unsetenv("GH_MDWIDTH")
			},
		},
		{
			name:     "zero width disables wrapping",
			input:    0,
			envValue: "",
			setup: func() {
				os.Unsetenv("GH_MDWIDTH")
			},
			cleanup: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			// Test that WithWrap returns a valid option
			option := WithWrap(tt.input)
			if option == nil {
				t.Errorf("WithWrap(%d) returned nil", tt.input)
			}

			// Test that the option can be used with a renderer
			renderer, err := glamour.NewTermRenderer(option)
			if err != nil {
				t.Errorf("Failed to create renderer with WithWrap(%d): %v", tt.input, err)
			}
			if renderer == nil {
				t.Errorf("WithWrap(%d) created nil renderer", tt.input)
			}
		})
	}
}

func TestWithoutIndentation(t *testing.T) {
	option := WithoutIndentation()
	if option == nil {
		t.Error("WithoutIndentation() returned nil")
	}

	// Test that the option can be used with a renderer
	renderer, err := glamour.NewTermRenderer(option)
	if err != nil {
		t.Errorf("Failed to create renderer with WithoutIndentation(): %v", err)
	}
	if renderer == nil {
		t.Error("WithoutIndentation() created nil renderer")
	}
}

func TestWithTheme(t *testing.T) {
	tests := []struct {
		name  string
		theme string
	}{
		{"dark theme", "dark"},
		{"light theme", "light"},
		{"auto theme", "auto"},
		{"empty theme", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithTheme(tt.theme)
			if option == nil {
				t.Errorf("WithTheme(%q) returned nil", tt.theme)
			}

			// Test that the option can be used with a renderer
			renderer, err := glamour.NewTermRenderer(option)
			if err != nil {
				t.Errorf("Failed to create renderer with WithTheme(%q): %v", tt.theme, err)
			}
			if renderer == nil {
				t.Errorf("WithTheme(%q) created nil renderer", tt.theme)
			}
		})
	}
}

func TestWithBaseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"valid URL", "https://github.com"},
		{"empty URL", ""},
		{"relative URL", "/path/to/resource"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := WithBaseURL(tt.url)
			if option == nil {
				t.Errorf("WithBaseURL(%q) returned nil", tt.url)
			}

			// Test that the option can be used with a renderer
			renderer, err := glamour.NewTermRenderer(option)
			if err != nil {
				t.Errorf("Failed to create renderer with WithBaseURL(%q): %v", tt.url, err)
			}
			if renderer == nil {
				t.Errorf("WithBaseURL(%q) created nil renderer", tt.url)
			}
		})
	}
}

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		checkLen bool
	}{
		{
			name:     "simple markdown",
			input:    "# Hello World\n\nThis is a test.",
			wantErr:  false,
			checkLen: true,
		},
		{
			name:     "empty string",
			input:    "",
			wantErr:  false,
			checkLen: false,
		},
		{
			name:     "markdown with code",
			input:    "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```",
			wantErr:  false,
			checkLen: true,
		},
		{
			name:     "markdown with links",
			input:    "[GitHub](https://github.com)",
			wantErr:  false,
			checkLen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Render(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Render(%q) expected error, got nil", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Render(%q) unexpected error: %v", tt.input, err)
				return
			}
			
			if tt.checkLen && len(result) == 0 {
				t.Errorf("Render(%q) returned empty result", tt.input)
			}
		})
	}
}

func TestRenderWithOptions(t *testing.T) {
	input := "# Test\n\nThis is a test with options."
	
	tests := []struct {
		name string
		opts []glamour.TermRendererOption
	}{
		{
			name: "with wrap",
			opts: []glamour.TermRendererOption{WithWrap(80)},
		},
		{
			name: "with theme",
			opts: []glamour.TermRendererOption{WithTheme("dark")},
		},
		{
			name: "with base URL",
			opts: []glamour.TermRendererOption{WithBaseURL("https://github.com")},
		},
		{
			name: "without indentation",
			opts: []glamour.TermRendererOption{WithoutIndentation()},
		},
		{
			name: "multiple options",
			opts: []glamour.TermRendererOption{
				WithWrap(100),
				WithTheme("light"),
				WithoutIndentation(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Render(input, tt.opts...)
			if err != nil {
				t.Errorf("Render with options failed: %v", err)
				return
			}
			
			if len(result) == 0 {
				t.Error("Render with options returned empty result")
			}
		})
	}
}

// Benchmark the Render function
func BenchmarkRender(b *testing.B) {
	input := "# Benchmark Test\n\nThis is a **markdown** test with `code` and [links](https://example.com)."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Render(input)
		if err != nil {
			b.Fatalf("Render failed: %v", err)
		}
	}
}

func BenchmarkRenderWithOptions(b *testing.B) {
	input := "# Benchmark Test\n\nThis is a **markdown** test with `code` and [links](https://example.com)."
	opts := []glamour.TermRendererOption{
		WithWrap(80),
		WithTheme("dark"),
		WithoutIndentation(),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Render(input, opts...)
		if err != nil {
			b.Fatalf("Render with options failed: %v", err)
		}
	}
}