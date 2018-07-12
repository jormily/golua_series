set GOARCH=amd64
set GOOS=windows
set GOROOT=c:\go
set GOBIN=%GOROOT%\bin
set GOPATH=%GOPATH%;C:\Users\XXX\Desktop\gohper-lua


go install tool
go build -o bin/lua.exe src/main.go

pause.