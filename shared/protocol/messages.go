package protocol

import (
	"strconv"
	"strings"
)

type MessageType int32

const (
	MessageTypePresentation MessageType = iota
	MessageTypeBetBatch
	MessageTypeBetAck
	MessageTypeAllBetsSent
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

func (m *MessageBet) GetSize() int {
	return len(m.FirstName) + len(m.LastName) + len(m.Document) + len(m.Birthdate) + len(m.Number) + 4
}

type MessageBetBatch struct {
	Bets []MessageBet
}

func (m *MessageBetBatch) GetMessageType() MessageType {
	return MessageTypeBetBatch
}

func (m *MessageBetBatch) Encode() string {
	encodedBets := make([]string, len(m.Bets))
	for i, bet := range m.Bets {
		encodedBets[i] = bet.Encode()
	}
	return strings.Join(encodedBets, "|")
}

func (m *MessageBetBatch) Decode(data string) error {
	values := strings.Split(data, "|")
	bets := make([]MessageBet, len(values))

	for i, value := range values {
		bet := MessageBet{}
		err := bet.Decode(value)
		if err != nil {
			return err
		}
		bets[i] = bet
	}

	m.Bets = bets
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

// MessageBetAck is a struct that represents a bet ack message
type MessageAllBetsSent struct {
}

func (m *MessageAllBetsSent) GetMessageType() MessageType {
	return MessageTypeAllBetsSent
}

func (m *MessageAllBetsSent) Encode() string {
	return ""
}

func (m *MessageAllBetsSent) Decode(data string) error {
	return nil
}

// MessageBetAck is a struct that represents a bet ack message
type MessagePresentation struct {
	Agency string
}

func (m *MessagePresentation) GetMessageType() MessageType {
	return MessageTypePresentation
}

func (m *MessagePresentation) Encode() string {
	return m.Agency
}

func (m *MessagePresentation) Decode(data string) error {
	m.Agency = data
	return nil
}
