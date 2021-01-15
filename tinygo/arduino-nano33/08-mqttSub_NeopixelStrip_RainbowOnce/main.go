//
// https://learn.adafruit.com/led-corset-with-circuit-playground-and-makecode/neopixel-strips
// https://www.makerguides.com/arduino-nano/#arduino-nano-33-iot
// https://github.com/solarwinds/tinygo-lessons/blob/5c01dc4fcda0/mqtt-consumer/main.go
//
// tinygo flash -target arduino-nano33 main.go
// screen /dev/ttyACM0
// mosquitto_pub -h 192.168.0.24 -t "nano33/input" -m "1"
//

package main

import (
	"fmt"
	"image/color"
	"machine"
	"math"
	"math/rand"
	"time"

	"tinygo.org/x/drivers/net/mqtt"
	"tinygo.org/x/drivers/wifinina"
	"tinygo.org/x/drivers/ws2812"
)

const ssid = "CHANGEME"
const pass = "CHANGEME"
const server = "tcp://192.168.0.24:1883"
const rstPin machine.Pin = machine.D2
const neo machine.Pin = machine.D3

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
	subTopic     = "nano33/input"
	hbTopic      = "nano33/ping"

	maxConnErrCount int = 5

	mqttToNeo chan int = make(chan int)

	leds       [60]color.RGBA
	neopixels  = len(leds)
	ledColor   = [3]byte{0xff, 0xff, 0x00}
	ws         ws2812.Device
	ledsStatus int = 0
)

func rainbowWheel(pos int) (int, int, int) {

	var r, g, b int = 0, 0, 0

	if (pos < 0) || (pos > 255) {
		r, g, b = 0, 0, 0
	} else if pos < 85 {
		r = pos * 3
		g = 255 - pos*3
		b = 0
	} else if pos < 170 {
		pos = pos - 85
		r = 255 - pos*3
		g = 0
		b = pos * 3
	} else {
		pos = pos - 170
		r = 0
		g = pos * 3
		b = 255 - pos*3
	}
	return r, g, b
}

func rainbowCycle(ws ws2812.Device) error {
	for j := 0; j <= 255; j++ {
		for i := 0; i < neopixels; i++ {
			it1 := float64(i * 256)
			it2 := math.Floor(it1 / float64(neopixels))
			neoIndex := int(it2) + j
			r, g, b := rainbowWheel(neoIndex & 255)
			r8, g8, b8 := uint8(r), uint8(g), uint8(b)
			leds[i] = color.RGBA{R: r8, G: g8, B: b8}
		}

		ws.WriteColors(leds[:])
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}

func ledsOff(ws ws2812.Device) error {
	for i := 0; i < neopixels; i++ {
		leds[i] = color.RGBA{R: 0, G: 0, B: 0}
	}
	ws.WriteColors(leds[:])
	return nil
}

func cycleLeds(ws ws2812.Device, ledsStatus int) {
	println("[cycleLeds] before loop...")
	if ledsStatus == 1 {
		println("[cycleLeds] starting new cycle...")
		rainbowCycle(ws)
		ledsOff(ws)
	}
}

func main() {
	time.Sleep(3 * time.Second)
	println("[main] NeopixelStrip + mqttSub starting...")

	// Neopixels strip settings
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws = ws2812.New(neo)
	ledsOff(ws)

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
	println("[main] Connecting to MQTT...")
	cl := mqtt.NewClient(opts)
	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		failMessage(token.Error().Error())
	} else {
		println("[main] connected to MQTT, starting heart beat check...")
		go mqttPubOrReset(cl, hbTopic)
	}

	mqttSub(cl, subTopic)

	for {
		ledsStatus = <-mqttToNeo
		fmt.Printf("[main] ledsStatus now is %d \n\r", ledsStatus)
		go cycleLeds(ws, ledsStatus)
	}

}

// connect to access point
func connectToAP() {
	connErrCount := 0
	time.Sleep(2 * time.Second)
	println("[ConnectToAp] connecting to " + ssid)
	adaptor.SetPassphrase(ssid, pass)
	for st, _ := adaptor.GetConnectionStatus(); st != wifinina.StatusConnected; {
		println("[ConnectToAp] connection status: " + st.String())
		time.Sleep(1 * time.Second)
		st, _ = adaptor.GetConnectionStatus()
		resetOnConnErrors(rstPin, connErrCount, maxConnErrCount)
		connErrCount++
	}
	println("[ConnectToAp] connected.")
	time.Sleep(2 * time.Second)
	ip, _, _, err := adaptor.GetIP()
	for ; err != nil; ip, _, _, err = adaptor.GetIP() {
		println(err.Error())
		time.Sleep(1 * time.Second)
		resetOnConnErrors(rstPin, connErrCount, maxConnErrCount)
		connErrCount++
	}
	println("[connectToAp] this device ip is", ip.String())
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

func subHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[mqttSub: %s] %s \r\n", msg.Topic(), msg.Payload())
	if string(msg.Payload()) == "1" {
		fmt.Print("[mqttSub] new ledStatus received \r\n")
		ledsStatus = 1
	} else {
		ledsStatus = 0
	}
	mqttToNeo <- ledsStatus
}

func mqttSub(c mqtt.Client, topic string) {
	if c.IsConnectionOpen() {
		token := c.Subscribe(topic, 1, subHandler)
		token.Wait()
		if token.Error() != nil {
			failMessage(token.Error().Error())
		} else {
			fmt.Printf("[mqttSub] Subscribed to topic: %s \r\n", topic)
		}
	}
}

func mqttPub(c mqtt.Client, topic string, msg string) {
	data := []byte(msg)
	token := c.Publish(topic, 0, false, data)
	token.Wait()
	if err := token.Error(); err != nil {
		switch t := err.(type) {
		case wifinina.Error:
			println(t.Error(), "attempting to reconnect")
			if token := c.Connect(); token.Wait() && token.Error() != nil {
				failMessage(token.Error().Error())
			}
		default:
			println(err.Error())
		}
	} else {
		println("[mqttPub] topic:", topic, " msg:", msg)
	}
}

func mqttPubOrReset(c mqtt.Client, topic string) {

	for {
		now := fmt.Sprintf("%v", time.Now().Unix())
		mqttPub(c, topic, now)
		time.Sleep(time.Second)
	}
}
