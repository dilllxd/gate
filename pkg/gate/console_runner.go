package gate

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/term"

	"go.minekube.com/common/minecraft/color"
	"go.minekube.com/common/minecraft/component"
	jproxy "go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/runtime/process"
)

// consoleRunner manages Gate's interactive console commands.
type consoleRunner struct {
	gate   *Gate
	reader io.ReadCloser
	writer io.Writer

	onceClose sync.Once
	isTTY     atomic.Bool
	stopping  atomic.Bool
}

var _ process.Runnable = (*consoleRunner)(nil)

var errStopRequested = errors.New("console stop requested")

func newConsoleRunner(g *Gate) process.Runnable {
	var reader io.ReadCloser
	if os.Stdin != nil {
		reader = os.Stdin
	} else {
		reader = io.NopCloser(strings.NewReader(""))
	}

	writer := io.Writer(os.Stdout)
	if writer == nil {
		writer = io.Discard
	}

	return &consoleRunner{
		gate:   g,
		reader: reader,
		writer: writer,
	}
}

func (r *consoleRunner) Start(ctx context.Context) error {
	log := logr.FromContextOrDiscard(ctx).WithName("console")

	if file, ok := r.reader.(*os.File); ok && file != nil {
		if fd := file.Fd(); fd >= 0 {
			r.isTTY.Store(term.IsTerminal(int(fd)))
		}
	}

	if !r.isTTY.Load() {
		log.V(1).Info("console disabled; stdin not a TTY")
		<-ctx.Done()
		return nil
	}

	fmt.Fprintln(r.writer, "Console ready. Type 'help' for available commands.")
	r.printPrompt()

	type incoming struct {
		line string
		err  error
	}

	lines := make(chan incoming, 1)
	var wg sync.WaitGroup
	wg.Add(2) // Track both scanner and context cancellation goroutines

	// Context cancellation goroutine - tracked by WaitGroup
	go func() {
		defer wg.Done()
		<-ctx.Done()
		r.closeReader()
	}()

	// Scanner goroutine
	go func() {
		defer wg.Done()
		defer close(lines)

		scanner := bufio.NewScanner(r.reader)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Non-blocking check for context cancellation
			}

			if !scanner.Scan() {
				// Scanner finished or errored
				if err := scanner.Err(); err != nil {
					select {
					case lines <- incoming{err: err}:
					case <-ctx.Done():
					}
				} else {
					select {
					case lines <- incoming{err: io.EOF}:
					case <-ctx.Done():
					}
				}
				return
			}

			text := scanner.Text()
			select {
			case lines <- incoming{line: text}:
			case <-ctx.Done():
				return
			}
		}
	}()

	defer func() {
		// Give goroutines a brief moment to clean up
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// All goroutines finished cleanly
		case <-time.After(1 * time.Second):
			// Timeout waiting for goroutines - they may be blocked on scanner.Scan()
			// This is acceptable as the process is shutting down anyway
		}
	}()

	for {
		select {
		case <-ctx.Done():
			r.closeReader()
			return nil
		case in, ok := <-lines:
			if !ok {
				r.closeReader()
				return nil
			}
			if in.err != nil {
				r.closeReader()
				if errors.Is(in.err, io.EOF) || errors.Is(in.err, os.ErrClosed) {
					return nil
				}
				return in.err
			}
			if err := r.execute(strings.TrimSpace(in.line)); err != nil {
				if errors.Is(err, errStopRequested) {
					r.closeReader()
					return nil
				}
				fmt.Fprintf(r.writer, "Error: %v\n", err)
				if r.stopping.Load() {
					r.closeReader()
					return nil
				}
			}
			if r.stopping.Load() {
				r.closeReader()
				return nil
			}
			r.printPrompt()
		}
	}
}

func (r *consoleRunner) printPrompt() {
	if !r.isTTY.Load() {
		return
	}
	fmt.Fprint(r.writer, "> ")
}

