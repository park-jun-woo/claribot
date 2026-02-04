package claude

import (
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	// Reset manager for test
	Init(Config{
		Timeout: 30 * time.Second,
		Max:     3,
	})

	result, err := Run(Options{
		SystemPrompt: "You are a helpful assistant. Keep responses very brief.",
		UserPrompt:   "Say 'Hello Claribot' and nothing else.",
	})

	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	t.Logf("Exit code: %d", result.ExitCode)
	t.Logf("Output:\n%s", result.Output)

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
}

func TestConcurrencyLimit(t *testing.T) {
	// Initialize with max 2
	globalManager = nil
	managerOnce = sync.Once{}
	Init(Config{
		Timeout: 30 * time.Second,
		Max:     2,
	})

	mgr := GetManager()

	// Check initial state
	if mgr.Available() != 2 {
		t.Errorf("Expected 2 available, got %d", mgr.Available())
	}

	t.Logf("Max: 2, Available: %d", mgr.Available())
}
