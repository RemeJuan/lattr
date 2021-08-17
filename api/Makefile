GO=go
GOCOVER=$(GO) tool cover

run:
	go run main.go

test:
	go test ./...

coverage:
	go test ./... -race -covermode=atomic -coverprofile=coverage.out
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out

