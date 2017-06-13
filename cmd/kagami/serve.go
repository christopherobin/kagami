package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/christopherobin/kagami"
)

func serveCommand(configPath string, logLevel string) {
	setLogLevel(logLevel)

	config, err := kagami.LoadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	kagami.NewServer(config)
}
