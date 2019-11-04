package main

import (
	"bufio"
	"fmt"
	"../../lora"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	device "github.com/d2r2/go-hd44780"
	i2c "github.com/d2r2/go-i2c"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

var lcd *device.Lcd
var isPi bool

func printLog(text string, line device.ShowOptions) {
	fmt.Println(text)
	if isPi {
		err := lcd.ShowMessage(text, line)
		checkErr(err)
	}
}
func cleanup() {
	println("Cleanup")
	// Turn off the backlight and exit
	err := lcd.BacklightOff()
	checkErr(err)
}
func main() {
	isPi = lora.IsPiInternal()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if isPi { // cleanup lcd
			cleanup()
		}
		os.Exit(1)
	}()

	fmt.Println("Lora Monitor")
	port := lora.GetSerial()

	// Make sure to close it later.
	defer port.Close()

	i2c, err := i2c.NewI2C(0x3f, 1)
	if err != nil {
		log.Fatal(err)
	}
	// Free I2C connection on exit
	defer i2c.Close()
	// Construct lcd-device connected via I2C connection

	lcd, err = device.NewLcd(i2c, device.LCD_16x2)
	checkErr(err)

	// Turn on the backlight
	err = lcd.BacklightOn()
	checkErr(err)

	lcd.Clear()

	printLog("-= LoraMon =-", device.SHOW_LINE_1)
	printLog("Waiting 4 beacon", device.SHOW_LINE_2)

	// Wait 5 secs
	time.Sleep(5 * time.Second)
	lcd.BacklightOff()

	r := bufio.NewReader(port)
	counter := 0
	for {
		res, err := r.ReadString('\n')
		checkErr(err)
		fields := strings.Fields(res)
		println(res)
		lcd.BacklightOn()
		line1 := "L: " + fields[0] + " " + fields[1]
		printLog(line1, device.SHOW_LINE_1)

		line2 := " " + fields[3] + " " + fields[4]
		printLog(line2, device.SHOW_LINE_2)

		counter++
		go func() {
			time.Sleep(10 * time.Second)
			lcd.BacklightOff()
		}()
		if counter == -1 {
			break
		}
	}
}
