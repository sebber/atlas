BINARY_NAME_ATLAS=atlas
BINARY_NAME_TERMINAL=terminal

all: build-atlas build-terminal

build-atlas:
	go build -o ./bin/$(BINARY_NAME_ATLAS) ./cmd/atlas/

build-terminal:
	go build -o ./bin/$(BINARY_NAME_TERMINAL) ./cmd/terminal/

test:
	go test ./...

clean:
	rm -f ./bin/$(BINARY_NAME_ATLAS) ./bin/$(<BINARY_NAME_TERMINAL)

run-atlas:
	go run ./cmd/atlas

run-terminal:
	go run ./cmd/terminal

fmt:
	go fmt ./...
