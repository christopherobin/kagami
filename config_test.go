package kagami

import (
	"testing"
)

var testConfig = []byte(`cache {
  path = "/tmp/kagami"
}

# Configure the HTTP server used by kagami
server {
  addr = ":5000"
}

# This is a provider, it tells kagami which provider to use
# with which credentials
provider "dummy" {
  type = "dummy"
}

# A mirror describes a project to mirror, it must have one source
# and at least 1 target
mirror "kagami" {
  source {
    provider = "dummy"
    path = "christopherobin/kagami"
    branches = ["develop", "master"]
  }

  target "dummy" {
    provider = "dummy"
    path = "christopherobin/kagami2"
  }
}
`)

func TestConfigLoad(t *testing.T) {
	RegisterProvider("dummy", dummyProviderFactory)

	_, err := LoadConfigFromBytes(testConfig)

	if err != nil {
		t.Error(err)
	}
}
