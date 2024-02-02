package hwconfig

// The MessageSender interface allows sending off any struct that can become a specific type of message
type MessageSender[T any] interface {
	Send(m Msger[T]) error
}

type BridgeSender[T any] interface {
	Send(m T) error
}

// A Msger is any struct that can become a specific type of message
type Msger[T any] interface {
	ToBridgeMsg() T
}
