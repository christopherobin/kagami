package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/christopherobin/kagami"
)

func checkCommand(configPath string, logLevel string) {
	setLogLevel(logLevel)

	_, err := kagami.LoadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}
}
