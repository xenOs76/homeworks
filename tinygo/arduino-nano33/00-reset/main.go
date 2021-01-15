/*
*  https://store.arduino.cc/arduino-nano-33-iot
*  http://www.getmicros.net/amp/a-look-at-the-arduino-nano-33-iot.php
*  https://www.makerguides.com/arduino-nano/#arduino-nano-33-iot
*  https://tinygo.org/microcontrollers/arduino-nano33-iot/
*  https://www.hackster.io/alankrantas/tinygo-on-arduino-uno-an-introduction-6130f6
*
*  tinygo flash -target=arduino-nano33 main.go
*  screen /dev/ttyACM0
 */

package main

import (
	"machine"
	"time"
)

func main() {
	time.Sleep(time.Second * 3)
	println("Reset example")
	var RST machine.Pin = machine.D2
	RST.Configure(machine.PinConfig{Mode: machine.PinOutput})
	RST.High()

	println("This board is going to reset in ... seconds")
	for i := 5; i > 0; i-- {
		println(i)
		time.Sleep(time.Second * 1)
	}

	RST.Low()
}
