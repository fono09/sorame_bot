DIRECTORY="github.com/fono09/sorame_bot"


dep:
	go get github.com/golang/dep/cmd/dep

bot: dep_status
	go build -v -o bot bot.go

dep_status: dep
	dep status

dep_ensure: dep
	dep ensure -v

depends:
	docker run --rm -v ${PWD}:/go/src/${DIRECTORY} -w /go/src/${DIRECTORY} golang:1.9 make dep_ensure

build:
	docker run --rm -v ${PWD}:/go/src/${DIRECTORY} -w /go/src/${DIRECTORY} golang:1.9 make bot


