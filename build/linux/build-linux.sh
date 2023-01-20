#!/bin/bash

GOOS=linux
GOARCH=amd64

cd ../../

go build -ldflags "-s -w" -o build/linux/compilex