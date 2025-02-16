.SILENT:

run:
	clear
	go run ./cmd/main.go | pplog

lint:
	golangci-lint run