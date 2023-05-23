package main

import (
	"audio-convert/models"
	"audio-convert/services"
	"fmt"
	"os"
)

func getInputs() models.UserInputs {
	var userInputs models.UserInputs

	// fmt.Print("INPUT path: ")
	// fmt.Scanln(&userInputs.InputPath)

	// fmt.Print("OUTPUT path: ")
	// fmt.Scanln(&userInputs.OutputPath)

	fmt.Print("MP3 bitrate: ")
	fmt.Scanln(&userInputs.BitRate)
	if userInputs.BitRate > 320 {
		userInputs.BitRate = 320
	} else if userInputs.BitRate < 128 {
		userInputs.BitRate = 128
	}

	return userInputs
}

// //////////////
func main() {

	// limiter := &cpulimit.Limiter{
	// 	MaxCPUUsage:     20.0,                   // throttle if current cpu usage is over 50%
	// 	MeasureInterval: time.Millisecond * 333, // measure cpu usage in an interval of 333 milliseconds
	// 	Measurements:    3,                      // use the average of the last 3 measurements for cpu usage calculation
	// }
	// limiter.Start()

	userInputs := getInputs()

	f, err := os.Open(services.InputPath)
	if err != nil {
		err := fmt.Errorf("Error: %q", err)
		fmt.Println(err)
		return
	}
	defer f.Close()

	services.HandleReadFiles(f, userInputs)
}
