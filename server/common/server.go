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
	agencies     []*Client
	cantAgencies int
}

func NewServer(port int, listenBacklog int, cantAgencies int) (*Server, error) {
	serverSocket, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		return nil, err
	}

	return &Server{serverSocket: serverSocket, cantAgencies: cantAgencies}, nil
}

func (s *Server) Close() {
	s.serverSocket.Close()
	for _, agency := range s.agencies {
		if agency != nil {
			agency.conn.Close()
		}
	}
}

// Dummy Server loop

// Server that accept a new connections and establishes a
// communication with a client. After client with communucation
// finishes, servers starts to accept new connections again
func (s *Server) Run() {
	for i := 0; i < s.cantAgencies; i++ {
		client, err := s.acceptNewConnection()
		s.agencies = append(s.agencies, client)
		if err != nil {
			return
		}

		s.handleConnection(client)
	}

	s.handleWinners()
}

func (s *Server) acceptNewConnection() (*Client, error) {
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

	receivedMessage, err := protocol.Receive(clientSocket)
	if err != nil || receivedMessage.MessageType != protocol.MessageTypePresentation {
		log.Errorf("action: accept_connections | result: fail | error: %s", err)
		return nil, err
	}

	presentation := protocol.MessagePresentation{}
	err = presentation.Decode(receivedMessage.Data)
	if err != nil {
		return nil, err
	}

	client := NewClient(clientSocket, presentation.Agency)

	log.Infof("action: accept_connections | result: success | agency: %d", client.agency)

	return client, nil
}

func (s *Server) handleConnection(client *Client) {
	for {
		msg, err := protocol.Receive(client.conn)
		// Si el reader devuelve EOF, el cliente se desconectó
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			log.Errorf("action: cliente_desconectado | agency: %d", client.agency)
			return
		}
		// Si ocurre otro error, lo registramos
		if err != nil {
			handleFailedBetBatch(client, protocol.MessageBetBatch{}, err)
			return
		}

		switch msg.MessageType {

		case protocol.MessageTypeBetBatch:
			s.handleNewBets(client, msg)

		case protocol.MessageTypeAllBetsSent:
			log.Info("action: total_apuestas_recibidas | result: success")
			return

		default:
			log.Errorf("action: handle_message | result: fail | error: mensaje no soportado %s", msg.MessageType)
			return
		}
	}
}

func (s *Server) handleWinners() {
	log.Infof("action: sorteo | result: in_progress")
	bets, err := LoadBets()
	if err != nil {
		log.Errorf("action: sorteo | result: fail | error: %s", err)
		return
	}

	all_winners := make([]*Bet, 0)
	for _, bet := range bets {
		if bet.HasWon() {
			all_winners = append(all_winners, bet)
		}
	}

	for _, client := range s.agencies {
		agency_winners := make([]string, 0)
		for _, winner := range all_winners {
			if winner.agency == client.agency {
				agency_winners = append(agency_winners, winner.document)
			}
		}

		winnersMsg := protocol.MessageWinners{Winners: agency_winners}
		err := protocol.Send(client.conn, &winnersMsg)
		if err != nil {
			log.Errorf("action: sorteo | result: fail | agency: %d | error: %s", client.agency, err)
			return
		}
	}

	log.Infof("action: sorteo | result: success")
}

func (s *Server) handleNewBets(client *Client, msg *protocol.ReceivedMessage) {

	if msg.MessageType != protocol.MessageTypeBetBatch {
		handleFailedBetBatch(client, protocol.MessageBetBatch{}, errors.New("Mensaje recibido no es una apuesta"))
		return
	}

	betsBatchMsg := protocol.MessageBetBatch{}
	err := betsBatchMsg.Decode(msg.Data)
	if err != nil {
		return
	}

	bets := make([]*Bet, 0)
	for _, bet := range betsBatchMsg.Bets {
		b, err := NewBet(strconv.Itoa(client.agency), bet.FirstName, bet.LastName, bet.Document, bet.Birthdate, bet.Number)
		if err != nil {
			handleFailedBetBatch(client, betsBatchMsg, err)
			return
		}
		bets = append(bets, b)
	}

	err = StoreBets(bets)

	if err != nil {
		handleFailedBetBatch(client, betsBatchMsg, err)
		return
	}

	betAck := protocol.MessageBetAck{Result: true}

	error := protocol.Send(client.conn, &betAck)
	if error != nil {
		return
	}

	log.Infof("action: apuesta_recibida | result: success | cantidad: %d", len(bets))
}

func handleFailedBetBatch(client *Client, bet protocol.MessageBetBatch, err error) {
	// Si el error es de conexión cerrada, se termina el programa, si no, se envía un mensaje de error
	if !errors.Is(err, io.EOF) || !errors.Is(err, net.ErrClosed) {
		log.Errorf("action: apuesta_recibida | result: fail | cantidad: %d | error: %s", len(bet.Bets), err)
	}

	betAck := protocol.MessageBetAck{Result: false}
	protocol.Send(client.conn, &betAck)
}

type Client struct {
	conn   *net.TCPConn
	agency int
}

func NewClient(conn *net.TCPConn, agency string) *Client {
	agencyID, _ := strconv.Atoi(agency)
	return &Client{conn: conn, agency: agencyID}
}
