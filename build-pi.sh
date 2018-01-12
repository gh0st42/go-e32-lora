#!/bin/sh

OUTPATH=out/pi
mkdir -p $OUTPATH

echo Building bcaster
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o $OUTPATH/bcaster cmds/bcaster/bcaster.go 

echo Building monitor
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w"  -o $OUTPATH/monitor cmds/monitor/monitor.go 

echo Building e32config
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w"  -o $OUTPATH/e32config cmds/e32config/e32config.go 

echo Building pktdump
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w"  -o $OUTPATH/pktdump cmds/pktdump/pktdump.go 

echo Building pktsend
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w"  -o $OUTPATH/pktsend cmds/pktsend/pktsend.go 

echo Building pktrecv
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w"  -o $OUTPATH/pktrecv cmds/pktrecv/pktrecv.go 

