run:
	@go run main.go

test:
	@go test ./...
t: test

build:
	@go build -o go-rest-api-blog main.go

lint:
	@go vet ./...
	@gofmt -d -w .
fix: lint

clean:
	@rm -f go-rest-api-blog

