package kagami

import (
	"net/http"
	"testing"

	"reflect"

	"github.com/hashicorp/hcl/hcl/ast"
)

// helpers
var providerType = reflect.TypeOf((*Provider)(nil)).Elem()

type dummyProvider struct {
	name string
}

func dummyProviderFactory(name string, _ ast.Node) Provider {
	return &dummyProvider{
		name: name,
	}
}

func (t *dummyProvider) Name() string {
	return t.name
}

func (t *dummyProvider) Type() string {
	return "dummy"
}

func (t *dummyProvider) Pull(repo *Repository, path string) error {
	return nil
}

func (t *dummyProvider) Push(repo *Repository, path string) error {
	return nil
}

func (t *dummyProvider) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(501)
}

func TestProviderRegister(t *testing.T) {
	RegisterProvider("dummy", dummyProviderFactory)

	// make sure we can create a valid instance
	instance := CreateProvider("dummy", nil)
	if !reflect.TypeOf(instance).Implements(providerType) {
		t.Fatal("Registering provider failed")
	}
}

func TestProviderRegisterInstance(t *testing.T) {
	RegisterProvider("dummy", dummyProviderFactory)

	instance := CreateProvider("dummy", nil)
	RegisterProviderInstance("dummy/bar", instance)
	if GetProviderInstance("dummy/bar") != instance {
		t.Fatal("Registering provider instance failed")
	}
}