func (r *consoleRunner) closeReader() {
	r.onceClose.Do(func() {
		if r.reader == nil {
			return
		}

		done := make(chan struct{})
		// Closing the console reader can block on Windows for tens of seconds,
		// so perform it asynchronously and only wait briefly.
		go func(reader io.Closer) {
			defer close(done)
			_ = reader.Close()
		}(r.reader)

		select {
		case <-done:
		case <-time.After(250 * time.Millisecond):
		}
	})
}

func (r *consoleRunner) execute(line string) error {
	if line == "" {
		return nil
	}

	args := strings.Fields(line)
	cmd := strings.ToLower(args[0])
	rest := args[1:]

	proxy := r.javaProxy()
	liteEnabled := proxy != nil && proxy.Config().Lite.Enabled

	if liteEnabled {
		switch cmd {
		case "help", "?":
			r.printHelp()
		case "routes":
			return r.listRoutes()
		case "stop", "shutdown":
			return r.stopGate(rest)
		default:
			fmt.Fprintf(r.writer, "Command '%s' is not available in Gate Lite mode. Type 'help' for supported commands.\n", cmd)
		}
		return nil
	}

	switch cmd {
	case "help", "?":
		r.printHelp()
	case "list":
		return r.cmdList(rest)
	case "glist":
		return r.cmdGList(rest)
	case "servers":
		return r.listServers()
	case "kick":
		return r.kickPlayer(rest)
	case "move", "send":
		return r.moveCommand(rest)
	case "alert", "broadcast":
		return r.alertCommand(rest)
	case "find":
		return r.findCommand(rest)
	case "ip":
		return r.ipCommand(rest)
	case "reload":
		return r.reloadCommand(rest)
	case "info", "version":
		return r.infoCommand(rest)
	case "stop", "shutdown":
		return r.stopGate(rest)
	case "routes":
		fmt.Fprintln(r.writer, "Command 'routes' is only available in Gate Lite mode.")
	default:
		fmt.Fprintf(r.writer, "Unknown command '%s'. Type 'help' for a list.\n", cmd)
	}
	return nil
}

func (r *consoleRunner) printHelp() {
	proxy := r.javaProxy()
	liteEnabled := proxy != nil && proxy.Config().Lite.Enabled

	if liteEnabled {
		fmt.Fprintln(r.writer, "Available commands:")
		fmt.Fprintln(r.writer, "help - Show available commands")
		fmt.Fprintln(r.writer, "routes - Show Gate Lite route configuration")
		fmt.Fprintln(r.writer, "stop [reason] - Gracefully stop Gate")
	} else {
		fmt.Fprintln(r.writer, "Available commands:")
		fmt.Fprintln(r.writer, "help - Show available commands")
		fmt.Fprintln(r.writer, "list - List all online players")
		fmt.Fprintln(r.writer, "glist [server|player] - List players by server or show player info")
		fmt.Fprintln(r.writer, "servers - List registered backend servers")
		fmt.Fprintln(r.writer, "kick <player> [reason] - Disconnect a player")
		fmt.Fprintln(r.writer, "move <player|server> <server> - Move a player or all players from a server")
		fmt.Fprintln(r.writer, "alert <message> - Send an alert message to all online players")
		fmt.Fprintln(r.writer, "find <player> - Find which server a player is connected to")
		fmt.Fprintln(r.writer, "ip <player> - Show a player's IP address")
		fmt.Fprintln(r.writer, "reload - Show information about automatic config reload")
		fmt.Fprintln(r.writer, "info - Show Gate version and system information")
		fmt.Fprintln(r.writer, "stop [reason] - Gracefully stop Gate")
	}
}

func (r *consoleRunner) cmdList(_ []string) error {
	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	cfg := proxy.Config()
	if cfg.Lite.Enabled {
		fmt.Fprintln(r.writer, "Player listing commands are not available in Gate Lite mode.")
		return nil
	}

	return r.listAllPlayers(proxy)
}

func (r *consoleRunner) cmdGList(args []string) error {
	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	cfg := proxy.Config()
	if cfg.Lite.Enabled {
		fmt.Fprintln(r.writer, "Player listing commands are not available in Gate Lite mode.")
		return nil
	}

	if len(args) == 0 {
		return r.listPlayersByServer(proxy)
	}

	targetName := args[0]

	// Try to find as server first
	if server := proxy.Server(targetName); server != nil {
		return r.listPlayersOnServer(server)
	}

	// Try to find as player
	if player := proxy.PlayerByName(targetName); player != nil {
		return r.showPlayerInfo(player)
	}

	fmt.Fprintf(r.writer, "No server or player named '%s' found.\n", targetName)
	return nil
}

