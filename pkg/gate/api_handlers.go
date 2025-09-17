package gate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"gopkg.in/yaml.v3"

	config2 "go.minekube.com/gate/pkg/edition/java/lite/config"
	jping "go.minekube.com/gate/pkg/edition/java/ping"
	jver "go.minekube.com/gate/pkg/edition/java/proto/version"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/gate/config"
	gproto "go.minekube.com/gate/pkg/gate/proto"
	"go.minekube.com/gate/pkg/internal/api"
	pb "go.minekube.com/gate/pkg/internal/api/gen/minekube/gate/v1"
	"go.minekube.com/gate/pkg/internal/reload"
	"go.minekube.com/gate/pkg/util/componentutil"
	"go.minekube.com/gate/pkg/util/configutil"
	"go.minekube.com/gate/pkg/util/favicon"
	"go.minekube.com/gate/pkg/version"
)

// ConfigHandlerImpl implements the ConfigHandler interface
type ConfigHandlerImpl struct {
	mu             *sync.Mutex
	cfg            *config.Config
	eventMgr       event.Manager
	proxy          *proxy.Proxy
	configFilePath string
}

func NewConfigHandler(mu *sync.Mutex, cfg *config.Config, eventMgr event.Manager, proxy *proxy.Proxy, configFilePath string) *ConfigHandlerImpl {
	return &ConfigHandlerImpl{
		mu:             mu,
		cfg:            cfg,
		eventMgr:       eventMgr,
		proxy:          proxy,
		configFilePath: configFilePath,
	}
}

func (h *ConfigHandlerImpl) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	h.mu.Lock()
	isLiteMode := h.cfg.Config.Lite.Enabled
	h.mu.Unlock()

	response := &pb.GetStatusResponse{
		Version: version.String(),
	}

	if isLiteMode {
		response.Mode = pb.ProxyMode_PROXY_MODE_LITE

		// Get lite mode statistics
		h.mu.Lock()
		routes := h.cfg.Config.Lite.Routes
		h.mu.Unlock()

		// Count total active connections across all backends
		sm := h.proxy.Lite().StrategyManager()
		var totalConnections int32
		for _, route := range routes {
			for _, backend := range route.Backend {
				if counter := sm.GetOrCreateCounter(backend); counter != nil {
					totalConnections += int32(counter.Load())
				}
			}
		}

		response.Stats = &pb.GetStatusResponse_Lite{
			Lite: &pb.LiteStats{
				Connections: totalConnections,
				Routes:      int32(len(routes)),
			},
		}
	} else {
		response.Mode = pb.ProxyMode_PROXY_MODE_CLASSIC

		// Count players in classic mode
		var players int32
		for _, s := range h.proxy.Servers() {
			s.Players().Range(func(proxy.Player) bool {
				players++
				return true
			})
		}

		response.Stats = &pb.GetStatusResponse_Classic{
			Classic: &pb.ClassicStats{
				Players: players,
				Servers: int32(len(h.proxy.Servers())),
			},
		}
	}

	return response, nil
}

func (h *ConfigHandlerImpl) GetConfig(ctx context.Context, req *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	h.mu.Lock()
	cfgCopy := *h.cfg
	h.mu.Unlock()

	data, err := yaml.Marshal(cfgCopy)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to encode config: %w", err))
	}
	return &pb.GetConfigResponse{Payload: string(data)}, nil
}

func (h *ConfigHandlerImpl) ValidateConfig(ctx context.Context, req *pb.ValidateConfigRequest) ([]string, error) {
	var newCfg config.Config
	if err := yaml.Unmarshal([]byte(req.GetConfig()), &newCfg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid YAML/JSON: %w", err))
	}
	warns, errs := newCfg.Validate()
	if len(errs) > 0 {
		errStrs := make([]string, len(errs))
		for i, err := range errs {
			errStrs[i] = err.Error()
		}
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("config validation failed: %s", strings.Join(errStrs, "; ")))
	}
	warnStrs := make([]string, len(warns))
	for i, warn := range warns {
		warnStrs[i] = warn.Error()
	}
	return warnStrs, nil
}

