build:
	@go build -o bin/cmd/api main.go
run: build
	@./bin/api
tests:
	@go test ./test/... -v
doc:
	@swag init --output etc/doc