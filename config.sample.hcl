cache {
  path = "/tmp/kagami"
}

# Configure the HTTP server used by kagami
server {
  addr = ":5000"
}

# This is a provider, it tells kagami which provider to use
# with which credentials
provider "github" {
  #type = "github"
  type = "dummy"

  deploy_key = "/etc/kagami/christopherobin.pem"
}

provider "gitlab" {
  #type = "gitlab"
  type = "dummy"

  deploy_key = "/etc/kagami/christopherobin.pem"
}

# TODO: Not supported yet
#provider "bitbucket" {
#  type = "gitlab"
#
#  deploy_key = "/etc/kagami/christopherobin.pem"
#}
#
# Polls a raw git repo on a regular basis
#provider "git" {
#  type = "git"
#
#  url = "git+ssh://my.custom.server.net/my/repo.git
#  interval = "1m"
#}

# A mirror describes a project to mirror, it must have one source
# and at least 1 target
mirror "kagami" {
  source {
    provider = "github"
    path = "christopherobin/kagami"
    branches = ["develop", "master"]
  }

  target "gitlab" {
    provider = "gitlab"
    path = "christopherobin/kagami"
  }

  target "bitbucket" {
    provider = "bitbucket"
    path = "christopherobin/kagami"
  }
}