func (h *ConfigHandlerImpl) ApplyConfig(ctx context.Context, req *pb.ApplyConfigRequest) ([]string, error) {
	var newCfg config.Config
	if err := yaml.Unmarshal([]byte(req.GetConfig()), &newCfg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid YAML/JSON: %w", err))
	}

	warns, errs := newCfg.Validate()
	if len(errs) > 0 {
		errStrs := make([]string, len(errs))
		for i, err := range errs {
			errStrs[i] = err.Error()
		}
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("config validation failed: %s", strings.Join(errStrs, "; ")))
	}

	h.mu.Lock()
	prev := *h.cfg
	*h.cfg = newCfg
	reload.FireConfigUpdate(h.eventMgr, h.cfg, &prev)
	h.mu.Unlock()
	logr.FromContextOrDiscard(ctx).Info("applied config via api")

	warnStrs := make([]string, len(warns))
	for i, warn := range warns {
		warnStrs[i] = warn.Error()
	}

	// If persist is enabled, try to write the config to disk
	if req.GetPersist() {
		if err := h.persistConfig(&newCfg); err != nil {
			logr.FromContextOrDiscard(ctx).Error(err, "failed to persist config to disk (config applied in-memory)")
			warnStrs = append(warnStrs, fmt.Sprintf("failed to persist config to disk: %v", err))
		} else {
			logr.FromContextOrDiscard(ctx).Info("config persisted to disk")
		}
	}

	return warnStrs, nil
}

func (h *ConfigHandlerImpl) persistConfig(cfg *config.Config) error {
	configFile := h.configFilePath
	if configFile == "" {
		return errors.New("config file path not available - cannot persist config")
	}

	// Determine format from file extension
	var (
		data []byte
		err  error
	)
	switch path.Ext(configFile) {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(cfg)
	default:
		return fmt.Errorf("unsupported config file format: %s (only .yml and .yaml are supported)", configFile)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configFile, err)
	}

	return nil
}

// LiteHandlerImpl implements the LiteHandler interface
type LiteHandlerImpl struct {
	mu       *sync.Mutex
	cfg      *config.Config
	eventMgr event.Manager
	proxy    *proxy.Proxy
}

func NewLiteHandler(mu *sync.Mutex, cfg *config.Config, eventMgr event.Manager, proxy *proxy.Proxy) *LiteHandlerImpl {
	return &LiteHandlerImpl{
		mu:       mu,
		cfg:      cfg,
		eventMgr: eventMgr,
		proxy:    proxy,
	}
}

func (h *LiteHandlerImpl) ListLiteRoutes(ctx context.Context, req *pb.ListLiteRoutesRequest) (*pb.ListLiteRoutesResponse, error) {
	h.mu.Lock()
	routes := make([]config2.Route, len(h.cfg.Config.Lite.Routes))
	copy(routes, h.cfg.Config.Lite.Routes)
	h.mu.Unlock()
	resp := &pb.ListLiteRoutesResponse{}
	for _, r := range routes {
		resp.Routes = append(resp.Routes, h.toProtoRoute(r))
	}
	return resp, nil
}

func (h *LiteHandlerImpl) GetLiteRoute(ctx context.Context, req *pb.GetLiteRouteRequest) (*pb.GetLiteRouteResponse, error) {
	host := strings.TrimSpace(req.GetHost())
	if host == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("host is required"))
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	idx := h.findRouteIdx(host)
	if idx < 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("route not found"))
	}
	return &pb.GetLiteRouteResponse{Route: h.toProtoRoute(h.cfg.Config.Lite.Routes[idx])}, nil
}

func (h *LiteHandlerImpl) UpdateLiteRouteStrategy(ctx context.Context, req *pb.UpdateLiteRouteStrategyRequest) ([]string, error) {
	host := strings.TrimSpace(req.GetHost())
	if host == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("host is required"))
	}
	strategy := strings.TrimSpace(api.ConvertStrategyToString(req.GetStrategy()))
	if strategy == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("strategy is required"))
	}
	h.mu.Lock()
	newCfg := *h.cfg
	h.mu.Unlock()
	idx := h.findRouteIdxInConfig(&newCfg, host)
	if idx < 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("route not found"))
	}
	old := string(newCfg.Config.Lite.Routes[idx].Strategy)
	newCfg.Config.Lite.Routes[idx].Strategy = config2.Strategy(strategy)
	return h.applyConfigUpdate(ctx, newCfg, "lite route strategy updated", "host", host, "old", old, "new", strategy)
}

