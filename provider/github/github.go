package provider

import (
	"github.com/christopherobin/kagami"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

func init() {
	kagami.RegisterProvider("github", NewGithub)
}

// GithubPushPayload represents a push payload from a github webhook, at least
// the field that we actually care about: https://developer.github.com/v3/activity/events/types/#pushevent
type GithubPushPayload struct {
	HeadCommit struct {
		ID string `json:"id"`
	} `json:"head_commit"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		SSHURL   string `json:"ssh_url"`
	} `json:"repository"`
}

// Github is the github provider and uses deploy keys as credentials
// TODO: we could provide just wrap the normal git provider
type Github struct {
	DeployKey         string `hcl:"deploy_key"`
	DeployKeyPassword string `hcl:"deploy_key_password"`

	name string
}

// NewGithub creates a new github instance
func NewGithub(name string, config ast.Node) kagami.Provider {
	var github Github
	hcl.DecodeObject(&github, config)
	github.name = name
	return &github
}

// Name is the name of the provider
func (gh *Github) Name() string {
	return gh.name
}

// Type is the type of the provider
func (gh *Github) Type() string {
	return "github"
}

// Pull pulls a repo from github locally
func (gh *Github) Pull(repo *kagami.Repository, path string) error {
	clonePath, err := repo.GetRepoPath()
	if err != nil {
		return err
	}

	log.Infof("Pulling repository %s from %s into path %s", path, gh.name, clonePath)

	auth, err := ssh.NewPublicKeysFromFile("git", gh.DeployKey, gh.DeployKeyPassword)
	if err != nil {
		return err
	}

	if !repo.Exists() {
		// TODO: instead of building the URL here, we should find a way to forward
		// the repo information all the way from the hook to here
		_, err = git.PlainClone(clonePath, true, &git.CloneOptions{
			URL:        "git@github.com:" + path + ".git",
			Progress:   os.Stdout,
			RemoteName: gh.name,
			Auth:       auth,
		})
	} else {
		gitRepo, err := git.PlainOpen(clonePath)
		if err != nil {
			return err
		}

		err = gitRepo.Pull(&git.PullOptions{
			RemoteName: gh.name,
			Auth:       auth,
		})
	}

	return err
}

// Push pushes a repo to a remote target
func (gh *Github) Push(repo *kagami.Repository, path string) error {
	return nil
}

// Handle takes care of github webhooks
func (gh *Github) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var payload GithubPushPayload

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)

	// we can't decode? ignore it
	if err != nil {
		w.WriteHeader(400)
		return
	}

	kagami.TrySync(gh, payload.Repository.FullName)
	w.WriteHeader(200)
}
