#!/bin/bash
if [ ! -d "./pb" ]; then
    mkdir ./pb
fi

rm -rf ./pb/*
mkdir ./pb/tr477 ./pb/control_relay

#used protoc version 3.12.4
# protoc --proto_path=proto --go_out=./pb/tr451 --go_opt=paths=source_relative --go-grpc_out=./pb/tr451 --go-grpc_opt=paths=source_relative ./proto/tr451*

protoc --proto_path=proto --go_out=./pb/tr477 --go_opt=paths=source_relative --go-grpc_out=./pb/tr477 --go-grpc_opt=paths=source_relative ./proto/tr477*

protoc --proto_path=proto --go_out=./pb/control_relay --go_opt=paths=source_relative --go-grpc_out=./pb/control_relay --go-grpc_opt=paths=source_relative ./proto/control_relay*