package kagami

import (
	"net/http"

	"github.com/hashicorp/hcl/hcl/ast"

	log "github.com/sirupsen/logrus"
)

// ProviderCreateFunc is a function that creates a provider from a HCL config node
type ProviderCreateFunc func(name string, config ast.Node) Provider

var providers map[string]ProviderCreateFunc

// RegisterProvider allows a provider to register itself with kagami
func RegisterProvider(name string, providerFunc ProviderCreateFunc) {
	log.Debugf("registering provider %s", name)

	if providers == nil {
		providers = make(map[string]ProviderCreateFunc)
	}

	providers[name] = providerFunc
}

// CreateProvider creates a new provider instance
func CreateProvider(name string, config ast.Node) Provider {
	if _, ok := providers[name]; !ok {
		log.Fatalf("Unknown provider %s", name)
	}

	return providers[name](name, config)
}

// Provider represents a valid provider with valid credentials
type Provider interface {
	Name() string
	Type() string
	Pull(repo *Repository, path string) error
	Push(repo *Repository, path string) error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

var providerInstances map[string]Provider

// RegisterProviderInstance saves a provider instance under a certain name
func RegisterProviderInstance(name string, provider Provider) {
	providerInstances[name] = provider
}

// GetProviderInstance returns a named provider instance
func GetProviderInstance(name string) Provider {
	return providerInstances[name]
}

// GetProviderInstances returns all provider instances
func GetProviderInstances() map[string]Provider {
	return providerInstances
}

func init() {
	providerInstances = make(map[string]Provider)
}
