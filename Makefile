.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: test
## test: runs `go test`
test:
	@go test -race ./lichess/... -coverprofile cover.out

.PHONY: lint
## lint: runs `golangci-lint`
lint:
	@golangci-lint run ./lichess/...

