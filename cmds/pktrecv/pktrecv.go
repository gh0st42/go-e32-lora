package main

import (
	"fmt"
	"github.com/gh0st42/go-e32-lora/lora"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func cleanup() {
	println("Cleanup")
}
func OnRecv(inp chan []byte) {
	for {
		select {
		case data := <-inp:
			fmt.Printf("incoming (%d): %X\n", len(data), data)
		}
	}
}

func main() {
	recvchan := make(chan []byte)
	go OnRecv(recvchan)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	fmt.Println("-= LoraMon - PktRecv =-")
	fmt.Println("Waiting 4 packets")

	lora.RecvLoop(recvchan)
}
