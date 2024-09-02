package common

import (
	"errors"
	"net"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/shared/protocol"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

type LoopConfig struct {
	Period time.Duration `mapstructure:"period"`
	Amount int           `mapstructure:"amount"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type Config struct {
	ID     string       `mapstructure:"id"`
	Server ServerConfig `mapstructure:"server"`
	Loop   LoopConfig   `mapstructure:"loop"`
	Log    LogConfig    `mapstructure:"log"`

	Bet protocol.MessageBet `mapstructure:"bet"`
}

// Client Entity that encapsulates how
type Client struct {
	config Config
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config Config) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.Server.Address)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

func (c *Client) cancel() {
	c.close()
	log.Debugf("action: close_connection | result: success | client_id: %v", c.config.ID)
}

func (c *Client) close() {
	c.conn.Close()
}

func (c *Client) SendBet(bet protocol.MessageBet) {
	c.createClientSocket()
	defer c.close()

	if err := protocol.Send(c.conn, &bet); err != nil {
		handleFailedBet(bet, err)
		return
	}

	msg, err := protocol.Receive(c.conn)
	if err != nil || msg.MessageType != protocol.MessageTypeBetAck {
		handleFailedBet(bet, err)
		return
	}

	betAck := protocol.MessageBetAck{}
	err = betAck.Decode(msg.Data)
	if err != nil {
		handleFailedBet(bet, err)
		return
	}

	if betAck.Result {
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", bet.Document, bet.Number)
	} else {
		handleFailedBet(bet, errors.New("El servidor no almacenó la apuesta"))
	}
}

func handleFailedBet(bet protocol.MessageBet, err error) {
	log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v | error: %v",
		bet.Document,
		bet.Number,
		err,
	)
}
