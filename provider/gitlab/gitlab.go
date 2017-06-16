package provider

import (
	"fmt"

	"github.com/christopherobin/kagami"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

func init() {
	kagami.RegisterProvider("gitlab", NewGitlab)
}

// Gitlab is the github provider and uses deploy keys as credentials
type Gitlab struct {
	DeployKey string `hcl:"deploy_key"`
}

// NewGitlab creates a new gitlab instance
func NewGitlab(config ast.Node) kagami.Provider {
	var gitlab Gitlab
	hcl.DecodeObject(&gitlab, config)
	return &gitlab
}

// Name is the name of the provider
func (gl *Gitlab) Name() string {
	return "github"
}

// Pull pulls a repo from github locally
func (gl *Gitlab) Pull(repo *kagami.Repository, path string) error {
	return fmt.Errorf("Not implemented")
}

// Push pushes a repo to a remote target
func (gl *Gitlab) Push(repo *kagami.Repository, path string) error {
	return fmt.Errorf("Not implemented")
}
