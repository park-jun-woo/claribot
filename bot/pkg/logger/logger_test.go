package logger

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func resetGlobal() {
	global = nil
	globalOnce = sync.Once{}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", DEBUG},
		{"info", INFO},
		{"warn", WARN},
		{"error", ERROR},
		{"", INFO},
		{"unknown", INFO},
	}

	for _, tt := range tests {
		result := parseLevel(tt.input)
		if result != tt.expected {
			t.Errorf("parseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestLoggerInit(t *testing.T) {
	resetGlobal()

	err := Init(Config{Level: "debug"})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	logger := GetLogger()
	if logger == nil {
		t.Fatal("GetLogger returned nil")
	}

	if logger.level != DEBUG {
		t.Errorf("Expected level DEBUG, got %v", logger.level)
	}
}

func TestLoggerFileOutput(t *testing.T) {
	resetGlobal()

	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	err := Init(Config{
		Level:    "info",
		FilePath: logFile,
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Close()

	Info("test message %d", 42)

	// Read log file
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Read log file failed: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "[INFO]") {
		t.Errorf("Log file missing [INFO]: %s", content)
	}
	if !strings.Contains(content, "test message 42") {
		t.Errorf("Log file missing message: %s", content)
	}
}

func TestLogLevelFiltering(t *testing.T) {
	resetGlobal()

	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	err := Init(Config{
		Level:    "warn",
		FilePath: logFile,
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer Close()

	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	// Read log file
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Read log file failed: %v", err)
	}

	content := string(data)
	if strings.Contains(content, "[DEBUG]") {
		t.Error("DEBUG should be filtered out")
	}
	if strings.Contains(content, "[INFO]") {
		t.Error("INFO should be filtered out")
	}
	if !strings.Contains(content, "[WARN]") {
		t.Error("WARN should be present")
	}
	if !strings.Contains(content, "[ERROR]") {
		t.Error("ERROR should be present")
	}
}
