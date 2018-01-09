package lora

import (
	"io"
	"log"
	"os"

	"github.com/jacobsa/go-serial/serial"
)

func GetSerial() io.ReadWriteCloser {
	serialport := "/dev/ttyS0"
	if os.Getenv("LORAPORT") != "" {
		serialport = os.Getenv("LORAPORT")
	}
	options := serial.OpenOptions{
		//PortName:        "/dev/tty.SLAB_USBtoUART",
		//PortName:        "/dev/ttyUSB0",
		PortName:        serialport,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
		os.Exit(1)
	}
	// Make sure to close it later.
	// defer port.Close()
	return port
}
