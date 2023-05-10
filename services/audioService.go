package services

import (
	"audio-convert/models"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const maxGoroutines = 10

var convertTo = "mp3"
var wg = new(sync.WaitGroup)

var outputPath = "./output"
var InputPath = "./music"

func handleConvert(folderName string, file fs.FileInfo, guard <-chan struct{}, userInputs models.UserInputs) {
	defer wg.Done()

	var pathFile string = InputPath + "/" + file.Name()
	if len(folderName) > 0 {
		pathFile = InputPath + "/" + folderName + "/" + file.Name()
	}
	originalExtension := file.Name()[strings.LastIndex(file.Name(), ".")+1:]
	newFileName := strings.Replace(file.Name(), originalExtension, convertTo, -1)
	newPathFile := outputPath + "/" + newFileName

	err := exec.Command("ffmpeg", "-i", pathFile, "-ab", strconv.Itoa(userInputs.BitRate)+"k", "-map_metadata", "0", "-id3v2_version", "3", newPathFile).Run()
	// err := exec.Command("ffmpeg", "-i", pathFile, "-c:a", "alac", "-c:v", "copy", newPathFile).Run() // use this for  FLAC to ALAC
	if err != nil {
		fmt.Println("Error:", file.Name(), err)
		<-guard
		return
	}

	<-guard
	fmt.Println("OK:", file.Name())
}

func HandleReadFiles(f *os.File, userInputs models.UserInputs) {
	guard := make(chan struct{}, maxGoroutines) // to limit Go-routines

	var files []fs.FileInfo

	folders, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, folder := range folders {
		if folder.IsDir() {
			_folder, err := os.Open(InputPath + "/" + folder.Name())
			if err != nil {
				fmt.Println(err)
				continue
			}

			_files, err := _folder.Readdir(0)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, file := range _files {
				guard <- struct{}{}

				wg.Add(1)
				go handleConvert(folder.Name(), file, guard, userInputs)
			}
		} else {
			files = append(files, folder) // now, folder is a file
		}
	}

	if len(files) > 0 {
		for _, file := range files {
			guard <- struct{}{}

			wg.Add(1)
			go handleConvert("", file, guard, userInputs)
		}
	}

	wg.Wait()
}
