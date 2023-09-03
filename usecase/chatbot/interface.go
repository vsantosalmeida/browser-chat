package chatbot

import "context"

// Reader interface to read messages from a queue.
type Reader interface {
	ReadMessage(ctx context.Context, msgCH chan<- []byte)
}

// Writer interface to send messages to a queue.
type Writer interface {
	WriteMessage(ctx context.Context, payload []byte) error
}

// Broker handle the read/write message flow.
type Broker interface {
	Reader
	Writer
}

// UseCase represents the chatbot starter.
type UseCase interface {
	Start(ctx context.Context)
}
