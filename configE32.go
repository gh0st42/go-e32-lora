package lora

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	ConfigPermanent = 0xC0
	ConfigTemporary = 0xC2

	SerialParity8N1  = 0 << 6
	SerialParity8O1  = 1 << 6
	SerialParity8E1  = 2 << 6
	SerialParity8N13 = 3 << 6

	SerialRate1200   = 0 << 3
	SerialRate2400   = 1 << 3
	SerialRate4800   = 2 << 3
	SerialRate9600   = 3 << 3
	SerialRate19200  = 4 << 3
	SerialRate38400  = 5 << 3
	SerialRate57600  = 6 << 3
	SerialRate115200 = 7 << 3

	WirelessRate1K  = 0
	WirelessRate2K  = 1
	WirelessRate5K  = 2
	WirelessRate8K  = 3
	WirelessRate10K = 4
	WirelessRate15K = 5
	WirelessRate20K = 6
	WirelessRate25K = 7

	Channel433 = 0x17

	TransmissionModeTransparent = 0 << 7
	TransmissionModeFixed       = 1 << 7

	IoModePushPull = 1 << 6
	IoModeOpen     = 0 << 6

	WakeUp250ms  = 0 << 3
	Wakeup500ms  = 1 << 3
	Wakeup750ms  = 2 << 3
	Wakeup1000ms = 3 << 3
	Wakeup1250ms = 4 << 3
	Wakeup1500ms = 5 << 3
	Wakeup1750ms = 6 << 3
	Wakeup2000ms = 7 << 3

	FecOff = 0 << 2
	FecOn  = 1 << 2

	Power30dBm = 0
	Power27dBm = 1
	Power24dBm = 2
	Power21dBm = 3
)

// E32config e32-ttl-100 configuration
type E32config struct {
	Head             byte
	Addr             uint16
	Parity           byte
	Ttlrate          byte
	Wirelessrate     byte
	Channel          byte
	Transmissionmode byte
	Iomode           byte
	Wakeuptime       byte
	Fec              byte
	Power            byte
}

func NewE32Config() E32config {
	var c E32config
	c.Head = ConfigTemporary
	c.Addr = 0xFFFF
	c.Parity = SerialParity8N1
	c.Ttlrate = SerialRate9600
	c.Wirelessrate = WirelessRate1K
	c.Channel = Channel433
	c.Transmissionmode = TransmissionModeTransparent
	c.Iomode = IoModePushPull
	c.Wakeuptime = WakeUp250ms
	c.Fec = FecOn
	c.Power = Power30dBm
	return c
}
func ParseHex(hexinput string) E32config {
	hexstr := strings.Replace(hexinput, " ", "", -1)
	decoded, err := hex.DecodeString(hexstr)
	if err != nil {
		log.Fatal(err)
	}
	if len(decoded) != 6 {
		log.Fatal("Hex string to short for e32 config!")
		os.Exit(1)
	}
	c := NewE32Config()
	c.Head = decoded[0]
	c.Addr = uint16(decoded[2]) | uint16(decoded[1])<<8

	c.Parity = (decoded[3]&(1<<7) | decoded[3]&(1<<6))
	c.Ttlrate = (decoded[3]&(1<<5) | decoded[3]&(1<<4) | decoded[3]&(1<<3))
	c.Wirelessrate = decoded[3]&(1<<2) | decoded[3]&(1<<1) | decoded[3]&(1<<0)
	c.Channel = decoded[4]

	c.Transmissionmode = (decoded[5] & (1 << 7))
	c.Iomode = (decoded[5] & (1 << 6))
	c.Wakeuptime = (decoded[5]&(1<<5) | decoded[5]&(1<<4) | decoded[5]&(1<<3))
	c.Fec = (decoded[5] & (1 << 2))
	c.Power = (decoded[5]&(1<<1) | decoded[5]&(1<<0))
	return c
}
func (config *E32config) Bytes() [6]byte {
	var c [6]byte
	c[0] = config.Head
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, config.Addr)
	c[1] = b[0]
	c[2] = b[1]
	c[3] = config.Parity | config.Ttlrate | config.Wirelessrate
	c[4] = config.Channel
	c[5] = config.Transmissionmode | config.Iomode | config.Wakeuptime | config.Fec | config.Power
	return c
}

func (config *E32config) Print() {
	cfgbytes := config.Bytes()

	for _, b := range cfgbytes {
		fmt.Printf("%02X %08b\n", b, b)
	}
	fmt.Println("\nConfig HEX bytes:")
	fmt.Printf("%X\n", cfgbytes)
}

func (config *E32config) PrintConfig() {
	SerialParity := []string{"8N1", "8O1", "8E1", "8N13"}
	SerialRate := []int{1200, 2400, 4800, 9600, 19200, 38400, 57600, 115200}
	WirelessRate := []int{1, 2, 5, 8, 10, 15, 20, 25}
	TranmissionMode := []string{"transparent", "fixed"}
	IoMode := []string{"open", "push-pull"}
	c := config.Bytes()
	fmt.Println("\nConfig HEX bytes:")
	fmt.Printf("%X\n", c)

	if config.Head == 0xC0 {
		fmt.Println("- Permanent")
	} else if config.Head == 0xC2 {
		fmt.Println("- Temporary")
	} else {
		fmt.Println("Unkown config head!")
	}
	fmt.Printf("- Addr: %04X\n", config.Addr)
	fmt.Printf("- Parity: %v\n", SerialParity[config.Parity>>6])
	fmt.Printf("- Serial Rate: %v\n", SerialRate[config.Ttlrate>>3])
	fmt.Printf("- Wireless Rate: %vK\n", WirelessRate[config.Wirelessrate])
	fmt.Printf("- Channel: %f\n", float32(425+float32(config.Channel)*0.1))
	fmt.Printf("- Transmission Mode: %v\n", TranmissionMode[config.Transmissionmode>>7])
	fmt.Printf("- IO Mode: %v\n", IoMode[config.Iomode>>6])
	fmt.Printf("- Response Time: %dms\n", int16(1+(config.Wakeuptime>>3))*250)
	fmt.Printf("- FEC: %v\n", config.Fec>>2)
	fmt.Printf("- Transmit Power: %vdBm\n", 30-config.Power*3)

}
