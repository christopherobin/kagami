package main

import (

	// import providers
	"os"

	"github.com/christopherobin/kagami"
	_ "github.com/christopherobin/kagami/provider/github"
	_ "github.com/christopherobin/kagami/provider/gitlab"

	"github.com/alecthomas/kingpin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

var (
	app        = kingpin.New("kagami", "Git mirroring agent")
	configPath = app.Flag("config", "Configuration file.").
			Short('c').
			Default("/etc/kagami.hcl").
			String()
	logLevel = app.Flag("loglevel", "Log level.").
			Short('l').
			Default("INFO").
			Enum("DEBUG", "INFO", "WARN", "ERROR", "FATAL")
	logFormat = app.Flag("log-format", "The format to use for logs.").
			Default("text").
			Enum("json", "text")
	logOutput = app.Flag("log-output", "Where to output logs, use - for stdin.").
			Default("-").
			String()
	testConfig = app.Flag("test-config", "Test the config and exit.").
			Short('t').
			Bool()
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

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *logOutput != "-" {
		logOutFd, err := os.OpenFile(*logOutput, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
		if err != nil {
			panic(err)
		}
		log.SetOutput(logOutFd)
	}

	setLogLevel(*logLevel)

	switch *logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}

	if *testConfig {
		_, err := kagami.LoadConfig(*configPath)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	kagami.Init(*configPath)
}
