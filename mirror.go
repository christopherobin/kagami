package kagami

import log "github.com/sirupsen/logrus"

var mirrors map[string]MirrorConfig

// SetMirrors sets the current active mirrors
func SetMirrors(config *Config) {
	mirrors = config.Mirrors
}

// TrySync will take path, remote and commitID and try to match mirrors
func TrySync(provider Provider, path string) {
	// match the path to one of the mirroring
	for name, mirror := range mirrors {
		if mirror.Source.Path == path && mirror.Source.Provider == provider.Name() {
			log.Infof("Received sync demand for mirror %s", name)
			Sync(mirror)
		}
	}
}

// Sync will synchronize a mirror
func Sync(mirror MirrorConfig) {
	provider := mirror.Source.providerInstance()
	repo := NewRepository(provider, mirror.Source.Path)

	// pull the repo from the source
	err := provider.Pull(repo, mirror.Source.Path)
	if err != nil {
		log.Error(err)
		return
	}

	// push to each target
	for _, target := range mirror.Targets {
		err = target.providerInstance().Push(repo, target.Path)

		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (sc SourceConfig) providerInstance() Provider {
	return GetProviderInstance(sc.Provider)
}

func (tc TargetConfig) providerInstance() Provider {
	return GetProviderInstance(tc.Provider)
}
