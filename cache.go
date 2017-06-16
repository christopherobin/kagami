package kagami

import (
	"path"

	"io/ioutil"

	"os"

	appdir "github.com/ProtonMail/go-appdir"
)

// Cache takes care of abstracting caching directory name generation and cleanup
type Cache struct {
	cachePath string
	disabled  bool
}

var cacheInstance *Cache

// NewCache creates a new cache instance
func NewCache(config *Config) *Cache {
	cachePath := config.Cache.Path

	if config.Cache.Disabled {
		dir, err := ioutil.TempDir("", "kagami")
		if err != nil {
			panic(err)
		}
		cachePath = dir
	}

	if cachePath == "" {
		dirs := appdir.New("kagami")
		cachePath = dirs.UserCache()
	}

	return &Cache{
		cachePath: cachePath,
		disabled:  config.Cache.Disabled,
	}
}

// Repository returns the cache path for a given repository
func (c *Cache) Repository(p Provider, r *Repository) (string, error) {
	// generate the path and create it
	dir := path.Join(c.cachePath, p.Name(), r.Path)
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return "", err
	}

	return dir, nil
}

// Cleanup destroys all temporary directories created by the cache system
func (c *Cache) Cleanup() error {
	if c.cachePath != "" {
		return os.RemoveAll(c.cachePath)
	}

	return nil
}

// GetCacheInstance returns the cache instance used by kagami
func GetCacheInstance() *Cache {
	return cacheInstance
}

// SetCacheInstance sets the cache instance used by kagami
func SetCacheInstance(cache *Cache) {
	cacheInstance = cache
}
