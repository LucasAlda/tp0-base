package common

import (
	"encoding/csv"
	"errors"
	"io"
	"net"
	"os"
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
	log.Debugf("action: cerrar_conexion | result: success | agency: %v", c.config.ID)
	c.Close()
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SendBets(betsFile *os.File) error {
	defer betsFile.Close()

	c.createClientSocket()
	betsReader := csv.NewReader(betsFile)

	println("Sending bets")

	batch := protocol.MessageBetBatch{
		Bets: make([]protocol.MessageBet, 0),
	}
	batchSize := 4 + 4 // 4 bytes for the size and 4 bytes for the type

	for {
		record, err := betsReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		bet := protocol.MessageBet{
			FirstName: record[0],
			LastName:  record[1],
			Document:  record[2],
			Birthdate: record[3],
			Number:    record[4],
		}

		if batchSize+bet.GetSize()+1 > 8*1024 || len(batch.Bets) >= c.config.Batch.MaxAmount {
			err := c.sendBetBatch(&batch)
			if err != nil {
				return err
			}

			batch.Bets = make([]protocol.MessageBet, 0)
			batchSize = 4 + 4 // 4 bytes for the size + 4 bytes for the type
			time.Sleep(c.config.Loop.Period)
		}

		batch.Bets = append(batch.Bets, bet)
		batchSize += bet.GetSize() + 1
	}

	return nil
}

func (c *Client) GetWinners() {
	log.Infof("action: consulta_ganadores | result: in_progress")

	allBetsSent := protocol.MessageAllBetsSent{}
	err := protocol.Send(c.conn, &allBetsSent)
	if err != nil {
		log.Errorf("action: servidor_desconectado")
		return
	}

	receivedMsg, err := protocol.Receive(c.conn)
	if err != nil {
		log.Errorf("action: servidor_desconectado")
		return
	}
	if receivedMsg.MessageType != protocol.MessageTypeWinners {
		log.Errorf("action: consulta_ganadores | result: fail | error: message type mismatch")
		return
	}

	winners := protocol.MessageWinners{}
	err = winners.Decode(receivedMsg.Data)
	if err != nil {
		log.Errorf("action: consulta_ganadores | result: fail | error: %s", err)
		return
	}

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners.Winners))
}

func (c *Client) sendBetBatch(batch *protocol.MessageBetBatch) error {

	err := protocol.Send(c.conn, batch)
	if err != nil {
		log.Errorf("action: servidor_desconectado")
		return err
	}

	msg, err := protocol.Receive(c.conn)
	if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
		log.Errorf("action: servidor_desconectado")
		return err
	}
	if err != nil || msg.MessageType != protocol.MessageTypeBetAck {
		log.Infof("action: apuesta_enviada | result: fail | cantidad: %d", len(batch.Bets))
		return err
	}

	betAck := protocol.MessageBetAck{}
	err = betAck.Decode(msg.Data)
	if err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | cantidad: %d", len(batch.Bets))
		return err
	}

	if betAck.Result {
		log.Infof("action: apuesta_enviada | result: success | cantidad: %d", len(batch.Bets))
	} else {
		log.Errorf("action: apuesta_enviada | result: fail | cantidad: %d", len(batch.Bets))
		return errors.New("El servidor no almacenó la apuesta")
	}

	return nil
}
