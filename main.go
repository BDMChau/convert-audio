package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type UserInputs struct {
	inputPath  string
	outputPath string
	bitRate    int
}

var convertTo = ".mp3"
var wg = new(sync.WaitGroup)

const maxGoroutines = 10

func handleConvert(folderName string, file fs.FileInfo, guard <-chan struct{}, userInputs UserInputs) {
	defer wg.Done()

	pathFile := userInputs.inputPath + "/" + folderName + "/" + file.Name()
	originalExtension := file.Name()[strings.LastIndex(file.Name(), ".")+1:]
	newFileName := strings.Replace(file.Name(), originalExtension, "."+convertTo, -1)
	newPathFile := userInputs.outputPath + "/" + newFileName

	err := exec.Command("ffmpeg", "-i", pathFile, "-ab", strconv.Itoa(userInputs.bitRate)+"k", "-map_metadata", "0", "-id3v2_version", "3", newPathFile).Run()
	if err != nil {
		fmt.Println("Error:", file.Name(), err)
		<-guard
		return
	}

	<-guard
	fmt.Println("OK:", file.Name())
	file.Name()
}

func handleReadFiles(f *os.File, guard chan struct{}, userInputs UserInputs) {
	folders, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, folder := range folders {
		if folder.IsDir() {
			_folder, err := os.Open(userInputs.inputPath + "/" + folder.Name())
			if err != nil {
				fmt.Println(err)
				continue
			}

			files, err := _folder.Readdir(0)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, file := range files {
				guard <- struct{}{}

				wg.Add(1)
				go handleConvert(folder.Name(), file, guard, userInputs)

			}

			wg.Wait()
		}
	}
}

func getInputs() UserInputs {
	var userInputs UserInputs

	fmt.Print("INPUT path: ")
	fmt.Scan(&userInputs.inputPath)

	fmt.Print("OUTPUT path: ")
	fmt.Scan(&userInputs.outputPath)

	fmt.Print("MP3 bitrate: ")
	fmt.Scan(&userInputs.bitRate)
	if userInputs.bitRate > 320 {
		userInputs.bitRate = 320
	} else if userInputs.bitRate < 128 {
		userInputs.bitRate = 128
	}

	return userInputs
}

// //////////////
func main() {
	guard := make(chan struct{}, maxGoroutines) // to limit Go-routines

	userInputs := getInputs()

	f, err := os.Open(userInputs.inputPath)

	if err != nil {
		err := fmt.Errorf("Error: %q", err)
		fmt.Println(err)
		return
	}

	handleReadFiles(f, guard, userInputs)
}
