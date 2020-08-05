ALL: lint run

.PHONY: lint
lint:
	golangci-lint run

.PHONY: run
run:
	go run example/main.go
