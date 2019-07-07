package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

// Reading represents a single reading by the device
type Reading struct {
	DeviceID string  `json:"device_id"`
	Value    float64 `json:"value"`
}

func main() {
	devID := "RaspberryPi-1" // device ID
	// var diff float64

	var err error

	host := flag.String("host", "https://earthworks.islandcivil.com", "hostname to connect to")
	flag.Parse()

	// set up the Raspberry Pi/Analog to digital converter (ADS1015)
	board := raspi.NewAdaptor()
	ads1015 := i2c.NewADS1015Driver(board)
	ads1015.DefaultGain, _ = ads1015.BestGainForVoltage(3.3)

	// work is a function that collects readings and sends them to a server to be stored
	work := func() {
		gobot.Every(1*time.Minute, func() {

			r, _ := ads1015.ReadWithDefaults(1)
			err := report(context.Background(), host, Reading{Value: r, DeviceID: devID})
			if err != nil {
				log.Print("Error sending data to server")
			} else {
				log.Print("Response from server: ok") // response had no errors
			}
		})
	}

	robot := gobot.NewRobot("thermBot",
		[]gobot.Connection{board},
		[]gobot.Device{ads1015},
		work,
	)

	err = robot.Start()
	if err != nil {
		log.Println(err)
	}

}

func report(ctx context.Context, host *string, v Reading) error {

	jsonReq, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		return err
	}

	resp, err := http.Post(*host+"/api/v1/instrument_data", "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	return err
}
