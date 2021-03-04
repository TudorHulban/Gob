package seri

// Message Structure used for serialization.
type Message struct {
	ID      int64
	Payload string
}

type ISerialization interface {
	Encode(Message) ([]byte, error)
	Decode(payload []byte) (Message, error)
}
