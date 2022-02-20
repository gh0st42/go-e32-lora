#!/bin/sh

GO111MODULE=auto

OUTPATH=out/native
mkdir -p $OUTPATH

echo Building bcaster
go build -ldflags="-s -w" -o $OUTPATH/bcaster cmds/bcaster/bcaster.go 

echo Building monitor
go build -ldflags="-s -w"  -o $OUTPATH/monitor cmds/monitor/monitor.go 

echo Building e32config
go build -ldflags="-s -w"  -o $OUTPATH/e32config cmds/e32config/e32config.go 

echo Building pktdump
go build -ldflags="-s -w"  -o $OUTPATH/pktdump cmds/pktdump/pktdump.go 

echo Building pktsend
go build -ldflags="-s -w"  -o $OUTPATH/pktsend cmds/pktsend/pktsend.go 

echo Building pktrecv
go build -ldflags="-s -w"  -o $OUTPATH/pktrecv cmds/pktrecv/pktrecv.go 

