package proxy

import (
	"errors"
	"net"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.minekube.com/gate/pkg/edition/java/config"
	"go.minekube.com/gate/pkg/util/configutil"
)

func TestIsConnectionRefused(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "syscall.ECONNREFUSED should be detected",
			err:      syscall.ECONNREFUSED,
			expected: true,
		},
		{
			name:     "error with 'connection refused' in message should be detected",
			err:      errors.New("dial tcp 127.0.0.1:25566: connect: connection refused"),
			expected: true,
		},
		{
			name:     "error with 'Connection Refused' (different case) should be detected",
			err:      errors.New("Connection Refused by server"),
			expected: true,
		},
		{
			name:     "wrapped ECONNREFUSED should be detected",
			err:      &testError{Inner: syscall.ECONNREFUSED},
			expected: true,
		},
		{
			name:     "timeout error should not be detected as connection refused",
			err:      errors.New("dial tcp 127.0.0.1:25566: i/o timeout"),
			expected: false,
		},
		{
			name:     "other network error should not be detected",
			err:      errors.New("dial tcp: missing address"),
			expected: false,
		},
		{
			name:     "nil error should not be detected",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsConnectionRefused(tt.err)
			assert.Equal(t, tt.expected, result, "IsConnectionRefused should correctly identify connection refused errors")
		})
	}
}

func TestGetErrorVerbosity(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		expectedVerbosity int
	}{
		{
			name:             "connection refused gets debug verbosity",
			err:              syscall.ECONNREFUSED,
			expectedVerbosity: 1,
		},
		{
			name:             "connection refused message gets debug verbosity",
			err:              errors.New("connection refused"),
			expectedVerbosity: 1,
		},
		{
			name:             "timeout error gets info verbosity",
			err:              errors.New("i/o timeout"),
			expectedVerbosity: 0,
		},
		{
			name:             "other error gets info verbosity",
			err:              errors.New("some other error"),
			expectedVerbosity: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getErrorVerbosity(tt.err)
			assert.Equal(t, tt.expectedVerbosity, result, "getErrorVerbosity should return correct verbosity level")
		})
	}
}

