FROM golang:1.10.1-alpine

RUN apk add --no-cache git

WORKDIR /go/src/gitlab.oit.duke.edu/rz78/ups
COPY . .

RUN go get -v ./cmd/ups_server ./cmd/web_server
