//
// To connect a Neopixel strip:
// https://learn.adafruit.com/led-corset-with-circuit-playground-and-makecode/neopixel-strips

// https://www.makerguides.com/arduino-nano/#arduino-nano-33-iot
//
// tinygo flash -target arduino-nano33 main.go
//

package main

import (
	"image/color"
	"machine"
	"math"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

var leds [60]color.RGBA
var neopixels = len(leds)
var ledColor = [3]byte{0xff, 0xff, 0x00}
var neo machine.Pin = machine.D3
var ws ws2812.Device
var ledsOn bool = true

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
			r8 := uint8(r)
			g8 := uint8(g)
			b8 := uint8(b)
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

func cycleLeds(ws ws2812.Device) {
	for {
		if ledsOn {
			println("starting new cycle...")
			rainbowCycle(ws)
			ledsOff(ws)
			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	time.Sleep(3 * time.Second)
	println("Rainbow starting")

	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws = ws2812.New(neo)

	go cycleLeds(ws)

	for {

		println("printing after cycleLeds")
		time.Sleep(1 * time.Second)
	}
}
