package main

import (
	"bufio"
	"fmt"
	"go-e32-lora/lora"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lobaro/slip"
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
func main() {

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	port := lora.GetSerial()

	// Make sure to close it later.
	defer port.Close()

	fmt.Println("-= LoraMon - PktDump=-")
	fmt.Println("Waiting 4 packets")

	r := bufio.NewReader(port)
	for {
		res, err := r.ReadByte()
		checkErr(err)
		fmt.Printf("%02X ", res)
		if res == slip.END {
			println()
		}
	}
}
