package main

import (
	"fmt"
	"time"

	rfid "github.com/firmom/go-rfid-rc522/rfid"
	rc522 "github.com/firmom/go-rfid-rc522/rfid/rc522"
	"github.com/stianeikeland/go-rpio"
)

const triggerPin int8 = 7
const echoPin int8 = 11

func main() {
	// Ultrasonic sensor
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	trigger := rpio.Pin(triggerPin)
	trigger.Output()

	echo := rpio.Pin(echoPin)
	echo.Input()

	trigger.Low()
	fmt.Println("Esperando para medir la distancia")
	time.Sleep(2 * time.Second)
	fmt.Println("Calculando distancia...")
	trigger.High()
	time.Sleep(1 * time.Microsecond)
	trigger.Low()

	var startTime int64
	var endTime int64

	for echo.Read() == 0 {
		startTime = time.Now().UnixNano()
	}
	for echo.Read() == 1 {
		endTime = time.Now().UnixNano()
	}
	pulseDuration := endTime - startTime

	distance := pulseDuration * 17150
	fmt.Println("distancia======")
	fmt.Println(distance)
	fmt.Println("===========")

	defer rpio.Close()

	// RFID sensor
	product := ""
	fmt.Println("RFID STARTED")
	reader, err := rc522.NewRfidReader()
	if err != nil {
		fmt.Println(err)
		return
	}
	readerChan, err := rfid.NewReaderChan(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	rfidChan := readerChan.GetChan()
	for product == "" {
		select {
		case id := <-rfidChan:
			product = id
		}
	}
	println(product)
}
