package provider

import (
	"github.com/christopherobin/kagami"
	"github.com/christopherobin/kagami/provider"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"

	"encoding/json"
	"net/http"
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
	provider.GitConfig

	name string
}

// NewGithub creates a new github instance
func NewGithub(name string, config ast.Node) kagami.Provider {
	var github Github
	hcl.DecodeObject(&github, config)
	github.name = name

	if github.Domain == "" {
		github.Domain = "github.com"
	}

	if github.User == "" {
		github.User = "git"
	}

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
	return provider.Pull(gh.name, gh.GitConfig, repo, path)
}

// Push pushes a repo to a remote target
func (gh *Github) Push(repo *kagami.Repository, path string) error {
	return provider.Push(gh.name, gh.GitConfig, repo, path)
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
