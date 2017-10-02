DIRECTORY="github.com/fono09/sorame_bot"

bot:
	go get github.com/golang/dep/cmd/dep
	dep status
	go build --ldflags '-s -w -linkmode external -extldflags -static' -v -o bot bot.go

build:
	docker run --rm -v ${PWD}:/go/src/${DIRECTORY} -w /go/src/${DIRECTORY} golang:1.9 make bot

