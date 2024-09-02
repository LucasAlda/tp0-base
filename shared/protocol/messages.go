package protocol

import "strings"

type MessageType int32

const (
	MessageTypeBet MessageType = iota
)

// Protocolo de comunicacion entre cliente y servidor
// Interfaz para los mensajes que se intercambian
type Message interface {
	GetSize() int32
	GetMessageType() MessageType
	Encode() string
	Decode(data string)
}

// MessageBet is a struct that represents a bet message
type MessageBet struct {
	FirstName string
	LastName  string
	Document  string
	Birthdate string
	Number    string
}

func (m *MessageBet) GetSize() int32 {
	return int32(len(m.FirstName) + len(m.LastName) + len(m.Document) + len(m.Birthdate) + len(m.Number))
}

func (m *MessageBet) GetMessageType() MessageType {
	return MessageTypeBet
}

func (m *MessageBet) Encode() string {
	return m.FirstName + "," + m.LastName + "," + m.Document + "," + m.Birthdate + "," + m.Number
}

func (m *MessageBet) Decode(data string) {
	values := strings.Split(data, ",")
	m.FirstName = values[0]
	m.LastName = values[1]
	m.Document = values[2]
	m.Birthdate = values[3]
	m.Number = values[4]
}
