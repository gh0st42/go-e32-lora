package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/gh0st42/go-e32-lora/lora"
	"log"
	"os"
	"strconv"

	dialog "github.com/weldpua2008/go-dialog"
)

func dialog_checkErr(err error) {
	if err != nil {
		d := dialog.New(dialog.AUTO, 0)
		d.SetBackTitle("E32-TTL-100 Configuration Wizard")
		d.Msgbox("Configuration wizard canceled!\nHave a nice day...")
		os.Exit(0)
	}
}
func dialog_new() {
	var cfg lora.E32config

	d := dialog.New(dialog.AUTO, 0)
	//d.Msgbox("Hello world!")
	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("Store configuration")
	res, err := d.Radiolist(0, "C0", "Permanent", "on", "C2", "Temporary", "off")
	dialog_checkErr(err)
	h, err := hex.DecodeString(res)
	cfg.Head = h[0]

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("2-byte node address - hex encoded")
	res, err = d.Inputbox("FFFF")
	dialog_checkErr(err)
	n, err := strconv.ParseUint(res, 16, 16)
	cfg.Addr = uint16(n)

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("Port Configuration")
	res, err = d.Radiolist(0, "0", "8N1", "on", "1", "8O1", "off", "2", "8E1", "off", "3", "8N1", "off")
	dialog_checkErr(err)
	n, err = strconv.ParseUint(res, 10, 8)
	cfg.Parity = uint8(n) << 6

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("Serial Rate")
	res, err = d.Radiolist(0, "0", "1200", "off", "1", "2400", "off",
		"2", "4800", "off", "3", "9600 (Default)", "on",
		"4", "19200", "off", "5", "38400", "off",
		"6", "57600", "off", "7", "115200", "off",
	)
	dialog_checkErr(err)
	n, err = strconv.ParseUint(res, 10, 8)
	cfg.Ttlrate = uint8(n) << 3

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("Wireless Rate")
	res, err = d.Radiolist(0, "0", "1K", "on", "1", "2K", "off",
		"2", "5K", "off", "3", "8K", "off",
		"4", "10K", "off", "5", "15K", "off",
		"6", "20K", "off", "7", "25K", "off",
	)
	dialog_checkErr(err)
	n, err = strconv.ParseUint(res, 10, 8)
	cfg.Wirelessrate = uint8(n)

	// TODO add slider
	cfg.Channel = 0x50

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("Select Frequency")
	res, err = d.Inputbox("433")
	dialog_checkErr(err)
	f, err := strconv.ParseFloat(res, 10)
	cfg.Channel = uint8((f - 425) / 0.1)

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("Radio Mode")
	res, err = d.Radiolist(0, "0", "Transparent (Default)", "on", "1", "FSK", "off")
	dialog_checkErr(err)
	n, err = strconv.ParseUint(res, 10, 8)
	cfg.Transmissionmode = uint8(n) << 7

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("IO Mode")
	res, err = d.Radiolist(0, "0", "Open", "off", "1", "Push-Pull (Default)", "on")
	dialog_checkErr(err)
	n, err = strconv.ParseUint(res, 10, 8)
	cfg.Iomode = uint8(n) << 6

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("Wakeup Time")
	res, err = d.Radiolist(0, "0", "250ms (Default)", "on", "1", "500ms", "off",
		"2", "750ms", "off", "3", "1000ms", "off",
		"4", "1250ms", "off", "5", "1500ms", "off",
		"6", "1750ms", "off", "7", "2000ms", "off",
	)
	dialog_checkErr(err)
	n, err = strconv.ParseUint(res, 10, 8)
	cfg.Wakeuptime = uint8(n) << 3

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("Forward Error Correction")
	res, err = d.Radiolist(0, "0", "off", "off", "1", "on (Default)", "on")
	dialog_checkErr(err)
	n, err = strconv.ParseUint(res, 10, 8)
	cfg.Fec = uint8(n) << 2

	d.SetBackTitle("E32-TTL-100 Configuration Wizard")
	d.SetTitle("Transmit Power")
	res, err = d.Radiolist(0, "0", "30dBm (Default)", "on", "1", "27dBm", "off",
		"2", "24dBm", "off", "3", "21dBm", "off",
		"4", "18dBm", "off", "5", "15dBm", "off",
		"6", "12dBm", "off", "7", "9dBm", "off",
	)
	dialog_checkErr(err)
	n, err = strconv.ParseUint(res, 10, 8)
	cfg.Power = uint8(n)

	cfg.PrintConfig()
}

