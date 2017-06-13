package main

import (

	// import providers
	"os"

	_ "github.com/christopherobin/kagami/provider/github"
	_ "github.com/christopherobin/kagami/provider/gitlab"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
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

	check = app.Command("check", "Check if configuration is valid.")
	serve = app.Command("serve", "Start the kagami server.")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case check.FullCommand():
		checkCommand(*configPath, *logLevel)
	case serve.FullCommand():
		serveCommand(*configPath, *logLevel)
	}
}
