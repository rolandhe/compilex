#!/bin/bash

GOOS=darwin
GOARCH=amd64

cd ../../

go build -ldflags "-s -w" -o build/mac/compilex