func (r *consoleRunner) showPlayerInfo(player jproxy.Player) error {
	// Player.Username() returns the exact username including Bedrock '*' prefix if present
	fmt.Fprintf(r.writer, "Player: %s\n", player.Username())

	if conn := player.CurrentServer(); conn != nil && conn.Server() != nil {
		server := conn.Server().ServerInfo()
		fmt.Fprintf(r.writer, "  Current server: %s (%s)\n", server.Name(), server.Addr().String())
	} else {
		fmt.Fprintln(r.writer, "  Current server: pending connection")
	}

	fmt.Fprintf(r.writer, "  Protocol version: %d\n", int(player.Protocol()))
	if player.RemoteAddr() != nil {
		fmt.Fprintf(r.writer, "  Remote address: %s\n", player.RemoteAddr().String())
	}

	return nil
}

func (r *consoleRunner) javaProxy() *jproxy.Proxy {
	if r.gate == nil {
		return nil
	}
	return r.gate.Java()
}

func (r *consoleRunner) listAllPlayers(proxy *jproxy.Proxy) error {
	players := proxy.Players()
	if len(players) == 0 {
		fmt.Fprintln(r.writer, "No players online.")
		return nil
	}

	names := make([]string, 0, len(players))
	for _, p := range players {
		names = append(names, p.Username())
	}
	sortStringsCaseInsensitive(names)

	fmt.Fprintf(r.writer, "Players online (%d):\n", len(names))
	fmt.Fprintf(r.writer, "  %s\n", strings.Join(names, ", "))
	return nil
}

func (r *consoleRunner) listPlayersByServer(proxy *jproxy.Proxy) error {
	players := proxy.Players()
	if len(players) == 0 {
		fmt.Fprintln(r.writer, "No players online.")
		return nil
	}

	perServer := map[string][]string{}
	var pending []string
	for _, p := range players {
		server := "pending"
		if conn := p.CurrentServer(); conn != nil && conn.Server() != nil {
			server = conn.Server().ServerInfo().Name()
		}
		if server == "pending" {
			pending = append(pending, p.Username())
		} else {
			perServer[server] = append(perServer[server], p.Username())
		}
	}

	serverNames := make([]string, 0, len(perServer))
	for name := range perServer {
		serverNames = append(serverNames, name)
	}
	sortStringsCaseInsensitive(serverNames)

	fmt.Fprintf(r.writer, "Players by server (%d total):\n", len(players))
	for _, name := range serverNames {
		plist := perServer[name]
		sortStringsCaseInsensitive(plist)
		fmt.Fprintf(r.writer, "- %s (%d): %s\n", name, len(plist), strings.Join(plist, ", "))
	}

	if len(pending) > 0 {
		sortStringsCaseInsensitive(pending)
		fmt.Fprintf(r.writer, "- pending (%d): %s\n", len(pending), strings.Join(pending, ", "))
	}

	return nil
}

func (r *consoleRunner) listPlayersOnServer(server jproxy.RegisteredServer) error {
	players := jproxy.PlayersToSlice[jproxy.Player](server.Players())
	if len(players) == 0 {
		fmt.Fprintf(r.writer, "No players on %s.\n", server.ServerInfo().Name())
		return nil
	}

	names := make([]string, 0, len(players))
	for _, p := range players {
		names = append(names, p.Username())
	}
	sortStringsCaseInsensitive(names)

	fmt.Fprintf(r.writer, "%s (%d players): %s\n", server.ServerInfo().Name(), len(names), strings.Join(names, ", "))
	return nil
}

