
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
dist: dist-darwin dist-linux dist-windows

.PHONY: dist-darwin
dist-darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X $(PKG).Version=$(VERSION)" -o ./dist/darwin/quickbase-cli
	upx -qqq -9 ./dist/darwin/quickbase-cli
	(cd ./dist/darwin/ && zip -X quickbase-cli.darwin.zip quickbase-cli)

.PHONY: dist-linux
dist-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X $(PKG).Version=$(VERSION)" -o ./dist/linux/quickbase-cli
	upx -qqq -9 ./dist/linux/quickbase-cli
	(cd ./dist/linux/ && tar -czvf quickbase-cli.linux.tgz quickbase-cli)

.PHONY: dist-windows
dist-windows:
	GOOS=windows GOARCH=386 go build -ldflags="-s -w -X $(PKG).Version=$(VERSION)" -o ./dist/windows/quickbase-cli.exe
	upx -qqq -9 ./dist/windows/quickbase-cli.exe
	(cd ./dist/windows/ && zip -X quickbase-cli.windows.zip quickbase-cli.exe)

.PHONY: clean
clean:
	rm -f quickbase-cli
	rm -f dist/darwin/*
	rm -f dist/linux/*
	rm -f dist/windows/*