func TestFindPassthroughServer(t *testing.T) {
	tests := []struct {
		name                  string
		setupConfig           func() *config.Config
		virtualHost           string // hostname for testing forced hosts
		expectedServerName    string
		expectedServerConfig  *config.ServerConfig
	}{
		{
			name: "finds server from forcedHosts when virtual host matches",
			setupConfig: func() *config.Config {
				return &config.Config{
					Servers: config.ServerConfigs{
						"vanilla":  {Address: "localhost:25566", PassthroughMOTD: true},
						"fabric":   {Address: "localhost:25567", PassthroughMOTD: true},
						"neoforge": {Address: "localhost:25568", PassthroughMOTD: true},
					},
					Try: []string{"fabric", "vanilla"},
					ForcedHosts: map[string][]string{
						"snapshot.example.test": {"vanilla"},
						"fabric.example.test":   {"fabric"},
						"forge.example.test":    {"neoforge"},
					},
				}
			},
			virtualHost:        "snapshot.example.test:25565",
			expectedServerName: "vanilla",
			expectedServerConfig: &config.ServerConfig{
				Address:         "localhost:25566",
				PassthroughMOTD: true,
			},
		},
		{
			name: "falls back to Try list when enabled and virtual host doesn't match forcedHosts",
			setupConfig: func() *config.Config {
				return &config.Config{
					Servers: config.ServerConfigs{
						"vanilla":  {Address: "localhost:25566", PassthroughMOTD: false},
						"fabric":   {Address: "localhost:25567", PassthroughMOTD: true},
						"neoforge": {Address: "localhost:25568", PassthroughMOTD: true},
					},
					Try:                []string{"fabric", "vanilla"},
					TryPassthroughMotd: true,
					ForcedHosts: map[string][]string{
						"snapshot.example.test": {"vanilla"},
					},
				}
			},
			virtualHost:        "unknown.host.com:25565",
			expectedServerName: "fabric",
			expectedServerConfig: &config.ServerConfig{
				Address:         "localhost:25567",
				PassthroughMOTD: true,
			},
		},
		{
			name: "finds first server in Try list with passthrough enabled when option is true",
			setupConfig: func() *config.Config {
				return &config.Config{
					Servers: config.ServerConfigs{
						"server1": {Address: "localhost:25561", PassthroughMOTD: false},
						"server2": {Address: "localhost:25562", PassthroughMOTD: true},
						"server3": {Address: "localhost:25563", PassthroughMOTD: true},
					},
					Try:                []string{"server2", "server3"},
					TryPassthroughMotd: true,
				}
			},
			virtualHost:        "", // No virtual host
			expectedServerName: "server2",
			expectedServerConfig: &config.ServerConfig{
				Address:         "localhost:25562",
				PassthroughMOTD: true,
			},
		},
		{
			name: "returns nil when Try list passthrough is disabled",
			setupConfig: func() *config.Config {
				return &config.Config{
					Servers: config.ServerConfigs{
						"server1": {Address: "localhost:25561", PassthroughMOTD: false},
						"server2": {Address: "localhost:25562", PassthroughMOTD: true},
					},
					Try:                []string{"server2"},
					TryPassthroughMotd: false,
				}
			},
			virtualHost:          "",
			expectedServerName:   "",
			expectedServerConfig: nil,
		},
		{
			name: "returns nil when no servers have passthrough enabled",
			setupConfig: func() *config.Config {
				return &config.Config{
					Servers: config.ServerConfigs{
						"server1": {Address: "localhost:25561", PassthroughMOTD: false},
						"server2": {Address: "localhost:25562", PassthroughMOTD: false},
					},
					Try: []string{"server1", "server2"},
				}
			},
			virtualHost:          "",
			expectedServerName:   "",
			expectedServerConfig: nil,
		},
		{
			name: "returns nil when only non-Try servers have passthrough",
			setupConfig: func() *config.Config {
				return &config.Config{
					Servers: config.ServerConfigs{
						"server1":  {Address: "localhost:25561", PassthroughMOTD: false},
						"server2":  {Address: "localhost:25562", PassthroughMOTD: false},
						"fallback": {Address: "localhost:25563", PassthroughMOTD: true},
					},
					Try:                []string{"server1", "server2"},
					TryPassthroughMotd: true,
				}
			},
			virtualHost:          "",
			expectedServerName:   "",
			expectedServerConfig: nil,
		},
		{
			name: "returns nil when no servers exist",
			setupConfig: func() *config.Config {
				return &config.Config{
					Servers: config.ServerConfigs{},
					Try:     []string{},
				}
			},
			virtualHost:          "",
			expectedServerName:   "",
			expectedServerConfig: nil,
		},
		{
			name: "skips forcedHost server when it doesn't have passthrough enabled",
			setupConfig: func() *config.Config {
				return &config.Config{
					Servers: config.ServerConfigs{
						"vanilla": {Address: "localhost:25566", PassthroughMOTD: false}, // No passthrough
						"fabric":  {Address: "localhost:25567", PassthroughMOTD: true},
					},
					Try:                []string{"fabric"},
					TryPassthroughMotd: true,
					ForcedHosts: map[string][]string{
						"snapshot.example.test": {"vanilla"}, // This server doesn't have passthrough
					},
				}
			},
			virtualHost:        "snapshot.example.test:25565",
			expectedServerName: "fabric", // Should fall back to Try list
			expectedServerConfig: &config.ServerConfig{
				Address:         "localhost:25567",
				PassthroughMOTD: true,
			},
		},
		{
			name: "returns nil when forcedHost lacks passthrough and Try passthrough disabled",
			setupConfig: func() *config.Config {
				return &config.Config{
					Servers: config.ServerConfigs{
						"vanilla": {Address: "localhost:25566", PassthroughMOTD: false},
						"fabric":  {Address: "localhost:25567", PassthroughMOTD: true},
					},
					Try:                []string{"fabric"},
					TryPassthroughMotd: false,
					ForcedHosts: map[string][]string{
						"snapshot.example.test": {"vanilla"},
					},
				}
			},
			virtualHost:          "snapshot.example.test:25565",
			expectedServerName:   "",
			expectedServerConfig: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test proxy with the configured settings
			cfg := tt.setupConfig()
			proxy := &Proxy{
				cfg: cfg,
			}

			// Create virtual host if provided
			var virtualHost net.Addr = nil
			if tt.virtualHost != "" {
				virtualHost = &testNetAddr{address: tt.virtualHost}
			}

			// Test the function
			serverConfig, serverName := proxy.findPassthroughServer(virtualHost)

			// Verify results
			assert.Equal(t, tt.expectedServerName, serverName, "Server name should match expected")
			if tt.expectedServerConfig == nil {
				assert.Nil(t, serverConfig, "Server config should be nil")
			} else {
				require.NotNil(t, serverConfig, "Server config should not be nil")
				assert.Equal(t, tt.expectedServerConfig.Address, serverConfig.Address, "Server address should match")
				assert.Equal(t, tt.expectedServerConfig.PassthroughMOTD, serverConfig.PassthroughMOTD, "PassthroughMOTD should match")
			}
		})
	}
}

