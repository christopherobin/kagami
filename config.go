package kagami

import (
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"

	log "github.com/sirupsen/logrus"
)

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
	Server  ServerConfig            `hcl:"server"`
	Mirrors map[string]MirrorConfig `hcl:"mirror"`
}

// LoadConfig loads the git mirroring configuration
func LoadConfig(name string) (*Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var config Config
	confAst, err := hcl.ParseBytes(b)
	if err != nil {
		return nil, err
	}

	root := confAst.Node.(*ast.ObjectList)

	// fetch the server config
	if serverConfig := root.Filter("server"); len(serverConfig.Items) == 1 {
		hcl.DecodeObject(&config.Server, serverConfig.Items[0])
	} else {
		if len(serverConfig.Items) > 1 {
			log.Fatalln("duplicate server configuration")
		}

		log.Fatalln("missing server configuration")
	}

	// take care of providers
	if providerConfig := root.Filter("provider"); len(providerConfig.Items) > 0 {
		for _, providerNode := range providerConfig.Items {
			LoadProvider(providerNode)
		}
	} else {
		log.Fatalln("no provider configured")
	}

	return &config, err
}

// LoadProvider creates a named provider instance
func LoadProvider(node *ast.ObjectItem) {
	if len(node.Keys) != 1 {
		log.Fatalln("invalid name for provider")
	}

	name := node.Keys[0].Token.Value().(string)

	// extract the provider type from the ast Node
	var providerType string
	if providerTypeNode := node.Val.(*ast.ObjectType).List.Filter("type"); len(providerTypeNode.Items) != 1 {
		log.Fatalf("missing type for provider %s", name)
	} else {
		providerType = providerTypeNode.Items[0].Val.(*ast.LiteralType).Token.Value().(string)
	}

	if _, ok := providers[providerType]; !ok {
		log.Fatalf("unknown provider %s", providerType)
	}

	providerInstances[name] = providers[providerType](node.Val)

	log.Debugf("loaded provider %s of type %s", name, providerType)
}
