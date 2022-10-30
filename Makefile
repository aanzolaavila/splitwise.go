.PHONY: setup
setup:
	go mod download && go mod verify

.PHONY: tests
test: setup
	go test .

.PHONY: coverage
coverage: setup
	go test -v -cover -coverprofile=c.out .
	go tool cover -html=c.out
	rm -f c.out

.PHONY: examples
examples: setup
	cd examples; \
	go run run.go
