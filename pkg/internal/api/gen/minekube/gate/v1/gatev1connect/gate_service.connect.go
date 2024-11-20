// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: minekube/gate/v1/gate_service.proto

package gatev1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "go.minekube.com/gate/pkg/internal/api/gen/minekube/gate/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// GateServiceName is the fully-qualified name of the GateService service.
	GateServiceName = "minekube.gate.v1.GateService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// GateServiceGetPlayerProcedure is the fully-qualified name of the GateService's GetPlayer RPC.
	GateServiceGetPlayerProcedure = "/minekube.gate.v1.GateService/GetPlayer"
	// GateServiceListPlayersProcedure is the fully-qualified name of the GateService's ListPlayers RPC.
	GateServiceListPlayersProcedure = "/minekube.gate.v1.GateService/ListPlayers"
	// GateServiceListServersProcedure is the fully-qualified name of the GateService's ListServers RPC.
	GateServiceListServersProcedure = "/minekube.gate.v1.GateService/ListServers"
	// GateServiceRegisterServerProcedure is the fully-qualified name of the GateService's
	// RegisterServer RPC.
	GateServiceRegisterServerProcedure = "/minekube.gate.v1.GateService/RegisterServer"
	// GateServiceUnregisterServerProcedure is the fully-qualified name of the GateService's
	// UnregisterServer RPC.
	GateServiceUnregisterServerProcedure = "/minekube.gate.v1.GateService/UnregisterServer"
	// GateServiceConnectPlayerProcedure is the fully-qualified name of the GateService's ConnectPlayer
	// RPC.
	GateServiceConnectPlayerProcedure = "/minekube.gate.v1.GateService/ConnectPlayer"
	// GateServiceDisconnectPlayerProcedure is the fully-qualified name of the GateService's
	// DisconnectPlayer RPC.
	GateServiceDisconnectPlayerProcedure = "/minekube.gate.v1.GateService/DisconnectPlayer"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	gateServiceServiceDescriptor                = v1.File_minekube_gate_v1_gate_service_proto.Services().ByName("GateService")
	gateServiceGetPlayerMethodDescriptor        = gateServiceServiceDescriptor.Methods().ByName("GetPlayer")
	gateServiceListPlayersMethodDescriptor      = gateServiceServiceDescriptor.Methods().ByName("ListPlayers")
	gateServiceListServersMethodDescriptor      = gateServiceServiceDescriptor.Methods().ByName("ListServers")
	gateServiceRegisterServerMethodDescriptor   = gateServiceServiceDescriptor.Methods().ByName("RegisterServer")
	gateServiceUnregisterServerMethodDescriptor = gateServiceServiceDescriptor.Methods().ByName("UnregisterServer")
	gateServiceConnectPlayerMethodDescriptor    = gateServiceServiceDescriptor.Methods().ByName("ConnectPlayer")
	gateServiceDisconnectPlayerMethodDescriptor = gateServiceServiceDescriptor.Methods().ByName("DisconnectPlayer")
)

// GateServiceClient is a client for the minekube.gate.v1.GateService service.
type GateServiceClient interface {
	// GetPlayer returns the player by the given id or username.
	// Returns NOT_FOUND if the player is not online.
	// Returns INVALID_ARGUMENT if neither id nor username is provided, or if the id format is invalid.
	GetPlayer(context.Context, *connect.Request[v1.GetPlayerRequest]) (*connect.Response[v1.GetPlayerResponse], error)
	// ListPlayers returns all online players.
	// If servers are specified in the request, only returns players on those servers.
	ListPlayers(context.Context, *connect.Request[v1.ListPlayersRequest]) (*connect.Response[v1.ListPlayersResponse], error)
	// ListServers returns all registered servers.
	ListServers(context.Context, *connect.Request[v1.ListServersRequest]) (*connect.Response[v1.ListServersResponse], error)
	// RegisterServer adds a server to the proxy.
	// Returns ALREADY_EXISTS if a server with the same name is already registered.
	// Returns INVALID_ARGUMENT if the server name or address is invalid.
	RegisterServer(context.Context, *connect.Request[v1.RegisterServerRequest]) (*connect.Response[v1.RegisterServerResponse], error)
	// UnregisterServer removes a server from the proxy.
	// Returns NOT_FOUND if no matching server is found.
	// Returns INVALID_ARGUMENT if neither name nor address is provided.
	UnregisterServer(context.Context, *connect.Request[v1.UnregisterServerRequest]) (*connect.Response[v1.UnregisterServerResponse], error)
	// ConnectPlayer connects a player to a specified server.
	// Returns NOT_FOUND if either the player or target server doesn't exist.
	// Returns FAILED_PRECONDITION if the connection attempt fails.
	ConnectPlayer(context.Context, *connect.Request[v1.ConnectPlayerRequest]) (*connect.Response[v1.ConnectPlayerResponse], error)
	// DisconnectPlayer disconnects a player from the proxy.
	// Returns NOT_FOUND if the player doesn't exist.
	// Returns INVALID_ARGUMENT if the reason text is malformed.
	DisconnectPlayer(context.Context, *connect.Request[v1.DisconnectPlayerRequest]) (*connect.Response[v1.DisconnectPlayerResponse], error)
}

