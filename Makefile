ALL: lint run

.PHONY: gomod
gomod:
	rm -f go.sum && go mod tidy

.PHONY: lint
lint:
	golangci-lint run

.PHONY: run
run:
	go run example/main.go
