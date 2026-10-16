.PHONY: run test build docker clean

run:
	go run ./cmd/server

test:
	go test -race ./...

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server ./cmd/server

docker:
	docker build -t hidden-coast-backend:latest .

clean:
	rm -rf bin/
