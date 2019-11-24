// Package customss020 allows interfacing with the US020 ultrasonic range finder.
package customss020

import (
	"log"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	pulseDelay  = 30000 * time.Nanosecond
	defaultTemp = 25
)

// Thermometer ok
type Thermometer interface {
	Temperature() (float64, error)
}

type nullThermometer struct {
}

// Temperature ok
func (*nullThermometer) Temperature() (float64, error) {
	return defaultTemp, nil
}

// NullThermometer l
var NullThermometer = &nullThermometer{}

// US020 represents a US020 ultrasonic range finder.
type US020 struct {
	EchoPinNumber, TriggerPinNumber int

	Thermometer Thermometer

	echoPin    rpio.Pin
	triggerPin rpio.Pin

	speedSound float64

	initialized bool
	mu          sync.RWMutex

	Debug bool
}

// New creates a new US020 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(e, t int, thermometer Thermometer) *US020 {
	return &US020{EchoPinNumber: e, TriggerPinNumber: t, Thermometer: thermometer}
}

func (d *US020) setup() (err error) {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err = rpio.Open(); err != nil {
		return
	}

	d.echoPin = rpio.Pin(d.EchoPinNumber)       // ECHO port on the US020
	d.triggerPin = rpio.Pin(d.TriggerPinNumber) // TRIGGER port on the US020

	d.echoPin.Input()
	d.triggerPin.Output()

	if d.Thermometer == nil {
		d.Thermometer = NullThermometer
	}

	if temp, err := d.Thermometer.Temperature(); err == nil {
		d.speedSound = 331.3 + 0.606*temp

		if d.Debug {
			log.Printf("read a temperature of %v, so speed of sound = %v", temp, d.speedSound)
		}
	} else {
		d.speedSound = 340
	}

	d.initialized = true

	return
}

// Distance computes the distance of the bot from the closest obstruction.
func (d *US020) Distance() (distance float64, err error) {
	if err = d.setup(); err != nil {
		return
	}

	if d.Debug {
		log.Print("us020: trigerring pulse")
	}

	// Generate a TRIGGER pulse
	d.triggerPin.High()
	time.Sleep(pulseDelay)
	d.triggerPin.Low()

	if d.Debug {
		log.Print("us020: waiting for echo to go high")
	}

	// Wait until ECHO goes high
	for d.echoPin.Read() == rpio.Low {
	}

	startTime := time.Now() // Record time when ECHO goes high

	if d.Debug {
		log.Print("us020: waiting for echo to go low")
	}

	// Wait until ECHO goes low
	for d.echoPin.Read() == rpio.High {
	}

	duration := time.Since(startTime) // Calculate time lapsed for ECHO to transition from high to low

	// Calculate the distance based on the time computed
	distance = float64(duration.Nanoseconds()) / 10000000 * (d.speedSound / 2)

	return
}

// Close function
func (d *US020) Close() {
	d.echoPin.Output()
	rpio.Close()
}
