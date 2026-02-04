package render

import (
	"strings"
	"testing"
)

func TestShouldRenderAsFile(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected bool
	}{
		{"short text", "Hello world", false},
		{"long text", strings.Repeat("a", 1001), true},
		{"exactly threshold", strings.Repeat("a", 1000), true},
		{"short with code block", "Hello ```code``` world", true},
		{"short without code", "Hello `inline` world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldRenderAsFile(tt.markdown)
			if result != tt.expected {
				t.Errorf("ShouldRenderAsFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToTelegramHTML(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		contains []string
	}{
		{
			name:     "bold",
			markdown: "**bold text**",
			contains: []string{"<b>bold text</b>"},
		},
		{
			name:     "italic",
			markdown: "*italic text*",
			contains: []string{"<i>italic text</i>"},
		},
		{
			name:     "inline code",
			markdown: "`code`",
			contains: []string{"<code>code</code>"},
		},
		{
			name:     "code block",
			markdown: "```\ncode block\n```",
			contains: []string{"<pre>"},
		},
		{
			name:     "link",
			markdown: "[link](https://example.com)",
			contains: []string{`<a href="https://example.com">link</a>`},
		},
		{
			name:     "escape html",
			markdown: "<script>alert('xss')</script>",
			contains: []string{"&lt;script&gt;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToTelegramHTML(tt.markdown)
			for _, c := range tt.contains {
				if !strings.Contains(result, c) {
					t.Errorf("ToTelegramHTML() = %q, should contain %q", result, c)
				}
			}
		})
	}
}

func TestToHTMLFile(t *testing.T) {
	markdown := "# Title\n\nSome **bold** text."
	title := "Test Report"

	html, err := ToHTMLFile(markdown, title)
	if err != nil {
		t.Fatalf("ToHTMLFile failed: %v", err)
	}

	// Check essential parts
	checks := []string{
		"<!DOCTYPE html>",
		"<title>Test Report</title>",
		"<h1",
		"<strong>bold</strong>",
		"@media (prefers-color-scheme: dark)",
	}

	for _, c := range checks {
		if !strings.Contains(html, c) {
			t.Errorf("HTML should contain %q", c)
		}
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		markdown string
		expected string
	}{
		{"# Main Title\n\nContent", "Main Title"},
		{"## Sub Title\n\nContent", "Sub Title"},
		{"No heading\n\nJust text", "No heading"},
		{"\n\n\n", "Report"},
		{strings.Repeat("a", 100), strings.Repeat("a", 80) + "..."},
		// Extract content after "## 요약"
		{"## 요약\n\n작업 완료됨", "작업 완료됨"},
		{"## Summary\n\n- Task done", "Task done"},
		{"## 요약\n로그인 기능을 구현했습니다.\n\n## 상세\n- 내용", "로그인 기능을 구현했습니다."},
		{"## 요약\n\n## 상세\n결과 내용", "상세"},
		// Non-summary heading
		{"## 실제 제목\n\nContent", "실제 제목"},
	}

	for _, tt := range tests {
		result := ExtractTitle(tt.markdown)
		if result != tt.expected {
			t.Errorf("ExtractTitle(%q) = %q, want %q", tt.markdown[:min(20, len(tt.markdown))], result, tt.expected)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
