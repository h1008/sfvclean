FROM golang:1.16-buster

RUN go get -v golang.org/x/tools/gopls && \
    go get -v github.com/go-delve/delve/cmd/dlv