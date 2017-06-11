FROM golang:1.8.3

ADD . /go/src/github.com/weirdtales/senor-rosado

RUN cd /go/src/github.com/weirdtales/senor-rosado && \
    go get -d ./... && \
    go test ./... && \
    go build -o /senor-rosado && \
    mkdir /plugins && \
    go build -buildmode=plugin -o /plugins/hello.so _cartridges/hello.go && \
    go build -buildmode=plugin -o /plugins/giphy.so _cartridges/giphy.go && \
    go build -buildmode=plugin -o /plugins/weather.so _cartridges/weather.go

WORKDIR /

ENTRYPOINT /senor-rosado
