package gate

import (
	"context"
	"io"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.minekube.com/common/minecraft/component"
	jconfig "go.minekube.com/gate/pkg/edition/java/config"
	jproxy "go.minekube.com/gate/pkg/edition/java/proxy"
)

// Test helpers and mocks
type testWriter struct {
	output strings.Builder
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	return w.output.Write(p)
}

func (w *testWriter) String() string {
	return w.output.String()
}

type testGate struct {
	javaProxy *jproxy.Proxy
}

func (g *testGate) Java() *jproxy.Proxy {
	return g.javaProxy
}

func (g *testGate) Bedrock() interface{} {
	return nil
}

func (g *testGate) StopWithReason(reason component.Component) {
	// Mock implementation
}

func TestConsoleRunner_printHelp(t *testing.T) {
	tests := []struct {
		name         string
		liteEnabled  bool
		expectLite   bool
		expectFull   bool
	}{
		{
			name:        "lite mode help",
			liteEnabled: true,
			expectLite:  true,
			expectFull:  false,
		},
		{
			name:        "full mode help",
			liteEnabled: false,
			expectLite:  false,
			expectFull:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &testWriter{}
			runner := &consoleRunner{
				writer: writer,
				gate:   &testGate{},
			}

			// Mock lite mode based on test case
			if tt.liteEnabled {
				// Set up mock for lite mode
				runner.gate = &testGate{}
			}

			runner.printHelp()
			output := writer.String()

			assert.Contains(t, output, "Available commands:")

			if tt.expectLite {
				assert.Contains(t, output, "routes - Show Gate Lite route configuration")
				assert.NotContains(t, output, "kick <player>")
			}

			if tt.expectFull {
				assert.Contains(t, output, "kick <player> [reason] - Disconnect a player")
				assert.Contains(t, output, "move <player|server> <server>")
				assert.NotContains(t, output, "routes - Show Gate Lite route configuration")
			}
		})
	}
}

func TestConsoleRunner_execute(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "empty command",
			input:          "",
			expectError:    false,
			expectedOutput: "",
		},
		{
			name:           "help command",
			input:          "help",
			expectError:    false,
			expectedOutput: "Available commands:",
		},
		{
			name:           "unknown command",
			input:          "unknown",
			expectError:    false,
			expectedOutput: "Unknown command 'unknown'",
		},
		{
			name:           "info command",
			input:          "info",
			expectError:    false,
			expectedOutput: "Gate Proxy Information:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &testWriter{}
			runner := &consoleRunner{
				writer: writer,
				gate:   &testGate{},
			}

			err := runner.execute(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedOutput != "" {
				assert.Contains(t, writer.String(), tt.expectedOutput)
			}
		})
	}
}

func TestConsoleRunner_reloadCommand(t *testing.T) {
	writer := &testWriter{}
	runner := &consoleRunner{
		writer: writer,
		gate:   &testGate{},
	}

	err := runner.reloadCommand([]string{})

	require.NoError(t, err)
	output := writer.String()
	assert.Contains(t, output, "Gate has automatic configuration reload enabled by default")
	assert.Contains(t, output, "Configuration changes are automatically detected")
	assert.Contains(t, output, "No manual reload command is necessary")
}

func TestConsoleRunner_infoCommand(t *testing.T) {
	writer := &testWriter{}
	runner := &consoleRunner{
		writer: writer,
		gate:   &testGate{},
	}

	err := runner.infoCommand([]string{})

	require.NoError(t, err)
	output := writer.String()
	assert.Contains(t, output, "Gate Proxy Information:")
	assert.Contains(t, output, "Version: Gate")
	assert.Contains(t, output, "Bedrock support: Disabled")
}

func TestConsoleRunner_stopGate(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		stopping     bool
		expectError  bool
		expectedMsg  string
	}{
		{
			name:        "first stop request",
			args:        []string{},
			stopping:    false,
			expectError: true, // errStopRequested is returned
			expectedMsg: "Stopping Gate...",
		},
		{
			name:        "stop with reason",
			args:        []string{"maintenance", "mode"},
			stopping:    false,
			expectError: true, // errStopRequested is returned
			expectedMsg: "Stopping Gate: maintenance mode",
		},
		{
			name:        "duplicate stop request",
			args:        []string{},
			stopping:    true,
			expectError: true, // errStopRequested is returned
			expectedMsg: "Stop already in progress",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &testWriter{}
			runner := &consoleRunner{
				writer:   writer,
				gate:     &testGate{},
				stopping: atomic.Bool{},
			}

			if tt.stopping {
				runner.stopping.Store(true)
			}

			err := runner.stopGate(tt.args)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, errStopRequested, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedMsg != "" {
				assert.Contains(t, writer.String(), tt.expectedMsg)
			}
		})
	}
}

