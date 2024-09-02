package common

import (
	"bufio"
	"context"
	"fmt"
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

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(ctx context.Context) {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.Loop.Amount; msgID++ {
		// Create the connection and send the message
		c.sendMessage(msgID)

		select {
		case <-ctx.Done():
			c.cancel()
			return
		case <-time.After(c.config.Loop.Period):
		}

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) sendMessage(id int) {
	c.createClientSocket()
	defer c.close()

	fmt.Fprintf(
		c.conn,
		"[CLIENT %v] Message N°%v\n",
		c.config.ID,
		id,
	)
	msg, err := bufio.NewReader(c.conn).ReadString('\n')

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}

	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
		c.config.ID,
		msg,
	)
}
