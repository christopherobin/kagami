package kagami

import (
	log "github.com/sirupsen/logrus"
)

// Init logs the config, sets the cache and mirrors, then start the http server
func Init(configPath string) {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	SetCacheInstance(NewCache(config))
	SetMirrors(config)
	NewServer(config)
}
