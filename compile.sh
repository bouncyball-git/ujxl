#!/bin/sh

echo "building binaries for all platforms"
echo

export GOOS=darwin
export GOARCH=arm64
go build --ldflags "-w -s" -o ujxl-0.9.1-arm64-darwin
echo "ujxl-0.9.1-arm64-darwin"

export GOOS=darwin
export GOARCH=amd64
go build --ldflags "-w -s" -o ujxl-0.9.1-amd64-darwin
echo "ujxl-0.9.1-amd64-darwin"

export GOOS=windows
export GOARCH=386
go build --ldflags "-w -s" -o ujxl-0.9.1-386-windows.exe
echo "ujxl-0.9.1-386-windows.exe"

export GOOS=windows
export GOARCH=amd64
go build --ldflags "-w -s" -o ujxl-0.9.1-amd64-windows.exe
echo "ujxl-0.9.1-amd64-windows.exe"

export GOOS=linux
export GOARCH=386
go build --ldflags "-w -s" -o ujxl-0.9.1-386-linux
echo "ujxl-0.9.1-386-linux"

export GOOS=linux
export GOARCH=amd64
go build --ldflags "-w -s" -o ujxl-0.9.1-amd64-linux
echo "ujxl-0.9.1-amd64-linux"
