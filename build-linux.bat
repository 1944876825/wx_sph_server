@echo off
set CGO_ENABLED=1
set GOOS=linux
set GOARCH=amd64
set CC=zig cc -target x86_64-linux
SET CXX=zig c++ -target x86_64-linux  
go build -o wx_video_help -tags netgo -ldflags "-s -w"