func PortMode(mode string) {
	if mode == "0" {
		fmt.Println("Setting lora modem to transmission mode.")
		lora.GpioSet(17, 27, 0)
	} else if mode == "3" {
		fmt.Println("Setting lora modem to configuration mode.")
		lora.GpioSet(17, 27, 1)
	} else {
		log.Fatalln("Unknown mode! \n 0 = transmission\n 3 = config\n")
		os.Exit(1)
	}
}

func ReadConfig() {
	println("Reading config..")
	var mode uint8
	if lora.IsPiInternal() {
		mode = lora.GpioGet(17, 27)
		if mode == 0 {
			lora.GpioSet(17, 27, 1)
		}
	} else {
		fmt.Println("Please ensure that config mode is set manually")
	}

	port := lora.GetSerial()
	// Make sure to close it later.
	defer port.Close()
	_, err := port.Write([]byte{0xC1, 0xC1, 0xC1})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	r := bufio.NewReader(port)
	buf := make([]byte, 6)
	b, _ := r.ReadByte()
	buf[0] = b
	b, _ = r.ReadByte()
	buf[1] = b
	b, _ = r.ReadByte()
	buf[2] = b
	b, _ = r.ReadByte()
	buf[3] = b
	b, _ = r.ReadByte()
	buf[4] = b
	b, _ = r.ReadByte()
	buf[5] = b
	fmt.Printf("Current Config:\n%X\n", buf)
	cfg := lora.ParseHex(hex.EncodeToString(buf))
	cfg.PrintConfig()
	if lora.IsPiInternal() {
		if mode == 0 {
			lora.GpioSet(17, 27, 0)
		}
	}
}
func ApplyConfig(c string) {
	if len(c) == 2*6 {
		var mode uint8
		if lora.IsPiInternal() {
			mode = lora.GpioGet(17, 27)
			if mode == 0 {
				lora.GpioSet(17, 27, 1)
			}
			lora.Wait()
		} else {
			fmt.Println("Please ensure that config mode is set manually")
		}

		decoded, err := hex.DecodeString(c)
		if err != nil {
			log.Fatal(err)
		}
		if len(decoded) != 6 {
			log.Fatal("Hex string too short for e32 config!")
		}
		port := lora.GetSerial()
		// Make sure to close it later.
		defer port.Close()

		fmt.Printf("Writing config (%s) ..\n", c)
		_, err = port.Write(decoded)
		if err != nil {
			log.Fatal(err)
		}
		lora.Wait()

		if lora.IsPiInternal() {
			if mode == 0 {
				lora.GpioSet(17, 27, 0)
			}
		}

	} else {
		log.Fatalln("Invalid config string!")
	}
}
func usage() {
	fmt.Printf("usage: %s <config_hex_string> - print decoded configuration\n", os.Args[0])
	fmt.Printf("       %s read - read configuration\n", os.Args[0])
	fmt.Printf("       %s apply <config_hex_string> - apply configuration\n", os.Args[0])
	fmt.Printf("                Example config: C2FFFF185078\n")
	fmt.Printf("       %s new - generate new configuration\n", os.Args[0])
	fmt.Printf("       %s mode <0|3> - set lora module to transmission or config mode\n", os.Args[0])
	fmt.Printf("                     GPIO 17 and 27 must be wired to m0 and m1 (pi only!)\n")
	fmt.Printf("       %s help - display this usage information\n", os.Args[0])
}
func main() {
	fmt.Println("E32-TTL-100 Config\n")

	if len(os.Args) >= 2 {
		if len(os.Args[1]) == 6*2 {
			cfg := lora.ParseHex(os.Args[1])
			cfg.Print()
			cfg.PrintConfig()
		} else if os.Args[1] == "help" {
			usage()
		} else if os.Args[1] == "new" {
			dialog_new()
		} else if os.Args[1] == "mode" {
			PortMode(os.Args[2])
		} else if os.Args[1] == "apply" {
			ApplyConfig(os.Args[2])
		} else if os.Args[1] == "read" {
			ReadConfig()
		} else {
			usage()
		}
	} else {
		usage()
	}

}
