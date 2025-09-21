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
}

var _ process.Runnable = (*consoleRunner)(nil)

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

	if file, ok := r.reader.(*os.File); ok {
		r.isTTY.Store(term.IsTerminal(int(file.Fd())))
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
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(r.reader)
		for scanner.Scan() {
			text := scanner.Text()
			lines <- incoming{line: text}
		}
		if err := scanner.Err(); err != nil {
			lines <- incoming{err: err}
		} else {
			lines <- incoming{err: io.EOF}
		}
		close(lines)
	}()
	defer wg.Wait()

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
				fmt.Fprintf(r.writer, "Error: %v\n", err)
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
		if r.reader != nil {
			_ = r.reader.Close()
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

	switch cmd {
	case "help", "?":
		r.printHelp()
	case "list":
		return r.handleList(rest)
	case "servers":
		return r.listServers()
	case "routes":
		return r.listRoutes()
	case "kick":
		return r.kickPlayer(rest)
	case "move", "send":
		return r.movePlayer(rest)
	default:
		fmt.Fprintf(r.writer, "Unknown command '%s'. Type 'help' for a list.\n", cmd)
	}
	return nil
}

func (r *consoleRunner) printHelp() {
	lines := []string{
		"help               - Show this message",
		"list [players]     - List online players",
		"servers            - List registered backend servers",
		"kick <player> [reason] - Disconnect a player with an optional reason",
		"move <player> <server> - Move a player to a different server",
		"routes             - Show Gate Lite route configuration",
	}
	for _, l := range lines {
		fmt.Fprintln(r.writer, l)
	}
}

func (r *consoleRunner) handleList(args []string) error {
	if len(args) == 0 || strings.EqualFold(args[0], "players") {
		return r.listPlayers()
	}
	switch strings.ToLower(args[0]) {
	case "servers":
		return r.listServers()
	case "routes":
		return r.listRoutes()
	default:
		fmt.Fprintf(r.writer, "Unknown list target '%s'. Try 'players', 'servers', or 'routes'.\n", args[0])
		return nil
	}
}

func (r *consoleRunner) javaProxy() *jproxy.Proxy {
	if r.gate == nil {
		return nil
	}
	return r.gate.Java()
}

func (r *consoleRunner) listPlayers() error {
	proxy := r.javaProxy()
	if proxy == nil {
		fmt.Fprintln(r.writer, "Java proxy not available yet.")
		return nil
	}

	cfg := proxy.Config()
	if cfg.Lite.Enabled {
		fmt.Fprintln(r.writer, "Player management commands are not available in Gate Lite mode.")
		return r.listRoutes()
	}

	players := proxy.Players()
	if len(players) == 0 {
		fmt.Fprintln(r.writer, "No players online.")
		return nil
	}

	sort.Slice(players, func(i, j int) bool {
		return strings.ToLower(players[i].Username()) < strings.ToLower(players[j].Username())
	})

	perServer := map[string][]string{}
	for _, p := range players {
		server := "pending"
		if conn := p.CurrentServer(); conn != nil && conn.Server() != nil {
			server = conn.Server().ServerInfo().Name()
		}
		perServer[server] = append(perServer[server], p.Username())
	}

	serverNames := make([]string, 0, len(perServer))
	for name := range perServer {
		serverNames = append(serverNames, name)
	}
	sort.Slice(serverNames, func(i, j int) bool {
		if serverNames[i] == "pending" {
			return false
		}
		if serverNames[j] == "pending" {
			return true
		}
		return strings.ToLower(serverNames[i]) < strings.ToLower(serverNames[j])
	})

	fmt.Fprintf(r.writer, "Players online (%d):\n", len(players))
	for _, name := range serverNames {
		plist := perServer[name]
		sort.Slice(plist, func(i, j int) bool { return strings.ToLower(plist[i]) < strings.ToLower(plist[j]) })
		fmt.Fprintf(r.writer, "- %s (%d): %s\n", name, len(plist), strings.Join(plist, ", "))
	}

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
		return r.listRoutes()
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

	cfg := proxy.Config()
	routes := cfg.Lite.Routes
	if len(routes) == 0 {
		fmt.Fprintln(r.writer, "No Gate Lite routes configured.")
		return nil
	}

	fmt.Fprintf(r.writer, "Lite routes (%d):\n", len(routes))
	for idx, route := range routes {
		hosts := route.Host.Multi()
		backends := route.Backend.Multi()
		hostStr := strings.Join(hosts, ", ")
		backendStr := strings.Join(backends, ", ")
		strategy := string(route.Strategy)
		if strategy == "" {
			strategy = "sequential"
		}
		fmt.Fprintf(r.writer, "%d. host=%s -> backends=[%s] (strategy: %s)\n", idx+1, hostStr, backendStr, strategy)
	}

	return nil
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

	cfg := proxy.Config()
	if cfg.Lite.Enabled {
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

func (r *consoleRunner) movePlayer(args []string) error {
	if len(args) < 2 {
		fmt.Fprintln(r.writer, "Usage: move <player> <server>")
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

	player := proxy.PlayerByName(args[0])
	if player == nil {
		fmt.Fprintf(r.writer, "Player '%s' is not online.\n", args[0])
		return nil
	}

	target := proxy.Server(args[1])
	if target == nil {
		fmt.Fprintf(r.writer, "Server '%s' is not registered.\n", args[1])
		return nil
	}

	timeout := time.Millisecond * time.Duration(cfg.ConnectionTimeout)
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if player.CreateConnectionRequest(target).ConnectWithIndication(ctx) {
		fmt.Fprintf(r.writer, "Moved %s to %s.\n", player.Username(), target.ServerInfo().Name())
	} else {
		fmt.Fprintf(r.writer, "Failed to move %s to %s. Check logs for details.\n", player.Username(), target.ServerInfo().Name())
	}

	return nil
}
