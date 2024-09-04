package common

import (
	"errors"
	"io"
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

type BatchConfig struct {
	MaxAmount int `mapstructure:"maxAmount"`
}

type Config struct {
	ID     string       `mapstructure:"id"`
	Server ServerConfig `mapstructure:"server"`
	Loop   LoopConfig   `mapstructure:"loop"`
	Log    LogConfig    `mapstructure:"log"`

	Batch BatchConfig `mapstructure:"batch"`
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
			"action: connect | result: fail | agency: %v | error: %v",
			c.config.ID,
			err,
		)
	}

	presentation := protocol.MessagePresentation{
		Agency: c.config.ID,
	}
	err = protocol.Send(conn, &presentation)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | agency: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}

	c.conn = conn
	return nil
}

func (c *Client) Cancel() {
	c.close()
	log.Debugf("action: close_connection | result: success | agency: %v", c.config.ID)
}

func (c *Client) close() {
	c.conn.Close()
}

func (c *Client) SendBets(betsStr [][]string) {
	c.createClientSocket()
	defer c.close()

	bets := c.createBets(betsStr)

	batchs := c.batchBets(bets)

	for _, betsBatch := range batchs {
		err := c.sendBetBatch(betsBatch)
		if err != nil {
			return
		}

		time.Sleep(c.config.Loop.Period)
	}

	allBetsSent := protocol.MessageAllBetsSent{}
	err := protocol.Send(c.conn, &allBetsSent)
	if err != nil {
		log.Errorf("action: server_disconnected")
		return
	}

}

func (c *Client) createBets(betsStr [][]string) []protocol.MessageBet {
	bets := make([]protocol.MessageBet, len(betsStr))
	for i, betStr := range betsStr {
		bets[i] = protocol.MessageBet{
			FirstName: betStr[0],
			LastName:  betStr[1],
			Document:  betStr[2],
			Birthdate: betStr[3],
			Number:    betStr[4],
		}
	}
	return bets
}

func (c *Client) batchBets(bets []protocol.MessageBet) []protocol.MessageBetBatch {
	maxBetSize := 0
	for _, bet := range bets {
		if bet.GetSize() > maxBetSize {
			maxBetSize = bet.GetSize()
		}
	}

	maxPayloadSize := 8*1024 - 4 - 4 // 8 kB - 4 bytes (size) - 4 bytes (type)
	batchSize := min(
		c.config.Batch.MaxAmount,
		(maxPayloadSize)/(maxBetSize+1), // +1 for the null terminator
	)

	batchs := make([]protocol.MessageBetBatch, 0)

	for i := 0; i < len(bets); i += batchSize {
		end := min(i+batchSize, len(bets))

		betsBatch := protocol.MessageBetBatch{
			Bets: bets[i:end],
		}

		batchs = append(batchs, betsBatch)
	}

	log.Debugf("Tamaño maximo de apuesta: %d, cantidad de apuestas x batch: %d, cantidad de batchs: %d", maxBetSize, batchSize, len(batchs))
	return batchs
}

func (c *Client) sendBetBatch(bet protocol.MessageBetBatch) error {

	err := protocol.Send(c.conn, &bet)
	if err != nil {
		log.Errorf("action: server_disconnected")
		return err
	}

	msg, err := protocol.Receive(c.conn)
	if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
		log.Errorf("action: server_disconnected")
		return err
	}
	if err != nil || msg.MessageType != protocol.MessageTypeBetAck {
		log.Infof("action: apuesta_enviada | result: fail | cantidad: %d", len(bet.Bets))
		return err
	}

	betAck := protocol.MessageBetAck{}
	err = betAck.Decode(msg.Data)
	if err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | cantidad: %d", len(bet.Bets))
		return err
	}

	if betAck.Result {
		log.Infof("action: apuesta_enviada | result: success | cantidad: %d", len(bet.Bets))
	} else {
		log.Errorf("action: apuesta_enviada | result: fail | cantidad: %d", len(bet.Bets))
		return errors.New("El servidor no almacenó la apuesta")
	}

	return nil
}
