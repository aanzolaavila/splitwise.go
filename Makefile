.PHONY: setup
setup:
	go mod download && go mod verify

.PHONY: tests
test: setup
	go test .

.PHONY: coverage
coverage: setup
	go test -v -cover -coverprofile=coverage.out .
	go tool cover -html=coverage.out
	rm -f coverage.out

.PHONY: examples
examples: setup
	cd examples; \
	go run run.go
