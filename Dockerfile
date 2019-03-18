FROM ubuntu:16.04
MAINTAINER Zi Wang  <isolationwyn@gmail.com>

RUN apt-get update -y -q && apt-get upgrade -y -q \
    && DEBIAN_FRONTEND=noninteractive \
    && apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git \
    aufs-tools linux-image-extra-virtual psmisc criu vim sudo tree net-tools\
    && curl -s https://storage.googleapis.com/golang/go1.11.1.linux-amd64.tar.gz| tar -v -C /usr/local -xz\
    && go get -u github.com/go-delve/delve/cmd/dlv


EXPOSE 2345


USER root
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $PATH:/usr/local/go/bin
WORKDIR ${GOPATH}

CMD ["dlv", "debug", "--headless", "--listen=:2345", "--log"]