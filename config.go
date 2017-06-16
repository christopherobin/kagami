package kagami

import (
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"

	"fmt"

	log "github.com/sirupsen/logrus"
)

// CacheConfig configures the caching policy of kagami
type CacheConfig struct {
	Path     string `hcl:"path"`
	Disabled bool   `hcl:"disabled"`
}

// ServerConfig configures the http server that listens to service hooks
type ServerConfig struct {
	Addr string `hcl:"addr"`
}

// SourceConfig tells the mirror where to pull the repository from
type SourceConfig struct {
	Provider string `hcl:"provider"`
	Path     string `hcl:"path"`
}

// TargetConfig is a push target, represented by a provider name and a path to
// push to
type TargetConfig struct {
	Provider string `hcl:"provider"`
	Path     string `hcl:"path"`
}

// MirrorConfig represents a git mirror, it should have one source and can have
// multiple targets
type MirrorConfig struct {
	Source  SourceConfig            `hcl:"source"`
	Targets map[string]TargetConfig `hcl:"target"`
}

// Config is the configuration structure for the kagami project
type Config struct {
	Cache   CacheConfig             `hcl:"cache"`
	Server  ServerConfig            `hcl:"server"`
	Mirrors map[string]MirrorConfig `hcl:"mirror"`
}

// LoadConfig loads the git mirroring configuration from a file name
func LoadConfig(name string) (*Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return LoadConfigFromBytes(bytes)
}

// LoadConfigFromBytes loads the git mirroring configuration from a byte array
func LoadConfigFromBytes(bytes []byte) (*Config, error) {
	var config Config
	confAst, err := hcl.ParseBytes(bytes)
	if err != nil {
		return nil, err
	}

	root := confAst.Node.(*ast.ObjectList)

	// cache configuration
	if cacheConfig := root.Filter("cache"); len(cacheConfig.Items) == 1 {
		hcl.DecodeObject(&config.Cache, cacheConfig.Items[0])
	} else {
		if len(cacheConfig.Items) > 1 {
			return nil, fmt.Errorf("duplicate cache configuration")
		}
	}

	// fetch the server config
	if serverConfig := root.Filter("server"); len(serverConfig.Items) == 1 {
		hcl.DecodeObject(&config.Server, serverConfig.Items[0])
	} else {
		if len(serverConfig.Items) > 1 {
			return nil, fmt.Errorf("duplicate server configuration")
		}

		return nil, fmt.Errorf("missing server configuration")
	}

	// take care of providers
	if providerConfig := root.Filter("provider"); len(providerConfig.Items) > 0 {
		for _, providerNode := range providerConfig.Items {
			err = LoadProvider(providerNode)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, fmt.Errorf("no provider configured")
	}

	// finally configure mirrors
	config.Mirrors = make(map[string]MirrorConfig)
	if mirrorConfig := root.Filter("mirror"); len(mirrorConfig.Items) > 0 {
		for _, mirrorNode := range mirrorConfig.Items {
			err = LoadMirror(&config, mirrorNode)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, fmt.Errorf("no provider configured")
	}

	return &config, err
}

// LoadProvider creates a named provider instance
func LoadProvider(node *ast.ObjectItem) error {
	if len(node.Keys) != 1 {
		return fmt.Errorf("invalid name for provider")
	}

	name := node.Keys[0].Token.Value().(string)

	// extract the provider type from the ast Node
	var providerType string
	var providerTypeNode *ast.ObjectList
	if providerTypeNode = node.Val.(*ast.ObjectType).List.Filter("type"); len(providerTypeNode.Items) != 1 {
		return fmt.Errorf("missing type for provider %s", name)
	}

	providerType = providerTypeNode.Items[0].Val.(*ast.LiteralType).Token.Value().(string)

	if _, ok := providers[providerType]; !ok {
		return fmt.Errorf("unknown provider %s", providerType)
	}

	RegisterProviderInstance(name, CreateProvider(providerType, node.Val))

	log.Debugf("loaded provider %s of type %s", name, providerType)

	return nil
}

// LoadMirror creates a new mirror
func LoadMirror(config *Config, node *ast.ObjectItem) error {
	if len(node.Keys) != 1 {
		return fmt.Errorf("invalid name for provider")
	}

	name := node.Keys[0].Token.Value().(string)

	// extract the provider type from the ast Node
	var mirrorConfig MirrorConfig
	err := hcl.DecodeObject(&mirrorConfig, node.Val)
	if err != nil {
		return fmt.Errorf("couldn't decode mirror \"%s\"'s configuration", name)
	}

	config.Mirrors[name] = mirrorConfig

	return nil
}
