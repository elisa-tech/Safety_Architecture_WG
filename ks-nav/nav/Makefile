.PHONY: all cgo upx test testcgo code/check code/fmt clean

all:	nav.go
	go build

cgo:    nav.go dot_parser/libdotparser.so libdotparser.so
	go build -ldflags="-r ." -tags CGO

upx:	nav
	upx nav

test:
	go test -cover -v

test_cgo:
	go test -v -cover  -ldflags="-r ." -tags CGO

code/check:
	@gofmt -l `find . -type f -name '*.go' -not -path "./vendor/*"`

code/fmt:
	@gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

dot_parser/libdotparser.so: dot_parser/dot.l dot_parser/dot.y
	$(MAKE) -C dot_parser/ libdotparser.so

libdotparser.so: dot_parser/libdotparser.so
	ln -s dot_parser/libdotparser.so libdotparser.so

clean:
	$(MAKE) -C dot_parser/ clean
	rm -f nav libdotparser.so
