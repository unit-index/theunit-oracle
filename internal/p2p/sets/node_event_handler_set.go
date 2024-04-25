package sets

type NodeConfiguredEvent struct{}
type NodeStartingEvent struct{}
type NodeHostStartedEvent struct{}
type NodePubSubStartedEvent struct{}
type NodeStartedEvent struct{}
type NodeTopicSubscribedEvent struct{ Topic string }
type NodeTopicUnsubscribedEvent struct{ Topic string }
type NodeStoppingEvent struct{}
type NodeStoppedEvent struct{}

// NodeEventHandlerFunc is a adapter for the NodeEventHandler interface.
type NodeEventHandlerFunc func(event interface{})

// Handle calls f(topic, event).
func (f NodeEventHandlerFunc) Handle(event interface{}) {
	f(event)
}

// NodeEventHandler can ba implemented by type that supports handling the Node
// system events.
type NodeEventHandler interface {
	// Handle is called on a new event.
	Handle(event interface{})
}

// NodeEventHandlerSet stores multiple instances of the NodeEventHandler interface.
type NodeEventHandlerSet struct {
	eventHandler []NodeEventHandler
}

// NewNodeEventHandlerSet creates new instance of the NodeEventHandlerSet.
func NewNodeEventHandlerSet() *NodeEventHandlerSet {
	return &NodeEventHandlerSet{}
}

// Add adds new NodeEventHandler to the set.
func (n *NodeEventHandlerSet) Add(eventHandler ...NodeEventHandler) {
	n.eventHandler = append(n.eventHandler, eventHandler...)
}

// Handle invokes all registered handlers for given topic.
func (n *NodeEventHandlerSet) Handle(event interface{}) {
	for _, eventHandler := range n.eventHandler {
		eventHandler.Handle(event)
	}
}

var _ NodeEventHandler = (*NodeEventHandlerSet)(nil)
