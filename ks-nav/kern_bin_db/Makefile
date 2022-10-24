all:	main.go
	go build
aarch64: main.go
	GOARCH="arm64" GOOS="linux" go build
amd64: main.go
	GOARCH="amd64" GOOS="linux" go build
upx:	kern_bin_db
	upx kern_bin_db