func (r *consoleRunner) listServers() error {
	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	cfg := proxy.Config()
	if cfg.Lite.Enabled {
		fmt.Fprintln(r.writer, "Registered server list is not available in Gate Lite mode.")
		return nil
	}

	servers := proxy.Servers()
	if len(servers) == 0 {
		fmt.Fprintln(r.writer, "No backend servers registered.")
		return nil
	}

	sort.Slice(servers, func(i, j int) bool {
		return strings.ToLower(servers[i].ServerInfo().Name()) < strings.ToLower(servers[j].ServerInfo().Name())
	})

	fmt.Fprintf(r.writer, "Registered servers (%d):\n", len(servers))
	for _, srv := range servers {
		info := srv.ServerInfo()
		count := srv.Players().Len()
		fmt.Fprintf(r.writer, "- %s (%s) - %d online\n", info.Name(), info.Addr().String(), count)
	}

	return nil
}

func (r *consoleRunner) listRoutes() error {
	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	routes := proxy.Config().Lite.Routes
	if len(routes) == 0 {
		fmt.Fprintln(r.writer, "No Gate Lite routes configured.")
		return nil
	}

	fmt.Fprintf(r.writer, "Lite routes (%d):\n", len(routes))
	for idx, route := range routes {
		var hosts, backends []string
		if route.Host != nil {
			hosts = route.Host.Multi()
		}
		if route.Backend != nil {
			backends = route.Backend.Multi()
		}

		hostStr := "<none>"
		if len(hosts) > 0 {
			hostStr = strings.Join(hosts, ", ")
		}
		backendStr := "<none>"
		if len(backends) > 0 {
			backendStr = strings.Join(backends, ", ")
		}
		strategy := string(route.Strategy)
		if strategy == "" {
			strategy = "sequential"
		}
		fmt.Fprintf(r.writer, "%d. host=%s -> backends=[%s] (strategy: %s)\n", idx+1, hostStr, backendStr, strategy)
	}

	return nil
}

func (r *consoleRunner) stopGate(args []string) error {
	if r.stopping.Swap(true) {
		fmt.Fprintln(r.writer, "Stop already in progress.")
		return errStopRequested
	}

	var (
		reason  component.Component
		message string
	)
	if len(args) > 0 {
		text := strings.Join(args, " ")
		reason = &component.Text{Content: text}
		message = fmt.Sprintf("Stopping Gate: %s", text)
	} else {
		message = "Stopping Gate..."
	}

	fmt.Fprintln(r.writer, message)

	if r.gate != nil {
		r.gate.StopWithReason(reason)
	}

	return errStopRequested
}

func (r *consoleRunner) kickPlayer(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(r.writer, "Usage: kick <player> [reason]")
		return nil
	}

	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	if cfg := proxy.Config(); cfg.Lite.Enabled {
		fmt.Fprintln(r.writer, "Kick command is not available in Gate Lite mode.")
		return nil
	}

	target := proxy.PlayerByName(args[0])
	if target == nil {
		fmt.Fprintf(r.writer, "Player '%s' is not online.\n", args[0])
		return nil
	}

	reason := "Kicked by console"
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}

	target.Disconnect(&component.Text{Content: reason})
	fmt.Fprintf(r.writer, "Kicked %s.\n", target.Username())
	return nil
}

