build:
	@go build -o sqlgostructure

run:
	@go run .

test:
	@go test -v ./...

testCover:
	@go test -v ./... -cover

swagger:
	@"$HOME/go/bin/swag init -g ./src/routes/routes.go"

dockerize:
	@docker build -t sqlgostructure:latest .

dockerrun:
	@docker run --name sqlgostructure -p 4000:4000 sqlgostructure:latest