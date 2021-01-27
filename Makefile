VERSION=$(shell git describe --tags)

.PHONY: build
build:
	go build -ldflags "-X github.com/QuickBase/quickbase-cli/qbclient.Version=$(VERSION)"

.PHONY: install
install:
	go mod download

.PHONY: test
test:
	go test -v ./...

.PHONY: dist
dist:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -f quickbase-cli
	rm -rf dist
