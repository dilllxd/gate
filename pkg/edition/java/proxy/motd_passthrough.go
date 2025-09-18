package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/jellydator/ttlcache/v3"
	"go.minekube.com/gate/pkg/edition/java/config"
	"go.minekube.com/gate/pkg/edition/java/netmc"
	"go.minekube.com/gate/pkg/edition/java/proto/codec"
	"go.minekube.com/gate/pkg/edition/java/proto/packet"
	"go.minekube.com/gate/pkg/edition/java/proto/state"
	"go.minekube.com/gate/pkg/gate/proto"
	"go.minekube.com/gate/pkg/util/errs"
	"go.minekube.com/gate/pkg/util/netutil"
	"golang.org/x/sync/singleflight"
)

var (
	motdPassthroughCache = ttlcache.New[motdCacheKey, *motdCacheResult]()
	motdSingleFlight     = new(singleflight.Group)
)

func init() {
	go motdPassthroughCache.Start() // start ttl eviction
}

type motdCacheKey struct {
	serverName string
	protocol   proto.Protocol
}

type motdCacheResult struct {
	res *packet.StatusResponse
	err error
}

// ResetMOTDPassthroughCache resets the MOTD passthrough cache.
func ResetMOTDPassthroughCache() {
	motdPassthroughCache.DeleteAll()
}

// IsConnectionRefused returns true if err indicates a connection refused error.
func IsConnectionRefused(err error) bool {
	return err != nil && (errors.Is(err, syscall.ECONNREFUSED) ||
		strings.Contains(strings.ToLower(err.Error()), "connection refused"))
}

// resolveMOTDPassthrough resolves the MOTD by forwarding the ping request to a backend server
// that has passthrough enabled.
func (p *Proxy) resolveMOTDPassthrough(
	log logr.Logger,
	statusRequestCtx *proto.PacketContext,
) (logr.Logger, *packet.StatusResponse, error) {
	// Find the first server with MOTD passthrough enabled
	serverConfig, serverName := p.findPassthroughServer()
	if serverConfig == nil {
		// No server has passthrough enabled, return proxy's own MOTD
		return log, nil, errors.New("no server with MOTD passthrough enabled")
	}

	key := motdCacheKey{serverName, statusRequestCtx.Protocol}

	// Fast path: use cache without loader
	if serverConfig.CachePingEnabled() {
		item := motdPassthroughCache.Get(key)
		if item != nil {
			log.V(1).Info("returning cached MOTD passthrough result", "server", serverName)
			val := item.Value()
			return log, val.res, val.err
		}
	}

	// Slow path: load cache, block many requests to same server
	load := func(ctx context.Context) (*packet.StatusResponse, error) {
		log.V(1).Info("resolving MOTD passthrough", "server", serverName)

		conn, netmcConn, err := p.dialPassthroughServer(ctx, serverConfig.Address, statusRequestCtx.Protocol)
		if err != nil {
			// Use proper verbosity for connection errors
			return nil, &errs.VerbosityError{
				Verbosity: getErrorVerbosity(err),
				Err:       fmt.Errorf("failed to dial passthrough server %s: %w", serverName, err),
			}
		}
		defer func() { _ = netmcConn.Close() }()

		log = log.WithValues("server", serverName, "backendAddr", netutil.Host(netmcConn.RemoteAddr()))
		return p.fetchMOTDFromBackend(log, conn, netmcConn, statusRequestCtx.Protocol)
	}

	if !serverConfig.CachePingEnabled() {
		res, err := load(context.Background())
		return log, res, err
	}

	// Use singleflight to prevent multiple concurrent requests to the same server
	opt := withMOTDLoader(motdSingleFlight, serverConfig.GetCachePingTTL(), func(key motdCacheKey) *motdCacheResult {
		res, err := load(context.Background())
		return &motdCacheResult{res: res, err: err}
	})

	resultChan := make(chan *motdCacheResult, 1)
	go func() { resultChan <- motdPassthroughCache.Get(key, opt).Value() }()

	select {
	case result := <-resultChan:
		return log, result.res, result.err
	case <-time.After(time.Second * 5): // Timeout after 5 seconds
		return log, nil, errors.New("MOTD passthrough request timed out")
	}
}

// findPassthroughServer finds the first server in the Try list that has MOTD passthrough enabled.
func (p *Proxy) findPassthroughServer() (*config.ServerConfig, string) {
	cfg := p.config()

	// Check servers in Try order first
	for _, serverName := range cfg.Try {
		if serverConfig, exists := cfg.Servers[serverName]; exists && serverConfig.PassthroughMOTD {
			return &serverConfig, serverName
		}
	}

	// If no Try servers have passthrough, check all servers
	for serverName, serverConfig := range cfg.Servers {
		if serverConfig.PassthroughMOTD {
			return &serverConfig, serverName
		}
	}

	return nil, ""
}

