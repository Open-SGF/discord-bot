#!/usr/bin/env -S bash -e

V=$1

echo "Releasing to $V"

sed -i "s/VERSION=.*$/VERSION=$V/" .env

git add .env
git commit -m "Release v$V"
git tag -a v$V -m "Release v$V"

docker-compose build build-amd64
docker-compose build pkg-debian
docker-compose run pkg-debian

dpkg --contents ../out/opensgf-discord-bot_${V}_amd64.deb