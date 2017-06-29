package provider

// this file implements standard operations that are shared among all providers

import (
	"os"

	"github.com/christopherobin/kagami"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

type GitConfig struct {
	Domain            string `hcl:"domain"`
	UseSSH            bool   `hcl:"use_ssh"`
	User              string `hcl:"user"`
	DeployKey         string `hcl:"deploy_key"`
	DeployKeyPassword string `hcl:"deploy_key_password"`
}

func Pull(name string, config GitConfig, repo *kagami.Repository, path string) error {
	clonePath, err := repo.GetRepoPath()
	if err != nil {
		return err
	}

	log.Infof("Pulling repository %s from %s into path %s", path, name, clonePath)

	var auth ssh.AuthMethod

	if config.DeployKey != "" {
		auth, err = ssh.NewPublicKeysFromFile("git", config.DeployKey, config.DeployKeyPassword)
		if err != nil {
			return err
		}
	}

	if !repo.Exists() {
		repoURL := "https://" + config.Domain + "/" + path + ".git"
		if config.UseSSH {
			repoURL = config.User + "@" + config.Domain + ":" + path + ".git"
		}

		_, err = git.PlainClone(clonePath, true, &git.CloneOptions{
			URL:        repoURL,
			Progress:   os.Stdout,
			RemoteName: name,
			Auth:       auth,
		})
	} else {
		gitRepo, err := git.PlainOpen(clonePath)
		if err != nil {
			return err
		}

		err = gitRepo.Pull(&git.PullOptions{
			RemoteName: name,
			Auth:       auth,
		})
	}

	return err
}

// Push pushes a repo to a remote target
func Push(name string, gitConfig GitConfig, repo *kagami.Repository, path string) error {
	clonePath, err := repo.GetRepoPath()
	if err != nil {
		return err
	}

	log.Infof("Pushing repository %s to %s", path, name)

	gitRepo, err := git.PlainOpen(clonePath)
	if err != nil {
		return err
	}

	if remote, _ := gitRepo.Remote(name); remote == nil {
		repoURL := "https://" + gitConfig.Domain + "/" + path + ".git"
		if gitConfig.UseSSH {
			repoURL = gitConfig.User + "@" + gitConfig.Domain + ":" + path + ".git"
		}

		_, err = gitRepo.CreateRemote(&config.RemoteConfig{
			Name: name,
			URL:  repoURL,
		})

		if err != nil {
			return err
		}
	}

	// create the SSH auth method
	var auth ssh.AuthMethod

	if gitConfig.DeployKey != "" {
		auth, err = ssh.NewPublicKeysFromFile(gitConfig.User, gitConfig.DeployKey, gitConfig.DeployKeyPassword)
		if err != nil {
			return err
		}
	}

	err = gitRepo.Push(&git.PushOptions{
		RemoteName: name,
		Auth:       auth,
	})

	return err
}
