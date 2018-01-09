package main

import (
	"bytes"
	"fmt"
	"go-e32-lora/lora"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	//"go.bug.st/serial.v1"
)

func getMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}

func main() {
	interval := 60 * time.Second
	var counter uint64 = 0

	if len(os.Args) == 2 {
		num, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalln("Error parsing interval!\n")
			os.Exit(1)
		}
		interval = time.Second * time.Duration(num)
	}
	fmt.Printf("Lora Static Broadcaster\n broadcast interval: %v\n\n", interval)

	port := lora.GetSerial()

	// Make sure to close it later.
	defer port.Close()

	hostname, err := os.Hostname()
	if err != nil {
		hostname = getMacAddr()
	}
	for {
		statusstring := hostname + " " + fmt.Sprintf("#%v", counter) + " " + time.Now().Format(time.RFC850) + "\n"
		output := []byte(statusstring)
		fmt.Printf("SENDING: %v", statusstring)
		n, err := port.Write(output)
		if err != nil {
			log.Fatalf("port.Write: %v", err)
		}

		fmt.Println("Wrote", n, "bytes.")
		counter++
		time.Sleep(interval)
	}
}
