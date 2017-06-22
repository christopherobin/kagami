package provider

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/christopherobin/kagami"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

func init() {
	kagami.RegisterProvider("gitlab", NewGitlab)
}

// GitlabPushPayload represents a push payload from a gitlab webhook, at least
// the fields that we actually care about: https://docs.gitlab.com/ce/user/project/integrations/webhooks.html#push-events
type GitlabPushPayload struct {
	Project struct {
		PathWithNamespace string `json:"path_with_namespace"`
	} `json:"project"`
	Repository struct {
		Name   string `json:"name"`
		SSHURL string `json:"git_ssh_url"`
	} `json:"repository"`
}

// Gitlab is the github provider and uses deploy keys as credentials
type Gitlab struct {
	UseSSH            bool   `hcl:"use_ssh"`
	DeployKey         string `hcl:"deploy_key"`
	DeployKeyPassword string `hcl:"deploy_key_password"`

	name string
}

// NewGitlab creates a new gitlab instance
func NewGitlab(name string, config ast.Node) kagami.Provider {
	var gitlab Gitlab
	hcl.DecodeObject(&gitlab, config)
	gitlab.name = name
	return &gitlab
}

// Name is the name of the provider
func (gl *Gitlab) Name() string {
	return gl.name
}

// Type is the type of the provider
func (gl *Gitlab) Type() string {
	return "github"
}

// Pull pulls a repo from github locally
func (gl *Gitlab) Pull(repo *kagami.Repository, path string) error {
	clonePath, err := repo.GetRepoPath()
	if err != nil {
		return err
	}

	log.Infof("Pulling repository %s from %s into path %s", path, gl.name, clonePath)

	var auth ssh.AuthMethod

	if gl.DeployKey != "" {
		auth, err = ssh.NewPublicKeysFromFile("git", gl.DeployKey, gl.DeployKeyPassword)
		if err != nil {
			return err
		}
	}

	if !repo.Exists() {
		repoURL := "https://gitlab.com/" + path + ".git"
		if gl.UseSSH {
			repoURL = "git@gitlab.com:" + path + ".git"
		}

		_, err = git.PlainClone(clonePath, true, &git.CloneOptions{
			URL:        repoURL,
			Progress:   os.Stdout,
			RemoteName: gl.name,
			Auth:       auth,
		})
	} else {
		gitRepo, err := git.PlainOpen(clonePath)
		if err != nil {
			return err
		}

		err = gitRepo.Pull(&git.PullOptions{
			RemoteName: gl.name,
			Auth:       auth,
		})
	}

	return err
}

// Push pushes a repo to a remote target
func (gl *Gitlab) Push(repo *kagami.Repository, path string) error {
	clonePath, err := repo.GetRepoPath()
	if err != nil {
		return err
	}

	log.Infof("Pushing repository %s to %s", path, gl.name)

	gitRepo, err := git.PlainOpen(clonePath)
	if err != nil {
		return err
	}

	if remote, _ := gitRepo.Remote(gl.name); remote == nil {
		repoURL := "https://gitlab.com/" + path + ".git"
		if gl.UseSSH {
			repoURL = "git@gitlab.com:" + path + ".git"
		}

		_, err = gitRepo.CreateRemote(&config.RemoteConfig{
			Name: gl.name,
			URL:  repoURL,
		})

		if err != nil {
			return err
		}
	}

	// create the SSH auth method
	var auth ssh.AuthMethod

	if gl.DeployKey != "" {
		auth, err = ssh.NewPublicKeysFromFile("git", gl.DeployKey, gl.DeployKeyPassword)
		if err != nil {
			return err
		}
	}

	err = gitRepo.Push(&git.PushOptions{
		RemoteName: gl.name,
		Auth:       auth,
	})

	return err
}

// Handle takes care of gitlab webhooks
func (gl *Gitlab) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var payload GitlabPushPayload

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)

	// we can't decode? ignore it
	if err != nil {
		w.WriteHeader(400)
		return
	}

	kagami.TrySync(gl, payload.Project.PathWithNamespace)
	w.WriteHeader(200)
}
