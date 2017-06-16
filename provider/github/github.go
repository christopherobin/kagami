package provider

import (
	"github.com/christopherobin/kagami"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"

	"os"

	git "gopkg.in/src-d/go-git.v4"
)

func init() {
	kagami.RegisterProvider("github", NewGithub)
}

// Github is the github provider and uses deploy keys as credentials
// TODO: we could provide just wrap the normal git provider
type Github struct {
	DeployKey string `hcl:"deploy_key"`
}

// NewGithub creates a new github instance
func NewGithub(config ast.Node) kagami.Provider {
	var github Github
	hcl.DecodeObject(&github, config)
	return &github
}

// Name is the name of the provider
func (gh *Github) Name() string {
	return "github"
}

// Pull pulls a repo from github locally
func (gh *Github) Pull(repo *kagami.Repository, path string) error {
	clonePath, err := repo.GetRepoPath(gh)
	if err != nil {
		return err
	}

	_, err = git.PlainClone(clonePath, true, &git.CloneOptions{
		URL:      "https://github.com/" + repo.Path,
		Progress: os.Stdout,
	})

	return err
}

// Push pushes a repo to a remote target
func (gh *Github) Push(repo *kagami.Repository, path string) error {
	return nil
}
