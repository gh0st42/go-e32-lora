#!/bin/sh

mkdir -p out/pi

echo Building bcaster
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o out/pi/bcaster cmds/bcaster/bcaster.go 

echo Building monitor
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w"  -o out/pi/monitor cmds/monitor/monitor.go 

echo Building e32config
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w"  -o out/pi/e32config cmds/e32config/e32config.go 
