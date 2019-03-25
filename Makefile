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
	./$(BINARY_NAME) run -ti sh

# Cross compilation
build-linux:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
drun:
	docker run -it  -v "/Users/wyn/Library/Mobile Documents/com~apple~CloudDocs/Daocloud/GOProject/src/github.com/IsolationWyn/paddle":/go/src/github.com/IsolationWyn/paddle  --privileged=true registry.cn-qingdao.aliyuncs.com/wisati/paddle   bash
dexec:
	docker exec -it $$(docker container ls | grep paddle | awk '{split($$0,arr," ");print arr[1]}') bash