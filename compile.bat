echo building binaries for all platforms

@set GOOS=linux
@set GOARCH=386
go build --ldflags "-w -s" -o ujxl-0.9.1-386-linux

@set GOOS=linux
@set GOARCH=amd64
go build --ldflags "-w -s" -o ujxl-0.9.1-amd64-linux

@set GOOS=darwin
@set GOARCH=arm64
go build --ldflags "-w -s" -o ujxl-0.9.1-arm64-darwin

@set GOOS=darwin
@set GOARCH=amd64
go build --ldflags "-w -s" -o ujxl-0.9.1-amd64-darwin

@set GOOS=windows
@set GOARCH=386
go build --ldflags "-w -s" -o ujxl-0.9.1-386-windows.exe

@set GOOS=windows
@set GOARCH=amd64
go build --ldflags "-w -s" -o ujxl-0.9.1-amd64-windows.exe
