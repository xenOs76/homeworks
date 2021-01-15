/*
*  https://store.arduino.cc/arduino-nano-33-iot
*  http://www.getmicros.net/amp/a-look-at-the-arduino-nano-33-iot.php
*  https://www.makerguides.com/arduino-nano/#arduino-nano-33-iot
*  https://tinygo.org/microcontrollers/arduino-nano33-iot/
*  https://www.hackster.io/alankrantas/tinygo-on-arduino-uno-an-introduction-6130f6
*  https://randomnerdtutorials.com/esp32-mqtt-publish-bme280-arduino/
*
*  tinygo flash -target=arduino-nano33 main.go
*  screen /dev/ttyACM0
 */

package main

import (
	"machine"
	"math/rand"
	"strconv"
	"time"

	"tinygo.org/x/drivers/bme280"
	"tinygo.org/x/drivers/net/mqtt"
	"tinygo.org/x/drivers/wifinina"
)

const ssid = "CHANGEME"
const pass = "CHANGEME"
const server = "tcp://192.168.0.24:1883"
const rstPin machine.Pin = machine.D2

var (

	// these are the default pins for the Arduino Nano33 IoT.
	uart = machine.UART2
	tx   = machine.NINA_TX
	rx   = machine.NINA_RX
	spi  = machine.NINA_SPI

	// this is the ESP chip that has the WIFININA firmware flashed on it
	adaptor = &wifinina.Device{
		SPI:   spi,
		CS:    machine.NINA_CS,
		ACK:   machine.NINA_ACK,
		GPIO0: machine.NINA_GPIO0,
		RESET: machine.NINA_RESETN,
	}

	console = machine.UART0

	mqttClientID = "nano33"
	displayTopic = "coconut/scrollphat"
	tempTopic    = "nano33/temp"
	humTopic     = "nano33/hum"
	pressTopic   = "nano33/press"

	maxConnErrCount int = 5
)

func main() {
	time.Sleep(3000 * time.Millisecond)

	// device reset settings
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin.High()

	// wifi connection settings
	uart.Configure(machine.UARTConfig{TX: tx, RX: rx})
	rand.Seed(time.Now().UnixNano())
	spi.Configure(machine.SPIConfig{
		Frequency: 8 * 1e6,
		SDO:       machine.NINA_SDO,
		SDI:       machine.NINA_SDI,
		SCK:       machine.NINA_SCK,
	})
	adaptor.Configure()
	connectToAP()

	// bme280 sensor settings
	machine.I2C0.Configure(machine.I2CConfig{})
	sensor := bme280.New(machine.I2C0)
	sensor.Configure()

	// Mqtt settings
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server).SetClientID(mqttClientID)
	println("Connecting to MQTT...")
	cl := mqtt.NewClient(opts)
	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		failMessage(token.Error().Error())
	}

	connErrCount := 0
	for i := 0; ; i++ {
		if !sensor.Connected() {
			resetOnConnErrors(rstPin, connErrCount, maxConnErrCount)
			connErrCount++
		} else {
			temp, _ := sensor.ReadTemperature()
			tempMsg := strconv.FormatFloat(float64(temp)/1000, 'f', 2, 64)
			tempMsgDisplay := "  Temp: " + tempMsg + "°C    "
			mqttPub(cl, tempTopic, tempMsg)
			mqttPub(cl, displayTopic, tempMsgDisplay)

			hum, _ := sensor.ReadHumidity()
			humMsg := strconv.FormatFloat(float64(hum)/100, 'f', 2, 64)
			humMsgDisplay := "  Hum: " + humMsg + "%    "
			mqttPub(cl, humTopic, humMsg)
			mqttPub(cl, displayTopic, humMsgDisplay)

			time.Sleep(10 * time.Second)
		}
	}

}

// connect to access point
func connectToAP() {
	connErrCount := 0
	time.Sleep(2 * time.Second)
	println("Connecting to " + ssid)
	adaptor.SetPassphrase(ssid, pass)
	for st, _ := adaptor.GetConnectionStatus(); st != wifinina.StatusConnected; {
		println("Connection status: " + st.String())
		time.Sleep(1 * time.Second)
		st, _ = adaptor.GetConnectionStatus()
		resetOnConnErrors(rstPin, connErrCount, maxConnErrCount)
		connErrCount++
	}
	println("Connected.")
	time.Sleep(2 * time.Second)
	ip, _, _, err := adaptor.GetIP()
	for ; err != nil; ip, _, _, err = adaptor.GetIP() {
		println(err.Error())
		time.Sleep(1 * time.Second)
		resetOnConnErrors(rstPin, connErrCount, maxConnErrCount)
		connErrCount++
	}
	println(ip.String())
}

func mqttPub(c mqtt.Client, topic string, msg string) {

	data := []byte(msg)
	token := c.Publish(topic, 0, false, data)
	token.Wait()
	if err := token.Error(); err != nil {
		switch t := err.(type) {
		case wifinina.Error:
			println(t.Error(), "attempting to reconnect")
			println("mqttPub wifinina err 1")
			if token := c.Connect(); token.Wait() && token.Error() != nil {
				println("mqttPub wifinina err 2")
				failMessage(token.Error().Error())
			}
		default:
			println("mqttPub default err case")
			println(err.Error())
		}
	} else {
		println("[mqttPub] topic:", topic, " msg:", msg)
	}

}

// // Returns an int >= min, < max
// func randomInt(min, max int) int {
// 	return min + rand.Intn(max-min)
// }

// // Generate a random string of A-Z chars with len = l
// func randomString(len int) string {
// 	bytes := make([]byte, len)
// 	for i := 0; i < len; i++ {
// 		bytes[i] = byte(randomInt(65, 90))
// 	}
// 	return string(bytes)
// }

func failMessage(msg string) {
	connErrCount := 0
	for {
		println(msg)
		resetOnConnErrors(rstPin, connErrCount, maxConnErrCount)
		connErrCount++
	}
}

// Resets the device if a counter reached its max allowed value
func resetOnConnErrors(RstPin machine.Pin, connErrCount int, maxConnErrCount int) {
	if connErrCount >= maxConnErrCount {
		println("reached maxConnErrCount, resetting device in 3 seconds")
		time.Sleep(3 * time.Second)
		rstPin.Low()
	}
}
