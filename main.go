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

func selectService() int8 {
	var serviceNumber int8 = 0

	fmt.Println("Enter service: ")
	fmt.Println("1: Music Converter")
	fmt.Println("2: Web Crawler")
	fmt.Scanln(&serviceNumber)

	if serviceNumber == 0 {
		fmt.Println("Error, invalid input!")
		return selectService()
	}

	return serviceNumber
}

// //////////////
func main() {
	switch selectService() {
	case 1:
		userInputs := getInputs()

		f, err := os.Open(services.InputPath)
		if err != nil {
			err := fmt.Errorf("Error: %q", err)
			fmt.Println(err)
			return
		}
		defer f.Close()

		services.HandleReadFiles(f, userInputs)

	case 2:

		services.Crawler()
	}
}
