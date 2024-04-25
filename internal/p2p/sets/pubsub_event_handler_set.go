""

package sets

import (
pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// PubSubEventHandlerFunc is a adapter for the PubSubEventHandler interface.
type PubSubEventHandlerFunc func (topic string, event pubsub.PeerEvent)

// Handle calls f(topic, event).
func (f PubSubEventHandlerFunc) Handle(topic string, event pubsub.PeerEvent) {
f(topic, event)
}

// PubSubEventHandler can ba implemented by type that supports handling the PubSub
// system events.
type PubSubEventHandler interface {
// Handle is called on a new event.
Handle(topic string, event pubsub.PeerEvent)
}

// PubSubEventHandlerSet stores multiple instances of the PubSubEventHandler interface.
type PubSubEventHandlerSet struct {
eventHandler []PubSubEventHandler
}

// NewPubSubEventHandlerSet creates new instance of the PubSubEventHandlerSet.
func NewPubSubEventHandlerSet() *PubSubEventHandlerSet {
return &PubSubEventHandlerSet{}
}

// Add adds new PubSubEventHandler to the set.
func (n *PubSubEventHandlerSet) Add(eventHandler ...PubSubEventHandler) {
n.eventHandler = append(n.eventHandler, eventHandler...)
}

// Handle invokes all registered handlers for given topic.
func (n *PubSubEventHandlerSet) Handle(topic string, event pubsub.PeerEvent) {
for _, eventHandler := range n.eventHandler {
eventHandler.Handle(topic, event)
}
}

var _ PubSubEventHandler = (*PubSubEventHandlerSet)(nil)
