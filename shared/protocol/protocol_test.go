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

		assert.Equal(t, received.MessageType, MessageTypeBetBatch)
		bet := MessageBetBatch{}
		bet.Decode(received.Data)

		println("Received bet message")
		println("BET: ", bet.Bets[0].FirstName, bet.Bets[0].LastName, bet.Bets[0].Document, bet.Bets[0].Birthdate, bet.Bets[0].Number)
		assert.Equal(t, len(bet.Bets), 1)
		assert.Equal(t, bet.Bets[0].FirstName, "Juan")
		assert.Equal(t, bet.Bets[0].LastName, "Perez")
		assert.Equal(t, bet.Bets[0].Document, "1234567890")
		assert.Equal(t, bet.Bets[0].Birthdate, "2000-01-01")
		assert.Equal(t, bet.Bets[0].Number, "123456")
	}()

	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		panic(err)
	}

	msg := MessageBet{
		FirstName: "Juan",
		LastName:  "Perez",
		Document:  "1234567890",
		Birthdate: "2000-01-01",
		Number:    "123456",
	}

	batch := MessageBetBatch{
		Bets: []MessageBet{msg},
	}

	err = Send(conn, &batch)
	if err != nil {
		panic(err)
	}
}
