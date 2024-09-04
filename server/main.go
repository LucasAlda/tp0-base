package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/server/common"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/shared"
	"github.com/op/go-logging"

	"github.com/spf13/viper"
)

type IniData struct {
	Default Config
}

type Config struct {
	ServerPort          int    `mapstructure:"SERVER_PORT"`
	ServerIp            string `mapstructure:"SERVER_IP"`
	ServerListenBacklog int    `mapstructure:"SERVER_LISTEN_BACKLOG"`
	LoggingLevel        string `mapstructure:"LOGGING_LEVEL"`
	CantAgencies        int    `mapstructure:"CANT_AGENCIES"`
}

var log = logging.MustGetLogger("log")

func initializeConfig() *Config {
	v := viper.New()
	_ = v.BindEnv("default.server_port", "SERVER_PORT")
	_ = v.BindEnv("default.server_ip", "SERVER_IP")
	_ = v.BindEnv("default.server_listen_backlog", "SERVER_LISTEN_BACKLOG")
	_ = v.BindEnv("default.logging_level", "LOGGING_LEVEL")
	_ = v.BindEnv("default.cant_agencies", "CANT_AGENCIES")

	v.SetConfigFile("config.ini")

	v.SetConfigType("ini")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	iniData := IniData{}
	if err := v.Unmarshal(&iniData); err != nil {
		log.Fatalf("Unable to decode into struct: %s", err)
	}

	if iniData.Default.LoggingLevel == "" {
		log.Fatal("LOGGING_LEVEL is not set")
	}

	if iniData.Default.ServerPort == 0 {
		log.Fatal("SERVER_PORT is not set")
	}

	if iniData.Default.ServerListenBacklog == 0 {
		log.Fatal("SERVER_LISTEN_BACKLOG is not set")
	}

	if iniData.Default.CantAgencies == 0 {
		log.Fatal("CANT_AGENCIES is not set")
	}

	return &iniData.Default
}

// PrintConfig Print all the configuration parameters of the program.
// For debugging purposes only
func PrintConfig(config *Config) {

	log.Debugf("action: config | result: success | port: %d | listen_backlog: %d | logging_level: %s", config.ServerPort, config.ServerListenBacklog, config.LoggingLevel)
}

func main() {
	env := initializeConfig()

	if err := shared.InitLogger(env.LoggingLevel); err != nil {
		log.Criticalf("%s", err)
	}

	PrintConfig(env)

	server, err := common.NewServer(env.ServerPort, env.ServerListenBacklog, env.CantAgencies)
	if err != nil {
		log.Criticalf("Error creating server: %s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	server.Run()

	log.Info("action: close_server | result: success")
}
