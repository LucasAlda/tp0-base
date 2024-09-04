package common

import (
	"errors"
	"io"
	"net"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/shared/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type Server struct {
	serverSocket *net.TCPListener
	agencies     []*net.TCPConn
}

func NewServer(port int, listenBacklog int) (*Server, error) {
	serverSocket, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		return nil, err
	}

	return &Server{serverSocket: serverSocket}, nil
}

func (s *Server) Close() {
	s.serverSocket.Close()
	for _, agency := range s.agencies {
		if agency != nil {
			agency.Close()
		}
	}
}

// Dummy Server loop

// Server that accept a new connections and establishes a
// communication with a client. After client with communucation
// finishes, servers starts to accept new connections again
func (s *Server) Run() {
	agencyId := 1
	for {
		conn, err := s.acceptNewConnection()
		s.agencies = append(s.agencies, conn)
		if err != nil {
			return
		}

		err = s.handleConnection(conn, agencyId)
		if err != nil {
			log.Errorf("action: handle_connection | result: fail | error: %s", err)
		}
		agencyId++
	}
}

func (s *Server) acceptNewConnection() (*net.TCPConn, error) {
	log.Info("action: accept_connections | result: in_progress")

	clientSocket, err := s.serverSocket.AcceptTCP()
	// Si el error es que el socket ya está cerrado, simplemente terminamos el programa
	if errors.Is(err, net.ErrClosed) {
		return nil, err
	}
	// Si ocurre otro error, lo registramos
	if err != nil {
		log.Errorf("action: accept_connections | result: fail | error: %s", err)
		return nil, err
	}

	log.Infof("action: accept_connections | result: success | ip: %s", clientSocket.RemoteAddr())

	return clientSocket, nil
}

func (s *Server) handleConnection(conn *net.TCPConn, agencyId int) error {
	defer conn.Close()

	for {
		msg, err := protocol.Receive(conn)
		// Si el reader devuelve EOF, el cliente se desconectó
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			log.Errorf("action: client_disconected | ip: %s", conn.RemoteAddr())
			return nil
		}
		// Si ocurre otro error, lo registramos
		if err != nil {
			handleFailedBetBatch(conn, protocol.MessageBetBatch{}, err)
			return err
		}

		switch msg.MessageType {

		case protocol.MessageTypeBetBatch:
			s.handleNewBets(conn, agencyId, msg)

		case protocol.MessageTypeAllBetsSent:
			log.Info("action: all_bets_received | result: success")
			return nil

		default:
			log.Errorf("action: handle_message | result: fail | error: mensaje no soportado %s", msg.MessageType)
			return errors.New("mensaje no soportado")
		}
	}
}

func (s *Server) handleNewBets(conn *net.TCPConn, agencyId int, msg *protocol.ReceivedMessage) {

	if msg.MessageType != protocol.MessageTypeBetBatch {
		handleFailedBetBatch(conn, protocol.MessageBetBatch{}, errors.New("Mensaje recibido no es una apuesta"))
		return
	}

	betsBatchMsg := protocol.MessageBetBatch{}
	err := betsBatchMsg.Decode(msg.Data)
	if err != nil {
		return
	}

	bets := make([]*Bet, 0)
	for _, bet := range betsBatchMsg.Bets {
		b, err := NewBet(strconv.Itoa(agencyId), bet.FirstName, bet.LastName, bet.Document, bet.Birthdate, bet.Number)
		if err != nil {
			handleFailedBetBatch(conn, betsBatchMsg, err)
			return
		}
		bets = append(bets, b)
	}

	err = StoreBets(bets)
	if err != nil {
		handleFailedBetBatch(conn, betsBatchMsg, err)
		return
	}

	betAck := protocol.MessageBetAck{Result: true}

	error := protocol.Send(conn, &betAck)
	if error != nil {
		return
	}

	log.Infof("action: apuesta_recibida | result: success | cantidad: %d", len(bets))
}

func handleFailedBetBatch(clientSocket *net.TCPConn, bet protocol.MessageBetBatch, err error) {
	// Si el error es de conexión cerrada, se termina el programa, si no, se envía un mensaje de error
	if !errors.Is(err, io.EOF) || !errors.Is(err, net.ErrClosed) {
		log.Errorf("action: apuesta_recibida | result: fail | cantidad: %d | error: %s", len(bet.Bets), err)
	}

	betAck := protocol.MessageBetAck{Result: false}
	protocol.Send(clientSocket, &betAck)
}
