#!/bin/bash

me=`realpath $0`

cd `dirname $me`/..
pwd

export GOPATH=`dirname $me`/../../../
echo $GOPATH

export GO111MODULE="on"

which protoc || exit 1

go mod tidy

go build -C third_party/protobuf-go/cmd/protoc-gen-go -v -o ../../../../build/protoc-gen-go . || exit 1

go build -C third_party/grpc-go/cmd/protoc-gen-go-grpc -v -o ../../../../build/protoc-gen-go-grpc . || exit 1

#go build -C third_party/nats-server -v -o ../../build/nats-server main.go || exit 1

./third_party/mockery/mockery

gofmt . >/dev/null || exit 1

find . -type f -name "*.go" | grep -v third_party | xargs gofmt -d
find . -type f -name "*.go" | grep -v third_party | xargs gofmt -w

protoc -I proto proto/*.proto --go_out=./gen/ --plugin=build/protoc-gen-go || exit 1

protoc -I proto proto/*.proto --go-grpc_out=./gen/ --plugin=build/protoc-gen-go-grpc || exit 1

go test -coverprofile=/tmp/coverage.out -count=2 ./... || exit 1
go tool cover -html=/tmp/coverage.out -o coverage.html

TOOLS="sender grpc-receiver metrics"
for tool in $TOOLS; do
    echo "Build $tool"
    go build -v -o "build/$tool" cmd/$tool/main.go || exit 1
done

tree -I third_party
