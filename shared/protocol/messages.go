package protocol

import (
	"strconv"
	"strings"
)

type MessageType int32

const (
	MessageTypeBet MessageType = iota
	MessageTypeBetAck
)

// Protocolo de comunicacion entre cliente y servidor
// Interfaz para los mensajes que se intercambian
type Message interface {
	GetMessageType() MessageType
	Encode() string
	Decode(data string) error
}

// MessageBet is a struct that represents a bet message
type MessageBet struct {
	FirstName string
	LastName  string
	Document  string
	Birthdate string
	Number    string
}

func (m *MessageBet) GetMessageType() MessageType {
	return MessageTypeBet
}

func (m *MessageBet) Encode() string {
	return m.FirstName + "," + m.LastName + "," + m.Document + "," + m.Birthdate + "," + m.Number
}

func (m *MessageBet) Decode(data string) error {
	values := strings.Split(data, ",")
	m.FirstName = values[0]
	m.LastName = values[1]
	m.Document = values[2]
	m.Birthdate = values[3]
	m.Number = values[4]

	return nil
}

// MessageBetAck is a struct that represents a bet ack message
type MessageBetAck struct {
	Result bool
}

func (m *MessageBetAck) GetMessageType() MessageType {
	return MessageTypeBetAck
}

func (m *MessageBetAck) Encode() string {
	return strconv.FormatBool(m.Result)
}

func (m *MessageBetAck) Decode(data string) error {
	result, err := strconv.ParseBool(data)
	if err != nil {
		return err
	}
	m.Result = result
	return nil
}
