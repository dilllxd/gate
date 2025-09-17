---
title: "Gate Protocol Documentation - Generated API Reference"
description: "Generated protocol documentation for Gate Minecraft proxy API. Complete reference for all API endpoints, messages, and data structures."
---

# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [minekube/gate/v1/gate_service.proto](#minekube_gate_v1_gate_service-proto)
    - [APIConfig](#minekube-gate-v1-APIConfig)
    - [AddLiteRouteBackendRequest](#minekube-gate-v1-AddLiteRouteBackendRequest)
    - [AddLiteRouteBackendResponse](#minekube-gate-v1-AddLiteRouteBackendResponse)
    - [ApplyConfigRequest](#minekube-gate-v1-ApplyConfigRequest)
    - [ApplyConfigResponse](#minekube-gate-v1-ApplyConfigResponse)
    - [BedrockPlayerData](#minekube-gate-v1-BedrockPlayerData)
    - [ClassicStats](#minekube-gate-v1-ClassicStats)
    - [ConfigUpdateEvent](#minekube-gate-v1-ConfigUpdateEvent)
    - [ConnectConfig](#minekube-gate-v1-ConnectConfig)
    - [ConnectPlayerRequest](#minekube-gate-v1-ConnectPlayerRequest)
    - [ConnectPlayerResponse](#minekube-gate-v1-ConnectPlayerResponse)
    - [DisconnectPlayerRequest](#minekube-gate-v1-DisconnectPlayerRequest)
    - [DisconnectPlayerResponse](#minekube-gate-v1-DisconnectPlayerResponse)
    - [ForwardingConfig](#minekube-gate-v1-ForwardingConfig)
    - [GateConfig](#minekube-gate-v1-GateConfig)
    - [GetConfigRequest](#minekube-gate-v1-GetConfigRequest)
    - [GetConfigResponse](#minekube-gate-v1-GetConfigResponse)
    - [GetLiteRouteRequest](#minekube-gate-v1-GetLiteRouteRequest)
    - [GetLiteRouteResponse](#minekube-gate-v1-GetLiteRouteResponse)
    - [GetPlayerRequest](#minekube-gate-v1-GetPlayerRequest)
    - [GetPlayerResponse](#minekube-gate-v1-GetPlayerResponse)
    - [GetStatusRequest](#minekube-gate-v1-GetStatusRequest)
    - [GetStatusResponse](#minekube-gate-v1-GetStatusResponse)
    - [JavaConfig](#minekube-gate-v1-JavaConfig)
    - [JavaConfig.ServersEntry](#minekube-gate-v1-JavaConfig-ServersEntry)
    - [ListLiteRoutesRequest](#minekube-gate-v1-ListLiteRoutesRequest)
    - [ListLiteRoutesResponse](#minekube-gate-v1-ListLiteRoutesResponse)
    - [ListPlayersRequest](#minekube-gate-v1-ListPlayersRequest)
    - [ListPlayersResponse](#minekube-gate-v1-ListPlayersResponse)
    - [ListServersRequest](#minekube-gate-v1-ListServersRequest)
    - [ListServersResponse](#minekube-gate-v1-ListServersResponse)
    - [LiteConfig](#minekube-gate-v1-LiteConfig)
    - [LiteRoute](#minekube-gate-v1-LiteRoute)
    - [LiteRouteBackend](#minekube-gate-v1-LiteRouteBackend)
    - [LiteRouteFallback](#minekube-gate-v1-LiteRouteFallback)
    - [LiteRouteFallbackPlayers](#minekube-gate-v1-LiteRouteFallbackPlayers)
    - [LiteRouteFallbackVersion](#minekube-gate-v1-LiteRouteFallbackVersion)
    - [LiteRouteOptions](#minekube-gate-v1-LiteRouteOptions)
    - [LiteRouteUpdateEvent](#minekube-gate-v1-LiteRouteUpdateEvent)
    - [LiteStats](#minekube-gate-v1-LiteStats)
    - [Player](#minekube-gate-v1-Player)
    - [PlayerConnectEvent](#minekube-gate-v1-PlayerConnectEvent)
    - [PlayerDisconnectEvent](#minekube-gate-v1-PlayerDisconnectEvent)
    - [PlayerServerSwitchEvent](#minekube-gate-v1-PlayerServerSwitchEvent)
    - [ProxyEvent](#minekube-gate-v1-ProxyEvent)
    - [ReadyEvent](#minekube-gate-v1-ReadyEvent)
    - [RegisterServerRequest](#minekube-gate-v1-RegisterServerRequest)
    - [RegisterServerResponse](#minekube-gate-v1-RegisterServerResponse)
    - [RemoveLiteRouteBackendRequest](#minekube-gate-v1-RemoveLiteRouteBackendRequest)
    - [RemoveLiteRouteBackendResponse](#minekube-gate-v1-RemoveLiteRouteBackendResponse)
    - [RequestCookieRequest](#minekube-gate-v1-RequestCookieRequest)
    - [RequestCookieResponse](#minekube-gate-v1-RequestCookieResponse)
    - [Server](#minekube-gate-v1-Server)
    - [ServerRegisterEvent](#minekube-gate-v1-ServerRegisterEvent)
    - [ServerUnregisterEvent](#minekube-gate-v1-ServerUnregisterEvent)
    - [ShutdownEvent](#minekube-gate-v1-ShutdownEvent)
    - [StatusConfig](#minekube-gate-v1-StatusConfig)
    - [StoreCookieRequest](#minekube-gate-v1-StoreCookieRequest)
    - [StoreCookieResponse](#minekube-gate-v1-StoreCookieResponse)
    - [StreamEventsRequest](#minekube-gate-v1-StreamEventsRequest)
    - [UnregisterServerRequest](#minekube-gate-v1-UnregisterServerRequest)
    - [UnregisterServerResponse](#minekube-gate-v1-UnregisterServerResponse)
    - [UpdateLiteRouteFallbackRequest](#minekube-gate-v1-UpdateLiteRouteFallbackRequest)
    - [UpdateLiteRouteFallbackResponse](#minekube-gate-v1-UpdateLiteRouteFallbackResponse)
    - [UpdateLiteRouteOptionsRequest](#minekube-gate-v1-UpdateLiteRouteOptionsRequest)
    - [UpdateLiteRouteOptionsResponse](#minekube-gate-v1-UpdateLiteRouteOptionsResponse)
    - [UpdateLiteRouteStrategyRequest](#minekube-gate-v1-UpdateLiteRouteStrategyRequest)
    - [UpdateLiteRouteStrategyResponse](#minekube-gate-v1-UpdateLiteRouteStrategyResponse)
    - [ValidateConfigRequest](#minekube-gate-v1-ValidateConfigRequest)
    - [ValidateConfigResponse](#minekube-gate-v1-ValidateConfigResponse)
  
    - [BedrockDeviceOS](#minekube-gate-v1-BedrockDeviceOS)
    - [BedrockInputMode](#minekube-gate-v1-BedrockInputMode)
    - [BedrockUIProfile](#minekube-gate-v1-BedrockUIProfile)
    - [EventType](#minekube-gate-v1-EventType)
    - [LiteRouteStrategy](#minekube-gate-v1-LiteRouteStrategy)
    - [LoginStatus](#minekube-gate-v1-LoginStatus)
    - [ProxyMode](#minekube-gate-v1-ProxyMode)
    - [RouteOperation](#minekube-gate-v1-RouteOperation)
  
    - [GateEventsService](#minekube-gate-v1-GateEventsService)
    - [GateLiteService](#minekube-gate-v1-GateLiteService)
    - [GateService](#minekube-gate-v1-GateService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="minekube_gate_v1_gate_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## minekube/gate/v1/gate_service.proto



<a name="minekube-gate-v1-APIConfig"></a>

### APIConfig
APIConfig represents the Gate API configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| enabled | [bool](#bool) |  | Whether the API is enabled |
| bind | [string](#string) |  | The address to bind the API server to (using a localhost address is recommended) |






<a name="minekube-gate-v1-AddLiteRouteBackendRequest"></a>

### AddLiteRouteBackendRequest
AddLiteRouteBackendRequest adds a backend to a route.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host | [string](#string) |  | Host matcher to update (case-insensitive) |
| backend | [string](#string) |  | Backend address to add (e.g., &#34;localhost:25565&#34;) |






<a name="minekube-gate-v1-AddLiteRouteBackendResponse"></a>

### AddLiteRouteBackendResponse
AddLiteRouteBackendResponse contains validation warnings.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| warnings | [string](#string) | repeated |  |






<a name="minekube-gate-v1-ApplyConfigRequest"></a>

### ApplyConfigRequest
ApplyConfigRequest is the request for ApplyConfig method.
The config payload is parsed with a YAML decoder (which supports JSON as YAML is a superset).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| config | [string](#string) |  | Configuration data as YAML or JSON string |
| persist | [bool](#bool) |  | Whether to persist the config to disk by overwriting the existing config file. Only works if a config file exists. Defaults to false (in-memory only). |






<a name="minekube-gate-v1-ApplyConfigResponse"></a>

### ApplyConfigResponse
ApplyConfigResponse contains validation warnings emitted while applying the config.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| warnings | [string](#string) | repeated |  |






<a name="minekube-gate-v1-BedrockPlayerData"></a>

### BedrockPlayerData
BedrockPlayerData contains information specific to Bedrock Edition players.
This data is only available for players connecting through Geyser/Floodgate.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| xuid | [int64](#int64) |  | Xbox User ID (XUID) - unique identifier for Bedrock players |
| device_os | [BedrockDeviceOS](#minekube-gate-v1-BedrockDeviceOS) |  | Device operating system the player is using |
| language | [string](#string) |  | Client language code (e.g., &#34;en_US&#34;) |
| ui_profile | [BedrockUIProfile](#minekube-gate-v1-BedrockUIProfile) |  | UI profile (Classic or Pocket) |
| input_mode | [BedrockInputMode](#minekube-gate-v1-BedrockInputMode) |  | Input method (mouse, touch, gamepad, etc.) |
| behind_proxy | [bool](#bool) |  | Whether the player is connecting through a proxy |
| linked_player | [string](#string) |  | Linked Java Edition username (if any) |






<a name="minekube-gate-v1-ClassicStats"></a>

### ClassicStats
ClassicStats contains statistics for classic proxy mode.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| players | [int32](#int32) |  | Number of online players |
| servers | [int32](#int32) |  | Number of registered servers |






<a name="minekube-gate-v1-ConfigUpdateEvent"></a>

### ConfigUpdateEvent
ConfigUpdateEvent fires when proxy configuration is updated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| config_source | [string](#string) |  | config_source indicates how the config was updated. |
| lite_mode_enabled | [bool](#bool) |  | lite_mode_enabled indicates if lite mode is enabled in the new config. |
| route_count | [int32](#int32) |  | route_count is the number of lite routes configured (lite mode only). |






<a name="minekube-gate-v1-ConnectConfig"></a>

### ConnectConfig
ConnectConfig represents the Connect network configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| enabled | [bool](#bool) |  | Whether to connect Gate to the WatchService |
| name | [string](#string) |  | Endpoint name |
| allow_offline_mode_players | [bool](#bool) |  | Allow offline mode players to join |






<a name="minekube-gate-v1-ConnectPlayerRequest"></a>

### ConnectPlayerRequest
ConnectPlayerRequest is the request for ConnectPlayer method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player | [string](#string) |  | The player&#39;s username or ID to connect |
| server | [string](#string) |  | The target server name to connect the player to |






<a name="minekube-gate-v1-ConnectPlayerResponse"></a>

### ConnectPlayerResponse
ConnectPlayerResponse is the response for ConnectPlayer method.






<a name="minekube-gate-v1-DisconnectPlayerRequest"></a>

### DisconnectPlayerRequest
DisconnectPlayerRequest is the request for DisconnectPlayer method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player | [string](#string) |  | The player&#39;s username or ID to disconnect |
| reason | [string](#string) |  | The reason displayed to the player when they are disconnected.

Formats:

- `{&#34;text&#34;:&#34;Hello, world!&#34;}` - JSON text component. See https://wiki.vg/Text_formatting for details.

- `§aHello,\n§bworld!` - Simple color codes. See https://wiki.vg/Text_formatting#Colors

Optional, if empty no reason will be shown. |






<a name="minekube-gate-v1-DisconnectPlayerResponse"></a>

### DisconnectPlayerResponse
DisconnectPlayerResponse is the response for DisconnectPlayer method.






<a name="minekube-gate-v1-ForwardingConfig"></a>

### ForwardingConfig
ForwardingConfig represents player info forwarding settings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| mode | [string](#string) |  | Forwarding mode (&#34;none&#34;, &#34;legacy&#34;, &#34;velocity&#34;, &#34;bungeeguard&#34;) |
| velocity_secret | [string](#string) |  | Secret used with &#34;velocity&#34; mode |
| bungee_guard_secret | [string](#string) |  | Secret used with &#34;bungeeguard&#34; mode |






<a name="minekube-gate-v1-GateConfig"></a>

### GateConfig
GateConfig represents the root configuration structure


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| api | [APIConfig](#minekube-gate-v1-APIConfig) |  |  |
| connect | [ConnectConfig](#minekube-gate-v1-ConnectConfig) |  |  |
| config | [JavaConfig](#minekube-gate-v1-JavaConfig) |  |  |






<a name="minekube-gate-v1-GetConfigRequest"></a>

### GetConfigRequest
GetConfigRequest is the request for GetConfig method.






<a name="minekube-gate-v1-GetConfigResponse"></a>

### GetConfigResponse
GetConfigResponse contains the serialized config payload.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| payload | [string](#string) |  | YAML-serialized configuration data |






<a name="minekube-gate-v1-GetLiteRouteRequest"></a>

### GetLiteRouteRequest
GetLiteRouteRequest is the request for GetLiteRoute method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host | [string](#string) |  | Host matcher to look up (case-insensitive). |






<a name="minekube-gate-v1-GetLiteRouteResponse"></a>

### GetLiteRouteResponse
GetLiteRouteResponse is the response for GetLiteRoute method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| route | [LiteRoute](#minekube-gate-v1-LiteRoute) |  |  |






<a name="minekube-gate-v1-GetPlayerRequest"></a>

### GetPlayerRequest
GetPlayerRequest is the request for GetPlayer method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | Gets the player by their Minecraft UUID. Optional, if not set the username will be used. If both id and username are set, the id will be used. Must be a valid Minecraft UUID format (e.g. &#34;550e8400-e29b-41d4-a716-446655440000&#34;) |
| username | [string](#string) |  | Gets the player by their username. Optional, if not set the id will be used. Case-sensitive. |






<a name="minekube-gate-v1-GetPlayerResponse"></a>

### GetPlayerResponse
GetPlayerResponse is the response for GetPlayer method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player | [Player](#minekube-gate-v1-Player) |  | The player matching the request criteria |






<a name="minekube-gate-v1-GetStatusRequest"></a>

### GetStatusRequest
GetStatusRequest is the request for GetStatus method.






<a name="minekube-gate-v1-GetStatusResponse"></a>

### GetStatusResponse
GetStatusResponse contains proxy runtime metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [string](#string) |  | Gate version string |
| mode | [ProxyMode](#minekube-gate-v1-ProxyMode) |  | Current operating mode (classic or lite) |
| classic | [ClassicStats](#minekube-gate-v1-ClassicStats) |  | Statistics for classic mode |
| lite | [LiteStats](#minekube-gate-v1-LiteStats) |  | Statistics for lite mode |






<a name="minekube-gate-v1-JavaConfig"></a>

### JavaConfig
JavaConfig represents the main Java edition configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bind | [string](#string) |  | The address to listen for connections |
| online_mode | [bool](#bool) |  | Whether to enable online mode |
| forwarding | [ForwardingConfig](#minekube-gate-v1-ForwardingConfig) |  | Player info forwarding settings |
| status | [StatusConfig](#minekube-gate-v1-StatusConfig) |  | Status response settings |
| servers | [JavaConfig.ServersEntry](#minekube-gate-v1-JavaConfig-ServersEntry) | repeated | Registered servers (name:address) |
| try | [string](#string) | repeated | Try server names order |
| forced_hosts_json | [string](#string) |  | Note: forced_hosts is represented as JSON string due to protobuf limitations with map&lt;string, []string&gt; |
| accept_transfers | [bool](#bool) |  | Whether to accept transfers from other hosts via transfer packet |
| bungee_plugin_channel_enabled | [bool](#bool) |  | Whether to enable BungeeCord plugin messaging |
| lite | [LiteConfig](#minekube-gate-v1-LiteConfig) |  | Lite mode settings |






<a name="minekube-gate-v1-JavaConfig-ServersEntry"></a>

### JavaConfig.ServersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="minekube-gate-v1-ListLiteRoutesRequest"></a>

### ListLiteRoutesRequest
ListLiteRoutesRequest is the request for ListLiteRoutes method.






<a name="minekube-gate-v1-ListLiteRoutesResponse"></a>

### ListLiteRoutesResponse
ListLiteRoutesResponse is the response for ListLiteRoutes method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| routes | [LiteRoute](#minekube-gate-v1-LiteRoute) | repeated |  |






<a name="minekube-gate-v1-ListPlayersRequest"></a>

### ListPlayersRequest
ListPlayersRequest is the request for ListPlayers method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| servers | [string](#string) | repeated | Filter players by server names. Optional, if empty all online players are returned. If specified, only returns players on the listed servers. |






<a name="minekube-gate-v1-ListPlayersResponse"></a>

### ListPlayersResponse
ListPlayersResponse is the response for ListPlayers method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| players | [Player](#minekube-gate-v1-Player) | repeated |  |






<a name="minekube-gate-v1-ListServersRequest"></a>

### ListServersRequest
ListServersRequest is the request for ListServers method.






<a name="minekube-gate-v1-ListServersResponse"></a>

### ListServersResponse
ListServersResponse is the response for ListServers method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| servers | [Server](#minekube-gate-v1-Server) | repeated |  |






<a name="minekube-gate-v1-LiteConfig"></a>

### LiteConfig
LiteConfig represents Gate Lite mode configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| enabled | [bool](#bool) |  | Whether Lite mode is enabled |
| routes | [LiteRoute](#minekube-gate-v1-LiteRoute) | repeated | Configured lite routes |






<a name="minekube-gate-v1-LiteRoute"></a>

### LiteRoute
LiteRoute represents a configured lite route and runtime state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hosts | [string](#string) | repeated | Host patterns this route matches |
| backends | [LiteRouteBackend](#minekube-gate-v1-LiteRouteBackend) | repeated | Backend servers for this route |
| strategy | [LiteRouteStrategy](#minekube-gate-v1-LiteRouteStrategy) |  | Load balancing strategy |
| options | [LiteRouteOptions](#minekube-gate-v1-LiteRouteOptions) |  | Proxy behavior options |
| fallback | [LiteRouteFallback](#minekube-gate-v1-LiteRouteFallback) |  | Fallback response when all backends fail |






<a name="minekube-gate-v1-LiteRouteBackend"></a>

### LiteRouteBackend
LiteRouteBackend represents a backend target for a lite route.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| address | [string](#string) |  | Network address of the backend server |
| active_connections | [uint32](#uint32) |  | Number of active connections to this backend |






<a name="minekube-gate-v1-LiteRouteFallback"></a>

### LiteRouteFallback
LiteRouteFallback contains fallback response data served when all backends fail.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| motd_json | [string](#string) |  | JSON-formatted MOTD for fallback response |
| version | [LiteRouteFallbackVersion](#minekube-gate-v1-LiteRouteFallbackVersion) |  | Version information for fallback response |
| players | [LiteRouteFallbackPlayers](#minekube-gate-v1-LiteRouteFallbackPlayers) |  | Player count information for fallback response |
| favicon | [string](#string) |  | Base64-encoded favicon for fallback response |






<a name="minekube-gate-v1-LiteRouteFallbackPlayers"></a>

### LiteRouteFallbackPlayers
LiteRouteFallbackPlayers contains fallback player counts.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| online | [int32](#int32) |  | Number of online players to display |
| max | [int32](#int32) |  | Maximum players to display |






<a name="minekube-gate-v1-LiteRouteFallbackVersion"></a>

### LiteRouteFallbackVersion
LiteRouteFallbackVersion contains display version metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Version name displayed to clients |
| protocol | [int32](#int32) |  | Protocol version number |






<a name="minekube-gate-v1-LiteRouteOptions"></a>

### LiteRouteOptions
LiteRouteOptions captures proxy behaviour flags for a lite route.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| proxy_protocol | [bool](#bool) |  | Whether to enable HA-Proxy protocol for this route |
| tcp_shield_real_ip | [bool](#bool) |  | Whether to enable TCPShield real IP support |
| modify_virtual_host | [bool](#bool) |  | Whether to modify the virtual host header |
| cache_ping_ttl_ms | [int64](#int64) |  | Cache TTL for ping responses in milliseconds |






<a name="minekube-gate-v1-LiteRouteUpdateEvent"></a>

### LiteRouteUpdateEvent
LiteRouteUpdateEvent fires when lite mode routes are modified (lite mode only).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host | [string](#string) |  | host is the hostname pattern of the updated route. |
| operation | [RouteOperation](#minekube-gate-v1-RouteOperation) |  | operation describes what operation was performed on the route. |
| backend_count | [int32](#int32) |  | backend_count is the current number of backends for this route. |






<a name="minekube-gate-v1-LiteStats"></a>

### LiteStats
LiteStats contains statistics for lite proxy mode.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| connections | [int32](#int32) |  | Number of active connections being proxied |
| routes | [int32](#int32) |  | Number of configured routes |






<a name="minekube-gate-v1-Player"></a>

### Player
Player represents an online player on the proxy.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The player&#39;s Minecraft UUID |
| username | [string](#string) |  | The player&#39;s username |
| bedrock | [BedrockPlayerData](#minekube-gate-v1-BedrockPlayerData) |  | Optional Bedrock player data (only present for Bedrock players) |






<a name="minekube-gate-v1-PlayerConnectEvent"></a>

### PlayerConnectEvent
PlayerConnectEvent fires when a player connects to the proxy.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player_id | [string](#string) |  | player_id is the UUID of the connecting player. |
| player_username | [string](#string) |  | player_username is the username of the connecting player. |
| remote_address | [string](#string) |  | remote_address is the IP address the player connected from. |
| protocol_version | [int32](#int32) |  | protocol_version is the Minecraft protocol version the player is using. |






<a name="minekube-gate-v1-PlayerDisconnectEvent"></a>

### PlayerDisconnectEvent
PlayerDisconnectEvent fires when a player disconnects from the proxy.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player_id | [string](#string) |  | player_id is the UUID of the disconnecting player. |
| player_username | [string](#string) |  | player_username is the username of the disconnecting player. |
| reason | [string](#string) |  | reason is the disconnect reason if available. |
| login_status | [LoginStatus](#minekube-gate-v1-LoginStatus) |  | login_status indicates the login state when the player disconnected. |






<a name="minekube-gate-v1-PlayerServerSwitchEvent"></a>

### PlayerServerSwitchEvent
PlayerServerSwitchEvent fires when a player switches between servers.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player_id | [string](#string) |  | player_id is the UUID of the player switching servers. |
| player_username | [string](#string) |  | player_username is the username of the player switching servers. |
| from_server | [string](#string) |  | from_server is the name of the server the player is leaving (empty if initial connect). |
| to_server | [string](#string) |  | to_server is the name of the server the player is connecting to. |






<a name="minekube-gate-v1-ProxyEvent"></a>

### ProxyEvent
ProxyEvent represents a single event occurrence with full details.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event_type | [EventType](#minekube-gate-v1-EventType) |  | event_type indicates the type of event that occurred. |
| timestamp_ms | [int64](#int64) |  | timestamp_ms is the Unix timestamp in milliseconds when the event occurred. |
| player_connect | [PlayerConnectEvent](#minekube-gate-v1-PlayerConnectEvent) |  |  |
| player_disconnect | [PlayerDisconnectEvent](#minekube-gate-v1-PlayerDisconnectEvent) |  |  |
| player_server_switch | [PlayerServerSwitchEvent](#minekube-gate-v1-PlayerServerSwitchEvent) |  |  |
| config_update | [ConfigUpdateEvent](#minekube-gate-v1-ConfigUpdateEvent) |  |  |
| server_register | [ServerRegisterEvent](#minekube-gate-v1-ServerRegisterEvent) |  |  |
| server_unregister | [ServerUnregisterEvent](#minekube-gate-v1-ServerUnregisterEvent) |  |  |
| lite_route_update | [LiteRouteUpdateEvent](#minekube-gate-v1-LiteRouteUpdateEvent) |  |  |
| ready | [ReadyEvent](#minekube-gate-v1-ReadyEvent) |  |  |
| shutdown | [ShutdownEvent](#minekube-gate-v1-ShutdownEvent) |  |  |






<a name="minekube-gate-v1-ReadyEvent"></a>

### ReadyEvent
ReadyEvent fires when the proxy is ready to accept connections.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bind_address | [string](#string) |  | bind_address is the address the proxy is listening on. |
| proxy_mode | [ProxyMode](#minekube-gate-v1-ProxyMode) |  | proxy_mode indicates whether the proxy is running in classic or lite mode. |






<a name="minekube-gate-v1-RegisterServerRequest"></a>

### RegisterServerRequest
RegisterServerRequest is the request for RegisterServer method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The unique name of the server |
| address | [string](#string) |  | The network address of the server (e.g. &#34;localhost:25565&#34;) |






<a name="minekube-gate-v1-RegisterServerResponse"></a>

### RegisterServerResponse
RegisterServerResponse is the response for RegisterServer method.






<a name="minekube-gate-v1-RemoveLiteRouteBackendRequest"></a>

### RemoveLiteRouteBackendRequest
RemoveLiteRouteBackendRequest removes a backend from a route.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host | [string](#string) |  | Host matcher to update (case-insensitive) |
| backend | [string](#string) |  | Backend address to remove |






<a name="minekube-gate-v1-RemoveLiteRouteBackendResponse"></a>

### RemoveLiteRouteBackendResponse
RemoveLiteRouteBackendResponse contains validation warnings.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| warnings | [string](#string) | repeated |  |






<a name="minekube-gate-v1-RequestCookieRequest"></a>

### RequestCookieRequest
RequestCookieRequest is the request for RequestCookie method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player | [string](#string) |  | The player&#39;s username or ID |
| key | [string](#string) |  | The key of the cookie in format `namespace:key` |






<a name="minekube-gate-v1-RequestCookieResponse"></a>

### RequestCookieResponse
RequestCookieResponse is the response for RequestCookie method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| payload | [bytes](#bytes) |  | The payload of the cookie. May be empty if the cookie is not found. |






<a name="minekube-gate-v1-Server"></a>

### Server
Server represents a backend server where Gate can connect players to.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The unique name of the server. |
| address | [string](#string) |  | The network address of the server. |
| players | [int32](#int32) |  | The number of players currently on the server. |






<a name="minekube-gate-v1-ServerRegisterEvent"></a>

### ServerRegisterEvent
ServerRegisterEvent fires when a new server is registered with the proxy.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| server_name | [string](#string) |  | server_name is the name of the registered server. |
| server_address | [string](#string) |  | server_address is the address of the registered server. |






<a name="minekube-gate-v1-ServerUnregisterEvent"></a>

### ServerUnregisterEvent
ServerUnregisterEvent fires when a server is unregistered from the proxy.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| server_name | [string](#string) |  | server_name is the name of the unregistered server. |
| server_address | [string](#string) |  | server_address is the address of the unregistered server. |






<a name="minekube-gate-v1-ShutdownEvent"></a>

### ShutdownEvent
ShutdownEvent fires when the proxy begins shutting down.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reason | [string](#string) |  | reason is the shutdown reason if available. |






<a name="minekube-gate-v1-StatusConfig"></a>

### StatusConfig
StatusConfig represents status response settings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| motd | [string](#string) |  | Message of the Day displayed in server list |
| show_max_players | [int32](#int32) |  | Maximum players shown in server list |
| favicon | [string](#string) |  | Base64-encoded favicon image |
| log_ping_requests | [bool](#bool) |  | Whether to log ping requests |
| announce_forge | [bool](#bool) |  | Whether the proxy should present itself as Forge/FML-compatible |






<a name="minekube-gate-v1-StoreCookieRequest"></a>

### StoreCookieRequest
StoreCookieRequest is the request for StoreCookie method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player | [string](#string) |  | The player&#39;s username or ID |
| key | [string](#string) |  | The key of the cookie in format `namespace:key` |
| payload | [bytes](#bytes) |  | The payload to store. Passing an empty payload will remove the cookie. |






<a name="minekube-gate-v1-StoreCookieResponse"></a>

### StoreCookieResponse
StoreCookieResponse is the response for StoreCookie method.






<a name="minekube-gate-v1-StreamEventsRequest"></a>

### StreamEventsRequest
StreamEventsRequest configures the event stream subscription.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event_types | [EventType](#minekube-gate-v1-EventType) | repeated | event_types filters which event types to receive. Empty means all events. |
| include_system_events | [bool](#bool) |  | include_system_events determines if system events (config, server changes, ready, shutdown) are included. Default: true |
| include_player_events | [bool](#bool) |  | include_player_events determines if player events (connect, disconnect, switch) are included. Default: true |






<a name="minekube-gate-v1-UnregisterServerRequest"></a>

### UnregisterServerRequest
UnregisterServerRequest is the request for UnregisterServer method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the server. Optional, if not set, the address will be used to match servers. |
| address | [string](#string) |  | The address of the server. Optional, if not set, the name will be used to match servers. If both name and address are set, only the server that matches both properties exactly will be unregistered. If only the address is set, the first server matching that address will be unregistered. |






<a name="minekube-gate-v1-UnregisterServerResponse"></a>

### UnregisterServerResponse
UnregisterServerResponse is the response for UnregisterServer method.






<a name="minekube-gate-v1-UpdateLiteRouteFallbackRequest"></a>

### UpdateLiteRouteFallbackRequest
UpdateLiteRouteFallbackRequest updates fallback metadata using a field mask.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host | [string](#string) |  | Host matcher to update (case-insensitive) |
| fallback | [LiteRouteFallback](#minekube-gate-v1-LiteRouteFallback) |  | New fallback configuration to apply |
| update_mask | [google.protobuf.FieldMask](#google-protobuf-FieldMask) |  | Field mask specifying which fallback fields to update |






<a name="minekube-gate-v1-UpdateLiteRouteFallbackResponse"></a>

### UpdateLiteRouteFallbackResponse
UpdateLiteRouteFallbackResponse contains validation warnings.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| warnings | [string](#string) | repeated |  |






<a name="minekube-gate-v1-UpdateLiteRouteOptionsRequest"></a>

### UpdateLiteRouteOptionsRequest
UpdateLiteRouteOptionsRequest updates per-route options using a field mask.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host | [string](#string) |  | Host matcher to update (case-insensitive) |
| options | [LiteRouteOptions](#minekube-gate-v1-LiteRouteOptions) |  | New options to apply |
| update_mask | [google.protobuf.FieldMask](#google-protobuf-FieldMask) |  | Field mask specifying which options to update |






<a name="minekube-gate-v1-UpdateLiteRouteOptionsResponse"></a>

### UpdateLiteRouteOptionsResponse
UpdateLiteRouteOptionsResponse contains validation warnings.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| warnings | [string](#string) | repeated |  |






<a name="minekube-gate-v1-UpdateLiteRouteStrategyRequest"></a>

### UpdateLiteRouteStrategyRequest
UpdateLiteRouteStrategyRequest updates the load-balancing strategy for a route.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host | [string](#string) |  | Host matcher to update (case-insensitive) |
| strategy | [LiteRouteStrategy](#minekube-gate-v1-LiteRouteStrategy) |  | New load balancing strategy to apply |






<a name="minekube-gate-v1-UpdateLiteRouteStrategyResponse"></a>

### UpdateLiteRouteStrategyResponse
UpdateLiteRouteStrategyResponse contains validation warnings.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| warnings | [string](#string) | repeated |  |






<a name="minekube-gate-v1-ValidateConfigRequest"></a>

### ValidateConfigRequest
ValidateConfigRequest is the request for ValidateConfig method.
The config payload is parsed with a YAML decoder (which supports JSON as YAML is a superset).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| config | [string](#string) |  | Configuration data as YAML or JSON string |






<a name="minekube-gate-v1-ValidateConfigResponse"></a>

### ValidateConfigResponse
ValidateConfigResponse contains validation results when the config is processed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| warnings | [string](#string) | repeated |  |
| errors | [string](#string) | repeated |  |





 


<a name="minekube-gate-v1-BedrockDeviceOS"></a>

### BedrockDeviceOS
BedrockDeviceOS represents the operating system of a Bedrock Edition player&#39;s device.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BEDROCK_DEVICE_OS_UNSPECIFIED | 0 |  |
| BEDROCK_DEVICE_OS_UNKNOWN | 1 |  |
| BEDROCK_DEVICE_OS_ANDROID | 2 |  |
| BEDROCK_DEVICE_OS_IOS | 3 |  |
| BEDROCK_DEVICE_OS_MACOS | 4 |  |
| BEDROCK_DEVICE_OS_AMAZON | 5 |  |
| BEDROCK_DEVICE_OS_GEAR_VR | 6 |  |
| BEDROCK_DEVICE_OS_HOLOLENS | 7 | Deprecated |
| BEDROCK_DEVICE_OS_WINDOWS_UWP | 8 |  |
| BEDROCK_DEVICE_OS_WINDOWS_X86 | 9 |  |
| BEDROCK_DEVICE_OS_DEDICATED | 10 |  |
| BEDROCK_DEVICE_OS_APPLE_TV | 11 | Deprecated |
| BEDROCK_DEVICE_OS_PLAYSTATION | 12 |  |
| BEDROCK_DEVICE_OS_SWITCH | 13 |  |
| BEDROCK_DEVICE_OS_XBOX | 14 |  |
| BEDROCK_DEVICE_OS_WINDOWS_PHONE | 15 | Deprecated |
| BEDROCK_DEVICE_OS_LINUX | 16 |  |



<a name="minekube-gate-v1-BedrockInputMode"></a>

### BedrockInputMode
BedrockInputMode represents the input method used by a Bedrock Edition player.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BEDROCK_INPUT_MODE_UNSPECIFIED | 0 |  |
| BEDROCK_INPUT_MODE_UNKNOWN | 1 |  |
| BEDROCK_INPUT_MODE_MOUSE | 2 |  |
| BEDROCK_INPUT_MODE_TOUCH | 3 |  |
| BEDROCK_INPUT_MODE_GAMEPAD | 4 |  |
| BEDROCK_INPUT_MODE_MOTION_CONTROLLER | 5 |  |



<a name="minekube-gate-v1-BedrockUIProfile"></a>

### BedrockUIProfile
BedrockUIProfile represents the UI profile used by a Bedrock Edition player.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BEDROCK_UI_PROFILE_UNSPECIFIED | 0 |  |
| BEDROCK_UI_PROFILE_CLASSIC | 1 |  |
| BEDROCK_UI_PROFILE_POCKET | 2 |  |



<a name="minekube-gate-v1-EventType"></a>

### EventType
EventType represents the different types of events that can be streamed.

| Name | Number | Description |
| ---- | ------ | ----------- |
| EVENT_TYPE_UNSPECIFIED | 0 | EVENT_TYPE_UNSPECIFIED is the default value, not used for filtering. |
| EVENT_TYPE_PLAYER_CONNECT | 1 | EVENT_TYPE_PLAYER_CONNECT fires when a player connects to the proxy. |
| EVENT_TYPE_PLAYER_DISCONNECT | 2 | EVENT_TYPE_PLAYER_DISCONNECT fires when a player disconnects from the proxy. |
| EVENT_TYPE_PLAYER_SERVER_SWITCH | 3 | EVENT_TYPE_PLAYER_SERVER_SWITCH fires when a player switches servers. |
| EVENT_TYPE_CONFIG_UPDATE | 4 | EVENT_TYPE_CONFIG_UPDATE fires when proxy configuration is updated. |
| EVENT_TYPE_SERVER_REGISTER | 5 | EVENT_TYPE_SERVER_REGISTER fires when a new server is registered. |
| EVENT_TYPE_SERVER_UNREGISTER | 6 | EVENT_TYPE_SERVER_UNREGISTER fires when a server is unregistered. |
| EVENT_TYPE_LITE_ROUTE_UPDATE | 7 | EVENT_TYPE_LITE_ROUTE_UPDATE fires when lite mode routes are modified (lite mode only). |
| EVENT_TYPE_READY | 8 | EVENT_TYPE_READY fires when the proxy is ready to accept connections. |
| EVENT_TYPE_SHUTDOWN | 9 | EVENT_TYPE_SHUTDOWN fires when the proxy begins shutting down. |



<a name="minekube-gate-v1-LiteRouteStrategy"></a>

### LiteRouteStrategy
LiteRouteStrategy represents load balancing strategies for lite routes.

| Name | Number | Description |
| ---- | ------ | ----------- |
| LITE_ROUTE_STRATEGY_UNSPECIFIED | 0 |  |
| LITE_ROUTE_STRATEGY_SEQUENTIAL | 1 |  |
| LITE_ROUTE_STRATEGY_RANDOM | 2 |  |
| LITE_ROUTE_STRATEGY_ROUND_ROBIN | 3 |  |
| LITE_ROUTE_STRATEGY_LEAST_CONNECTIONS | 4 |  |
| LITE_ROUTE_STRATEGY_LOWEST_LATENCY | 5 |  |



<a name="minekube-gate-v1-LoginStatus"></a>

### LoginStatus
LoginStatus represents the login state when a player disconnected.

| Name | Number | Description |
| ---- | ------ | ----------- |
| LOGIN_STATUS_UNSPECIFIED | 0 | LOGIN_STATUS_UNSPECIFIED is the default value. |
| LOGIN_STATUS_SUCCESSFUL | 1 | LOGIN_STATUS_SUCCESSFUL indicates the player was successfully logged in. |
| LOGIN_STATUS_CONFLICTING | 2 | LOGIN_STATUS_CONFLICTING indicates there was a login conflict. |
| LOGIN_STATUS_CANCELED_BY_USER | 3 | LOGIN_STATUS_CANCELED_BY_USER indicates the user canceled the login. |
| LOGIN_STATUS_CANCELED_BY_PROXY | 4 | LOGIN_STATUS_CANCELED_BY_PROXY indicates the proxy canceled the login. |
| LOGIN_STATUS_CANCELED_BEFORE_COMPLETE | 5 | LOGIN_STATUS_CANCELED_BEFORE_COMPLETE indicates the login was canceled before completion. |



<a name="minekube-gate-v1-ProxyMode"></a>

### ProxyMode
ProxyMode enumerates the current operating mode of Gate.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PROXY_MODE_UNSPECIFIED | 0 |  |
| PROXY_MODE_CLASSIC | 1 |  |
| PROXY_MODE_LITE | 2 |  |



<a name="minekube-gate-v1-RouteOperation"></a>

### RouteOperation
RouteOperation describes operations performed on lite routes.

| Name | Number | Description |
| ---- | ------ | ----------- |
| ROUTE_OPERATION_UNSPECIFIED | 0 | ROUTE_OPERATION_UNSPECIFIED is the default value. |
| ROUTE_OPERATION_BACKEND_ADDED | 1 | ROUTE_OPERATION_BACKEND_ADDED indicates a backend was added to the route. |
| ROUTE_OPERATION_BACKEND_REMOVED | 2 | ROUTE_OPERATION_BACKEND_REMOVED indicates a backend was removed from the route. |
| ROUTE_OPERATION_STRATEGY_UPDATED | 3 | ROUTE_OPERATION_STRATEGY_UPDATED indicates the load balancing strategy was changed. |
| ROUTE_OPERATION_OPTIONS_UPDATED | 4 | ROUTE_OPERATION_OPTIONS_UPDATED indicates route options were modified. |
| ROUTE_OPERATION_FALLBACK_UPDATED | 5 | ROUTE_OPERATION_FALLBACK_UPDATED indicates fallback settings were modified. |


 

 


<a name="minekube-gate-v1-GateEventsService"></a>

### GateEventsService
GateEventsService provides real-time event streaming for monitoring proxy activity.
This service streams events as they occur, allowing external tools to monitor
proxy state changes, player activities, and system events in real-time.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StreamEvents | [StreamEventsRequest](#minekube-gate-v1-StreamEventsRequest) | [ProxyEvent](#minekube-gate-v1-ProxyEvent) stream | StreamEvents subscribes to real-time proxy events with optional filtering. The stream remains open until the client disconnects or an error occurs. This is a server streaming RPC where the client sends one request and receives multiple events. |


<a name="minekube-gate-v1-GateLiteService"></a>

### GateLiteService
GateLiteService provides APIs for managing Gate Lite mode routes and backends.
This service is only available when Gate is running in Lite mode.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListLiteRoutes | [ListLiteRoutesRequest](#minekube-gate-v1-ListLiteRoutesRequest) | [ListLiteRoutesResponse](#minekube-gate-v1-ListLiteRoutesResponse) | ListLiteRoutes returns lite routes and their active connection counters. |
| GetLiteRoute | [GetLiteRouteRequest](#minekube-gate-v1-GetLiteRouteRequest) | [GetLiteRouteResponse](#minekube-gate-v1-GetLiteRouteResponse) | GetLiteRoute returns detailed information about a single lite route. |
| UpdateLiteRouteStrategy | [UpdateLiteRouteStrategyRequest](#minekube-gate-v1-UpdateLiteRouteStrategyRequest) | [UpdateLiteRouteStrategyResponse](#minekube-gate-v1-UpdateLiteRouteStrategyResponse) | UpdateLiteRouteStrategy updates the load-balancing strategy for a lite route. |
| AddLiteRouteBackend | [AddLiteRouteBackendRequest](#minekube-gate-v1-AddLiteRouteBackendRequest) | [AddLiteRouteBackendResponse](#minekube-gate-v1-AddLiteRouteBackendResponse) | AddLiteRouteBackend adds a backend target to a lite route. |
| RemoveLiteRouteBackend | [RemoveLiteRouteBackendRequest](#minekube-gate-v1-RemoveLiteRouteBackendRequest) | [RemoveLiteRouteBackendResponse](#minekube-gate-v1-RemoveLiteRouteBackendResponse) | RemoveLiteRouteBackend removes a backend target from a lite route. |
| UpdateLiteRouteOptions | [UpdateLiteRouteOptionsRequest](#minekube-gate-v1-UpdateLiteRouteOptionsRequest) | [UpdateLiteRouteOptionsResponse](#minekube-gate-v1-UpdateLiteRouteOptionsResponse) | UpdateLiteRouteOptions updates proxy options for a lite route using a field mask. |
| UpdateLiteRouteFallback | [UpdateLiteRouteFallbackRequest](#minekube-gate-v1-UpdateLiteRouteFallbackRequest) | [UpdateLiteRouteFallbackResponse](#minekube-gate-v1-UpdateLiteRouteFallbackResponse) | UpdateLiteRouteFallback updates fallback metadata for a lite route using a field mask. |


<a name="minekube-gate-v1-GateService"></a>

### GateService
GateService is the service API for managing a Gate proxy instance.
It provides methods for managing players and servers.
All methods follow standard gRPC error codes and include detailed error messages.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetPlayer | [GetPlayerRequest](#minekube-gate-v1-GetPlayerRequest) | [GetPlayerResponse](#minekube-gate-v1-GetPlayerResponse) | GetPlayer returns the player by the given id or username. Returns NOT_FOUND if the player is not online. Returns INVALID_ARGUMENT if neither id nor username is provided, or if the id format is invalid. |
| ListPlayers | [ListPlayersRequest](#minekube-gate-v1-ListPlayersRequest) | [ListPlayersResponse](#minekube-gate-v1-ListPlayersResponse) | ListPlayers returns all online players. If servers are specified in the request, only returns players on those servers. |
| ListServers | [ListServersRequest](#minekube-gate-v1-ListServersRequest) | [ListServersResponse](#minekube-gate-v1-ListServersResponse) | ListServers returns all registered servers. |
| RegisterServer | [RegisterServerRequest](#minekube-gate-v1-RegisterServerRequest) | [RegisterServerResponse](#minekube-gate-v1-RegisterServerResponse) | RegisterServer adds a server to the proxy. Returns ALREADY_EXISTS if a server with the same name is already registered. Returns INVALID_ARGUMENT if the server name or address is invalid. |
| UnregisterServer | [UnregisterServerRequest](#minekube-gate-v1-UnregisterServerRequest) | [UnregisterServerResponse](#minekube-gate-v1-UnregisterServerResponse) | UnregisterServer removes a server from the proxy. Returns NOT_FOUND if no matching server is found. Returns INVALID_ARGUMENT if neither name nor address is provided. |
| ConnectPlayer | [ConnectPlayerRequest](#minekube-gate-v1-ConnectPlayerRequest) | [ConnectPlayerResponse](#minekube-gate-v1-ConnectPlayerResponse) | ConnectPlayer connects a player to a specified server. Returns NOT_FOUND if either the player or target server doesn&#39;t exist. Returns FAILED_PRECONDITION if the connection attempt fails. |
| DisconnectPlayer | [DisconnectPlayerRequest](#minekube-gate-v1-DisconnectPlayerRequest) | [DisconnectPlayerResponse](#minekube-gate-v1-DisconnectPlayerResponse) | DisconnectPlayer disconnects a player from the proxy. Returns NOT_FOUND if the player doesn&#39;t exist. Returns INVALID_ARGUMENT if the reason text is malformed. |
| StoreCookie | [StoreCookieRequest](#minekube-gate-v1-StoreCookieRequest) | [StoreCookieResponse](#minekube-gate-v1-StoreCookieResponse) | StoreCookie stores a cookie on a player&#39;s client. Returns NOT_FOUND if the player doesn&#39;t exist. Passing an empty payload will remove the cookie. |
| RequestCookie | [RequestCookieRequest](#minekube-gate-v1-RequestCookieRequest) | [RequestCookieResponse](#minekube-gate-v1-RequestCookieResponse) | RequestCookie requests a cookie from a player&#39;s client. The payload in RequestCookieResponse may be empty if the cookie is not found. |
| GetStatus | [GetStatusRequest](#minekube-gate-v1-GetStatusRequest) | [GetStatusResponse](#minekube-gate-v1-GetStatusResponse) | GetStatus returns current proxy metadata including version, mode, players and servers. |
| GetConfig | [GetConfigRequest](#minekube-gate-v1-GetConfigRequest) | [GetConfigResponse](#minekube-gate-v1-GetConfigResponse) | GetConfig returns the current effective config. |
| ValidateConfig | [ValidateConfigRequest](#minekube-gate-v1-ValidateConfigRequest) | [ValidateConfigResponse](#minekube-gate-v1-ValidateConfigResponse) | ValidateConfig parses and validates a config payload without applying it. |
| ApplyConfig | [ApplyConfigRequest](#minekube-gate-v1-ApplyConfigRequest) | [ApplyConfigResponse](#minekube-gate-v1-ApplyConfigResponse) | ApplyConfig parses, validates, and applies a new config payload. |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

