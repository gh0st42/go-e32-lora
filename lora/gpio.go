package lora

import (
	"fmt"
	"log"
	"os"
	"runtime"

	rpio "github.com/stianeikeland/go-rpio"
)

func IsPiInternal() bool {
	if runtime.GOARCH != "arm" {
		return false
	}
	lp := os.Getenv("LORAPORT")
	if lp == "" || lp == "/dev/ttyS0" {
		return true
	}
	return false

}
func GpioSet(p1, p2, level uint8) {
	if !IsPiInternal() {
		fmt.Println("WARNING: External mode setting needed!")
		return
	}
	err := rpio.Open()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	m0 := rpio.Pin(p1)
	m1 := rpio.Pin(p2)

	m0.Output()
	m1.Output()

	if level == 0 {
		m0.Low()
		m1.Low()
	} else {
		m0.High()
		m1.High()
	}
	Wait()
	rpio.Close()
}
func GpioGet(p1, p2 uint8) uint8 {
	if !IsPiInternal() {
		fmt.Println("WARNING: External mode setting needed!")
		return 255
	}
	err := rpio.Open()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	m0 := rpio.Pin(p1)
	m1 := rpio.Pin(p2)

	m0.Input()
	m1.Input()

	res := uint8(m0.Read()) + uint8(m1.Read())
	Wait()

	rpio.Close()
	return res
}
