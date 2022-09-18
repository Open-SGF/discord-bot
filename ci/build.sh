#!/bin/sh -e

BIN_DIR=$BUILD_DIR/opt/$PROJNAME/bin

rm -rf $BUILD_DIR
mkdir -p $BIN_DIR $SYSTEMD_DIR
cp example.config.json $BIN_DIR/config.json

go mod tidy
CGO_ENABLED=0 go build -o $BIN_DIR/$PROJNAME main.go