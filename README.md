# go-E32-lora: A E32-TTL-100 Lora Toolkit

This project makes the cheap E32-TTL-100 lora chips accessible through an easy library and some command line tools.


## Installation

```
go get ./...
./build-native.sh
```

For cross compilation to raspi

```
./build-pi.sh
```

Binaries are build to ./out/\<platform>

## Direct Pi Hookup

GPIO pins 17 and 27 (BCM Numbering) must be connected to M0 and M1 on E32. VCC and GND should be connected to 3.3V and GND on Pi as well as RX and TX to the corresponding serial pins on the Pi. Aux is left dangling.

## Serial Port Selection

Default serial device is `/dev/ttyS0` as this is the internal serial port on a raspberry pi.

This can be overwritten by the `LORAPORT` environment variable.

On mac for example:
```
user@mymac$ LORAPORT=/dev/tty.SLAB_USBtoUART ./bcaster 15
```

Or linux with an attached usb serial converter:
```
user@mymac$ LORAPORT=/dev/ttyUSB0 ./bcaster 15
```
## Included tools

### e32config

Quick tool for decoding e32config strings, generating new ones using a dialog gui and changing modes and configs on properly configured devices.

Execute `e32config help` for more information. Beware some feature work only on E32 chips directly connected to a pi.

### bcaster
A simple program to broadcast the hostname and current date-time via lora. Useful for range checks.

Usage: `bcaster [sleeptimeinseconds]`

### monitor
Receives announces from `bcaster` and displays them on the console and in case of an i2c display connected to a pi also outputs to the display.

**WARNING** This code is currently tighlty coupled to the presence of the lcd display, some cleanup needed!