BINARY := converse
PLATFORMS := windows/amd64 linux/amd64 darwin/amd64

.PHONY: build test install cross

build:
	go build -o $(BINARY) ./main.go

test:
	go test ./... -v -cover

vet:
	go vet ./...

install:
	go install $(LDFLAGS) ./main.go

cross: $(addprefix dist/,$(PLATFORMS))

dist/%:
	GOOS=$(word 1, $(subst /, ,$@))
	GOARCH=$(word 2, $(subst /, ,$@))
	GOARM=$(word 3, $(subst /, ,$@))
	CGO_ENABLED=0 GOOS=$${GOOS} GOARCH=$${GOARCH} GOARM=$${GOARM} go build $(LDFLAGS) -o $@/$(BINARY) ./main.go
