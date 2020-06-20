go mod init github.com/hitzhangjie/gorpc-cli

go build -o gorpc.exe

copy gorpc.exe %GOPATH%\bin /y

mkdir C:\users\%username%\.gorpc

xcopy install C:\users\%username%\.gorpc\ /s /e /y /d

