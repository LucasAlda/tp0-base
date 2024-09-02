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
	for {
		err := s.acceptNewConnection()

		if s.cancelled {
			log.Debug("action: cancel_server | result: success")
			return
		}

		if err != nil {
			log.Errorf("action: accept_connections | result: fail | error: %s", err)
			continue
		}

	}
}

// Read message from a specific client socket and closes the socket

// If a problem arises in the communication with the client, the
// client socket will also be closed
func (s *Server) handleClientConnection(clientSocket *net.TCPConn) {
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

func (s *Server) acceptNewConnection() error {
	log.Info("action: accept_connections | result: in_progress")

	clientSocket, err := s.serverSocket.AcceptTCP()
	if err != nil {
		return err
	}

	log.Infof("action: accept_connections | result: success | ip: %s", clientSocket.RemoteAddr())

	defer clientSocket.Close()
	s.handleClientConnection(clientSocket)
	return nil
}
