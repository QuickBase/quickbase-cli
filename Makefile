
PKG=github.com/QuickBase/quickbase-cli/qbclient
VERSION=$(shell git describe --tags)

.PHONY: build
build:
	go build -ldflags "-X $(PKG).Version=$(VERSION)"

.PHONY: install
install:
	go get ./...

.PHONY: update
update:
	go get -u ./...

.PHONY: test
test: install
	go test -v ./qbclient

.PHONY: dist
dist:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X $(PKG).Version=$(VERSION)" -o ./dist/darwin/quickbase-cli
	upx -qqq -9 ./dist/darwin/quickbase-cli
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X $(PKG).Version=$(VERSION)" -o ./dist/linux/quickbase-cli
	upx -qqq -9 ./dist/linux/quickbase-cli
	GOOS=windows GOARCH=386 go build -ldflags="-s -w -X $(PKG).Version=$(VERSION)" -o ./dist/windows/quickbase-cli.exe
	upx -qqq -9 ./dist/windows/quickbase-cli.exe

.PHONY: clean
clean:
	rm -f quickbase-cli
	rm -f dist/darwin/*
	rm -f dist/linux/*
	rm -f dist/windows/*
