BINARY_NAME=products-api

build:
	go build -o bin/${BINARY_NAME} main.go

run:
	go run main.go

dep: 
	go mod download

tidy:
	go mod tidy

lint: 
	golangci-lint run --enable-all
 	
