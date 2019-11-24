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
	trigger := rpio.Pin(7)
	trigger.Output()

	echo := rpio.Pin(11)
	echo.Input()

	trigger.Low()
	fmt.Println("Esperando para medir la distancia")
	time.Sleep(2 * time.Second)
	fmt.Println("Calculando distancia...")
	trigger.High()
	time.Sleep(30000 * time.Nanosecond)
	trigger.Low()

	var startTime = time.Now()
	var endTime = time.Now()

	fmt.Println(rpio.Low)
	fmt.Println(rpio.High)

	for {
		val := echo.Read()
		startTime = time.Now()

		if val == rpio.Low {
			fmt.Printf("value: %v \n", val)
			continue
		}

		break

	}

	for {
		val2 := echo.Read()
		endTime = time.Now()

		if val2 == rpio.High {
			fmt.Printf("second value: %v \n", val2)
			continue
		}

		break

	}
	fmt.Println("bueno jeje")
	duration := endTime.Sub(startTime)
	durationAsInt64 := int64(duration)
	distance := duration.Seconds() * 34300
	distance = distance / 2 //one way travel time
	fmt.Println("distancia======")
	fmt.Printf("Distance : %v | duration: %v | raw: %v \n", distance, duration.Seconds(), durationAsInt64)
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
