package ipc

type MessageType int

const (
	READ MessageType = iota
	WRITE
	RESULT
)

type Message struct {
	Type MessageType
	Key string
	Payload string
}