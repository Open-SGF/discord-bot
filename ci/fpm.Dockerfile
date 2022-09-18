FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get -y install \
        binutils \
        rubygems \
        squashfs-tools && \
    gem install fpm

COPY ./ci/linux/ /work
COPY ./ci/pkg-deb.sh /work
COPY ./build /work/build

WORKDIR /work