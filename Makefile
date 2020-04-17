# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=quanta_lab_aip
BINARY_WIN=$(BINARY_NAME).exe
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
	$(GOBUILD) -o ./bin/$(BINARY_WIN) -v
tests:
	${GOTEST} -v ./...
clean:
	$(GOCLEAN)
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)
cleanw:
	$(GOCLEAN)
	del bin\$(BINARY_WIN)
run:
	$(GOBUILD) -o ./bin/$(BINARY_WIN) -v
	./bin/$(BINARY_WIN)
deps:
	$(GOGET) github.com/gin-gonic/gin
	$(GOGET) github.com/sirupsen/logrus
	$(GOGET) github.com/pborman/uuid
	$(GOGET) github.com/stretchr/testify
	$(GOGET) github.com/joho/godotenv
	$(GOGET) github.com/robfig/cron
	$(GOGET) gopkg.in/mgo.v2/bson
	$(GOGET) github.com/golang/protobuf
	$(GOGET) github.com/google/uuid
	$(GOGET) github.com/gorilla/websocket
	$(GOGET) github.com/mattn/go-isatty
	$(GOGET) github.com/ugorji/go/codec

# todo: Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
