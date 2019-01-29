FROM golang:1.11

ARG VERSION

WORKDIR /go/src/github.com/laszlocph/tsdbinfo/

ADD . /go/src/github.com/laszlocph/tsdbinfo/

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN dep ensure

RUN env GOOS=darwin GOARCH=amd64 go build \
  && mv tsdbinfo tsdbinfo-$VERSION-mac-amd64

RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
  && cp tsdbinfo tsdbinfo-$VERSION-linux-amd64
