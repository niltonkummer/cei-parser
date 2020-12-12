
build-cli:
	@go mod tidy
	@go build -o ./bin/cei-parser ./bin
