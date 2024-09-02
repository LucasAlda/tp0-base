package protocol

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendAndReceiveBetMessage(t *testing.T) {
	server, err := net.Listen("tcp", "localhost:8080")

	go func() {
		if err != nil {
			panic(err)
		}

		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}

		// time.Sleep(200 * time.Millisecond)
		received, err := Receive(conn)
		if err != nil {
			println("Error receiving message: ", err.Error())
			panic(err)
		}

		assert.Equal(t, received.MessageType, MessageTypeBet)
		bet := MessageBet{}
		bet.Decode(received.Data)

		println("Received bet message")
		println("BET: ", bet.FirstName, bet.LastName, bet.Document, bet.Birthdate, bet.Number)
		assert.Equal(t, bet.FirstName, "Juan")
		assert.Equal(t, bet.LastName, "Perez")
		assert.Equal(t, bet.Document, "1234567890")
		assert.Equal(t, bet.Birthdate, "2000-01-01")
		assert.Equal(t, bet.Number, "123456")
	}()

	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		panic(err)
	}

	msg := &MessageBet{
		FirstName: "Juan",
		LastName:  "Perez",
		Document:  "1234567890",
		Birthdate: "2000-01-01",
		Number:    "123456",
	}

	err = Send(conn, msg)
	if err != nil {
		panic(err)
	}
}
