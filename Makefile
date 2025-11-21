run:
	@go run main.go

test:
	@go test ./...

build:
	@go build -o go-rest-api-blog main.go

clean:
	@rm -f go-rest-api-blog

