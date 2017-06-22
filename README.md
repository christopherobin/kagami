# kagami

[![Build Status](https://travis-ci.org/christopherobin/kagami.svg)](https://travis-ci.org/christopherobin/kagami)

Small daemon that helps mirror git repository between various
providers such as github/bitbucket/etc...

Work in progress, pulling should work from github, pushing only works on SSH
providers with specifying a key.

## Usage

Copy `config.sample.hcl` to `/etc/kagami.hcl`, change based on
your needs.

```
usage: kagami [<flags>]

Git mirroring agent

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
  -c, --config="/etc/kagami.hcl"
                         Configuration file.
  -l, --loglevel=INFO    Log level.
      --log-format=text  The format to use for logs.
      --log-output="-"   Where to output logs, use - for stdin.
  -t, --test-config      Test the config and exit.
```

## Config

```hcl
# optional: cache configuration for kagami
cache {
  # optional: sets the path of caching, otherwise uses system defaults
  # path = "/var/cache/kagami"

  # optional: disable caching (uses /tmp instead)
  # disabled = false
}

# required: http server configuration
server {
  # required: the listen address of the http server
  addr = ":5050"
}

# required (at least one): sets up a named provider for use in mirrors
provider "github" {
  # required: the provider type, see the providers section
  type = "github"
}

provider "gitlab" {
  type = "gitlab"
}

# required (at least one): describes a git mirror
mirror "kagami" {
  # required: the mirror source
  source {
    # required: the provider to use for pulling/cloning
    provider = "github"

    # required: the path to the repository on the given provider
    path = "christopherobin/kagami"
  }

  # required (at least one target): the actual mirror
  target "gitlab" {
    # required: which provider to use for pushing
    provider = "gitlab"

    # required: the path to the repository on the given provider
    path = "crobin/kagami"
  }

}

```

## Providers

### github

```hcl
provider "github" {
  type = "github"

  # optional: use SSH instead of HTTPS
  # use_ssh: true

  # optional: the SSH key to use for SSH
  # deploy_key: /home/kagami/.ssh/id_ecdsa

  # optional: if the SSH key has a passphrase, enter it here
  # deploy_key_password: "passphrase"
}
```

### gitlab

```hcl
provider "gitlab" {
  type = "gitlab"

  # optional: use SSH instead of HTTPS
  # use_ssh: true

  # optional: the SSH key to use for SSH
  # deploy_key: /home/kagami/.ssh/id_ecdsa

  # optional: if the SSH key has a passphrase, enter it here
  # deploy_key_password: "passphrase"
}
```
