package kagami

// Repository represents a git repo on the disk
type Repository struct {
	Path string

	cache *Cache
}

// NewRepository creates a new repository instance
func NewRepository(path string) *Repository {
	return &Repository{
		Path: path,
	}
}

// Exists checks if the repo currently exists on the disk
func (r *Repository) Exists() bool {
	return false
}

// GetRepoPath creates and return a valid path where to clone the repository
func (r *Repository) GetRepoPath(p Provider) (string, error) {
	return GetCacheInstance().Repository(p, r)
}
