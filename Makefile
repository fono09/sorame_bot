DIRECTORY="github.com/fono09/sorame_bot"

bot:
	go get -v github.com/tools/godep
	godep restore -v
	go build --ldflags '-s -w -linkmode external -extldflags -static' -v -o bot bot.go

build:
	docker run --rm -v ${PWD}:/go/src/${DIRECTORY} -w /go/src/${DIRECTORY} golang:1.9 make bot
