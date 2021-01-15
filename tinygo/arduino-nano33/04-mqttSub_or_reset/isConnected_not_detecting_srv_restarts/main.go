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
*  while true; do mosquitto_pub -t "nano33/input" -h xor -m "$(pwgen)" && sleep 1; done
 */

package main

import (
	"fmt"
	"machine"
	"math/rand"
	"time"

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
	subTopic     = "nano33/input"

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

	// Mqtt settings
	opts := mqtt.NewClientOptions()
	opts.AddBroker(server).SetClientID(mqttClientID)
	println("Connecting to MQTT...")
	cl := mqtt.NewClient(opts)
	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		failMessage(token.Error().Error())
	} else {
		println("connected to MQTT")
	}

	mqttSub(cl, subTopic)

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

func mqttOrReset(client mqtt.Client, rstPin machine.Pin) {

	for {
		if client.IsConnected() {
			println("[mqttOrReset] client is connected")
		} else {
			println("[mqttOrReset] about to reset the device")
			time.Sleep(time.Second)
			rstPin.High()
		}
		time.Sleep(time.Second * 2)
	}

}

func subHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[mqttSub: %s] %s \r\n", msg.Topic(), msg.Payload())
	if string(msg.Payload()) == "sample1" {
		fmt.Print("[mqttSub] sample1 message received \r\n")
	}
}

func mqttSub(c mqtt.Client, topic string) {
	go mqttOrReset(c, rstPin)
	if c.IsConnectionOpen() {
		token := c.Subscribe(topic, 1, subHandler)
		token.Wait()
		if token.Error() != nil {
			failMessage(token.Error().Error())
		} else {
			fmt.Printf("[mqttSub] Subscribed to topic: %s \r\n", topic)
		}
		select {}
	}
}
