package provider

import (
	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/christopherobin/kagami"
)

func init() {
	kagami.RegisterProvider("gitlab", NewGitlab)
}

// Gitlab is the github provider and uses deploy keys as credentials
type Gitlab struct {
	Name      string `hcl:"name"`
	DeployKey string `hcl:"deploy_key"`
}

// NewGitlab creates a new gitlab instance
func NewGitlab(config ast.Node) kagami.Provider {
	var gitlab Gitlab
	hcl.DecodeObject(&gitlab, config)
	return &gitlab
}

// Pull pulls a repo from github locally
func (gl *Gitlab) Pull(repo *kagami.Repository, path string) error {
	return fmt.Errorf("Not implemented")
}

// Push pushes a repo to a remote target
func (gl *Gitlab) Push(repo *kagami.Repository, path string) error {
	return fmt.Errorf("Not implemented")
}