func (h *LiteHandlerImpl) AddLiteRouteBackend(ctx context.Context, req *pb.AddLiteRouteBackendRequest) ([]string, error) {
	host := strings.TrimSpace(req.GetHost())
	backend := strings.TrimSpace(req.GetBackend())
	if host == "" || backend == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("host and backend are required"))
	}
	h.mu.Lock()
	newCfg := *h.cfg
	h.mu.Unlock()
	idx := h.findRouteIdxInConfig(&newCfg, host)
	if idx < 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("route not found"))
	}
	exists := false
	for _, b := range newCfg.Config.Lite.Routes[idx].Backend {
		if strings.EqualFold(b, backend) {
			exists = true
			break
		}
	}
	if !exists {
		newCfg.Config.Lite.Routes[idx].Backend = append(newCfg.Config.Lite.Routes[idx].Backend, backend)
	}
	return h.applyConfigUpdate(ctx, newCfg, "lite route backend added", "host", host, "backend", backend, "alreadyExisted", exists)
}

func (h *LiteHandlerImpl) RemoveLiteRouteBackend(ctx context.Context, req *pb.RemoveLiteRouteBackendRequest) ([]string, error) {
	host := strings.TrimSpace(req.GetHost())
	backend := strings.TrimSpace(req.GetBackend())
	if host == "" || backend == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("host and backend are required"))
	}
	h.mu.Lock()
	newCfg := *h.cfg
	h.mu.Unlock()
	idx := h.findRouteIdxInConfig(&newCfg, host)
	if idx < 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("route not found"))
	}
	bs := newCfg.Config.Lite.Routes[idx].Backend
	filtered := make([]string, 0, len(bs))
	removed := false
	for _, b := range bs {
		if strings.EqualFold(b, backend) {
			removed = true
			continue
		}
		filtered = append(filtered, b)
	}
	newCfg.Config.Lite.Routes[idx].Backend = filtered
	return h.applyConfigUpdate(ctx, newCfg, "lite route backend removed", "host", host, "backend", backend, "removed", removed)
}

func (h *LiteHandlerImpl) UpdateLiteRouteOptions(ctx context.Context, req *pb.UpdateLiteRouteOptionsRequest) ([]string, error) {
	if req.GetOptions() == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("options payload is required"))
	}
	host := strings.TrimSpace(req.GetHost())
	if host == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("host is required"))
	}
	h.mu.Lock()
	newCfg := *h.cfg
	h.mu.Unlock()
	idx := h.findRouteIdxInConfig(&newCfg, host)
	if idx < 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("route not found"))
	}
	opts := req.GetOptions()
	paths := req.GetUpdateMask().GetPaths()
	if len(paths) == 0 {
		paths = []string{"proxy_protocol", "tcp_shield_real_ip", "modify_virtual_host", "cache_ping_ttl_ms"}
	}
	for _, path := range paths {
		switch path {
		case "proxy_protocol":
			newCfg.Config.Lite.Routes[idx].ProxyProtocol = opts.GetProxyProtocol()
		case "tcp_shield_real_ip":
			newCfg.Config.Lite.Routes[idx].TCPShieldRealIP = opts.GetTcpShieldRealIp()
		case "modify_virtual_host":
			newCfg.Config.Lite.Routes[idx].ModifyVirtualHost = opts.GetModifyVirtualHost()
		case "cache_ping_ttl_ms":
			newCfg.Config.Lite.Routes[idx].CachePingTTL = configutil.Duration(opts.GetCachePingTtlMs())
		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unsupported field mask path %q", path))
		}
	}
	return h.applyConfigUpdate(ctx, newCfg, "lite route options updated", "host", host)
}

