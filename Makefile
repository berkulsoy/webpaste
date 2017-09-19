BINARY=webpaste

all: clean deps test linux darwin windows docker

deps:
	go get "github.com/dustinkirkland/golang-petname"

linux:
	GOOS=linux GOARCH=amd64 go build  -o ${BINARY}-linux-amd64 . ;

darwin:
	GOOS=darwin GOARCH=amd64 go build  -o ${BINARY}-darwin-amd64 . ;

windows:
	GOOS=windows GOARCH=amd64 go build  -o ${BINARY}-windows-amd64.exe . ;

test:
	go test

docker:
	docker build -t webpaste:0.1 -t webpaste:latest .

clean:
	-rm -f ${BINARY}-*

.PHONY: clean deps test linux darwin windows docker