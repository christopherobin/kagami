package main

import (
	log "github.com/sirupsen/logrus"
)

var levelMap = map[string]log.Level{
	"DEBUG": log.DebugLevel,
	"INFO":  log.InfoLevel,
	"WARN":  log.WarnLevel,
	"ERROR": log.ErrorLevel,
	"FATAL": log.FatalLevel,
}

func setLogLevel(logLevel string) {
	log.SetLevel(levelMap[logLevel])
}
