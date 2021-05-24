#!/bin/sh

mkdir -p dist
go build -o dist -ldflags "-s -w" ./...