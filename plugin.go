package trelay

type Plugin interface {
	// Used for logging
	Name() string

	// Executed when plugin is loaded by server. Returned object will be used by the server in case plugin wants to force unique instance of self per server
	OnLoad(s Server) Plugin

	// Executed when server starts
	OnServerStart(s Server)

	// Executed when server stops
	OnServerStop(s Server)

	// Executed when new session is opened (before server connection is opened and server->client loop starts)
	OnSessionOpen(s Session)

	// Executed when session is closed (after client->server loop stopped)
	OnSessionClose(s Session)

	// Executed when client sends a packet.
	// This method is not goroutine-safe, multiple sessions may invoke it simultaneously.
	// If this function returns `true`, handling will stop and other plugins won't see it, otherwise packet will get forwarded to each plugin and then to server.
	OnClientPacket(pid PacketType, packet *Packet, session Session) (handled bool)

	// Executed when server sends a packet.
	// This method is not goroutine-safe, multiple sessions may invoke it simultaneously.
	// If this function returns `true`, handling will stop and other plugins won't see it, otherwise packet will get forwarded to each plugin and then to client.
	OnServerPacket(pid PacketType, packet *Packet, session Session) (handled bool)
}
