all: tools deps test

tools:
	@echo "Installing tools"
	@echo " * gox"
	@go get -u github.com/mitchellh/gox
	@echo " * dep"
	@go get -u github.com/golang/dep/cmd/dep
	@echo " * gocov"
	@go get -u github.com/axw/gocov/gocov
.PHONY: tools

deps:
	@echo "Installing dependencies"
	@dep ensure
.PHONY: deps

test:
	@echo "Running tests"
	@go test -cover $(go list ./... | grep -v /vendor/)
.PHONY: test

install:
	@echo "Installing kagami"
	@go install ./cmd/...
.PHONY: install

build:
	@echo "Build distribution binaries"
	@gox -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" -os="linux darwin windows openbsd" ./cmd/...
.PHONY: build
