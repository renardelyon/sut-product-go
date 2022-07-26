FROM golang:alpine AS builder

WORKDIR /workspace
ENV GO111MODULE=on
ARG SERVICE
ARG PORT

RUN apk update && apk add --no-cache protobuf bash ca-certificates git gcc g++ libc-dev curl 

COPY . .

RUN /usr/local/go/bin/go install github.com/golang/protobuf/protoc-gen-go@latest
RUN mkdir pb && protoc --proto_path=proto/ --go_out=paths=source_relative,plugins=grpc:./pb proto/*/*.proto
RUN /usr/local/go/bin/go mod init ${SERVICE} && /usr/local/go/bin/go mod tidy && mkdir /tmp/result && /usr/local/go/bin/go get github.com/gin-gonic/gin@v1.7.7
RUN /usr/local/go/bin/go build -o . ./cmd/main.go

EXPOSE ${PORT}

CMD ./main