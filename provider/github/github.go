package provider

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/christopherobin/kagami"
)

func init() {
	kagami.RegisterProvider("github", NewGithub)
}

// Github is the github provider and uses deploy keys as credentials
type Github struct {
	Name      string `hcl:"name"`
	DeployKey string `hcl:"deploy_key"`
}

func NewGithub(config ast.Node) kagami.Provider {
	var github Github
	hcl.DecodeObject(&github, config)
	return &github
}

// Pull pulls a repo from github locally
func (gh *Github) Pull(repo *kagami.Repository, path string) error {
	return nil
}

// Push pushes a repo to a remote target
func (gh *Github) Push(repo *kagami.Repository, path string) error {
	return nil
}