func (r *consoleRunner) moveCommand(args []string) error {
	if len(args) < 2 {
		fmt.Fprintln(r.writer, "Usage: move <player|server> <server>")
		fmt.Fprintln(r.writer, "       move server:<server_name> <destination>  (force server mode)")
		fmt.Fprintln(r.writer, "       move player:<player_name> <destination>  (force player mode)")
		fmt.Fprintln(r.writer, "Note: Bedrock player names may start with '*' - include the full name")
		return nil
	}

	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	cfg := proxy.Config()
	if cfg.Lite.Enabled {
		fmt.Fprintln(r.writer, "Move command is not available in Gate Lite mode.")
		return nil
	}

	destination := proxy.Server(args[len(args)-1])
	if destination == nil {
		fmt.Fprintf(r.writer, "Server '%s' is not registered.\n", args[len(args)-1])
		return nil
	}

	timeout := time.Millisecond * time.Duration(cfg.ConnectionTimeout)
	if timeout <= 100*time.Millisecond {
		timeout = 5 * time.Second
	}

	subject := args[0]

	// Handle explicit prefixes for disambiguation
	if strings.HasPrefix(subject, "server:") {
		serverName := strings.TrimPrefix(subject, "server:")
		r.moveServerPlayers(proxy, serverName, destination, timeout)
		return nil
	}
	if strings.HasPrefix(subject, "player:") {
		playerName := strings.TrimPrefix(subject, "player:")
		if player := proxy.PlayerByName(playerName); player != nil {
			r.moveSinglePlayer(player, destination, timeout)
			return nil
		}
		fmt.Fprintf(r.writer, "Player '%s' is not online.\n", playerName)
		return nil
	}

	// Smart resolution: check both player and server
	player := proxy.PlayerByName(subject)
	server := proxy.Server(subject)

	if player != nil && server != nil {
		// Ambiguous case - both exist
		fmt.Fprintf(r.writer, "Ambiguous: both player '%s' and server '%s' exist.\n", subject, subject)
		fmt.Fprintln(r.writer, "Use 'move player:"+subject+"' or 'move server:"+subject+"' to specify.")
		return nil
	}

	if player != nil {
		r.moveSinglePlayer(player, destination, timeout)
		return nil
	}

	if server != nil {
		r.moveServerPlayers(proxy, subject, destination, timeout)
		return nil
	}

	fmt.Fprintf(r.writer, "No player or server named '%s' found.\n", subject)
	return nil
}

func (r *consoleRunner) moveSinglePlayer(player jproxy.Player, destination jproxy.RegisteredServer, timeout time.Duration) {
	if current := player.CurrentServer(); current != nil && jproxy.RegisteredServerEqual(current.Server(), destination) {
		fmt.Fprintf(r.writer, "%s is already connected to %s.\n", player.Username(), destination.ServerInfo().Name())
		return
	}

	if r.movePlayerWithTimeout(player, destination, timeout) {
		fmt.Fprintf(r.writer, "Moved %s to %s.\n", player.Username(), destination.ServerInfo().Name())
		return
	}

	fmt.Fprintf(r.writer, "Failed to move %s to %s. Check logs for details.\n", player.Username(), destination.ServerInfo().Name())
}

func (r *consoleRunner) moveServerPlayers(proxy *jproxy.Proxy, sourceName string, destination jproxy.RegisteredServer, timeout time.Duration) {
	source := proxy.Server(sourceName)
	if source == nil {
		fmt.Fprintf(r.writer, "Server '%s' is not registered.\n", sourceName)
		return
	}

	if jproxy.RegisteredServerEqual(source, destination) {
		fmt.Fprintf(r.writer, "Source and destination are both %s.\n", destination.ServerInfo().Name())
		return
	}

	players := jproxy.PlayersToSlice[jproxy.Player](source.Players())
	if len(players) == 0 {
		fmt.Fprintf(r.writer, "No players connected to %s.\n", source.ServerInfo().Name())
		return
	}

	moved := make([]string, 0, len(players))
	failed := make([]string, 0)
	for _, player := range players {
		if r.movePlayerWithTimeout(player, destination, timeout) {
			moved = append(moved, player.Username())
		} else {
			failed = append(failed, player.Username())
		}
	}

	sortStringsCaseInsensitive(moved)
	sortStringsCaseInsensitive(failed)

	fmt.Fprintf(r.writer, "Attempted to move %d player(s) from %s to %s.\n", len(players), source.ServerInfo().Name(), destination.ServerInfo().Name())
	if len(moved) > 0 {
		fmt.Fprintf(r.writer, "Moved: %s\n", strings.Join(moved, ", "))
	}
	if len(failed) > 0 {
		fmt.Fprintf(r.writer, "Failed: %s\n", strings.Join(failed, ", "))
	}
}

func (r *consoleRunner) movePlayerWithTimeout(player jproxy.Player, destination jproxy.RegisteredServer, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return player.CreateConnectionRequest(destination).ConnectWithIndication(ctx)
}

func sortStringsCaseInsensitive(values []string) {
	sort.Slice(values, func(i, j int) bool {
		return strings.ToLower(values[i]) < strings.ToLower(values[j])
	})
}

