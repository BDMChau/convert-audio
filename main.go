package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func handleReadFiles(f *os.File) {
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	convertTo := ".mp3"

	for _, file := range files {
		fileName := "./music/" + file.Name()
		splitName := strings.Split(file.Name(), ".")
		splitName[len(splitName)-1] =
			strings.Replace(splitName[len(splitName)-1], splitName[len(splitName)-1], convertTo, -1)

		newFileName := strings.Join(splitName, "")

		// ffmpeg -i input.mp3 output.flac
		cmd := "ffmpeg -i " + fileName + " " + newFileName

		err = exec.Command(cmd).Run()
		if err != nil {

		}

		fmt.Println(exec.Command("ffmpeg -i ", "D:/School/Web/Go/convert-audio/music/04._Work_it_out.flac", " 04_Work_it_out.mp3").Run())
	}
}

// //////////////
func main() {
	f, err := os.Open("./music")

	if err != nil {
		err := fmt.Errorf("Error: %q", err)
		fmt.Println(err)
		return
	}

	handleReadFiles(f)
}