func TestServerConfig_CachePingEnabled(t *testing.T) {
	tests := []struct {
		name     string
		config   config.ServerConfig
		expected bool
	}{
		{
			name: "cache enabled with TTL set and passthrough enabled",
			config: config.ServerConfig{
				PassthroughMOTD: true,                           // Required for cache to be enabled
				CachePingTTL:    configutil.Duration(5_000_000_000), // 5 seconds
			},
			expected: true,
		},
		{
			name: "cache disabled with zero TTL",
			config: config.ServerConfig{
				CachePingTTL: configutil.Duration(0),
			},
			expected: false,
		},
		{
			name: "cache disabled with default (zero) TTL",
			config: config.ServerConfig{
				// CachePingTTL not set (defaults to 0)
			},
			expected: false,
		},
		{
			name: "cache disabled when passthrough is false even with TTL set",
			config: config.ServerConfig{
				PassthroughMOTD: false, // Disabled
				CachePingTTL:    configutil.Duration(5_000_000_000), // 5 seconds
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.CachePingEnabled()
			assert.Equal(t, tt.expected, result, "CachePingEnabled should return correct value")
		})
	}
}

// testError is a test error type that wraps another error for testing error unwrapping
type testError struct {
	Inner error
}

func (e *testError) Error() string {
	return e.Inner.Error()
}

func (e *testError) Unwrap() error {
	return e.Inner
}

// testNetAddr is a test implementation of net.Addr for testing virtual host processing
type testNetAddr struct {
	address string
}

func (t *testNetAddr) Network() string {
	return "tcp"
}

func (t *testNetAddr) String() string {
	return t.address
}

func TestGetVirtualHostnameFromAddr(t *testing.T) {
	tests := []struct {
		name         string
		virtualHost  net.Addr
		expectedHost string
	}{
		{
			name:         "extracts hostname from address with port",
			virtualHost:  &testNetAddr{address: "snapshot.example.test:25565"},
			expectedHost: "snapshot.example.test",
		},
		{
			name:         "extracts hostname from address without port",
			virtualHost:  &testNetAddr{address: "fabric.example.test"},
			expectedHost: "fabric.example.test",
		},
		{
			name:         "handles localhost",
			virtualHost:  &testNetAddr{address: "localhost:25565"},
			expectedHost: "localhost",
		},
		{
			name:         "returns empty string for nil virtual host",
			virtualHost:  nil,
			expectedHost: "",
		},
		{
			name:         "converts to lowercase",
			virtualHost:  &testNetAddr{address: "EXAMPLE.COM:25565"},
			expectedHost: "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxy := &Proxy{}
			result := proxy.getVirtualHostnameFromAddr(tt.virtualHost)
			assert.Equal(t, tt.expectedHost, result, "getVirtualHostnameFromAddr should extract correct hostname")
		})
	}
}
