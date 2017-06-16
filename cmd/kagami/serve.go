package main

import (
	"github.com/christopherobin/kagami"
)

func serveCommand(configPath string, logLevel string) {
	setLogLevel(logLevel)
	kagami.Init(configPath)

	/*config, err := kagami.LoadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	kagami.NewServer(config)*/
}
