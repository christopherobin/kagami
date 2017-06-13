test:
	@go test -cover ./...
.PHONY: test

install:
	@go install ./cmd/...
.PHONY: install

build:
	@gox -os="linux darwin windows openbsd" ./...
.PHONY: build