package transport

// ReceivedMessage contains a Message received from a Transport with
// an additional data.
type ReceivedMessage struct {
	// Message contains the message content. It is nil when the Error field
	// is not nil.
	Message Message
	// Data contains an optional data associated with the message. A type of
	// the data is different depending on a transport implementation.
	Data interface{}
	// Error contains an optional error returned by a transport.
	Error error
}

type Message interface {
	Marshall() ([]byte, error)
	Unmarshall([]byte) error
}

// Transport is the interface for different implementations of a
// publish–subscribe messaging solutions for the Oracle network.
type Transport interface {
	Broadcast(topic string, message Message) error
	// Messages returns a channel that will deliver incoming messages. Note,
	// that only messages for subscribed topics will be supported by this
	// method, for unsubscribed topic nil will be returned. In case of an
	// error, error will be returned in a ReceivedMessage structure.
	Messages(topic string) chan ReceivedMessage
	// Start starts listening for messages.
	Start() error
	// Wait waits until transport's context is cancelled.
	Wait()
}
