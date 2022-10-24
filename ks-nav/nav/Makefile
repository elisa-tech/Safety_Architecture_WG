all:	nav.go
	go build
aarch64: nav.go
	GOARCH="arm64" GOOS="linux" go build
amd64: nav.go
	GOARCH="amd64" GOOS="linux" go build
upx:	nav
	upx nav

