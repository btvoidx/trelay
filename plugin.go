package trelay

type Plugin interface {
	// Used for logging.
	Name() string

	// Executed when server starts.
	OnServerStart()

	// Executed when server stops.
	OnServerStop()

	// Executed when new session is opened (before server connection is opened and server->client loop starts).
	OnSessionOpen(s Session)

	// Executed when session is closed (after client->server loop stopped).
	OnSessionClose(s Session)

	// Executed when client sends a packet.
	// If this function returns `true`, handling will stop and other plugins won't see it, otherwise packet will get forwarded to each plugin and then to server.
	//
	// This method is not goroutine-safe, multiple sessions may invoke it simultaneously.
	OnClientPacket(pid PacketType, packet *Packet, session Session) (handled bool)

	// Executed when server sends a packet.
	// If this function returns `true`, handling will stop and other plugins won't see it, otherwise packet will get forwarded to each plugin and then to client.
	//
	// This method is not goroutine-safe, multiple sessions may invoke it simultaneously.
	OnServerPacket(pid PacketType, packet *Packet, session Session) (handled bool)
}
