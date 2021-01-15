/*
*  http://www.getmicros.net/a-look-at-the-arduino-nano-33-iot.php
*  https://randomnerdtutorials.com/esp32-mqtt-publish-bme280-arduino/
*
*  tinygo flash -target=arduino-nano33 main.go
*  screen /dev/ttyACM0
 */

package main

import (
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/bme280"
)

func main() {
	time.Sleep(5 * time.Second)

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	machine.I2C0.Configure(machine.I2CConfig{})
	sensor := bme280.New(machine.I2C0)
	sensor.Configure()

	connected := sensor.Connected()
	if !connected {
		println("BME280 not detected")
		led.High()
	}
	println("BME280 detected")

	for {
		led.High()
		temp, _ := sensor.ReadTemperature()
		println("Temperature:", strconv.FormatFloat(float64(temp)/1000, 'f', 2, 64), "°C")
		press, _ := sensor.ReadPressure()
		println("Pressure:", strconv.FormatFloat(float64(press)/100000, 'f', 2, 64), "hPa")
		hum, _ := sensor.ReadHumidity()
		println("Humidity:", strconv.FormatFloat(float64(hum)/100, 'f', 2, 64), "%")
		alt, _ := sensor.ReadAltitude()
		println("Altitude:", alt, "m")
		println("##############################")
		led.Low()

		time.Sleep(2 * time.Second)
	}
}
