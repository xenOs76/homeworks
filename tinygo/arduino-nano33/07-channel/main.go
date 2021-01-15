/*
*
*
*  tinygo flash -target arduino-nano33 main.go
 */

package main

import (
	"fmt"
	"time"
)

func getTime(outCh chan int) {
	start := int(time.Now().Unix())
	fmt.Printf("[getTime] start at %d", int(start))
	for {
		time.Sleep(1 * time.Second)
		now := int(time.Now().Unix())
		fmt.Printf("[getTime] now is %d \n\r", int(now))
		outCh <- int(now)
	}
}

func showTime(inCh chan int) {

	for {
		t := <-inCh
		fmt.Printf("[showTime] now is %d \n\r", t)
	}
}

func loopAndPrint() {
	for {
		fmt.Printf("[loopAndPring] hello \r\n")
		time.Sleep(700 * time.Millisecond)
	}
}

func main() {
	time.Sleep(3 * time.Second)
	fmt.Printf("[main] Channel example starting...\r\n")

	mainStart := int(time.Now().Unix())
	fmt.Printf("[main] start at %d \r\n", int(mainStart))

	timeExchCh := make(chan int, 1)

	fmt.Printf("[main] invoking getTime...\n\r")
	go getTime(timeExchCh)

	fmt.Printf("[main] invoking showTime...\n\r")
	go showTime(timeExchCh)

	fmt.Printf("[main] invoking loopAndPrint...\n\r")
	go loopAndPrint()

	select {}
}
