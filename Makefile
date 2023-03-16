build:
	@go build -o sqlgostructure

run:
	@go run .

test:
	@go test -v ./...

testCover:
	@go test -v ./... -cover

dockerize:
	@docker build -t sqlgostructure:latest .

dockerrun:
	@docker run --name sqlgostructure -p 4000:4000 sqlgostructure:latest