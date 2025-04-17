.PHONY: build clean dev

build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/app cmd/server/main.go

clean:
	rm -rf bin/ tmp/

dev:
	go run github.com/air-verse/air@latest