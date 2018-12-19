# Go parameters

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=paddle
BINARY_UNIX=$(BINARY_NAME)_unix

all: 
	test build
build: 
	$(GOBUILD) .
clean:
	$(GOCLEAN) 
	rm -f $(BINARY_NAME) 
	rm -f $(BINARY_UNIX)
run:
	./$(BINARY_NAME) run -ti bash

# Cross compilation
build-linux:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-run:
	docker run -it  -v "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle":/go/src/github.com/IsolationWyn/paddle  --privileged=true --net=bridge registry.cn-qingdao.aliyuncs.com/wisati/paddle   bash

