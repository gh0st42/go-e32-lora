module go-e32-lora

go 1.16

require (
	github.com/d2r2/go-hd44780 v0.0.0-20181002113701-74cc28c83a3e // indirect
	github.com/d2r2/go-i2c v0.0.0-20191123181816-73a8a799d6bc // indirect
	github.com/d2r2/go-logger v0.0.0-20210606094344-60e9d1233e22 // indirect
	github.com/gh0st42/go-e32-lora/lora v0.0.0
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/weldpua2008/go-dialog v0.0.0-20160901232730-3f70232df8fe // indirect
)

replace github.com/gh0st42/go-e32-lora/lora v0.0.0 => ./lora
