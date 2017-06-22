package kagami

import (
	"log"
)

// Init initialize kagami
func Init(configPath string) {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	SetCacheInstance(NewCache(config))
	SetMirrors(config)
	NewServer(config)
}
