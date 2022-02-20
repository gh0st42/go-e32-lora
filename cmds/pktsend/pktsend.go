package main

import (
	"github.com/gh0st42/go-e32-lora/lora"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	println("send pkt")

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <data>\n", os.Args[0])
	}
	data := []byte(os.Args[1])

	println("DATA LEN:", len(data))

	buf := lora.WriteBytes(data)

	port := lora.GetSerial()
	// Make sure to close it later.
	defer port.Close()

	_, err := port.Write(buf)
	if err != nil {
		log.Fatal(err)
	}
	for _, _ = range buf {
		lora.Wait()
	}
}
