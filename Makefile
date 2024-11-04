build:
	go build -o app

lint:
	golangci-lint run

test:
	go test ./...

run:
	go run main.go
