package main

import (
	"fmt"
	"math"
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
	time.Sleep(2 * time.Millisecond)
	fmt.Println("Calculando distancia...")
	trigger.High()
	time.Sleep(1 * time.Microsecond)
	trigger.Low()

	var pulse_start_time time
	var pulse_end_time time

	for echo.Read() == 0 {
		pulse_start_time = time.Now()
	}
	for echo.Read() == 1 {
		pulse_end_time = time.Now()
	}
	pulse_duration := pulse_end_time - pulse_start_time

	distance := math.Round(pulse_duration*17150, 2)
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
