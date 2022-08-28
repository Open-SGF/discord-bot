#!/bin/bash

mkdir -p bin
go build -o bin/opensgf-discord-bot main.go
chown -R 1000:1000 ./bin