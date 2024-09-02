package main

import (
	"github.com/7574-sistemas-distribuidos/docker-compose-init/server/common"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/shared"
	"github.com/op/go-logging"

	"github.com/spf13/viper"
)

type IniData struct {
	Default Env
}

type Env struct {
	ServerPort          int    `mapstructure:"SERVER_PORT"`
	ServerIp            string `mapstructure:"SERVER_IP"`
	ServerListenBacklog int    `mapstructure:"SERVER_LISTEN_BACKLOG"`
	LoggingLevel        string `mapstructure:"LOGGING_LEVEL"`
}

var log = logging.MustGetLogger("log")

func initializeConfig() *Env {
	v := viper.New()
	_ = v.BindEnv("default.server_port", "SERVER_PORT")
	_ = v.BindEnv("default.server_ip", "SERVER_IP")
	_ = v.BindEnv("default.server_listen_backlog", "SERVER_LISTEN_BACKLOG")
	_ = v.BindEnv("default.logging_level", "LOGGING_LEVEL")

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

	return &iniData.Default
}

// PrintConfig Print all the configuration parameters of the program.
// For debugging purposes only
func PrintConfig(env *Env) {

	log.Debugf("action: config | result: success | port: %d | listen_backlog: %d | logging_level: %s", env.ServerPort, env.ServerListenBacklog, env.LoggingLevel)
}

func main() {
	env := initializeConfig()

	if err := shared.InitLogger(env.LoggingLevel); err != nil {
		log.Criticalf("%s", err)
	}

	PrintConfig(env)

	server, err := common.NewServer(env.ServerPort, env.ServerListenBacklog)
	if err != nil {
		log.Criticalf("Error creating server: %s", err)
	}

	server.Run()
}
