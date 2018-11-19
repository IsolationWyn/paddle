FROM ubuntu:16.04

COPY os-requirement.sh .

RUN cp /dev/null /etc/apt/source.list \
    # && bash os-requirement.sh \
    && apt-get update -y -q \
    && apt-get upgrade -y -q \
    && mkdir -p /goroot \
    && curl https://storage.googleapis.com/golang/go1.11.1.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

ENV GOROOT /goroot
ENV GOPATH /gopath
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
WORKDIR /gopath
COPY . .
CMD ["bash"]


