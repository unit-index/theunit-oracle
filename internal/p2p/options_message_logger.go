""

package p2p

import (
pubsub "github.com/libp2p/go-libp2p-pubsub"

"github.com/toknowwhy/theunit-oracle/pkg/log"
"github.com/toknowwhy/theunit-oracle/pkg/transport"
)

// MessageLogger logs published and received messages.
func MessageLogger() Options {
return func (n *Node) error {
mlh := &messageLoggerHandler{n: n}
n.AddMessageHandler(mlh)
return nil
}
}

type messageLoggerHandler struct {
n *Node
}

func (m *messageLoggerHandler) Published(topic string, raw []byte, _ transport.Message) {
m.n.tsLog.get().
WithFields(log.Fields{
"topic":   topic,
"message": string(raw),
}).
Debug("Published a new message")
}

func (m *messageLoggerHandler) Received(topic string, msg *pubsub.Message, _ pubsub.ValidationResult) {
m.n.tsLog.get().
WithFields(log.Fields{
"topic":              topic,
"message":            string(msg.Data),
"peerID":             msg.GetFrom().String(),
"receivedFromPeerID": msg.ReceivedFrom.String(),
}).
Debug("Received a new message")
}

func (m *messageLoggerHandler) Broken(topic string, msg *pubsub.Message, err error) {
m.n.tsLog.get().
WithError(err).
WithFields(log.Fields{
"topic":              topic,
"peerID":             msg.GetFrom().String(),
"receivedFromPeerID": msg.ReceivedFrom.String(),
}).
Debug("Unable to unmarshall received message")
}