func (r *consoleRunner) alertCommand(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(r.writer, "Usage: alert <message>")
		return nil
	}

	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	if cfg := proxy.Config(); cfg.Lite.Enabled {
		fmt.Fprintln(r.writer, "Alert command is not available in Gate Lite mode.")
		return nil
	}

	players := proxy.Players()
	if len(players) == 0 {
		fmt.Fprintln(r.writer, "No players online to receive the alert.")
		return nil
	}

	message := strings.Join(args, " ")
	alertMsg := &component.Text{
		Content: "[ALERT] " + message,
		S:       component.Style{Color: color.Red, Bold: component.True},
	}

	count := 0
	for _, player := range players {
		if err := player.SendMessage(alertMsg); err == nil {
			count++
		}
	}

	fmt.Fprintf(r.writer, "Alert sent to %d/%d online players: %s\n", count, len(players), message)
	return nil
}

func (r *consoleRunner) findCommand(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(r.writer, "Usage: find <player>")
		return nil
	}

	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	if cfg := proxy.Config(); cfg.Lite.Enabled {
		fmt.Fprintln(r.writer, "Find command is not available in Gate Lite mode.")
		return nil
	}

	player := proxy.PlayerByName(args[0])
	if player == nil {
		fmt.Fprintf(r.writer, "Player '%s' is not online.\n", args[0])
		return nil
	}

	if conn := player.CurrentServer(); conn != nil && conn.Server() != nil {
		server := conn.Server().ServerInfo()
		fmt.Fprintf(r.writer, "Player '%s' is connected to server '%s' (%s)\n",
			player.Username(), server.Name(), server.Addr().String())
	} else {
		fmt.Fprintf(r.writer, "Player '%s' is online but not connected to any server (pending).\n", player.Username())
	}

	return nil
}

func (r *consoleRunner) ipCommand(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(r.writer, "Usage: ip <player>")
		return nil
	}

	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	if cfg := proxy.Config(); cfg.Lite.Enabled {
		fmt.Fprintln(r.writer, "IP command is not available in Gate Lite mode.")
		return nil
	}

	player := proxy.PlayerByName(args[0])
	if player == nil {
		fmt.Fprintf(r.writer, "Player '%s' is not online.\n", args[0])
		return nil
	}

	if addr := player.RemoteAddr(); addr != nil {
		fmt.Fprintf(r.writer, "Player '%s' IP address: %s\n", player.Username(), addr.String())
	} else {
		fmt.Fprintf(r.writer, "Unable to retrieve IP address for player '%s'.\n", player.Username())
	}

	return nil
}

func (r *consoleRunner) reloadCommand(args []string) error {
	fmt.Fprintln(r.writer, "Gate has automatic configuration reload enabled by default.")
	fmt.Fprintln(r.writer, "Configuration changes are automatically detected and applied.")
	fmt.Fprintln(r.writer, "No manual reload command is necessary.")
	return nil
}

func (r *consoleRunner) infoCommand(args []string) error {
	fmt.Fprintln(r.writer, "Gate Proxy Information:")
	fmt.Fprintln(r.writer, "  Version: Gate (build info not available)")

	proxy := r.javaProxy()
	if proxy != nil {
		players := proxy.Players()
		servers := proxy.Servers()
		fmt.Fprintf(r.writer, "  Players online: %d\n", len(players))
		fmt.Fprintf(r.writer, "  Registered servers: %d\n", len(servers))

		cfg := proxy.Config()
		if cfg.Lite.Enabled {
			fmt.Fprintln(r.writer, "  Mode: Gate Lite")
			fmt.Fprintf(r.writer, "  Lite routes: %d\n", len(cfg.Lite.Routes))
		} else {
			fmt.Fprintln(r.writer, "  Mode: Full proxy")
		}
	}

	if r.gate != nil {
		if bedrock := r.gate.Bedrock(); bedrock != nil {
			fmt.Fprintln(r.writer, "  Bedrock support: Enabled")
		} else {
			fmt.Fprintln(r.writer, "  Bedrock support: Disabled")
		}
	}

	return nil
}
