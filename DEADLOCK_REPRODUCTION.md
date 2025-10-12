# Gate Lite Mode Deadlock Issue - Reproduction Guide

## Issue Description

When calling `AddLiteRouteBackend` (and other lite-mode configuration modification operations), the request blocks and hangs forever, preventing other lite-mode requests including `ListLiteRoutes` and `GetLiteRoute`. This affects "pretty much any route for lite mode that isn't just simply getting routes or getting a specific route."

## Root Cause Analysis

### The Problem
The issue is a **mutex deadlock** that occurs in the configuration reload sequence. The deadlock happens in this order:

1. `AddLiteRouteBackend()` calls `applyConfigUpdate()`
2. `applyConfigUpdate()` acquires the shared mutex lock
3. `applyConfigUpdate()` calls `reload.FireConfigUpdate()` **WHILE STILL HOLDING THE LOCK**
4. Config reload event handler fires and calls proxy re-initialization
5. Any other concurrent lite-mode requests (`ListLiteRoutes`, `GetLiteRoute`) try to acquire the same mutex
6. **DEADLOCK**: All requests hang indefinitely waiting for the mutex

### Code Location

**File**: `pkg/gate/api_handlers.go`
**Function**: `LiteHandlerImpl.applyConfigUpdate`
**Lines**: ~537-541

### Affected Operations

All lite-mode operations that modify configuration were affected:
- `AddLiteRouteBackend`
- `RemoveLiteRouteBackend`
- `UpdateLiteRouteStrategy`
- `UpdateLiteRouteOptions`
- `UpdateLiteRouteFallback`

Read-only operations would hang when called concurrently with any modification operation:
- `ListLiteRoutes`
- `GetLiteRoute`

## Bug vs Fix Comparison

### 🐛 BUGGY CODE (pre-fix version - commit 8c25f9a)
```go
func (h *LiteHandlerImpl) applyConfigUpdate(ctx context.Context, newCfg config.Config, logMsg string, kv ...any) ([]string, error) {
    // ... validation code ...

    h.mu.Lock()
    prev := *h.cfg
    *h.cfg = newCfg
    reload.FireConfigUpdate(h.eventMgr, h.cfg, &prev)  // ❌ CALLED WHILE HOLDING LOCK
    h.mu.Unlock()

    // ... logging code ...
}
```

### ✅ FIXED CODE (post-fix version - commit 7d869fa)
```go
func (h *LiteHandlerImpl) applyConfigUpdate(ctx context.Context, newCfg config.Config, logMsg string, kv ...any) ([]string, error) {
    // ... validation code ...

    h.mu.Lock()
    prev := *h.cfg
    *h.cfg = newCfg
    h.mu.Unlock()
    reload.FireConfigUpdate(h.eventMgr, h.cfg, &prev)  // ✅ CALLED AFTER RELEASING LOCK

    // ... logging code ...
}
```

## Reproduction Steps

### Prerequisites
- Go development environment
- Gate repository cloned
- Lite mode configuration file

### Method 1: Using Git and Manual Testing

#### 1. Test the Buggy Version (Should Reproduce Hanging)

```bash
# Checkout the pre-fix version
git checkout 8c25f9a

# Build Gate
go build -o gate-buggy ./cmd/gate

# Create a minimal lite config
cat > lite-config.yml << EOF
lite:
  enabled: true
  routes:
    - host: ["example.com"]
      strategy: sequential
      backend: ["backend1.example.com:25565"]
api:
  enabled: true
  listen: "0.0.0.0:8080"
EOF

# Start Gate in background
./gate-buggy run --config lite-config.yml &
GATE_PID=$!

# Wait for Gate to start
sleep 3

# Terminal 1: Start AddLiteRouteBackend (this will hang)
echo "Starting AddLiteRouteBackend request..."
curl -X POST http://localhost:8080/v1/lite/routes/backends/add \
     -H 'Content-Type: application/json' \
     -d '{"host": "example.com", "backend": "backend2.example.com:25565"}' &

# Wait a moment for the first request to acquire the lock
sleep 1

# Terminal 2: Try ListLiteRoutes (this should be blocked)
echo "Starting ListLiteRoutes request..."
timeout 10s curl -X GET http://localhost:8080/v1/lite/routes || echo "ListLiteRoutes timed out (expected)"

# Terminal 3: Try GetLiteRoute (this should be blocked)
echo "Starting GetLiteRoute request..."
timeout 10s curl -X GET http://localhost:8080/v1/lite/routes/example.com || echo "GetLiteRoute timed out (expected)"

# Clean up
kill $GATE_PID
```

