#!/bin/bash

PROJNAME=opensgf-discord-bot
ROOT_DIR=build
BIN_DIR=$ROOT_DIR/opt/$PROJNAME/bin

rm -rf $ROOT_DIR
mkdir -p $BIN_DIR $SYSTEMD_DIR
cp example.config.json $BIN_DIR/config.json

CGO_ENABLED=0 go build -o $BIN_DIR/$PROJNAME main.go

rm -f $PROJNAME*.deb
fpm \
    -s dir \
    -t deb \
    -v 0.1.0 \
    -n $PROJNAME \
    --description "Open SGF Discord bot" \
    --deb-systemd-enable \
    --deb-systemd opensgf-discord-bot.service \
    --deb-systemd-restart-after-upgrade \
    -C $ROOT_DIR \
    --deb-user opensgf \
    --deb-group opensgf \
    --after-install postinstall \
    --before-install preinstall \
    .