// dialPassthroughServer dials a backend server for MOTD passthrough.
func (p *Proxy) dialPassthroughServer(
	ctx context.Context,
	serverAddr string,
	protocol proto.Protocol,
) (net.Conn, netmc.MinecraftConn, error) {
	cfg := p.config()

	addr, err := netutil.Parse(serverAddr, "tcp")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse server address %s: %w", serverAddr, err)
	}

	conn, err := net.DialTimeout("tcp", addr.String(), time.Duration(cfg.ConnectionTimeout))
	if err != nil {
		// Use proper verbosity for connection errors (matching lite module pattern)
		return nil, nil, &errs.VerbosityError{
			Verbosity: getErrorVerbosity(err),
			Err:       fmt.Errorf("failed to connect to backend %s: %w", serverAddr, err),
		}
	}

	// Create a proper handshake packet for status connection
	host := netutil.Host(addr)
	port := netutil.Port(addr)
	if port == 0 {
		port = 25565 // Default Minecraft port
	}

	handshake := &packet.Handshake{
		ProtocolVersion: int(protocol),
		ServerAddress:   host,
		Port:            int(port),
		NextStatus:      1, // 1 = Status, 2 = Login
	}

	// Create and send handshake packet
	netmcConn, err := p.sendHandshakeToBackend(conn, handshake, protocol)
	if err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("failed to send handshake packet: %w", err)
	}

	return conn, netmcConn, nil
}

// fetchMOTDFromBackend fetches the MOTD from a backend server.
func (p *Proxy) fetchMOTDFromBackend(
	log logr.Logger,
	conn net.Conn,
	netmcConn netmc.MinecraftConn,
	protocol proto.Protocol,
) (*packet.StatusResponse, error) {
	// Send status request packet
	statusRequest := &packet.StatusRequest{}
	if err := netmcConn.WritePacket(statusRequest); err != nil {
		return nil, fmt.Errorf("failed to write status request packet to backend: %w", err)
	}

	dec := codec.NewDecoder(conn, proto.ClientBound, log.V(2))
	dec.SetProtocol(protocol)
	dec.SetState(state.Status)

	for {
		pktCtx, err := dec.Decode()
		if err != nil {
			return nil, fmt.Errorf("failed to read status response from backend: %w", err)
		}

		if statusResponse, ok := pktCtx.Packet.(*packet.StatusResponse); ok {
			return statusResponse, nil
		}
		// Ignore other packets and continue reading
	}
}

// sendHandshakeToBackend sends a handshake packet to the backend server and returns the netmc connection.
func (p *Proxy) sendHandshakeToBackend(conn net.Conn, handshake *packet.Handshake, protocol proto.Protocol) (netmc.MinecraftConn, error) {
	cfg := p.config()

	// Create a temporary netmc connection for easier packet encoding
	// We're acting as a client to the backend server, so use ClientBound direction
	netmcConn, _ := netmc.NewMinecraftConn(
		context.Background(),
		conn,
		proto.ClientBound,
		time.Duration(cfg.ReadTimeout),
		time.Duration(cfg.ConnectionTimeout),
		0, // no compression
	)
	netmcConn.SetProtocol(protocol)

	// Send handshake packet
	err := netmcConn.WritePacket(handshake)
	if err != nil {
		return nil, err
	}

	// Update state to Status after handshake
	netmcConn.SetState(state.Status)

	return netmcConn, nil
}


// getErrorVerbosity returns appropriate verbosity level for different error types.
// Connection refused errors get debug level to reduce spam.
func getErrorVerbosity(err error) int {
	if IsConnectionRefused(err) {
		return 1 // Debug level for connection refused
	}
	return 0 // Info level for other errors
}

// withMOTDLoader returns a ttlcache option that uses the given load function to load a value for a key
// if it is not already cached.
func withMOTDLoader[K comparable, V any](group *singleflight.Group, ttl time.Duration, load func(key K) V) ttlcache.Option[K, V] {
	loader := ttlcache.LoaderFunc[K, V](
		func(c *ttlcache.Cache[K, V], key K) *ttlcache.Item[K, V] {
			v := load(key)
			return c.Set(key, v, ttl)
		},
	)
	return ttlcache.WithLoader[K, V](loader)
}