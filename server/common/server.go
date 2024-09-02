package common

import (
	"bufio"
	"fmt"
	"net"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type Server struct {
	serverSocket *net.TCPListener
}

func NewServer(port int, listenBacklog int) (*Server, error) {
	serverSocket, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		return nil, err
	}

	return &Server{serverSocket: serverSocket}, nil
}

// Dummy Server loop

// Server that accept a new connections and establishes a
// communication with a client. After client with communucation
// finishes, servers starts to accept new connections again
func (s *Server) Run() {
	for {
		clientSocket, err := s.acceptNewConnection()
		if err != nil {
			log.Errorf("action: accept_connections | result: fail | error: %s", err)
			continue
		}
		s.handleClientConnection(clientSocket)
	}
}

// Read message from a specific client socket and closes the socket

// If a problem arises in the communication with the client, the
// client socket will also be closed
func (s *Server) handleClientConnection(clientSocket *net.TCPConn) {
	defer clientSocket.Close()

	msg, err := bufio.NewReader(clientSocket).ReadString('\n')
	if err != nil {
		log.Error("Error reading client message: %s", err)
		return
	}
	log.Infof("action: receive_message | result: success | ip: %s | msg: %s", clientSocket.RemoteAddr(), msg)

	// TODO: Modify the send to avoid short-write
	fmt.Fprintf(
		clientSocket,
		msg,
	)

	// clientSocket.Write(msg)
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