func (h *LiteHandlerImpl) UpdateLiteRouteFallback(ctx context.Context, req *pb.UpdateLiteRouteFallbackRequest) ([]string, error) {
	host := strings.TrimSpace(req.GetHost())
	if host == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("host is required"))
	}
	h.mu.Lock()
	newCfg := *h.cfg
	h.mu.Unlock()
	idx := h.findRouteIdxInConfig(&newCfg, host)
	if idx < 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("route not found"))
	}
	if newCfg.Config.Lite.Routes[idx].Fallback == nil {
		newCfg.Config.Lite.Routes[idx].Fallback = &config2.Status{}
	}
	fb := newCfg.Config.Lite.Routes[idx].Fallback
	paths := req.GetUpdateMask().GetPaths()
	if len(paths) == 0 {
		paths = []string{"motd_json", "version", "players", "favicon"}
	}
	for _, path := range paths {
		switch path {
		case "motd_json":
			if req.GetFallback() == nil || strings.TrimSpace(req.GetFallback().GetMotdJson()) == "" {
				fb.MOTD = nil
			} else {
				motd, err := h.parseMOTD(req.GetFallback().GetMotdJson())
				if err != nil {
					return nil, err
				}
				fb.MOTD = motd
			}
		case "version":
			if req.GetFallback() == nil || req.GetFallback().GetVersion() == nil {
				fb.Version = jping.Version{}
			} else {
				fb.Version.Name = req.GetFallback().GetVersion().GetName()
				fb.Version.Protocol = gproto.Protocol(req.GetFallback().GetVersion().GetProtocol())
			}
		case "players":
			if req.GetFallback() == nil || req.GetFallback().GetPlayers() == nil {
				fb.Players = nil
			} else {
				if fb.Players == nil {
					fb.Players = &jping.Players{}
				}
				fb.Players.Online = int(req.GetFallback().GetPlayers().GetOnline())
				fb.Players.Max = int(req.GetFallback().GetPlayers().GetMax())
			}
		case "favicon":
			if req.GetFallback() == nil || strings.TrimSpace(req.GetFallback().GetFavicon()) == "" {
				fb.Favicon = ""
			} else {
				fb.Favicon = favicon.Favicon(req.GetFallback().GetFavicon())
			}
		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unsupported field mask path %q", path))
		}
	}
	return h.applyConfigUpdate(ctx, newCfg, "lite route fallback updated", "host", host)
}

func (h *LiteHandlerImpl) toProtoFallback(src *config2.Status) *pb.LiteRouteFallback {
	if src == nil {
		return nil
	}
	pbFallback := &pb.LiteRouteFallback{}
	if src.MOTD != nil {
		if data, err := json.Marshal(src.MOTD); err == nil {
			pbFallback.MotdJson = string(data)
		}
	}
	if src.Version.Name != "" || src.Version.Protocol != 0 {
		pbFallback.Version = &pb.LiteRouteFallbackVersion{
			Name:     src.Version.Name,
			Protocol: int32(src.Version.Protocol),
		}
	}
	if src.Players != nil {
		pbFallback.Players = &pb.LiteRouteFallbackPlayers{
			Online: int32(src.Players.Online),
			Max:    int32(src.Players.Max),
		}
	}
	if src.Favicon != "" {
		pbFallback.Favicon = string(src.Favicon)
	}
	return pbFallback
}

func (h *LiteHandlerImpl) toProtoRoute(route config2.Route) *pb.LiteRoute {
	sm := h.proxy.Lite().StrategyManager()
	pbRoute := &pb.LiteRoute{
		Hosts:    route.Host,
		Strategy: api.ConvertStrategyFromString(string(route.Strategy)),
		Options: &pb.LiteRouteOptions{
			ProxyProtocol:     route.ProxyProtocol,
			TcpShieldRealIp:   route.GetTCPShieldRealIP(),
			ModifyVirtualHost: route.ModifyVirtualHost,
			CachePingTtlMs:    int64(route.CachePingTTL),
		},
	}
	for _, backend := range route.Backend {
		var active uint32
		if counter := sm.GetOrCreateCounter(backend); counter != nil {
			active = counter.Load()
		}
		pbRoute.Backends = append(pbRoute.Backends, &pb.LiteRouteBackend{
			Address:           backend,
			ActiveConnections: active,
		})
	}
	pbRoute.Fallback = h.toProtoFallback(route.Fallback)
	return pbRoute
}

