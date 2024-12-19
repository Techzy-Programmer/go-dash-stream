package handler

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

var handler func(float64)

func SubscribeForDistance(distHandler func(float64)) {
	handler = distHandler
	startReadLoop()
}

func startReadLoop() {
	// Initialize periph.io
	if _, err := host.Init(); err != nil {
		fmt.Printf("Failed to initialize periph: %v\n", err)
		return
	}

	// Get handles for trigger and echo pins
	// Adjust these pin numbers according to your wiring
	triggerPin := gpioreg.ByName("GPIO23") // Use the GPIO number you connected TRIG to
	echoPin := gpioreg.ByName("GPIO24")    // Use the GPIO number you connected ECHO to

	if triggerPin == nil || echoPin == nil {
		fmt.Printf("Failed to find pins\n")
		return
	}

	// Configure trigger as output and echo as input
	if err := triggerPin.Out(gpio.Low); err != nil {
		fmt.Printf("Failed to configure trigger pin: %v\n", err)
		return
	}

	if err := echoPin.In(gpio.Float, gpio.NoEdge); err != nil {
		fmt.Printf("Failed to configure echo pin: %v\n", err)
		return
	}

	for {
		distance, err := measureDistance(triggerPin, echoPin)
		if err != nil {
			fmt.Printf("Error measuring distance: %v\n", err)
		} else {
			fmt.Printf("Distance: %.2f cm\n", distance)
		}
		time.Sleep(time.Second)
		handler(distance)
	}
}

func measureDistance(triggerPin gpio.PinOut, echoPin gpio.PinIn) (float64, error) {
	// Ensure trigger is low
	triggerPin.Out(gpio.Low)
	time.Sleep(2 * time.Microsecond)

	// Send trigger pulse
	triggerPin.Out(gpio.High)
	time.Sleep(10 * time.Microsecond)
	triggerPin.Out(gpio.Low)

	// Wait for echo to go high
	startTime := time.Now()
	timeout := startTime.Add(time.Second)
	for echoPin.Read() == gpio.Low {
		if time.Now().After(timeout) {
			return 0, fmt.Errorf("timeout waiting for echo start")
		}
	}
	pulseStart := time.Now()

	// Wait for echo to go low
	for echoPin.Read() == gpio.High {
		if time.Now().After(timeout) {
			return 0, fmt.Errorf("timeout waiting for echo end")
		}
	}
	pulseEnd := time.Now()

	// Calculate distance
	// Speed of sound = 34300 cm/s
	// Distance = (Time × Speed of Sound) ÷ 2
	// Divide by 2 because the sound wave travels to the object and back

	pulseDuration := pulseEnd.Sub(pulseStart)
	distance := float64(pulseDuration.Nanoseconds()) / 1000000000 * 34300 / 2

	return distance, nil
}
