build:
	@go build -o bin/roundtable

run: build
	@./bin/roundtable

test:
	@go test -v ./...
