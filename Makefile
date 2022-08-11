build: tools-install generate
	go mod vendor

run:
	go run cmd/main.go