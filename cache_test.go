package kagami

import (
	"testing"
)

func TestCache(t *testing.T) {
	// setup a test provider
	RegisterProvider("dummy", dummyProviderFactory)
	provider := CreateProvider("dummy", nil)

	// create our cache instance
	cache := NewCache(&Config{
		Cache: CacheConfig{
			Disabled: true,
		},
	})

	cache.Repository(provider, &Repository{
		Path: "org/foobar",
	})

}