func (h *LiteHandlerImpl) findRouteIdx(host string) int {
	for i, r := range h.cfg.Config.Lite.Routes {
		for _, h := range r.Host {
			if strings.EqualFold(h, host) {
				return i
			}
		}
	}
	return -1
}

func (h *LiteHandlerImpl) findRouteIdxInConfig(c *config.Config, host string) int {
	for i, r := range c.Config.Lite.Routes {
		for _, h := range r.Host {
			if strings.EqualFold(h, host) {
				return i
			}
		}
	}
	return -1
}

func (h *LiteHandlerImpl) parseMOTD(s string) (*configutil.TextComponent, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	tc, err := componentutil.ParseTextComponent(jver.MinimumVersion.Protocol, s)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid motd: %w", err))
	}
	return (*configutil.TextComponent)(tc), nil
}

func (h *LiteHandlerImpl) applyConfigUpdate(ctx context.Context, newCfg config.Config, logMsg string, kv ...any) ([]string, error) {
	warns, errs := newCfg.Validate()
	if len(errs) > 0 {
		errStrs := make([]string, len(errs))
		for i, err := range errs {
			errStrs[i] = err.Error()
		}
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("config validation failed: %s", strings.Join(errStrs, "; ")))
	}
	warnStrs := make([]string, len(warns))
	for i, warn := range warns {
		warnStrs[i] = warn.Error()
	}
	h.mu.Lock()
	prev := *h.cfg
	*h.cfg = newCfg
	reload.FireConfigUpdate(h.eventMgr, h.cfg, &prev)
	h.mu.Unlock()
	logr.FromContextOrDiscard(ctx).Info(logMsg, kv...)
	return warnStrs, nil
}

// EventsHandlerImpl implements the EventsHandler interface for real-time event streaming
type EventsHandlerImpl struct {
	mu       *sync.Mutex
	cfg      *config.Config
	eventMgr event.Manager
	proxy    *proxy.Proxy
}

func NewEventsHandler(mu *sync.Mutex, cfg *config.Config, eventMgr event.Manager, proxy *proxy.Proxy) *EventsHandlerImpl {
	return &EventsHandlerImpl{
		mu:       mu,
		cfg:      cfg,
		eventMgr: eventMgr,
		proxy:    proxy,
	}
}

// StreamEvents implements the server streaming RPC for real-time event monitoring
func (h *EventsHandlerImpl) StreamEvents(ctx context.Context, req *connect.Request[pb.StreamEventsRequest], stream *connect.ServerStream[pb.ProxyEvent]) error {
	log := logr.FromContextOrDiscard(ctx)
	log.Info("client connected to event stream")

	// Create buffered channel for events with reasonable buffer size
	eventChan := make(chan *pb.ProxyEvent, 100)
	defer close(eventChan)

	// Set up event filters
	eventTypes := make(map[pb.EventType]bool)
	if len(req.Msg.EventTypes) > 0 {
		for _, eventType := range req.Msg.EventTypes {
			eventTypes[eventType] = true
		}
	} else {
		// If no specific types requested, include all
		for i := pb.EventType_EVENT_TYPE_PLAYER_CONNECT; i <= pb.EventType_EVENT_TYPE_SHUTDOWN; i++ {
			eventTypes[i] = true
		}
	}

	// Handle defaults for event categories
	// Since protobuf bools default to false, but our API documentation says default is true,
	// we treat the case where both are false as "include everything"
	includePlayerEvents := req.Msg.IncludePlayerEvents
	includeSystemEvents := req.Msg.IncludeSystemEvents

	// If client sends empty request {}, include all events by default
	if !includePlayerEvents && !includeSystemEvents {
		includePlayerEvents = true
		includeSystemEvents = true
	}

	// Subscribe to various events and create conversions
	unsubscribers := h.subscribeToEvents(eventChan, eventTypes, includePlayerEvents, includeSystemEvents)
	defer func() {
		for _, unsub := range unsubscribers {
			unsub()
		}
	}()

	// Stream events until client disconnects or context is canceled
	for {
		select {
		case <-ctx.Done():
			log.Info("client disconnected from event stream")
			return nil
		case proxyEvent, ok := <-eventChan:
			if !ok {
				return nil
			}
			if err := stream.Send(proxyEvent); err != nil {
				log.Error(err, "failed to send event to client")
				return err
			}
		}
	}
}

