all: build

build:
	@go build -o 9stream

build-arm:
	@GOOS=linux GOARCH=arm go build -o 9stream-arm