#### 2. Test the Fixed Version (Should Work Correctly)

```bash
# Checkout the fixed version
git checkout 7d869fa  # or master

# Build Gate
go build -o gate-fixed ./cmd/gate

# Start Gate in background
./gate-fixed run --config lite-config.yml &
GATE_PID=$!

# Wait for Gate to start
sleep 3

# Terminal 1: Start AddLiteRouteBackend (this should complete)
echo "Starting AddLiteRouteBackend request..."
curl -X POST http://localhost:8080/v1/lite/routes/backends/add \
     -H 'Content-Type: application/json' \
     -d '{"host": "example.com", "backend": "backend2.example.com:25565"}' &

# Wait a moment for the first request to complete
sleep 1

# Terminal 2: Try ListLiteRoutes (this should work)
echo "Starting ListLiteRoutes request..."
curl -X GET http://localhost:8080/v1/lite/routes && echo "ListLiteRoutes completed successfully"

# Terminal 3: Try GetLiteRoute (this should work)
echo "Starting GetLiteRoute request..."
curl -X GET http://localhost:8080/v1/lite/routes/example.com && echo "GetLiteRoute completed successfully"

# Clean up
kill $GATE_PID
```

## Expected Results

### With the Bug (pre-fix version)
- `AddLiteRouteBackend` hangs indefinitely
- `ListLiteRoutes` and `GetLiteRoute` are blocked and timeout
- Multiple concurrent requests result in complete deadlock
- No API responses are received

### With the Fix (post-fix version)
- All requests complete successfully
- `AddLiteRouteBackend` completes and releases the lock immediately
- `ListLiteRoutes` and `GetLiteRoute` can proceed normally
- No blocking or hanging occurs

## Technical Details

### Mutex Sharing
Both `ConfigHandlerImpl` and `LiteHandlerImpl` share the same mutex instance:
```go
// In pkg/gate/api.go
configHandler := NewConfigHandler(&mu, cfg, eventMgr, initialEnable, configFilePath)
liteHandler := NewLiteHandler(&mu, cfg, eventMgr, initialEnable)
//                                      ^^^ Same mutex
```

### Event Firing During Config Reload
The `reload.FireConfigUpdate()` triggers synchronous event handlers that can:
- Re-initialize proxy components
- Perform server synchronization
- Reset caches
- Execute other time-consuming operations

When these operations are executed while holding the config mutex, any other config access blocks indefinitely.

### Fix Rationale
Moving `reload.FireConfigUpdate()` outside the mutex lock ensures:
1. Config updates are atomic and fast
2. Event handlers can access the new configuration without blocking
3. Concurrent API requests can proceed normally
4. No deadlock scenarios can occur

## Verification

To verify the fix is working:

1. **Code Inspection**: Ensure `reload.FireConfigUpdate()` is called after `h.mu.Unlock()`
2. **Runtime Testing**: Run concurrent API requests and verify they complete
3. **Load Testing**: Stress test with multiple concurrent lite-mode operations

## Files Changed

- **Primary Fix**: `pkg/gate/api_handlers.go:541` - Moved `reload.FireConfigUpdate()` outside mutex lock
- **Related Fix**: `pkg/gate/api_handlers.go:167` - Same fix applied to `ConfigHandlerImpl.ApplyConfig`

## Commit History

- **Fix Applied**: `7d869fa` - "fix(api): ensure config reload fires after mutex unlock"
- **Bug Present**: `8c25f9a` and earlier versions

## Git Commands to Test Both Versions

```bash
# Test the buggy version (should reproduce hanging):
git checkout 8c25f9a
go build && ./gate run --config lite-config.yml
# Then run the curl commands above

# Test the fixed version (should work correctly):
git checkout 7d869fa  # or master
go build && ./gate run --config lite-config.yml
# Then run the curl commands above
```