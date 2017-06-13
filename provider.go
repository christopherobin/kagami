package kagami

import (
	"github.com/hashicorp/hcl/hcl/ast"

	log "github.com/sirupsen/logrus"
)

// ProviderCreateFunc is a function that creates a provider from a HCL config node
type ProviderCreateFunc func(config ast.Node) Provider

var providers map[string]ProviderCreateFunc

// RegisterProvider allows a provider to register itself with kagami
func RegisterProvider(name string, providerFunc ProviderCreateFunc) {
	log.Debugf("registering provider %s", name)

	if providers == nil {
		providers = make(map[string]ProviderCreateFunc)
	}

	providers[name] = providerFunc
}

// Provider represents a valid provider with valid credentials
type Provider interface {
	Pull(repo *Repository, path string) error
	Push(repo *Repository, path string) error
}

var providerInstances map[string]Provider

func init() {
	providerInstances = make(map[string]Provider)
}
