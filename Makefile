BINARY_NAME_ATLAS=atlas

all: build-atlas

build-atlas:
	go build -o ./bin/$(BINARY_NAME_ATLAS) ./cmd/atlas/

test:
	go test ./...

clean:
	rm -f ./bin/$(BINARY_NAME_ATLAS)

run-atlas:
	go run ./cmd/atlas

fmt:
	go fmt ./...
