BINARY := converse
PLATFORMS := windows/amd64 linux/amd64 darwin/amd64

.PHONY: build test install cross dist

build:
	go build -o $(BINARY) ./main.go

test:
	go test ./... -v -cover

vet:
	go vet ./...

install:
	go install $(LDFLAGS) ./main.go

xcompile:
	./xcompile.sh
