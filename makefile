# build things with docker cross compiling for arm linux

GOOS=linux
GOARCH=arm


all: hello

hello: src/hello/hello.go
	docker run -it --rm -v $(PWD):/darpi -w /darpi -e GOOS=${GOOS} -e GOARCH=${GOARCH} golang:1.6 go build src/hello/hello.go
sftp-hello:
	@echo sftp pi@darpi

get:
	docker run -it --rm -v $(PWD):/darpi -w /darpi golang:1.6 go get src/darpi/.go

build:
	docker run -it --rm -v $(PWD):/darpi -w /darpi -e GOOS=${GOOS} -e GOARCH=${GOARCH} golang:1.6 go build src/darpi/hello.go