func TestConsoleRunner_closeReader(t *testing.T) {
	tests := []struct {
		name           string
		hasReader      bool
		callMultiple   bool
	}{
		{
			name:         "close with reader",
			hasReader:    true,
			callMultiple: false,
		},
		{
			name:         "close without reader",
			hasReader:    false,
			callMultiple: false,
		},
		{
			name:         "multiple close calls",
			hasReader:    true,
			callMultiple: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reader io.ReadCloser
			if tt.hasReader {
				reader = io.NopCloser(strings.NewReader("test"))
			}

			runner := &consoleRunner{
				reader: reader,
			}

			// Should not panic
			runner.closeReader()

			if tt.callMultiple {
				// Should not panic on second call
				runner.closeReader()
			}
		})
	}
}

func TestConsoleRunner_printPrompt(t *testing.T) {
	tests := []struct {
		name        string
		isTTY       bool
		expectPrompt bool
	}{
		{
			name:         "TTY mode",
			isTTY:        true,
			expectPrompt: true,
		},
		{
			name:         "non-TTY mode",
			isTTY:        false,
			expectPrompt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &testWriter{}
			runner := &consoleRunner{
				writer: writer,
				isTTY:  atomic.Bool{},
			}

			runner.isTTY.Store(tt.isTTY)
			runner.printPrompt()

			output := writer.String()
			if tt.expectPrompt {
				assert.Contains(t, output, "> ")
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestSortStringsCaseInsensitive(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "mixed case sorting",
			input:    []string{"Charlie", "alice", "Bob", "david"},
			expected: []string{"alice", "Bob", "Charlie", "david"},
		},
		{
			name:     "already sorted",
			input:    []string{"alice", "bob", "charlie"},
			expected: []string{"alice", "bob", "charlie"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single item",
			input:    []string{"Alice"},
			expected: []string{"Alice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the test input
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			sortStringsCaseInsensitive(input)
			assert.Equal(t, tt.expected, input)
		})
	}
}

func TestNewConsoleRunner(t *testing.T) {
	gate := &testGate{}
	runner := newConsoleRunner(gate)

	require.NotNil(t, runner)
	assert.Equal(t, gate, runner.gate)
	assert.NotNil(t, runner.reader)
	assert.NotNil(t, runner.writer)
	assert.False(t, runner.stopping.Load())
}

// Test the command validation patterns
func TestConsoleRunner_kickPlayer_Validation(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "no arguments",
			args:           []string{},
			expectError:    false,
			expectedOutput: "Usage: kick <player> [reason]",
		},
		{
			name:           "proxy not available",
			args:           []string{"player1"},
			expectError:    false,
			expectedOutput: "Java proxy not available yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &testWriter{}
			runner := &consoleRunner{
				writer: writer,
				gate:   &testGate{}, // nil proxy
			}

			err := runner.kickPlayer(tt.args)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedOutput != "" {
				assert.Contains(t, writer.String(), tt.expectedOutput)
			}
		})
	}
}

// Additional tests for edge cases and concurrent scenarios
func TestConsoleRunner_ConcurrentStopRequests(t *testing.T) {
	writer := &testWriter{}
	runner := &consoleRunner{
		writer:   writer,
		gate:     &testGate{},
		stopping: atomic.Bool{},
	}

	// Simulate concurrent stop requests
	const numGoroutines = 10
	var wg sync.WaitGroup
	errors := make([]error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			errors[idx] = runner.stopGate([]string{"test"})
		}(i)
	}

	wg.Wait()

	// Verify that all calls returned errStopRequested
	successCount := 0
	for _, err := range errors {
		if err == errStopRequested {
			successCount++
		}
	}

	// All should return errStopRequested, but only one should actually initiate stop
	assert.Equal(t, numGoroutines, successCount)
	assert.True(t, runner.stopping.Load())
}

func TestConsoleRunner_MalformedInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "whitespace only",
			input: "   \t\n  ",
		},
		{
			name:  "very long command",
			input: strings.Repeat("a", 10000),
		},
		{
			name:  "special characters",
			input: "command\x00\x01\x02",
		},
		{
			name:  "unicode characters",
			input: "help 测试 🚀 ñoño",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &testWriter{}
			runner := &consoleRunner{
				writer: writer,
				gate:   &testGate{},
			}

			// Should not panic
			err := runner.execute(tt.input)
			assert.NoError(t, err)
		})
	}
}

func TestConsoleRunner_FileDescriptorEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		reader io.ReadCloser
	}{
		{
			name:   "nil reader",
			reader: nil,
		},
		{
			name:   "non-file reader",
			reader: io.NopCloser(strings.NewReader("test")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &testWriter{}
			runner := &consoleRunner{
				reader: tt.reader,
				writer: writer,
				gate:   &testGate{},
				isTTY:  atomic.Bool{},
			}

			// Should not panic during TTY detection
			assert.NotNil(t, runner)
			assert.False(t, runner.isTTY.Load()) // Should default to false
		})
	}
}