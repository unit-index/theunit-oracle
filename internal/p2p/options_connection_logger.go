package p2p

import (
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/multiformats/go-multiaddr"

	"github.com/toknowwhy/theunit-oracle/pkg/log"
)

// ConnectionLogger logs connected and disconnected hosts,
func ConnectionLogger() Options {
	return func(n *Node) error {
		n.AddNotifee(&connectionLoggerNotifee{n: n})
		return nil
	}
}

type connectionLoggerNotifee struct {
	n *Node
}

// Listen implements the network.Notifiee interface.
func (n *connectionLoggerNotifee) Listen(network.Network, multiaddr.Multiaddr) {}

// ListenClose implements the network.Notifiee interface.
func (n *connectionLoggerNotifee) ListenClose(network.Network, multiaddr.Multiaddr) {}

// Connected implements the network.Notifiee interface.
func (n *connectionLoggerNotifee) Connected(_ network.Network, conn network.Conn) {
	n.n.tsLog.get().
		WithFields(log.Fields{
			"peerID": conn.RemotePeer().String(),
			"addr":   conn.RemoteMultiaddr().String(),
		}).
		Info("Connected to a host")
}

// Disconnected implements the network.Notifiee interface.
func (n *connectionLoggerNotifee) Disconnected(_ network.Network, conn network.Conn) {
	n.n.tsLog.get().
		WithFields(log.Fields{
			"peerID": conn.RemotePeer().String(),
			"addr":   conn.RemoteMultiaddr().String(),
		}).
		Info("Disconnected from a host")
}

// OpenedStream implements the network.Notifiee interface.
func (n *connectionLoggerNotifee) OpenedStream(network.Network, network.Stream) {}

// ClosedStream implements the network.Notifiee interface.
func (n *connectionLoggerNotifee) ClosedStream(network.Network, network.Stream) {}
