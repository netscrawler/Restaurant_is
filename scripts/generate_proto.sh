#!/bin/bash

PROTO_ROOT=./proto/v1
GEN_DIR=./gen
GO_OPTS="paths=source_relative"

if [ ! -d "./third_party" ]; then
    git clone --depth 1 https://github.com/protocolbuffers/protobuf.git third_party
fi

protoc -I=$PROTO_ROOT \
    -I=./third_party/protobuf/src \
    --go_out=$GEN_DIR \
    --go_opt=$GO_OPTS \
    --go-grpc_out=$GEN_DIR \
    --go-grpc_opt=$GO_OPTS \
    $PROTO_ROOT/*.proto