// subscribeToEvents sets up event listeners and returns unsubscribe functions
func (h *EventsHandlerImpl) subscribeToEvents(eventChan chan<- *pb.ProxyEvent, eventTypes map[pb.EventType]bool, includePlayer, includeSystem bool) []func() {
	var unsubscribers []func()

	// Subscribe to config update events
	if includeSystem && eventTypes[pb.EventType_EVENT_TYPE_CONFIG_UPDATE] {
		unsub := reload.Subscribe(h.eventMgr, func(e *reload.ConfigUpdateEvent[config.Config]) {
			h.mu.Lock()
			isLiteMode := e.Config.Config.Lite.Enabled
			routeCount := len(e.Config.Config.Lite.Routes)
			h.mu.Unlock()

			proxyEvent := &pb.ProxyEvent{
				EventType:   pb.EventType_EVENT_TYPE_CONFIG_UPDATE,
				TimestampMs: time.Now().UnixMilli(),
				EventData: &pb.ProxyEvent_ConfigUpdate{
					ConfigUpdate: &pb.ConfigUpdateEvent{
						ConfigSource:     "api",
						LiteModeEnabled:  isLiteMode,
						RouteCount:       int32(routeCount),
					},
				},
			}
			select {
			case eventChan <- proxyEvent:
			default:
				// Channel full, drop event to prevent blocking
			}
		})
		unsubscribers = append(unsubscribers, unsub)
	}

	// Subscribe to player connection events
	if includePlayer && eventTypes[pb.EventType_EVENT_TYPE_PLAYER_CONNECT] {
		unsub := event.Subscribe(h.eventMgr, 0, func(e *proxy.PostLoginEvent) {
			player := e.Player()
			var remoteAddr string
			if addr := player.RemoteAddr(); addr != nil {
				remoteAddr = addr.String()
			}

			proxyEvent := &pb.ProxyEvent{
				EventType:   pb.EventType_EVENT_TYPE_PLAYER_CONNECT,
				TimestampMs: time.Now().UnixMilli(),
				EventData: &pb.ProxyEvent_PlayerConnect{
					PlayerConnect: &pb.PlayerConnectEvent{
						PlayerId:        player.ID().String(),
						PlayerUsername:  player.Username(),
						RemoteAddress:   remoteAddr,
						ProtocolVersion: int32(player.Protocol()),
					},
				},
			}
			select {
			case eventChan <- proxyEvent:
			default:
			}
		})
		unsubscribers = append(unsubscribers, unsub)
	}

	// Subscribe to player disconnect events
	if includePlayer && eventTypes[pb.EventType_EVENT_TYPE_PLAYER_DISCONNECT] {
		unsub := event.Subscribe(h.eventMgr, 0, func(e *proxy.DisconnectEvent) {
			player := e.Player()
			loginStatus := h.convertLoginStatus(e.LoginStatus())

			proxyEvent := &pb.ProxyEvent{
				EventType:   pb.EventType_EVENT_TYPE_PLAYER_DISCONNECT,
				TimestampMs: time.Now().UnixMilli(),
				EventData: &pb.ProxyEvent_PlayerDisconnect{
					PlayerDisconnect: &pb.PlayerDisconnectEvent{
						PlayerId:       player.ID().String(),
						PlayerUsername: player.Username(),
						Reason:         "", // DisconnectEvent doesn't provide reason
						LoginStatus:    loginStatus,
					},
				},
			}
			select {
			case eventChan <- proxyEvent:
			default:
			}
		})
		unsubscribers = append(unsubscribers, unsub)
	}

	// Subscribe to server switch events
	if includePlayer && eventTypes[pb.EventType_EVENT_TYPE_PLAYER_SERVER_SWITCH] {
		unsub := event.Subscribe(h.eventMgr, 0, func(e *proxy.ServerConnectedEvent) {
			player := e.Player()
			var fromServer string
			if prevServer := e.PreviousServer(); prevServer != nil {
				fromServer = prevServer.ServerInfo().Name()
			}

			proxyEvent := &pb.ProxyEvent{
				EventType:   pb.EventType_EVENT_TYPE_PLAYER_SERVER_SWITCH,
				TimestampMs: time.Now().UnixMilli(),
				EventData: &pb.ProxyEvent_PlayerServerSwitch{
					PlayerServerSwitch: &pb.PlayerServerSwitchEvent{
						PlayerId:       player.ID().String(),
						PlayerUsername: player.Username(),
						FromServer:     fromServer,
						ToServer:       e.Server().ServerInfo().Name(),
					},
				},
			}
			select {
			case eventChan <- proxyEvent:
			default:
			}
		})
		unsubscribers = append(unsubscribers, unsub)
	}

	// Subscribe to ready events
	if includeSystem && eventTypes[pb.EventType_EVENT_TYPE_READY] {
		unsub := event.Subscribe(h.eventMgr, 0, func(e *proxy.ReadyEvent) {
			h.mu.Lock()
			isLiteMode := h.cfg.Config.Lite.Enabled
			h.mu.Unlock()

			var proxyMode pb.ProxyMode
			if isLiteMode {
				proxyMode = pb.ProxyMode_PROXY_MODE_LITE
			} else {
				proxyMode = pb.ProxyMode_PROXY_MODE_CLASSIC
			}

			proxyEvent := &pb.ProxyEvent{
				EventType:   pb.EventType_EVENT_TYPE_READY,
				TimestampMs: time.Now().UnixMilli(),
				EventData: &pb.ProxyEvent_Ready{
					Ready: &pb.ReadyEvent{
						BindAddress: e.Addr(),
						ProxyMode:   proxyMode,
					},
				},
			}
			select {
			case eventChan <- proxyEvent:
			default:
			}
		})
		unsubscribers = append(unsubscribers, unsub)
	}

	// Subscribe to shutdown events
	if includeSystem && eventTypes[pb.EventType_EVENT_TYPE_SHUTDOWN] {
		unsub := event.Subscribe(h.eventMgr, 0, func(e *proxy.PreShutdownEvent) {
			var reason string
			if e.Reason() != nil {
				// Convert component to JSON representation as string
				if data, err := json.Marshal(e.Reason()); err == nil {
					reason = string(data)
				}
			}

			proxyEvent := &pb.ProxyEvent{
				EventType:   pb.EventType_EVENT_TYPE_SHUTDOWN,
				TimestampMs: time.Now().UnixMilli(),
				EventData: &pb.ProxyEvent_Shutdown{
					Shutdown: &pb.ShutdownEvent{
						Reason: reason,
					},
				},
			}
			select {
			case eventChan <- proxyEvent:
			default:
			}
		})
		unsubscribers = append(unsubscribers, unsub)
	}

	return unsubscribers
}

// convertLoginStatus converts internal LoginStatus to protobuf LoginStatus
func (h *EventsHandlerImpl) convertLoginStatus(status proxy.LoginStatus) pb.LoginStatus {
	switch status {
	case proxy.SuccessfulLoginStatus:
		return pb.LoginStatus_LOGIN_STATUS_SUCCESSFUL
	case proxy.ConflictingLoginStatus:
		return pb.LoginStatus_LOGIN_STATUS_CONFLICTING
	case proxy.CanceledByUserLoginStatus:
		return pb.LoginStatus_LOGIN_STATUS_CANCELED_BY_USER
	case proxy.CanceledByProxyLoginStatus:
		return pb.LoginStatus_LOGIN_STATUS_CANCELED_BY_PROXY
	case proxy.CanceledByUserBeforeCompleteLoginStatus:
		return pb.LoginStatus_LOGIN_STATUS_CANCELED_BEFORE_COMPLETE
	default:
		return pb.LoginStatus_LOGIN_STATUS_UNSPECIFIED
	}
}