![Logo](docs/images/cover3.png)

# The Minecraft Proxy _(alpha)_

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/minekube/gate?sort=semver)
[![Doc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go)](https://pkg.go.dev/go.minekube.com/gate)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/minekube/gate?logo=go)
[![Go Report Card](https://goreportcard.com/badge/go.minekube.com/gate)](https://goreportcard.com/report/go.minekube.com/gate)
![test](https://github.com/minekube/gate/workflows/test/badge.svg)
![Discord](https://img.shields.io/discord/633708750032863232?logo=discord)

**Gate is an extensible Minecraft proxy written in Go**

> This project is in its early stages, not yet ready for production and
> subject to have breaking changes,
> but you can already start playing around with it!
>
> It will be tested & used in production by the Minekube network when
> there is a "stable enough" release.

### Features

- [**Fast**](#benchmarks)
- High performant parallelism (see benchmarks)
- [Quick installation](#quick-sample-setup)
    - simply pick a download from the releases
    - support windows/macOS/linux/...
    - single executable binary
    - (No Java runtime needed)
- Benefits from Go's awesome language features
    - simple, reliable, efficient
    - [and much more](https://golang.org/)

### Target audiences
- advanced networks wanting performance while operating at a high scale
- simple Minecraft network admins when [scripting languages](#script-languages)
are supported


## What Gate does

The whole point of a Minecraft proxy is to be able to
move players between servers without fully disconnecting them,
like switching the world but server-wise.

Similar to the proxies
[Velocity](https://github.com/VelocityPowered/Velocity)
_(where much of the knowledge and ideas for this proxy comes from)_,
[BungeeCord](https://github.com/SpigotMC/BungeeCord),
[Waterfall](https://github.com/PaperMC/Waterfall) etc.
Gate delivers rich interfaces to interact with connected players
on a network of Minecraft servers.

Therefore, Gate reads all packets sent between
players (Minecraft client) and servers (e.g. Minecraft spigot, paper, sponge, ...),
logs state changes and emits different events that 
custom plugins/code can react to.

## Benchmarks

> TODO

## Quick sample setup

This is a simple setup of a Minecraft network using Gate proxy,
a Paper 1.16.1 (server1) & Paper 1.8.8 (server2).

**You will only need a JRE (Java Runtime Environment, 8 or higher).**

1. `git clone https://github.com/minekube/gate.git`
2. Go to `docs/sample`
3. Open a terminal in this directory
and run the pre-compiled executable
    - linux: `./gate`
    - windows: `gate.exe`
    - Or build an executable yourself. ;)
4. Open two terminals, one for each `server1` & `server2`.
    - Run `java -jar <the server jar>` to run both servers.
    
Now you can connect to the network on `localhost:25565`
with a Minecraft version 1.16.1 and 1.8.x.
Gate tries to connect you to one of the servers as specified in the configuration.

Checkout the `/server` command
(temporarily hardcoded until commands system has been added).

> There will be an expressive documentation website for Gate later on!

## Known Bugs / Limitations

Remember that Gate is in a very early stage and there are many
more things planned (e.g. command system, more events, stability etc.).

- Can't login sometime in online mode.
  - _Dev note: Is this due to Mojang API database latencies when
 replicating data for the `hasJoined` endpoint?_

## Extending Gate with custom code

- [Native Go](#native-go)
- [Script languages](#script-languages)

### Native Go

You can import Gate as a Go module and use it like a framework
in your own external projects.

Go get it into your project
- `go get -u go.minekube.com/gate`


create a `proxy.New(...)`, register servers and event subscribers and start
`proxy.Listen(...)` for connections.

> TODO: code examples

### Script languages

To simplify and accelerate customization of Gate there
will be added support for scripting languages such as
[Tengo](https://github.com/d5/tengo) and
[Lua](https://github.com/yuin/gopher-lua).

> This feature will be added when highly requested.

## Anticipated future of Gate

- Gate will be a high performance & cloud native Minecraft proxy
ready for massive scale in the cloud!

- Players can always join and will never be kicked if there is
no available server to connect to, or the network is too full.
Instead, players will be moved to an empty virtual room simulated
by the proxy to queue players to wait.

- _Distant future, or maybe not too far?_ A proxy for Java & Bedrock edition to mix and match players & servers of all kinds.
(protocol translation back and forth...)