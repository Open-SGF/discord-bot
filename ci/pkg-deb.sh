#!/bin/sh -e

mkdir -p /out

rm -f $PROJNAME*.deb
fpm \
    -s dir \
    -t deb \
    -v $VERSION \
    -n $PROJNAME \
    --description "Open SGF Discord bot" \
    --deb-systemd-enable \
    --deb-systemd opensgf-discord-bot.service \
    --deb-systemd-restart-after-upgrade \
    -C $BUILD_DIR \
    --deb-user opensgf \
    --deb-group opensgf \
    --after-install postinstall \
    --before-install preinstall \
    .

cp *.deb /out