package trelay

// TServer represents a terraria server. However, there is no underlying connection
// and it is not possible to send packets to TServer directly, use TPlayer.SendPacketTS
type TServer interface {
}
