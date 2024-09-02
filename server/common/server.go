package common

import (
	"errors"
	"net"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/shared/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type Server struct {
	serverSocket *net.TCPListener
	cancelled    bool
}

func NewServer(port int, listenBacklog int) (*Server, error) {
	serverSocket, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		return nil, err
	}

	return &Server{serverSocket: serverSocket, cancelled: false}, nil
}

func (s *Server) Close() {
	s.cancelled = true
	s.serverSocket.Close()
}

// Dummy Server loop

// Server that accept a new connections and establishes a
// communication with a client. After client with communucation
// finishes, servers starts to accept new connections again
func (s *Server) Run() {
	agencyId := 1
	for {
		conn, err := s.acceptNewConnection()

		if s.cancelled {
			log.Debug("action: cancel_server | result: success")
			return
		}

		if err != nil {
			log.Errorf("action: accept_connections | result: fail | error: %s", err)
			continue
		}

		s.handleNewBet(conn, agencyId)
		agencyId++
	}
}

// Read message from a specific client socket and closes the socket

// If a problem arises in the communication with the client, the
// client socket will also be closed
func (s *Server) handleNewBet(clientSocket *net.TCPConn, agencyId int) {
	msg, err := protocol.Receive(clientSocket)
	if err != nil {
		handleFailedBet(clientSocket, protocol.MessageBet{}, err)
		return
	}

	if msg.MessageType != protocol.MessageTypeBet {
		handleFailedBet(clientSocket, protocol.MessageBet{}, errors.New("Mensaje recibido no es una apuesta"))
		return
	}

	betMsg := protocol.MessageBet{}
	err = betMsg.Decode(msg.Data)
	if err != nil {
		handleFailedBet(clientSocket, protocol.MessageBet{}, errors.New("Error al decodificar la apuesta"))
		return
	}

	bet, err := NewBet(strconv.Itoa(agencyId), betMsg.FirstName, betMsg.LastName, betMsg.Document, betMsg.Birthdate, betMsg.Number)
	if err != nil {
		handleFailedBet(clientSocket, betMsg, err)
		return
	}

	err = StoreBets([]*Bet{bet})
	if err != nil {
		handleFailedBet(clientSocket, betMsg, err)
		return
	}

	log.Infof("action: apuesta_almacenada | result: success | dni: %s | numero: %s", betMsg.Document, betMsg.Number)

	betAck := protocol.MessageBetAck{Result: true}
	protocol.Send(clientSocket, &betAck)
}

func (s *Server) acceptNewConnection() (*net.TCPConn, error) {
	log.Info("action: accept_connections | result: in_progress")

	clientSocket, err := s.serverSocket.AcceptTCP()
	if err != nil {
		return nil, err
	}

	log.Infof("action: accept_connections | result: success | ip: %s", clientSocket.RemoteAddr())

	return clientSocket, nil
}

func handleFailedBet(clientSocket *net.TCPConn, bet protocol.MessageBet, err error) {
	log.Errorf("action: apuesta_almacenada | result: fail | ip: %s | error: %s", clientSocket.RemoteAddr(), err)
	betAck := protocol.MessageBetAck{Result: false}
	protocol.Send(clientSocket, &betAck)
}
