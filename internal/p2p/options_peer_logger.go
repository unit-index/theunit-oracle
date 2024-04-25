package p2p

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/toknowwhy/theunit-oracle/internal/p2p/sets"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
)

// PeerLogger logs all peers handled by libp2p's pubsub system.
func PeerLogger() Options {
	return func(n *Node) error {
		n.AddPubSubEventHandler(sets.PubSubEventHandlerFunc(func(topic string, event pubsub.PeerEvent) {
			addrs := n.Peerstore().PeerInfo(event.Peer).Addrs
			ua := getPeerUserAgent(n.Peerstore(), event.Peer)
			pp := getPeerProtocols(n.Peerstore(), event.Peer)
			pv := getPeerProtocolVersion(n.Peerstore(), event.Peer)

			switch event.Type {
			case pubsub.PeerJoin:
				n.tsLog.get().
					WithFields(log.Fields{
						"peerID":          event.Peer.String(),
						"topic":           topic,
						"listenAddrs":     log.Format(addrs),
						"userAgent":       ua,
						"protocolVersion": pv,
						"protocols":       log.Format(pp),
					}).
					Info("Connected to a peer")
			case pubsub.PeerLeave:
				n.tsLog.get().
					WithFields(log.Fields{
						"peerID":      event.Peer.String(),
						"topic":       topic,
						"listenAddrs": log.Format(addrs),
					}).
					Info("Disconnected from a peer")
			}
		}))
		return nil
	}
}

func getPeerProtocols(ps peerstore.Peerstore, pid peer.ID) []string {
	pp, _ := ps.GetProtocols(pid)
	return pp
}

func getPeerUserAgent(ps peerstore.Peerstore, pid peer.ID) string {
	av, _ := ps.Get(pid, "AgentVersion")
	if s, ok := av.(string); ok {
		return s
	}
	return ""
}

func getPeerProtocolVersion(ps peerstore.Peerstore, pid peer.ID) string {
	av, _ := ps.Get(pid, "ProtocolVersion")
	if s, ok := av.(string); ok {
		return s
	}
	return ""
}
