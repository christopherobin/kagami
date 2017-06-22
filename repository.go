package kagami

import (
	"os"
	"path"
)

// Repository represents a git repo on the disk
type Repository struct {
	Provider Provider
	Path     string

	cache *Cache
}

// NewRepository creates a new repository instance
func NewRepository(p Provider, path string) *Repository {
	return &Repository{
		Path:     path,
		Provider: p,
		cache:    GetCacheInstance(),
	}
}

// Exists checks if the repo currently exists on the disk
func (r *Repository) Exists() bool {
	repoPath, err := r.GetRepoPath()
	if err != nil {
		return false
	}
	configPath := path.Join(repoPath, "config")
	_, err = os.Stat(configPath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// GetRepoPath creates and return a valid path where to clone the repository
func (r *Repository) GetRepoPath() (string, error) {
	return r.cache.Repository(r.Provider, r)
}
