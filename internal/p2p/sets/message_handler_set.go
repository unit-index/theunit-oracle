package sets

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/toknowwhy/theunit-oracle/pkg/transport"
)

// MessageHandler can ba implemented by type that supports handling the PubSub
// system messages.
type MessageHandler interface {
	// Published is called when new message is published.
	Published(topic string, raw []byte, msg transport.Message)
	// Received is called when new message is received.
	Received(topic string, msg *pubsub.Message, result pubsub.ValidationResult)
	// Broken is called when it is impossible to unmarshall message,
	Broken(topic string, msg *pubsub.Message, err error)
}

// MessageHandlerSet stores multiple instances of the MessageHandler interface.
type MessageHandlerSet struct {
	messageHandler []MessageHandler
}

// NewMessageHandlerSet creates new instance of the MessageHandlerSet.
func NewMessageHandlerSet() *MessageHandlerSet {
	return &MessageHandlerSet{}
}

// Add adds new MessageHandler to the set.
func (n *MessageHandlerSet) Add(messageHandler ...MessageHandler) {
	n.messageHandler = append(n.messageHandler, messageHandler...)
}

// Published invokes all registered handlers.
func (n *MessageHandlerSet) Published(topic string, raw []byte, msg transport.Message) {
	for _, messageHandler := range n.messageHandler {
		messageHandler.Published(topic, raw, msg)
	}
}

// Received invokes all registered handlers.
func (n *MessageHandlerSet) Received(topic string, msg *pubsub.Message, result pubsub.ValidationResult) {
	for _, messageHandler := range n.messageHandler {
		messageHandler.Received(topic, msg, result)
	}
}

// Broken invokes all registered handlers.
func (n *MessageHandlerSet) Broken(topic string, msg *pubsub.Message, err error) {
	for _, messageHandler := range n.messageHandler {
		messageHandler.Broken(topic, msg, err)
	}
}

var _ MessageHandler = (*MessageHandlerSet)(nil)