// NewGateServiceClient constructs a client for the minekube.gate.v1.GateService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewGateServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) GateServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &gateServiceClient{
		getPlayer: connect.NewClient[v1.GetPlayerRequest, v1.GetPlayerResponse](
			httpClient,
			baseURL+GateServiceGetPlayerProcedure,
			connect.WithSchema(gateServiceGetPlayerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listPlayers: connect.NewClient[v1.ListPlayersRequest, v1.ListPlayersResponse](
			httpClient,
			baseURL+GateServiceListPlayersProcedure,
			connect.WithSchema(gateServiceListPlayersMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listServers: connect.NewClient[v1.ListServersRequest, v1.ListServersResponse](
			httpClient,
			baseURL+GateServiceListServersProcedure,
			connect.WithSchema(gateServiceListServersMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		registerServer: connect.NewClient[v1.RegisterServerRequest, v1.RegisterServerResponse](
			httpClient,
			baseURL+GateServiceRegisterServerProcedure,
			connect.WithSchema(gateServiceRegisterServerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		unregisterServer: connect.NewClient[v1.UnregisterServerRequest, v1.UnregisterServerResponse](
			httpClient,
			baseURL+GateServiceUnregisterServerProcedure,
			connect.WithSchema(gateServiceUnregisterServerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		connectPlayer: connect.NewClient[v1.ConnectPlayerRequest, v1.ConnectPlayerResponse](
			httpClient,
			baseURL+GateServiceConnectPlayerProcedure,
			connect.WithSchema(gateServiceConnectPlayerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		disconnectPlayer: connect.NewClient[v1.DisconnectPlayerRequest, v1.DisconnectPlayerResponse](
			httpClient,
			baseURL+GateServiceDisconnectPlayerProcedure,
			connect.WithSchema(gateServiceDisconnectPlayerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// gateServiceClient implements GateServiceClient.
type gateServiceClient struct {
	getPlayer        *connect.Client[v1.GetPlayerRequest, v1.GetPlayerResponse]
	listPlayers      *connect.Client[v1.ListPlayersRequest, v1.ListPlayersResponse]
	listServers      *connect.Client[v1.ListServersRequest, v1.ListServersResponse]
	registerServer   *connect.Client[v1.RegisterServerRequest, v1.RegisterServerResponse]
	unregisterServer *connect.Client[v1.UnregisterServerRequest, v1.UnregisterServerResponse]
	connectPlayer    *connect.Client[v1.ConnectPlayerRequest, v1.ConnectPlayerResponse]
	disconnectPlayer *connect.Client[v1.DisconnectPlayerRequest, v1.DisconnectPlayerResponse]
}

// GetPlayer calls minekube.gate.v1.GateService.GetPlayer.
func (c *gateServiceClient) GetPlayer(ctx context.Context, req *connect.Request[v1.GetPlayerRequest]) (*connect.Response[v1.GetPlayerResponse], error) {
	return c.getPlayer.CallUnary(ctx, req)
}

// ListPlayers calls minekube.gate.v1.GateService.ListPlayers.
func (c *gateServiceClient) ListPlayers(ctx context.Context, req *connect.Request[v1.ListPlayersRequest]) (*connect.Response[v1.ListPlayersResponse], error) {
	return c.listPlayers.CallUnary(ctx, req)
}

// ListServers calls minekube.gate.v1.GateService.ListServers.
func (c *gateServiceClient) ListServers(ctx context.Context, req *connect.Request[v1.ListServersRequest]) (*connect.Response[v1.ListServersResponse], error) {
	return c.listServers.CallUnary(ctx, req)
}

// RegisterServer calls minekube.gate.v1.GateService.RegisterServer.
func (c *gateServiceClient) RegisterServer(ctx context.Context, req *connect.Request[v1.RegisterServerRequest]) (*connect.Response[v1.RegisterServerResponse], error) {
	return c.registerServer.CallUnary(ctx, req)
}

// UnregisterServer calls minekube.gate.v1.GateService.UnregisterServer.
func (c *gateServiceClient) UnregisterServer(ctx context.Context, req *connect.Request[v1.UnregisterServerRequest]) (*connect.Response[v1.UnregisterServerResponse], error) {
	return c.unregisterServer.CallUnary(ctx, req)
}

// ConnectPlayer calls minekube.gate.v1.GateService.ConnectPlayer.
func (c *gateServiceClient) ConnectPlayer(ctx context.Context, req *connect.Request[v1.ConnectPlayerRequest]) (*connect.Response[v1.ConnectPlayerResponse], error) {
	return c.connectPlayer.CallUnary(ctx, req)
}

// DisconnectPlayer calls minekube.gate.v1.GateService.DisconnectPlayer.
func (c *gateServiceClient) DisconnectPlayer(ctx context.Context, req *connect.Request[v1.DisconnectPlayerRequest]) (*connect.Response[v1.DisconnectPlayerResponse], error) {
	return c.disconnectPlayer.CallUnary(ctx, req)
}

// GateServiceHandler is an implementation of the minekube.gate.v1.GateService service.
type GateServiceHandler interface {
	// GetPlayer returns the player by the given id or username.
	// Returns NOT_FOUND if the player is not online.
	// Returns INVALID_ARGUMENT if neither id nor username is provided, or if the id format is invalid.
	GetPlayer(context.Context, *connect.Request[v1.GetPlayerRequest]) (*connect.Response[v1.GetPlayerResponse], error)
	// ListPlayers returns all online players.
	// If servers are specified in the request, only returns players on those servers.
	ListPlayers(context.Context, *connect.Request[v1.ListPlayersRequest]) (*connect.Response[v1.ListPlayersResponse], error)
	// ListServers returns all registered servers.
	ListServers(context.Context, *connect.Request[v1.ListServersRequest]) (*connect.Response[v1.ListServersResponse], error)
	// RegisterServer adds a server to the proxy.
	// Returns ALREADY_EXISTS if a server with the same name is already registered.
	// Returns INVALID_ARGUMENT if the server name or address is invalid.
	RegisterServer(context.Context, *connect.Request[v1.RegisterServerRequest]) (*connect.Response[v1.RegisterServerResponse], error)
	// UnregisterServer removes a server from the proxy.
	// Returns NOT_FOUND if no matching server is found.
	// Returns INVALID_ARGUMENT if neither name nor address is provided.
	UnregisterServer(context.Context, *connect.Request[v1.UnregisterServerRequest]) (*connect.Response[v1.UnregisterServerResponse], error)
	// ConnectPlayer connects a player to a specified server.
	// Returns NOT_FOUND if either the player or target server doesn't exist.
	// Returns FAILED_PRECONDITION if the connection attempt fails.
	ConnectPlayer(context.Context, *connect.Request[v1.ConnectPlayerRequest]) (*connect.Response[v1.ConnectPlayerResponse], error)
	// DisconnectPlayer disconnects a player from the proxy.
	// Returns NOT_FOUND if the player doesn't exist.
	// Returns INVALID_ARGUMENT if the reason text is malformed.
	DisconnectPlayer(context.Context, *connect.Request[v1.DisconnectPlayerRequest]) (*connect.Response[v1.DisconnectPlayerResponse], error)
}

// NewGateServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewGateServiceHandler(svc GateServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	gateServiceGetPlayerHandler := connect.NewUnaryHandler(
		GateServiceGetPlayerProcedure,
		svc.GetPlayer,
		connect.WithSchema(gateServiceGetPlayerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gateServiceListPlayersHandler := connect.NewUnaryHandler(
		GateServiceListPlayersProcedure,
		svc.ListPlayers,
		connect.WithSchema(gateServiceListPlayersMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gateServiceListServersHandler := connect.NewUnaryHandler(
		GateServiceListServersProcedure,
		svc.ListServers,
		connect.WithSchema(gateServiceListServersMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gateServiceRegisterServerHandler := connect.NewUnaryHandler(
		GateServiceRegisterServerProcedure,
		svc.RegisterServer,
		connect.WithSchema(gateServiceRegisterServerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gateServiceUnregisterServerHandler := connect.NewUnaryHandler(
		GateServiceUnregisterServerProcedure,
		svc.UnregisterServer,
		connect.WithSchema(gateServiceUnregisterServerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gateServiceConnectPlayerHandler := connect.NewUnaryHandler(
		GateServiceConnectPlayerProcedure,
		svc.ConnectPlayer,
		connect.WithSchema(gateServiceConnectPlayerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	gateServiceDisconnectPlayerHandler := connect.NewUnaryHandler(
		GateServiceDisconnectPlayerProcedure,
		svc.DisconnectPlayer,
		connect.WithSchema(gateServiceDisconnectPlayerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/minekube.gate.v1.GateService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case GateServiceGetPlayerProcedure:
			gateServiceGetPlayerHandler.ServeHTTP(w, r)
		case GateServiceListPlayersProcedure:
			gateServiceListPlayersHandler.ServeHTTP(w, r)
		case GateServiceListServersProcedure:
			gateServiceListServersHandler.ServeHTTP(w, r)
		case GateServiceRegisterServerProcedure:
			gateServiceRegisterServerHandler.ServeHTTP(w, r)
		case GateServiceUnregisterServerProcedure:
			gateServiceUnregisterServerHandler.ServeHTTP(w, r)
		case GateServiceConnectPlayerProcedure:
			gateServiceConnectPlayerHandler.ServeHTTP(w, r)
		case GateServiceDisconnectPlayerProcedure:
			gateServiceDisconnectPlayerHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedGateServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedGateServiceHandler struct{}

func (UnimplementedGateServiceHandler) GetPlayer(context.Context, *connect.Request[v1.GetPlayerRequest]) (*connect.Response[v1.GetPlayerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("minekube.gate.v1.GateService.GetPlayer is not implemented"))
}

func (UnimplementedGateServiceHandler) ListPlayers(context.Context, *connect.Request[v1.ListPlayersRequest]) (*connect.Response[v1.ListPlayersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("minekube.gate.v1.GateService.ListPlayers is not implemented"))
}

func (UnimplementedGateServiceHandler) ListServers(context.Context, *connect.Request[v1.ListServersRequest]) (*connect.Response[v1.ListServersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("minekube.gate.v1.GateService.ListServers is not implemented"))
}

func (UnimplementedGateServiceHandler) RegisterServer(context.Context, *connect.Request[v1.RegisterServerRequest]) (*connect.Response[v1.RegisterServerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("minekube.gate.v1.GateService.RegisterServer is not implemented"))
}

func (UnimplementedGateServiceHandler) UnregisterServer(context.Context, *connect.Request[v1.UnregisterServerRequest]) (*connect.Response[v1.UnregisterServerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("minekube.gate.v1.GateService.UnregisterServer is not implemented"))
}

func (UnimplementedGateServiceHandler) ConnectPlayer(context.Context, *connect.Request[v1.ConnectPlayerRequest]) (*connect.Response[v1.ConnectPlayerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("minekube.gate.v1.GateService.ConnectPlayer is not implemented"))
}

func (UnimplementedGateServiceHandler) DisconnectPlayer(context.Context, *connect.Request[v1.DisconnectPlayerRequest]) (*connect.Response[v1.DisconnectPlayerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("minekube.gate.v1.GateService.DisconnectPlayer is not implemented"))
}