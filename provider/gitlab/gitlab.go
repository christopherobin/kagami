package provider

import (
	"encoding/json"
	"net/http"

	"github.com/christopherobin/kagami"
	"github.com/christopherobin/kagami/provider"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
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
	provider.GitConfig

	name string
}

// NewGitlab creates a new gitlab instance
func NewGitlab(name string, config ast.Node) kagami.Provider {
	var gitlab Gitlab
	hcl.DecodeObject(&gitlab, config)
	gitlab.name = name

	if gitlab.Domain == "" {
		gitlab.Domain = "gitlab.com"
	}

	if gitlab.User == "" {
		gitlab.User = "git"
	}

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
	return provider.Pull(gl.name, gl.GitConfig, repo, path)
}

// Push pushes a repo to a remote target
func (gl *Gitlab) Push(repo *kagami.Repository, path string) error {
	return provider.Push(gl.name, gl.GitConfig, repo, path)
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
