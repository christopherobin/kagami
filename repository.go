package kagami

// Repository represents a git repo on the disk
type Repository struct {
	Name string
	Path string
}

// Exists checks if the repo currently exists on the disk
func (r *Repository) Exists() bool {
	return false
}
