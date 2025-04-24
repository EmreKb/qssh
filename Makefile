.PHONY: build

build:
	go build -o qssh cmd/qssh/main.go

.PHONY: run
run:
	go run cmd/qssh/main.go
