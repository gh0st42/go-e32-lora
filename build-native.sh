#!/bin/sh

mkdir -p out/native

echo Building bcaster
go build -ldflags="-s -w" -o out/native/bcaster cmds/bcaster/bcaster.go 

echo Building monitor
go build -ldflags="-s -w"  -o out/native/monitor cmds/monitor/monitor.go 

echo Building e32config
go build -ldflags="-s -w"  -o out/native/e32config cmds/e32config/e32config.go 
