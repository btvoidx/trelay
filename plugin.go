package trelay

type Plugin interface {
	// Used for logging.
	Name() string

	// Executed when server starts.
	OnServerStart()

	// Executed when server stops.
	OnServerStop()

	// Executed when session is opened (before server connection is opened and server->client loop starts).
	OnSessionOpen(*Session)

	// Executed when session is closed (after client->server loop is stopped).
	OnSessionClose(*Session)

	// Executed when client sends a packet.
	// Call `ctx.SetHandled()` to prevent forwarding this packet to server.
	//
	// This method is not goroutine-safe, multiple sessions may invoke it simultaneously.
	OnClientPacket(ctx *PacketContext)

	// Executed when server sends a packet.
	// Call `ctx.SetHandled()` to prevent forwarding this packet to client.
	//
	// This method is not goroutine-safe, multiple sessions may invoke it simultaneously.
	OnServerPacket(ctx *PacketContext)
}
